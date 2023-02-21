---
layout: default
title: Authentication & Authorization
parent: Advanced Configuration
nav_order: 5
---

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

To get help with creating certificates and your CA, see the [Getting Started](./gettingstarted) section.

Each certificate holds four pieces of information:

- a public encryption key
- a Common Name (CN)
- Hosts in the form of Subject Alternative Names (SAN) that the certificate is valid for
- key usages that the certificate may be used for

The Common Name can be set arbitrarily in most cases, but it is recommended to have this be unique between clients.
In the case of client authorization, the CN is used as the username.
So if a client makes a request with a certificate that has the CN "Client1", roles for the user "Client1" will be applied to authorize any operation.

The SAN are validated automatically as well.
This means that the IP address entered as a SAN must match the IP that makes a request.
If required, multiple IP addresses can be entered.
In our testing, it was always required to add the loopback "127.0.0.1" address in as well.
It might even be possible to use wildcards here, but who knows.
In the future, we might run into issues with mobile clients, but we will cross that bridge when we get to it.

The key usages include `keyEncipherment` and `dataEncipherment`.
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

Please note that just because all users can access this information, it is not considered public: users must still be authenticated with a certificate in order to talk to FReD at all.

If a user is the only user with the `ConfigureKeygroups` role for a particular keygroup, that user is, in theory, able to remove this permission from itself.
This would lead the keygroup to become non-configurable, thus it is not recommended.
