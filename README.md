# FReD: Fog Replicated Data

[![pipeline status](https://git.tu-berlin.de/mcc-fred/fred/badges/main/pipeline.svg)](https://git.tu-berlin.de/mcc-fred/fred/-/commits/main)
[![coverage report](https://git.tu-berlin.de/mcc-fred/fred/badges/main/coverage.svg)](https://git.tu-berlin.de/mcc-fred/fred/-/commits/main)
[![License MIT](https://img.shields.io/badge/License-MIT-brightgreen.svg)](https://img.shields.io/badge/License-MIT-brightgreen.svg)
[![Go Report Card](https://goreportcard.com/badge/git.tu-berlin.de/mcc-fred/fred)](https://goreportcard.com/report/git.tu-berlin.de/mcc-fred/fred)
[![Godoc](http://img.shields.io/badge/go-documentation-blue.svg)](https://pkg.go.dev/git.tu-berlin.de/mcc-fred/fred)

**FReD** is a distributed middleware for **F**og **Re**plicated **D**ata.
It abstracts data management for fog-based applications by grouping data into _keygroups_, each keygroup a set of key-value pairs that can be managed independently.
Applications have full control over keygroup replication: replicate your data where you need it.

FReD is maintained by [Tobias Pfandzelter, Trever Schirmer, and Nils Japke of the Mobile Cloud Computing research group at Technische Universität Berlin and Einstein Center Digital Future in the scope of the FogStore project](https://www.mcc.tu-berlin.de/).
Funded by the Deutsche Forschungsgemeinschaft (DFG, German Research Foundation) -- 415899119.

FReD is open-source software and contributions are welcome.
All contributions should be submitted as merge requests on the [main repository on the TU Berlin GitLab](https://git.tu-berlin.de/mcc-fred/fred) and are subject to review by the maintainers.
Check out the [Contributing](#contributing) section for more information.

## Architecture

A FReD deployment comprises a FReD node on all or a subset of fog nodes in the available fog infrastructure.
This can range from edge nodes all the way to the cloud, different node sizes are possible.

![fred architecture](./docs/assets/architecture.png)

Each FReD node consists of a number of machines running the `fred` software and storage backend.
If a node has multiple servers, a load balancer distributes requests among the nodes.
The storage backend can be a cloud database like DynamoDB, a dedicated database server, or embedded on a single machine with the `fred` software.

Additionally, a centralized "Naming Service" ("NaSe") based on `etcd` keeps track of system configuration and helps nodes find each other (a decentralized version is in progress.)

Clients interact with the FReD system through the _Application Level Extension to Allow Node Discovery and Replica Appointment_ (_ALExANDRA_) middleware.
Applications can create their keygroups to store data as needed and instruct FReD to replicate individual keygroups to where that data is needed.

Additionally, Trigger Nodes can be set up for every keygroup.
Trigger nodes are sent all updates for all data items in a particular keygroup.
This can be useful to transform data to write it back into FReD, event processing, and more.

## Getting Started

The smallest possible FReD deployment comprises a single `etcd` node as a NaSe and a `fred` node with integrated storage.
This is a useful starting point to understand how FReD works but does not offer any useful functionality as data can only be replicated to a single node.

### Minimum FReD Deployment

The simplest way to get started is to use Docker containers to deploy the needed software.

If you're not familiar with Docker yet, take the time to learn about it [here](https://www.docker.com/101-tutorial).
Alternatively, you are of course able to run every component manually.

#### Certificates

Before deploying `etcd` and `fred` software, certificates have to be created to secure the connection between all components and authenticate the individual services.
We provide tooling to create new certificates.
More information about authentication and authorization with certificates can be found further down in this introduction.
This guide is adapted from [scriptcrunch.com](https://scriptcrunch.com/create-ca-tls-ssl-certificates-keys/).

##### Certificate Authority

First, a certificate authority (CA) has to be created to issue new certificates.

```bash
# generate a CA private key
openssl genrsa -out ca.key 2048

# generate a CA certificate
openssl req -x509 -new -nodes \
     -key ca.key -sha512 \
     -days 1825 -out ca.crt

$ You are about to be asked to enter information that will be incorporated
$ into your certificate request.
$ What you are about to enter is what is called a Distinguished Name or a DN.
$ There are quite a few fields but you can leave some blank
$ For some fields there will be a default value,
$ If you enter '.', the field will be left blank.
$ -----
$ Country Name (2 letter code) []:DE
$ State or Province Name (full name) []:Berlin
$ Locality Name (eg, city) []:Berlin
$ Organization Name (eg, company) []:MCC
$ Organizational Unit Name (eg, section) []:FRED
$ Common Name (eg, fully qualified host name) []:
$ Email Address []:
```

You must distribute the `ca.crt` public certificate with all services so its issued certificates can be verified.
However, you must keep the private key `ca.key` private at all times!
**Access to your private key allows anyone to issue arbitrary certificates from your certificate authority!**

##### Generating Certificates

With your CA private key in hand you may now generate certificates for `etcd`, `fred`, storage server, client, and any other software components.
Use the included `gen-cert.sh` script like this: `gen-cert.sh {NAME} {IP}`, where `{NAME}` is the name for your certificate and `{IP}` is the IP address of your service.
The name you assign to the certificate for a FReD node should be the ID of the node.
It is the common name (CN) of this certificate.

Although not strictly required, it is recommended to generate unique certificates for every software component.
In this example case, the following commands are required:

```bash
./gen-cert.sh etcdnase 172.26.1.1
./gen-cert.sh fredNodeA 172.26.1.2
./gen-cert.sh alexandra 172.26.1.3
./gen-cert.sh fredClient 172.26.1.4
```

#### Network

If you run this example in Docker, you must first create a simple network for the individual services to talk to each other:

```bash
docker network create fredwork --gateway 172.26.0.1 --subnet 172.26.0.0/16
```

#### NaSe

To start a simple `etcd` instance in Docker with our certificates mounted as volumes, you can use this command:

```bash
docker pull gcr.io/etcd-development/etcd:v3.5.7
docker run -d \
-v $(pwd)/etcdnase.crt:/cert/etcdnase.crt \
-v $(pwd)/etcdnase.key:/cert/etcdnase.key \
-v $(pwd)/ca.crt:/cert/ca.crt \
--network=fredwork \
--ip=172.26.1.1 \
gcr.io/etcd-development/etcd:v3.5.7 \
etcd --name s-1 \
--data-dir /tmp/etcd/s-1 \
--listen-client-urls https://172.26.1.1:2379 \
--advertise-client-urls https://172.26.1.1:2379 \
--listen-peer-urls http://172.26.1.1:2380 \
--initial-advertise-peer-urls http://172.26.1.1:2380 \
--initial-cluster s-1=http://172.26.1.1:2380 \
--initial-cluster-token tkn \
--initial-cluster-state new \
--cert-file=/cert/etcdnase.crt \
--key-file=/cert/etcdnase.key \
--client-cert-auth \
--trusted-ca-file=/cert/ca.crt
```

This runs a `etcd` cluster with a single machine (hence no fault tolerance), gives it the name `s-1` (arbitrary), adds the certificate files to the container, enables client certificate authentication, and connects the container to the `fredwork` network.

#### FReD

Running the `fred` software is similar to running `etcd`, we make container images available on our GitLab container registry.

```bash
docker pull git.tu-berlin.de:5000/mcc-fred/fred/fred:latest
docker run -d \
-v $(pwd)/fredNodeA.crt:/cert/fredNodeA.crt \
-v $(pwd)/fredNodeA.key:/cert/fredNodeA.key \
-v $(pwd)/ca.crt:/cert/ca.crt \
--network=fredwork \
--ip=172.26.1.2 \
git.tu-berlin.de:5000/mcc-fred/fred/fred:latest \
--log-level info \
--handler dev \
--nodeID fredNodeA \
--host 172.26.1.2:9001 \
--peer-host 172.26.1.2:5555 \
--adaptor badgerdb \
--badgerdb-path ./db \
--nase-host https://172.26.1.1:2379 \
--nase-cert /cert/fredNodeA.crt \
--nase-key /cert/fredNodeA.key \
--nase-ca /cert/ca.crt \
--trigger-cert /cert/fredNodeA.crt \
--trigger-key /cert/fredNodeA.key \
--trigger-ca /cert/ca.crt \
--peer-cert /cert/fredNodeA.crt \
--peer-key /cert/fredNodeA.key \
--peer-ca /cert/ca.crt \
--cert /cert/fredNodeA.crt \
--key /cert/fredNodeA.key \
--ca-file /cert/ca.crt
```

This starts an instance of the `fred` software with the `info` log level using the `dev` log handler.
The ID of this node is `fredNodeA`.
It also uses an embedded BadgerDB database as a storage backend.

#### ALExANDRA

We will use the `alexandra` middleware for handling client requests. `alexandra` can be started with the following command.

```bash
docker pull git.tu-berlin.de:5000/mcc-fred/fred/alexandra:latest
docker run -d \
-v $(pwd)/alexandra.crt:/cert/alexandra.crt \
-v $(pwd)/alexandra.key:/cert/alexandra.key \
-v $(pwd)/ca.crt:/cert/ca.crt \
--network=fredwork \
--ip=172.26.1.3 \
-p 10000:10000 \
git.tu-berlin.de:5000/mcc-fred/fred/alexandra:latest \
--address :10000 \
--lighthouse 172.26.1.2:9001 \
--ca-cert /cert/ca.crt \
--alexandra-key /cert/alexandra.key \
--alexandra-cert /cert/alexandra.crt \
--clients-key /cert/alexandra.key \
--clients-cert /cert/alexandra.crt \
--experimental
```

This starts `alexandra` in a Docker container, connects it to `fred` and exposes the port 10000 on localhost with port forwarding, so that
clients can easily connect to it from the same machine.

#### Using FReD

Your initial FReD deployment is now complete!
If you want to try it out, use the `middleware.proto` in `./proto/middleware` to build a client or use [`grpcc`](https://github.com/njpatel/grpcc) to get a REPL interface:

```bash
docker build -t grpcc -f grpcc.Dockerfile .
docker run \
-v $(pwd)/fredClient.crt:/cert/fredClient.crt \
-v $(pwd)/fredClient.key:/cert/fredClient.key \
-v $(pwd)/ca.crt:/cert/ca.crt \
-v $(pwd)/proto/middleware/middleware.proto:/middleware.proto \
--network=fredwork \
--ip=172.26.1.4 \
-it \
grpcc \
grpcc -p middleware.proto \
-a 172.26.1.3:10000 \
--root_cert /cert/ca.crt \
--private_key /cert/fredClient.key \
--cert_chain /cert/fredClient.crt
```

Alternatively, you may use [`grpcui`](https://github.com/fullstorydev/grpcui), which gives you a webinterface to interactively call ALExANDRA.
After building `grpcui`, you can run it with the following command.

```bash
grpcui -open-browser -proto $(pwd)/proto/middleware/middleware.proto -cacert ca.crt -cert fredClient.crt -key fredClient.key 127.0.0.1:10000
```

You may now also add new FReD nodes, different storage backends, Trigger nodes, and more to extend your FReD deployment.

### Storage Adaptors

Each FReD node has its own local storage adaptor to persist data.
Regardless of the data format exposed to applications, internal storage is based on key-value stores and can thus easily be extended with new storage backends.

Currently, four adaptors are supported:

- In-Memory (recommended for testing-only)
- Local Filesystem (via BadgerDB)
- AWS DynamoDB
- Remote (uses a custom storage backend executable)

To choose between these backends, set the `--adaptor` flag when starting FReD.

#### BadgerDB

[BadgerDB is a key-value store developed by DGraph](https://github.com/dgraph-io/badger).
It creates a database backed by the local file system (or, optionally, in memory) with support for key expiry.

#### DynamoDB

DynamoDB is a distributed NoSQL column-family datastore by Amazon, available as-a-Service on AWS.

To use the DynamoDB storage backend, a table must already exist in DynamoDB.
It should have a composite key with the String Hash Key "Keygroup" and String Range Key "Key", and a [Number field "Expiry" that is enabled as the TTL attribute](https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/time-to-live-ttl-how-to.html).
Furthermore, the `fred` process that talks to DynamoDB should have IAM keys configured as environment variables and the corresponding IAM user must have permission to access the table.
To create a table named `fred` (this must be passed in as command-line parameter `--dynamo-table=fred`) using the AWS CLI (feel free to adapt provisioned throughput to suit your needs):

```bash
export AWS_PAGER=""
aws dynamodb create-table --table-name fred --attribute-definitions "AttributeName=Keygroup,AttributeType=S AttributeName=Key,AttributeType=S" --key-schema "AttributeName=Keygroup,KeyType=HASH AttributeName=Key,KeyType=RANGE" --provisioned-throughput "ReadCapacityUnits=1,WriteCapacityUnits=1"
aws dynamodb update-time-to-live --table-name fred --time-to-live-specification "Enabled=true, AttributeName=Expiry"
```

To delete the table:

```bash
export AWS_PAGER=""
aws dynamodb delete-table --table-name fred
```

For debugging purposes, you can also set up a local instance of DynamoDB with the `amazon/dynamodb-local:1.16.0` Docker image.
To use that, use the `--dynamodb-endpoint` flag to point to your local endpoint.
You can also use the `--dynamodb-create-table` flag to have FReD create your DynamoDB table, yet that is not recommended, e.g., when multiple FReD machines share a table.

#### Remote

Instead of accessing the storage backends directly, all storage backends can also be accessed through gRPC, which makes it possible to run the storage backend service separately.
This is useful if you want multiple `fred` instances acting as a single FReD node to share one datastore.
To use this feature, use the `storageserver` executable (or the `storage.Dockerfile` in Docker) and configure it like you would configure the storage settings in `fred`.
It supports BadgerDB, but uses the same interface as `fred` and can thus be easily extended.

If you want to use this backend, you will need to generate certificates for the storage server as well in order to secure the gRPC connection.

### API

**Clients should preferably use the ALExANDRA middleware! All API methods are available here, plus some nice-to-haves such as client-centric consistency when using version vectors!**

FReD uses a gRPC API.
The main benefit is that language libraries can easily be generated by using the `protoc` compiler.
Compare that to JSON + HTTP + REST, where you need to roll your language-specific code or we would need to provide different libraries.

For test purposes, the compiled Go and Python language files are already available in `./proto/middleware`.
To compile your own language bindings, read more about protobuf [here](https://developers.google.com/protocol-buffers).

The API provides endpoints to interact with FReD on five different aspects:

- Data
- Keygroup Management
- Replica Management
- User Management
- Trigger Node Management

#### Data

Data in FReD keygroups can be read, updated/appended and deleted.

Keygroups can be configured to be either mutable or immutable.
Data in immutable keygroups can only be appended, existing keys cannot be updated or deleted.
You need to specifiy a key in the form of a 64-bit unsigned integer, e.g., a timestamp.

Data in mutable keygroups cannot be appended as there is no concept of key incrementation, but you can insert and update data at specified keys with the update operation.
If you update a key that does not exist yet it will be created.
You can use the delete operation to delete a key from the store.

Data keys can be any string that matches the RegEx pattern `^[a-zA-Z0-9]+$`, i.e., they must be alphanumeric. Data values can be any string, although the protobuf encoding is limited to UTF-8 (AFAWK).

#### Keygroups

Keygroups can be created and deleted.
Keygroups cannot be updated in their configuration as they have no mutable attributes: the mutability of their data can only be specified when creating a new keygroup.

Keygroup names must be alphanumeric.
Additionally, you can specify whether a new keygroup should be mutable or not, i.e., if the data in the keygroup can be modified or deleted.

Furthermore, keygroups support data expiry.
Each appended, updated, or created key-value pair in the keygroup will be deleted after expiration of the data.
Note that we cannot guarantee that the data will be deleted exactly after expiring as this depends on the garbage collection of the underlying storage backends.

Expiration durations are set **per keygroup and per replica FReD node**.
The same keygroup can have different expiration duration on different FReD nodes.
For example, you may want your data to be removed on lightweight edge nodes but to be persisted in the Cloud.
Set an expiry of 0 to disable expiration.

If you set an expiry during keygroup creation, this expiry will only apply to the FReD node you are asking to create the keygroup.
This node will also automatically become the first replica node for that keygroup.

#### Replica Management

Every FReD node is a possible replica node for every keygroup.
A key concept of FReD is that this replication is configurable by application.

The API allows the `GetAllReplica` and `GetReplica` endpoints so that clients can access information about all available FReD nodes in the cluster.
Each FReD node can be uniquely identified with its ID.

The `GetKeygroupReplica` command returns a list of replica nodes for a given keygroup including the expiry settings for each replica node.
Replica nodes for keygroups can also be added by providing a node ID, keygroup name, and expiry.
Additionally, replicas can also be removed from keygroups.

#### User Management

FReD supports a simple authentication and authorization procedure.
More information about this system can be found below.

The API allows adding and removing roles for users for keygroups with the `AddUser` and `RemoveUser` endpoints.
You can specify a user ID, keygroup name, and role for that user in that keygroup.

#### Trigger Node Management

The API also allows adding, reading, and removing triggers for keygroups.

Triggers are defined **per keygroup per FReD node**.
If you add a trigger for a keygroup by sending the corresponding API request to a FReD node, only this node will send data updates to this trigger node.
As data is replicated between all FReD members of a keygroup, all updates will eventually also reach the trigger node.
This makes it easier to reason about communication between FReD nodes and trigger nodes and foregoes duplicate transmissions.

However, if you remove a replica node for a keygroup, all keygroup triggers configured on that FReD node will also be removed.

## Adding More FReD Nodes

FReD nodes communicate directly over gRPC.
It is thus vital that all nodes can communicate with each other.
You can specify your ports and addresses during FReD deployment, make sure that any proxy or firewall settings allow communication.

FReD nodes find each other through the NaSe.
As such, supplying the address of a shared NaSe during deployment is sufficient to add a new node to a FReD cluster and all nodes with access to this NaSe will consider each other as peers.
No further configuration is necessary to allow for this communication.

Keep in mind that unique identifiers should be supplied to each FReD node.
As each `fred` instance registers at the NaSe with its ID, old entries will be overwritten.
This is useful behaviour if you need to restart a failed `fred` instance with an existing ID, it will simply pick up its old operation.
Similarly, if you have a FReD node with several `fred` machines, you will need to give all of the same node ID so they can cooperate.
In this case, all other parameters should be equal so they behave equally.

## Creating a Multi-Machine Node

For horizontal scalability, a FReD node is not constrained to only a single host.
A single FReD node can instead be distributed over several hosts by sharing a common database.
This common database is where all state is kept (in addition to the NaSe that keeps additional configuration data).
For this use-case it is thus recommended to use a powerful database instance, either a remote store on a powerful machine or DynamoDB.

For consistency reasons, it's required that each keygroup in this node is bound to a particular machine within this node.
Our included `fredproxy` makes sure of that with consistent hashing.
Start the individual instances with the same FReD node identifier, their personal address (to bind interfaces correctly) and the proxy's public address (to propagate to other FReD nodes).
After all your machines have been started, add a `fredproxy` instance in front of these machines (you can actually add several proxies if you want and balance them with a load balancing L3/L4 proxy if you want).

An example of this can be found in the 3 node test.

## Authentication & Authorization

FReD uses certificates for both authentication and authorization.
This is inspired by the way `etcd` works.

Communication between any two services in FReD is peer-to-peer over gRPC, which itself uses HTTP/2 over TLS.
For example, there is the communication between two FReD nodes, from a FReD node to the NaSe, from client to a FReD node, from FReD node to a storage server, from FReD node to trigger node, and more.

### Certificates

In every such interaction, both parties each hold a private key and public certificate.
This allows for a secure connection between these parties.
To validate that the parties can be trusted, we introduce a certificate authority (CA) that can issue the certificates.
A separate CA can be used for every type of communication but in most cases it is sufficient to have a single CA in a FReD cluster and a single certificate and private key for each machine.

To get help with creating certficiates and your CA, see the [Generating Certificates](#generating-certificates) section.

Each certificate holds four pieces of information:

- a public encryption key
- a Common Name (CN)
- Hosts in the form of Subject Alternative Names (SAN) that the certificate is valid for
- key usages that the certificate may be used for

The Common Name can be set arbitrarily in most cases, but it is recommended to have this be unique between clients.
In the case of client authorization, the CN is used as the user name.
So if a client makes a request with a certificate that has the CN "Client1", roles for the user "Client1" will be applied to authorize any operation.

The SAN are validated automatically as well.
This means that the IP address entered as a SAN must match the IP that makes a request.
If required, multiple IP addresses can be entered.
In our testing, it was always required to add the loopback "127.0.0.1" address in as well.
It might even be possible to use wildcards here, but who knows.
In the future, we might run into issues with mobile clients, but we will cross that bridge when we get to it.

The key usages incldue `keyEncipherment` and `dataEncipherment`.
Additionally, the self-explanatory `serverAuth` and `clientAuth` extended key usages must be set if you want to use this certificate for server and client authentication, respectively.

### RBAC

FReD uses a lightweight role-based access control (RBAC) to enable multi-tenancy and multiple users managing data on a single FReD deployment.
Roles can be configured per user per keygroup.
That means that a user can have access to one keygroup but not to another.
Roles for a user in a keygroup are stored in the NaSe and thus consistent across all FReD nodes.

The following permissions exist:

- `Read`: read data items from a keygroup
- `Update`: update and append data items in a keygroup
- `Delete`: remove data items from a keygroup
- `AddReplica`: add a FReD node as a replica node to a keygroup
- `GetReplica`: retrieve a list of replica nodes for a keygroup along with their configuration
- `RemoveReplica`: remove a replica node from a keygroup
- `DeleteKeygroup`: delete a keygroup along with its data
- `AddUser`: add a user with a given role to a keygroup
- `RemoveUser`: remove a user's permission from a keygroup
- `GetTrigger`: get the trigger nodes for a keygroup on a replica node
- `AddTrigger`: add a trigger node as a trigger for a keygroup on a replica node
- `RemoveTrigger`: remove an existing trigger node from a keygroup on a replica node

These individual permissions are grouped into roles, which each have a unique identifier:

| Role               | Permissions                                 | Identifier           |
| ------------------ | ------------------------------------------- | -------------------- |
| Read Keygroup      | `Read`                                      | `ReadKeygroup`       |
| Write Keygroup     | `Update`, `Delete`                          | `WriteKeygroup`      |
| Configure Replica  | `AddReplica`, `GetReplica`, `RemoveReplica` | `ConfigureReplica`   |
| Configure Trigger  | `GetTrigger`, `AddTrigger`, `RemoveTrigger` | `ConfigureTrigger`   |
| Configure Keygroup | `DeleteKeygroup`, `AddUser`, `RemoveUser`   | `ConfigureKeygroups` |

When a user creates a keygroup, that user automatically receives all roles for that keygroup.

All users have the following permissions because they're not specific to a keygroup:

- retrieve a list of all available FReD nodes
- retrieve information about a particular FReD node
- create a new keygroup

Please note that just because all users can access this information, it is not considered public: users must still be authenticated with a certficate in order to talk to FReD at all.

If a user is the only user with the `ConfigureKeygroups` role for a particular keygroup, that user is, in theory, able to remove this permission from itself.
This would lead the keygroup to become unconfigurable, thus it is not recommended.

## Trigger Nodes

Trigger nodes enable getting data out of FReD automatically, for example to easily build distributed fog applications or to transform data automatically.
The basic trigger node interface is simple, yet powerful: an endpoint is set up to receive any updates to data in a keygroup.
The code behind that endpoint can do anything, including filtering data, reacting to events, transforming data into new formats, replicating data to persistent storage, and more.

A trigger node is not managed by FReD and not technically part of a FReD deployment.
Instead, a trigger node is any software service that supports the gRPC interface in `./proto/trigger/trigger.proto`.
The trigger node software thus exposes a gRPC server with this interface to any address reachable by the FReD node it is supposed to be used with.

Two types of messages will be sent to the trigger node:

1. `PutItemTriggerRequest`: includes the keygroup name, data key, and (updated) value of the data item
2. `DeleteItemTriggerRequest`: includes a keygroup name and key of the data item that was deleted

Trigger nodes can respond with an `OK` or an `ERROR`, including an optional error message.
This error message is used only for debugging and logged to FReD, not passed to any client.

## Caching in Nameservice

A CLI flag has been added to optionally enable caching for the nameservice.
Pass `--nase-cached` to your `fred` instance to activate caching.
This improves performance for requests to `fred` but may lead to data inconsistency if configuration changes often.
By default it is turned off.

## Contributing

For development, it is recommended to install [GoLand](https://www.jetbrains.com/go/).

### Git Workflow

Setup git environment with `sh ./ci/env-setup.sh` (installs git hooks). Be sure to have Go (>1.16) installed.

The `main` branch is protected and only approved pull requests can push to it.
Most important part of the workflow is `rebase`, [here's](https://www.atlassian.com/git/tutorials/merging-vs-rebasing) a refresher on merging vs rebasing.

1. Switch to `main` -> `git checkout main`
2. Update `main` -> `git pull --rebase` (ALWAYS use `rebase` when pulling!!!)
3. Create new branch from `main` -> `git checkout -b tp/new-feature` (where 'tp' is your own name/abbreviation)
4. Work on branch and push changes
5. Rebase `main` onto branch to not have merge conflicts later -> `git pull origin main --rebase` (AGAIN use`--rebase`)
6. Push branch again, this time force push to include rebased `main` (`git push --force`)
7. Create a pull request from `git.tu-berlin.de`
8. Get pull request reviewed and merge it into `main`

Some last words, keep pull requests small (not 100 files changed etc :D), so they are easier to review and rather create a lot of small pull requests than one big.

### Code Quality and Testing

In order to keep our code clean and working, we provide a number of test suites and support a number of code quality tools.

#### Static Analysis

Static analysis analyses the code in the repository without actually executing any of it.

##### Compiling

Before anything else, code must of course compile to be considered valid.
To compile the main code, use `make`, which builds the main `fred` software in the `./cmd/frednode` folder.
It first fetches dependencies and then builds the software.
Any compiler warnings and errors produced by the build process make this code invalid.

##### Linting

We use the [`golangci-lint`](https://github.com/golangci/golangci-lint) to run a number of linting tasks on our code.
Check [the documentation](https://golangci-lint.run/usage/install/#local-installation) to install it on your machine.

Once the utility is available to you, run `make lint` to lint all Go code files in the repository.
This uses the [default list of linters](https://golangci-lint.run/usage/linters/#enabled-by-default-linters), finds some basic errors and style faults, and marks them.
These linting errors mostly don't make the code invalid but ignoring them can lead to bad code quality.

##### Mega-Linting

Additionally, we also provide the `make megalint` command in our Makefile.
This runs a number of additional checks on the code, yet passing these checks is not mandatory for code to be merged into the repository.

#### Unit Tests

For a number of packages we also include unit tests to test some functionality.
This allows testing basic functionality, test for race conditions, and memory sanitation.
Use `make test`, `make race`, and `make msan` to execute all tests.

To get a coverage report, you can use `make coverage` and `make coverhtml`, depending on your preference.

#### System Tests

There are two system tests to test functionality of a FReD deployment as a whole.
This is part of a TDD approach where tests can be defined first and the software is refined until it completes all tests.
All system tests can be found in `./tests`.
All tests require Docker and Docker Compose to work.

##### 3 Node Test

The "3 node test" starts a FReD deployment of three FReD nodes and runs a client against the FReD cluster that validates different functionalities.
It uses Docker compose and can thus easily be started with `make 3n-all`.

The deployment comprises a single `etcd` Docker container as a NaSe, a simple trigger node, two FReD nodes that each comprise only a single machine (node _B_ and _C_) with a storage server, and a distributed FReD node _A_ that comprises three individual FReD machines behind a `fredproxy` sharing a single DynamoDB storage server.
All machines are connected over a Docker network.

The test client runs a number of operations against the FReD deployment and outputs a list of errors.
The complete code for the test client can be found in `./tests/3NodeTest`.

When the debug log output of the individual nodes is not enough to debug an issue, it is also possible to connect a `dlv` debugger directly to FReD node _B_ to set breakpoints or step through code.
This is currently configured to use the included debugger in the GoLand IDE.
Further information can be found in the 3 node test documentation.

##### Failing Node Test

As FReD is a distributed system, it is important to also test the impact of a failing node in the deployment.
The "Failing Node Test" allows this.
It starts the same deployment as in the 3 node test and runs a number of queries before killing one of the nodes and starting it back up.
It uses the Docker API to destroy and start the corresponding containers.

The code can be found in `./tests/FailingNodeTest` but can be started with `make failtest` in `./tests/3NodeTest/` after a deployment has been created with `make fred`.

##### ALExANDRA Test

The ALExANDRA test tests a limited amount of middleware functionality.
Use `make alexandratest` to run it.
The complete code can be found and extended in `./tests/AlexandraTest`.

##### Consistency Test

The consistency test tests consistency guarantees provided by the middleware.
In `./tests/consistency` run `bash ./run-cluster.sh [NUM_NODES] [NUM_CLIENTS]` and specify the number of FReD nodes and clients.

#### Cluster

You can easily set up a cluster of FReD nodes by using the `run-cluster.sh` script in the `cluster/` folder.
Simply run `bash run-cluster.sh [NUM_NODES]` to spawn up to 263 FReD nodes.

#### Profiling

FReD supports CPU and memory profiling for the main `frednode` binary.
Use the `--cpuprofile` and `--memprofile` flags in addition to your other flags to enable profiling.
Keep in mind that this may have an impact on performance in some cases.

```sh
# start fred with profiling
$ ./frednode --cpuprofile fredcpu.pprof --memprof fredmem.pprof [ALL_YOUR_OTHER_FLAGS...]

# run tests, benchmarks, your application, etc.
# then quit fred with SIGINT or SIGKILL and your files will be written
# open pprof files and convert to pdf (note that you need graphviz installed)
# you also need to provide the path to your frednode binary
$ go tool pprof --pdf ./frednode fredcpu.pprof > cpu.pdf
$ go tool pprof --pdf ./frednode fredmem.pprof > mem.pdf

```
