package replication

import (
	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/data"
	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/keygroup"
)

// Service is an interface that encapsulates the needed methods to replicate data across nodes.
type Service interface {
	CreateKeygroup(k keygroup.Keygroup) error
	DeleteKeygroup(k keygroup.Keygroup) error
	RelayDeleteKeygroup(k keygroup.Keygroup) error
	RelayUpdate(i data.Item) error
	RelayDelete(i data.Item) error
	AddNode(n Node, relay bool) error
	RemoveNode(n Node, relay bool) error
	AddReplica(k keygroup.Keygroup, n Node, i []data.Item, relay bool) error
	RemoveReplica(k keygroup.Keygroup, n Node, relay bool) error
	GetNode(n Node) (Node, error)
	GetNodes() ([]Node, error)
	GetReplica(k keygroup.Keygroup) ([]Node, error)
	Seed(n Node) error
	Unseed() error
}
