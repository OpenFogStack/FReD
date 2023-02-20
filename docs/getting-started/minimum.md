---
layout: default
title: Minimum FReD Deployment
parent: Getting Started
nav_order: 1
---

## Minimum FReD Deployment

The simplest way to get started is to use Docker containers to deploy the needed software.

If you're not familiar with Docker yet, take the time to learn about it [here](https://www.docker.com/101-tutorial).
Alternatively, you are of course able to run every component manually.

### Certificates

Before deploying `etcd` and `fred` software, certificates have to be created to secure the connection between all components and authenticate the individual services.
We provide tooling to create new certificates.
More information about authentication and authorization with certificates can be found further down in this introduction.
This guide is adapted from [scriptcrunch.com](https://scriptcrunch.com/create-ca-tls-ssl-certificates-keys/).

#### Certificate Authority

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

#### Generating Certificates

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

### Network

If you run this example in Docker, you must first create a simple network for the individual services to talk to each other:

```bash
docker network create fredwork --gateway 172.26.0.1 --subnet 172.26.0.0/16
```

### NaSe

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

### FReD

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

### ALExANDRA

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

### Using FReD

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

Alternatively, you may use [`grpcui`](https://github.com/fullstorydev/grpcui), which gives you a web interface to interactively call ALExANDRA.
After building `grpcui`, you can run it with the following command.

```bash
grpcui -open-browser -proto $(pwd)/proto/middleware/middleware.proto -cacert ca.crt -cert fredClient.crt -key fredClient.key 127.0.0.1:10000
```

You may now also add new FReD nodes, different storage backends, Trigger nodes, and more to extend your FReD deployment.
