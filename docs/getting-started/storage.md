---
layout: default
title: Storage Adaptors
parent: Getting Started
nav_order: 2
---

## Storage Adaptors

Each FReD node has its own local storage adaptor to persist data.
Regardless of the data format exposed to applications, internal storage is based on key-value stores and can thus easily be extended with new storage backends.

Currently, four adaptors are supported:

- In-memory (recommended for testing-only)
- Local file system (via BadgerDB)
- AWS DynamoDB
- Remote (uses a custom storage backend executable)

To choose between these backends, set the `--adaptor` flag when starting FReD.

### BadgerDB

[BadgerDB is a key-value store developed by DGraph](https://github.com/dgraph-io/badger).
It creates a database backed by the local file system (or, optionally, in memory) with support for key expiry.

### DynamoDB

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

### Remote

Instead of accessing the storage backends directly, all storage backends can also be accessed through gRPC, which makes it possible to run the storage backend service separately.
This is useful if you want multiple `fred` instances acting as a single FReD node to share one datastore.
To use this feature, use the `storageserver` executable (or the `storage.Dockerfile` in Docker) and configure it like you would configure the storage settings in `fred`.
It supports BadgerDB, but uses the same interface as `fred` and can thus be easily extended.

If you want to use this backend, you will need to generate certificates for the storage server as well in order to secure the gRPC connection.
