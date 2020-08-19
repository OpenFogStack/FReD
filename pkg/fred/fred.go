package fred

import (
	"github.com/go-errors/errors"
	"github.com/rs/zerolog/log"
)

// Config holds configuration parameters for an instance of FReD.
type Config struct {
	Store       Store
	Client      Client
	PeeringHost string
	NodeID      string
	NaSeHosts   []string
}

// Fred is an instance of FReD.
type Fred struct {
	E ExtHandler
	I IntHandler
}

// IntHandler is an interface that abstracts the methods of the handler that handles internal requests.
type IntHandler interface {
	HandleCreateKeygroup(k Keygroup, nodes []Node) error
	HandleDeleteKeygroup(k Keygroup) error
	HandleUpdate(i Item) error
	HandleDelete(i Item) error
	HandleAddReplica(k Keygroup, n Node) error
	HandleRemoveReplica(k Keygroup, n Node) error
}

// ExtHandler is an interface that abstracts the methods of the handler that handles external requests.
type ExtHandler interface {
	HandleCreateKeygroup(k Keygroup) error
	HandleDeleteKeygroup(k Keygroup) error
	HandleRead(i Item) (Item, error)
	HandleUpdate(i Item) error
	HandleDelete(i Item) error
	HandleAddReplica(k Keygroup, n Node) error
	HandleGetKeygroupReplica(k Keygroup) ([]Node, error)
	HandleRemoveReplica(k Keygroup, n Node) error
	HandleAddNode(n []Node) error
	HandleGetReplica(n Node) (Node, error)
	HandleGetAllReplica() ([]Node, error)
	HandleRemoveNode(n Node) error
}

// New creates a new FReD instance.
func New(config *Config) (f Fred) {
	s := newStoreService(config.Store)

	n, err := newNameService(config.NodeID, config.NaSeHosts)

	if err != nil {
		log.Err(err).Msg(err.(*errors.Error).ErrorStack())
		panic(err)
	}

	err = n.registerSelf(Address{Addr: config.PeeringHost})

	if err != nil {
		log.Err(err).Msg(err.(*errors.Error).ErrorStack())
		panic(err)
	}

	r := newReplicationService(config.Store, config.Client, n)

	return Fred{
		E: newExthandler(s, r),
		I: newInthandler(s, r),
	}
}
