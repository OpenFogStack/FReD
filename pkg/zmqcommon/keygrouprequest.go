package zmqcommon

import (
	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/commons"
	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/replication"
)

// KeygroupRequest has all data for a ZMQ request to add a new node to the keygroup.
type KeygroupRequest struct {
	Keygroup commons.KeygroupName
	Nodes    []replication.Node
}
