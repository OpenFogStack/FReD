package zmq

// Different types of messages that the zmqclient can send to other nodes
// Currently no messages are received
const (
	CreateKeygroup byte = 0x10
	DeleteKeygroup byte = 0x11
	//GET_ITEM byte = 0x12
	PutItem       byte = 0x13
	DeleteItem    byte = 0x14
	AddNode       byte = 0x15
	RemoveNode    byte = 0x16
	AddReplica    byte = 0x17
	RemoveReplica byte = 0x18
)
