package zmq

// MessageHandler provides all methods to handle an incoming ZMQ request.
type MessageHandler interface {
	HandlePutValueIntoKeygroup(req *DataRequest, from string)
	HandleDeleteFromKeygroup(req *DataRequest, from string)
	HandleDeleteKeygroup(req *KeygroupRequest, from string)
	HandleCreateKeygroup(req *KeygroupRequest, src string)
	HandleAddReplica(req *ReplicationRequest, src string)
	HandleRemoveReplica(req *ReplicationRequest, src string)
}
