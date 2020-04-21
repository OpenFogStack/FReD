package zmqserver

import (
	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/zmqcommon"
)

// MessageHandler provides all methods to handle an incoming ZMQ request.
type MessageHandler interface {
	HandlePutValueIntoKeygroup(req *zmqcommon.DataRequest, from string)
	HandleDeleteFromKeygroup(req *zmqcommon.DataRequest, from string)
	HandleDeleteKeygroup(req *zmqcommon.KeygroupRequest, from string)
	HandleCreateKeygroup(req *zmqcommon.KeygroupRequest, src string)
	HandleAddReplica(req *zmqcommon.ReplicationRequest, src string)
	HandleRemoveReplica(req *zmqcommon.ReplicationRequest, src string)
}
