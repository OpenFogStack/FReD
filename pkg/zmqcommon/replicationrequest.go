package zmqcommon

import (
	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/commons"
	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/replication"
)

// ReplicationRequest has all data for a ZMQ request for changing replication.
type ReplicationRequest struct {
	Keygroup commons.KeygroupName
	Node replication.Node
}