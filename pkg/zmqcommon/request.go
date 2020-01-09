package zmqcommon

// Request has all data for a ZMQ request.
// Currently its the same as item.go, but this might change so we use two interfaces
type Request struct {
	Keygroup string
	ID       string
	Value    string
}

// ReplicationRequest has all data for a ZMQ request for changing replication.
type ReplicationRequest struct {
	Keygroup string
	Node struct {
		ID string
		IP string
		Port string
	}
}
