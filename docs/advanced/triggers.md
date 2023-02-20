---
layout: default
title: Trigger Nodes
parent: Advanced Configuration
nav_order: 4
---

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
