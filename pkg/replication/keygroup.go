package replication

import (
	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/commons"
)

// Keygroup has a name and a list of replica nodes.
type Keygroup struct {
	Name    commons.KeygroupName
	Replica map[ID]struct{}
}
