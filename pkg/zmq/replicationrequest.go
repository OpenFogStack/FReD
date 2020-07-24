package zmq

import (
	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/fred"
)

// ReplicationRequest has all data for a ZMQ request for changing replication.
type ReplicationRequest struct {
	Keygroup fred.KeygroupName
	Node     fred.Node
}
