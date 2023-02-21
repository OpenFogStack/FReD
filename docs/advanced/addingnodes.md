---
layout: default
title: Adding More FReD Nodes
parent: Advanced Configuration
nav_order: 1
---

## Adding More FReD Nodes

FReD nodes communicate directly over gRPC.
It is thus vital that all nodes can communicate with each other.
You can specify your ports and addresses during FReD deployment, make sure that any proxy or firewall settings allow communication.

FReD nodes find each other through the NaSe.
As such, supplying the address of a shared NaSe during deployment is sufficient to add a new node to a FReD cluster and all nodes with access to this NaSe will consider each other as peers.
No further configuration is necessary to allow for this communication.

Keep in mind that unique identifiers should be supplied to each FReD node.
As each `fred` instance registers at the NaSe with its ID, old entries will be overwritten.
This is useful behavior if you need to restart a failed `fred` instance with an existing ID, it will simply pick up its old operation.
Similarly, if you have a FReD node with several `fred` machines, you will need to give all of them the same node ID, so they can cooperate.
In this case, all other parameters should be equal, so they behave equally.
