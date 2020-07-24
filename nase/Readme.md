# FReD Naming Service

The naming service uses etcd.io for distributed storage.

The provided docker-compose starts a local instance of etcd.io. To access the CLI enter the docker container `docker exec -it etcd-1 /bin/sh
` and execute `etcdctl`

Clients can reach etcd on port 2379, peers can reach it on port 2380

Prefix Search is also possible, see implementation in nameservice

## Data Representation

Etcd is a key-value store that supports a range parameter for GET request, which could be used to store additional information in the key.

We need to store:

- What Keygroups exist
- All the replicas of a keygroup (and their rights (read/write) ??)
- The current status of a keygoup (running, thombstoned)

Key `kg-[keygroupname]-node-[nodeID]` has: Current status of the replica on this node (initializing, full, readOnly, ...)

Key `kg-[keygroupname]-status` has: String with current status

See if kg exists: Get range `kg-[keygroupname]-`

Get all Replicas: Get range `kg-[keygroupname]-node-`

Every Node will also store information about itself in `node-[nodeid]-status` and `node-nodeid-adress` (format: `ip:port`)

__Problem__: Keygroup names can't end in `-node` or `-status` (could also use other delimitor)

