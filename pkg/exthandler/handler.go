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
	HandleAddKeygroupReplica(k keygroup.Keygroup, n replication.Node) error
	HandleGetKeygroupReplica(k keygroup.Keygroup) ([]replication.Node, error)
	HandleRemoveKeygroupReplica(k keygroup.Keygroup, n replication.Node) error
	HandleAddReplica(n []replication.Node) error
	HandleGetReplica() ([]replication.Node, error)
	HandleRemoveReplica(n replication.Node) error
}