package exthandler

import (
	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/data"
	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/keygroup"
	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/replication"
)

// Handler is an interface that abstracts the methods of the handler that handles external requests.
type Handler interface {
	HandleCreateKeygroup(k keygroup.Keygroup) error
	HandleDeleteKeygroup(k keygroup.Keygroup) error
	HandleRead(i data.Item) (data.Item, error)
	HandleUpdate(i data.Item) error
	HandleDelete(i data.Item) error
	HandleAddReplica(k keygroup.Keygroup, n replication.Node) error
	HandleGetKeygroupReplica(k keygroup.Keygroup) ([]replication.Node, error)
	HandleRemoveReplica(k keygroup.Keygroup, n replication.Node) error
	HandleAddNode(n []replication.Node) error
	HandleGetReplica(n replication.Node) (replication.Node, error)
	HandleGetAllReplica() ([]replication.Node, error)
	HandleRemoveNode(n replication.Node) error
	HandleSeed(n replication.Node) error
}
