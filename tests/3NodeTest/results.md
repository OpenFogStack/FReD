| Description 	| Expected 	| Result 	|
|-------------	|----------	|--------	|
|Sending a GET for nonexistent item in existing KG|404 and error message|404 and no error message|
|Telling nodeA about nodeB, then adding nodeA as a replica node with nodeB | nodeB knows about nodeA and can use nodeA as Replicanode| 409 because nodeB doesn't have nodeA as Replicanode (note: Introduction Message doesnt log anything, please add logs to this handler) (1)|
| | | |
| | | |
| | | |
| | | |

```
(1)
[telling nodeA about nodeB]
nodeA    | 9:39AM DBG GetNodes from memoryrs: found 0 nodes
nodeA    | 9:39AM DBG Created a new Socket to send to node 172.26.0.11:5555 
nodeA    | 
nodeA    | 9:39AM DBG Sender has created a dealer to tcp://172.26.0.11:5555
nodeA    | 
nodeA    | 9:39AM DBG ZMQClient is sending a new message: addr=172.26.0.11, msType=25 msg="{\"Self\":{\"ID\":\"nodeA\",\"Addr\":{\"Addr\":\"172.26.0.10\",\"IsIP\":true},\"Port\":5555},\"Other\":{\"ID\":\"nodeB\",\"Addr\":{\"Addr\":\"172.26.0.11\",\"IsIP\":true},\"Port\":5555},\"Node\":[]}"
nodeA    | 9:39AM DBG Sending message type 0x19 to 
nodeA    | 9:39AM DBG CreateNode from memoryrs: in replication.Node{ID:"nodeB", Addr:replication.Address{Addr:"172.26.0.11", IsIP:true}, Port:5555}
nodeA    | 9:39AM INF Request ip=172.26.0.1 latency=0.619102 method=POST path=/v0/replica status=200 user-agent=Go-http-client/1.1

[nodeB receives message]
nodeB    | 9:39AM DBG ZMQServer received a request: msgType=25, msg=[]byte{0x7b, 0x22, 0x53, 0x65, 0x6c, 0x66, 0x22, 0x3a, 0x7b, 0x22, 0x49, 0x44, 0x22, 0x3a, 0x22, 0x6e, 0x6f, 0x64, 0x65, 0x41, 0x22, 0x2c, 0x22, 0x41, 0x64, 0x64, 0x72, 0x22, 0x3a, 0x7b, 0x22, 0x41, 0x64, 0x64, 0x72, 0x22, 0x3a, 0x22, 0x31, 0x37, 0x32, 0x2e, 0x32, 0x36, 0x2e, 0x30, 0x2e, 0x31, 0x30, 0x22, 0x2c, 0x22, 0x49, 0x73, 0x49, 0x50, 0x22, 0x3a, 0x74, 0x72, 0x75, 0x65, 0x7d, 0x2c, 0x22, 0x50, 0x6f, 0x72, 0x74, 0x22, 0x3a, 0x35, 0x35, 0x35, 0x35, 0x7d, 0x2c, 0x22, 0x4f, 0x74, 0x68, 0x65, 0x72, 0x22, 0x3a, 0x7b, 0x22, 0x49, 0x44, 0x22, 0x3a, 0x22, 0x6e, 0x6f, 0x64, 0x65, 0x42, 0x22, 0x2c, 0x22, 0x41, 0x64, 0x64, 0x72, 0x22, 0x3a, 0x7b, 0x22, 0x41, 0x64, 0x64, 0x72, 0x22, 0x3a, 0x22, 0x31, 0x37, 0x32, 0x2e, 0x32, 0x36, 0x2e, 0x30, 0x2e, 0x31, 0x31, 0x22, 0x2c, 0x22, 0x49, 0x73, 0x49, 0x50, 0x22, 0x3a, 0x74, 0x72, 0x75, 0x65, 0x7d, 0x2c, 0x22, 0x50, 0x6f, 0x72, 0x74, 0x22, 0x3a, 0x35, 0x35, 0x35, 0x35, 0x7d, 0x2c, 0x22, 0x4e, 0x6f, 0x64, 0x65, 0x22, 0x3a, 0x5b, 0x5d, 0x7d}
nodeB    | 9:39AM INF 
nodeB    | 9:39AM DBG ZMQServer received a request: msgType=21, msg=[]byte{0x7b, 0x22, 0x4b, 0x65, 0x79, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x22, 0x3a, 0x22, 0x22, 0x2c, 0x22, 0x4e, 0x6f, 0x64, 0x65, 0x22, 0x3a, 0x7b, 0x22, 0x49, 0x44, 0x22, 0x3a, 0x22, 0x6e, 0x6f, 0x64, 0x65, 0x43, 0x22, 0x2c, 0x22, 0x41, 0x64, 0x64, 0x72, 0x22, 0x3a, 0x7b, 0x22, 0x41, 0x64, 0x64, 0x72, 0x22, 0x3a, 0x22, 0x31, 0x37, 0x32, 0x2e, 0x32, 0x36, 0x2e, 0x30, 0x2e, 0x31, 0x32, 0x22, 0x2c, 0x22, 0x49, 0x73, 0x49, 0x50, 0x22, 0x3a, 0x74, 0x72, 0x75, 0x65, 0x7d, 0x2c, 0x22, 0x50, 0x6f, 0x72, 0x74, 0x22, 0x3a, 0x35, 0x35, 0x35, 0x35, 0x7d, 0x7d}
nodeB    | 9:39AM DBG CreateNode from memoryrs: in replication.Node{ID:"nodeC", Addr:replication.Address{Addr:"172.26.0.12", IsIP:true}, Port:5555}
nodeB    | 9:39AM INF 

[telling nodeB to create a Keygroup]
nodeB    | 9:39AM DBG CreateKeygroup from memorykg: in keygroup.Keygroup{Name:"KGN"}
nodeB    | 9:39AM DBG CreateKeygroup from replservice: in keygroup.Keygroup{Name:"KGN"}
nodeB    | 9:39AM DBG ExistsKeygroup from memoryrs: in replication.Keygroup{Name:"KGN", Replica:map[replication.ID]struct {}(nil)}, out false
nodeB    | 9:39AM DBG CreateKeygroup from memoryrs: in replication.Keygroup{Name:"KGN", Replica:map[replication.ID]struct {}(nil)}
nodeB    | 9:39AM INF Request ip=172.26.0.1 latency=0.22158 method=POST path=/v0/keygroup/KGN status=200 user-agent=Go-http-client/1.1

[telling nodeB to use nodeA as ReplicaNode]
nodeB    | 9:39AM DBG Exists from memorykg: in keygroup.Keygroup{Name:"KGN"}, out true
nodeB    | 9:39AM DBG AddReplica from replservice: in kg=keygroup.Keygroup{Name:"KGN"} no=replication.Node{ID:"nodeA", Addr:replication.Address{Addr:"", IsIP:false}, Port:0}
nodeB    | 9:39AM DBG ExistsKeygroup from memoryrs: in replication.Keygroup{Name:"KGN", Replica:map[replication.ID]struct {}(nil)}, out true
nodeB    | 9:39AM DBG GetKeygroup from memoryrs: in replication.Keygroup{Name:"KGN", Replica:map[replication.ID]struct {}(nil)}
nodeB    | 9:39AM DBG GetNode from memoryrs: in replication.Node{ID:"nodeA", Addr:replication.Address{Addr:"", IsIP:false}, Port:0}
nodeB    | 9:39AM ERR Exthandler cannot add a new keygroup replica error="memoryrs: no such node"
nodeB    | 9:39AM WRN Error #01: memoryrs: no such node
nodeB    |  ip=172.26.0.1 latency=0.251957 method=POST path=/v0/keygroup/KGN/replica/nodeA status=409 user-agent=Go-http-client/1.1

[putting something into keygroup in nodeB]
nodeB    | 9:39AM DBG Exists from memorykg: in keygroup.Keygroup{Name:"KGN"}, out true
nodeB    | 9:39AM DBG Update from levedbsd: in data.Item{Keygroup:"KGN", ID:"Item", Data:"Value"}
nodeB    | 9:39AM DBG RelayUpdate from replservice: in data.Item{Keygroup:"KGN", ID:"Item", Data:"Value"}
nodeB    | 9:39AM DBG ExistsKeygroup from memoryrs: in replication.Keygroup{Name:"KGN", Replica:map[replication.ID]struct {}(nil)}, out true
nodeB    | 9:39AM DBG GetKeygroup from memoryrs: in replication.Keygroup{Name:"KGN", Replica:map[replication.ID]struct {}(nil)}
nodeB    | 9:39AM DBG RelayUpdate sending to: in map[replication.ID]struct {}{}
nodeB    | 9:39AM INF Request ip=172.26.0.1 latency=0.39164 method=PUT path=/v0/keygroup/KGN/data/Item status=200 user-agent=Go-http-client/1.1

```