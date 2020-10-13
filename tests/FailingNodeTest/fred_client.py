import grpc

import client_pb2
import client_pb2_grpc


class FredClient:
    def __init__(self, addr: str):
        with open('certs/client.crt') as f:
            crt = f.read().encode()

        with open('certs/client.key') as f:
            key = f.read().encode()

        with open('certs/ca.crt') as f:
            ca = f.read().encode()

        creds = grpc.ssl_channel_credentials(ca, key, crt)
        channel = grpc.secure_channel(addr, creds)
        self.stub = client_pb2_grpc.ClientStub(channel)

    def create_keygroup(self, name: str, mutable: bool, expiry: int):
        r = client_pb2.CreateKeygroupRequest()
        r.keygroup = name
        r.mutable = mutable
        r.expiry = expiry
        return self.stub.CreateKeygroup(r)

    def add_replica(self, keygroup: str, nodeId: str, expiry: int):
        r = client_pb2.AddReplicaRequest()
        r.keygroup = keygroup
        r.nodeId = nodeId
        r.expiry = expiry
        return self.stub.AddReplica(r)

    def read(self, keygroup: str, id: str):
        r = client_pb2.ReadRequest()
        r.keygroup = keygroup
        r.id = id
        return self.stub.Read(r)

    def update(self, keygroup: str, id: str, data: str):
        r = client_pb2.UpdateRequest()
        r.keygroup = keygroup
        r.id = id
        r.data = data
        return self.stub.Update(r)
