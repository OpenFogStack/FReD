package zmqcommon

import (
	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/replication"
)

// IntroductionRequest has all data for a ZMQ request for introduction.
type IntroductionRequest struct {
	Self  replication.Node
	Other replication.Node
	Node  []replication.Node
}
