package zmqcommon

// Request has all data for a ZMQ request.
// Currently its the same as item.go, but this might change so we use two interfaces
type Request struct {
	Keygroup string
	ID       string
	Value    string
}