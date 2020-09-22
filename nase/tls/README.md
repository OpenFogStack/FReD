# Generating etcd certificates

You need certificates to properly talk to etcd with authentication and encryption.
Use the `gen-cert.sh` script to generate client and server certificates with `gen-cert.sh server 172.26.1.1` to generate a `server.crt` and `server.key` for a server with IP address `172.26.1.1`.
You can use the same script to generate client certificates and peer certificates, it doesn't matter.
But please down delete the CA files.
pls.