package zmqserver

// Request has all data for a ZMQ request.
type Request struct {
	Keygroup string
	ID       string
	Value    string
}
