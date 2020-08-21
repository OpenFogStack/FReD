package fred

import (
	"github.com/go-errors/errors"
	"github.com/rs/zerolog/log"
)

// Config holds configuration parameters for an instance of FReD.
type Config struct {
	Store            Store
	Client           Client
	PeeringHost      string
	PeeringHostProxy string
	NodeID           string
	NaSeHosts        []string
	NaSeCert         string
	NaSeKey          string
	NaSeCA           string
	TriggerCert      string
	TriggerKey       string
}

// Fred is an instance of FReD.
type Fred struct {
	E ExtHandler
	I IntHandler
}

// IntHandler is an interface that abstracts the methods of the handler that handles internal requests.
type IntHandler interface {
	HandleCreateKeygroup(k Keygroup) error
	HandleDeleteKeygroup(k Keygroup) error
	HandleUpdate(i Item) error
	HandleDelete(i Item) error
	HandleAddReplica(k Keygroup, n Node) error
	HandleRemoveReplica(k Keygroup, n Node) error
	HandleGet(i Item) (Item, error)
	HandleGetAllItems(k Keygroup) ([]Item, error)
}

// ExtHandler is an interface that abstracts the methods of the handler that handles client requests.
type ExtHandler interface {
	HandleCreateKeygroup(user string, k Keygroup) error
	HandleDeleteKeygroup(user string, k Keygroup) error
	HandleRead(user string, i Item) (Item, error)
	HandleUpdate(user string, i Item) error
	HandleDelete(user string, i Item) error
	HandleAddReplica(user string, k Keygroup, n Node) error
	HandleGetKeygroupReplica(user string, k Keygroup) ([]Node, map[NodeID]int, error)
	HandleRemoveReplica(user string, k Keygroup, n Node) error
	HandleGetReplica(user string, n Node) (Node, error)
	HandleGetAllReplica(user string) ([]Node, error)
	HandleGetKeygroupTriggers(user string, keygroup Keygroup) ([]Trigger, error)
	HandleAddTrigger(user string, keygroup Keygroup, t Trigger) error
	HandleRemoveTrigger(user string, keygroup Keygroup, t Trigger) error
	HandleAddUser(user string, newuser string, keygroup Keygroup, role Role) error
	HandleRemoveUser(user string, newuser string, keygroup Keygroup, role Role) error
}

// New creates a new FReD instance.
func New(config *Config) (f Fred) {
	s := newStoreService(config.Store)

	n, err := newNameService(config.NodeID, config.NaSeHosts, config.NaSeCert, config.NaSeKey, config.NaSeCA)

	if err != nil {
		log.Err(err).Msg(err.(*errors.Error).ErrorStack())
		panic(err)
	}

	if config.PeeringHostProxy != "" && config.PeeringHost != config.PeeringHostProxy {
		// we are behind a proxy: register with the proxy address for everyone else to find us
		err = n.registerSelf(Address{Addr: config.PeeringHostProxy})

		if err != nil {
			log.Err(err).Msg(err.(*errors.Error).ErrorStack())
			panic(err)
		}
	} else {
		// not behind a proxy: register with the local bind address
		err = n.registerSelf(Address{Addr: config.PeeringHost})

		if err != nil {
			log.Err(err).Msg(err.(*errors.Error).ErrorStack())
			panic(err)
		}
	}

	r := newReplicationService(s, config.Client, n)

	t := newTriggerService(s, config.TriggerCert, config.TriggerKey)

	a := newAuthService(n)

	return Fred{
		E: newExthandler(s, r, t, n, a),
		I: newInthandler(s, r, t, n),
	}
}
