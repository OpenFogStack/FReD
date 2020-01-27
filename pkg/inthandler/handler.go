package inthandler

import (
	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/data"
	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/keygroup"
	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/replication"
)

// Handler is an interface that abstracts the methods of the handler that handles internal requests.
type Handler interface {
	HandleCreateKeygroup(k keygroup.Keygroup, nodes []replication.Node) error
	HandleDeleteKeygroup(k keygroup.Keygroup) error
	HandleUpdate(i data.Item) error
	HandleDelete(i data.Item) error
	HandleAddReplica(k keygroup.Keygroup, n replication.Node) error
	HandleRemoveReplica(k keygroup.Keygroup, n replication.Node) error
	HandleAddNode(n replication.Node) error
	HandleRemoveNode(n replication.Node) error
	HandleIntroduction(introducer replication.Node, self replication.Node, node []replication.Node) error
	HandleDetroduction() error
}
