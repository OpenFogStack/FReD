---
layout: default
title: Multi-Machine Node
parent: Advanced Configuration
nav_order: 2
---

## Creating a Multi-Machine Node

For horizontal scalability, a FReD node is not constrained to only a single host.
A single FReD node can instead be distributed over several hosts by sharing a common database.
This common database is where all state is kept (in addition to the NaSe that keeps additional configuration data).
For this use-case it is thus recommended using a powerful database instance, either a remote store on a powerful machine or DynamoDB.

For consistency reasons, it's required that each keygroup in this node is bound to a particular machine within this node.
Our included `fredproxy` makes sure of that with consistent hashing.
Start the individual instances with the same FReD node identifier, their personal address (to bind interfaces correctly) and the proxy's public address (to propagate to other FReD nodes).
After all your machines have been started, add a `fredproxy` instance in front of these machines (you can actually add several proxies if you want and balance them with a load balancing L3/L4 proxy if you want).

An example of this can be found in the 3 node test.
