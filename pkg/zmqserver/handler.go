package zmqserver

// MessageHandler provides all methods to handle an incoming ZMQ request.
type MessageHandler interface {
	// We dont send get requests to other nodes, they will push their updates
	//HandleGetValueFromKeygroup(req *Request, from string)
	HandlePutValueIntoKeygroup(req *Request, from string)
	HandleDeleteFromKeygroup(req *Request, from string)
	HandleDeleteKeygroup(req *Request, from string)
	HandleCreateKeygroup(req *Request, src string)
}
