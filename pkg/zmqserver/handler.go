package zmqserver

import (
	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/zmqcommon"
)

// MessageHandler provides all methods to handle an incoming ZMQ request.
type MessageHandler interface {
	// We dont send get requests to other nodes, they will push their updates
	//HandleGetValueFromKeygroup(req *Request, from string)
	HandlePutValueIntoKeygroup(req *zmqcommon.Request, from string)
	HandleDeleteFromKeygroup(req *zmqcommon.Request, from string)
	HandleDeleteKeygroup(req *zmqcommon.Request, from string)
	HandleCreateKeygroup(req *zmqcommon.Request, src string)
	HandleAddNode(req *zmqcommon.ReplicationRequest, src string)
	HandleRemoveNode(req *zmqcommon.ReplicationRequest, src string)
	HandleAddReplica(req *zmqcommon.ReplicationRequest, src string)
	HandleRemoveReplica(req *zmqcommon.ReplicationRequest, src string)
}
