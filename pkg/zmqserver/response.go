package zmqserver

// Response has all data for a ZMQ response.
type Response struct {
	Keygroup string
	ID       string
	Value    string
}
