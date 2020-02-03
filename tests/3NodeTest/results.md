| Description 	| Expected 	| Result 	|
|-------------	|----------	|--------	|
|Adding a Replica for a nonexistent Keygroup| Error Code 409| Error Code 404|
|Response in Error case should be json|something that can be parsed as json| plain text|
|Nodes should gossip new nodes to eachother|after intoducing a new node to nodeA this node should tell all connected nodes about the newly connected nodes|When Registering nodeC the existing nodeB does not receive a new message informing it about the new node|
||||
||||
||||