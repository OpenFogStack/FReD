package zmq

import (
	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/fred"
)

// KeygroupRequest has all data for a ZMQ request to add a new node to the keygroup.
type KeygroupRequest struct {
	Keygroup fred.KeygroupName
	Nodes    []fred.Node
}
