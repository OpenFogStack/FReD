package zmq

// DataRequest has all data for a ZMQ request.
// Currently its the same as item.go, but this might change so we use two interfaces
type DataRequest struct {
	Keygroup string
	ID       string
	Value    string
}
