package fred

import (
	"github.com/go-errors/errors"
	"github.com/rs/zerolog/log"
)

// Config holds configuration parameters for an instance of FReD.
type Config struct {
	Store             Store
	Client            Client
	NaSe              NameService
	PeeringHost       string
	PeeringHostProxy  string
	ExternalHost      string
	ExternalHostProxy string
	NodeID            string
	TriggerCert       string
	TriggerKey        string
	TriggerCA         []string
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
	HandleAppend(user string, i Item) (Item, error)
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

	if config.PeeringHostProxy != "" && config.PeeringHost != config.PeeringHostProxy {
		// we are behind a proxy: register with the proxy address for everyone else to find us
		err := config.NaSe.RegisterSelf(config.PeeringHostProxy, config.ExternalHostProxy)

		if err != nil {
			log.Err(err).Msg(err.(*errors.Error).ErrorStack())
			panic(err)
		}
	} else {
		// not behind a proxy: register with the local bind address
		err := config.NaSe.RegisterSelf(config.PeeringHost, config.ExternalHost)

		if err != nil {
			log.Err(err).Msg(err.(*errors.Error).ErrorStack())
			panic(err)
		}
	}

	s := newStoreService(config.Store)

	r := newReplicationService(s, config.Client, config.NaSe)

	t := newTriggerService(s, config.TriggerCert, config.TriggerKey, config.TriggerCA)

	a := newAuthService(config.NaSe)

	// TODO this code should live somewhere where it is called every n seconds, but for testing purposes the easiest way
	// TODO to simulate an internet shutdown is via killing a node, so testing once at startup should be enough
	missedItems := config.NaSe.RequestNodeStatus(config.NaSe.GetNodeID())
	if missedItems != nil {
		log.Warn().Msg("NodeStatus: This node was offline has missed some updates, getting them from other nodes")
		for _, item := range missedItems {
			nodeID, addr := config.NaSe.GetNodeWithBiggerExpiry(item.Keygroup)
			if addr == "" {
				log.Error().Msgf("NodeStatus: Was not able to find node that can provide item %s, skipping it...", item.Keygroup)
				continue
			}
			log.Info().Msgf("Getting item of KG %s ID %s from Node %s @ %s", string(item.Keygroup), item.ID, string(nodeID), addr)
			item, err := config.Client.SendGetItem(addr, item.Keygroup, item.ID)
			if err != nil {
				log.Err(err).Msg("Was not able to get Items from node")
			}
			expiry, _ := config.NaSe.GetExpiry(item.Keygroup)
			err = s.update(item, expiry)
			if err != nil {
				log.Error().Msgf("Could not update missed item %s", item.ID)
			}
		}
	} else {
		log.Debug().Msg("NodeStatus: No updates were missed by this node.")
	}

	return Fred{
		E: newExthandler(s, r, t, a, config.NaSe),
		I: newInthandler(s, r, t, config.NaSe),
	}
}
