import argparse
from os import read
import typing
import grpc
import sys
import random
import json
import time

import proto.middleware.middleware_pb2 as middleware_pb2
import proto.middleware.middleware_pb2_grpc as middleware_pb2_grpc

KEYGROUPNAME = "forumkeygroup"
FORUMKEY = "forumkey"
UPDATE_PERCENTAGE = 10
SLEEP_TIME = 1

class Client():
    def __init__(
        self,
        id: str,
        node_id: str,
        keygroupname: str,
        forumkey: str,
        stub: middleware_pb2_grpc.MiddlewareStub,
    ):
        self.id = id
        self.client = stub
        self.keygroupname = keygroupname
        self.forumkey = forumkey

        self.create_keygroup(self.keygroupname, node_id)

    def create_keygroup(self, keygroup: str, node: str) -> None:
        # try to create the keygroup and add an empty set as a value
        # if someone else did so before, that's ok
        try:
            r = middleware_pb2.CreateKeygroupRequest()

            r.keygroup = keygroup
            r.mutable = True
            r.firstNodeId = node

            print("\033[93mtrying to create a keygroup...\033[0m")
            self.client.CreateKeygroup(r)
            r = middleware_pb2.UpdateRequest()

            r.keygroup = keygroup
            r.id = self.forumkey
            r.data = json.dumps(list(set()))

            print("\033[93mtrying to add first item...\033[0m")
            self.client.Update(r)
            return

        except grpc.RpcError as rpc_error:
            print(rpc_error.details())

        r = middleware_pb2.AddReplicaRequest()

        r.keygroup = keygroup
        r.nodeId = node

        try:
            print("\033[93mtrying to add our node as replica to keygroup...\033[0m")
            self.client.AddReplica(r)
            return
        except grpc.RpcError as rpc_error:
            print(rpc_error.details())

    def update(self, keygroup: str, id: str, value: typing.Set[str]) -> None:
        r = middleware_pb2.UpdateRequest()

        r.keygroup = keygroup
        r.id = id
        r.data = json.dumps(sorted(value))

        self.client.Update(r)

    def add(self, keygroup: str, id: str, value: str) -> None:
        # here we try to add a new value to the items
        while True:
            # do an update
            # 1. read current values
            vals: typing.List[typing.Set[str]] = []

            # try to read values until we have at least one version
            while len(vals) == 0:
                try:
                    vals = self.read(keygroup, id)
                except grpc.RpcError:
                    time.sleep(0.01)

            # 2. add a new value
            # if we have more than one version, we need to merge first
            read_vals: typing.Set[str] = set()

            for v in vals:
                read_vals = read_vals.union(v)

            new_vals = read_vals

            # then we can also add our new version
            new_vals.add(value)

            # 3. write the new values
            # and try to write
            # if our write is not accepted, we need to start from the top
            try:
                self.update(keygroup, id, new_vals)
                print("\033[92mupdate success!\033[0m")
                return
            except grpc.RpcError as rpc_error:
                time.sleep(0.01)
                print("\033[91merror updating: %s\033[0m" % rpc_error.details())

    def read(self, keygroup: str, id: str) -> typing.List[typing.Set[str]]:
        r = middleware_pb2.ReadRequest()

        r.keygroup = keygroup
        r.id = id

        data = self.client.Read(r)

        values: typing.List[typing.Set[str]] = []

        for item in data.items:
            values.append(set(json.loads(item.val)))

        return values

    def run_tests(self, ops: int, sleep_time: int=100, update_percentage: int=10) -> int:

        mrcErrors = 0
        rywcErrors = 0

        seen: typing.Set[str] = set()
        written: typing.Set[str] = set()

        for x in range(1, ops):
            if random.randint(0, 99) < update_percentage:
                new_entry = self.id + "-" + str(x)
                print("\033[93m%d: adding %s\033[0m" % (x, new_entry))
                # seen.add(new_entry)
                written.add(new_entry)

                self.add(self.keygroupname, self.forumkey, new_entry)
                print("\033[92m%d: ADD success\033[0m" % x)

            # do a read
            vals: typing.List[typing.Set[str]] = []

            # if more than one set was seen, try to merge them
            while len(vals) != 1:
                vals = []
                while len(vals) == 0:
                    try:
                        vals = self.read(self.keygroupname, self.forumkey)
                    except grpc.RpcError:
                        print("\033[93m%d: error reading, trying again\033[0m" % x)
                        time.sleep(0.01)

                print("\033[93m%d: got %d different versions, trying to merge...\033[0m" % (x, len(vals)))

                # make a new set of values by merging all of the seen ones
                read_vals: typing.Set[str] = set()

                for v in vals:
                    read_vals = read_vals.union(v)

                # try to put that upated list
                try:
                    self.update(self.keygroupname, self.forumkey, read_vals)

                    # continue with our merged version
                    vals = [read_vals]
                    print("\033[92m%d: merged %s into %s\033[0m" % (x, str(vals), str(sorted(read_vals))))
                    break

                except grpc.RpcError:
                    # merge was not succesful? then we had read outdated stuff
                    print("\033[91m%d: merge was not successful, trying again...\033[0m" % x)
                    time.sleep(0.01)

            read_vals = vals[0]

            if not seen.issubset(read_vals):
                mrcErrors += 1
                print("\033[91m%d: MRC error reading: missing %s from updates %s (instead have %s)\033[0m" % (x, str(sorted(seen.difference(read_vals))), str(sorted(read_vals)), str(sorted(read_vals.difference(seen)))))
                time.sleep(10)
                exit()
            else:
                print("\033[92m%d: MRC success\033[0m" % x)

            if not written.issubset(read_vals):
                rywcErrors += 1
                print("\033[91m%d: RYWC error reading: missing %s from updates %s\033[0m" % (x, str(sorted(written.difference(read_vals))), str(sorted(read_vals))))
                time.sleep(10)
                exit()
            else:
                print("\033[92m%d: RYWC success\033[0m" % x)

            seen = seen.union(read_vals)

            time.sleep(sleep_time / 1000)

        print("\033[9" + ("1" if mrcErrors > 0 else "2") + "mtotal mrc errors\033[0m: %d/%d" % (mrcErrors, ops))
        print("\033[9" + ("1" if rywcErrors > 0 else "2") + "mtotal rywc errors\033[0m: %d/%d" % (rywcErrors, ops))

        return mrcErrors + rywcErrors

if __name__ == "__main__":
    parser = argparse.ArgumentParser()

    parser.add_argument("--id", help="id for this client", required=True, type=str)
    parser.add_argument("--ops", help="number of operations to perform", required=True, type=int)
    parser.add_argument("--host", help="host where the middleware is running", required=True, type=str)
    parser.add_argument("--node-id", help="node to use", required=True, type=str)
    parser.add_argument("--cert", help="certificate to use for the middleware connection", required=True, type=str)
    parser.add_argument("--key", help="key to use for middleware connection", required=True, type=str)
    parser.add_argument("--ca", help="certificate of the ca", required=True, type=str)

    args = parser.parse_args()

    print("ID: %s" % args.id)
    print("ops: %d" % args.ops)
    print("host: %s" % args.host)
    print("node: %s" % args.node_id)
    print("cert: %s" % args.cert)
    print("key: %s" % args.key)
    print("ca: %s" % args.ca)

    with open(args.cert) as f:
        crt = f.read().encode()

    with open(args.key) as f:
        key = f.read().encode()

    with open(args.ca) as f:
        ca = f.read().encode()

    creds = grpc.ssl_channel_credentials(ca, key, crt)
    channel = grpc.secure_channel(args.host, creds)

    client = Client(args.id, args.node_id, KEYGROUPNAME, FORUMKEY, middleware_pb2_grpc.MiddlewareStub(channel))

    print("running experiments\n")
    sys.exit(client.run_tests(ops=args.ops, sleep_time=SLEEP_TIME, update_percentage=UPDATE_PERCENTAGE))
