package fred

import (
	"github.com/go-errors/errors"
	"github.com/rs/zerolog/log"
)

// Config holds configuration parameters for an instance of FReD.
type Config struct {
	Store                   Store
	Client                  Client
	NaSe                    NameService
	PeeringHost             string
	PeeringHostProxy        string
	PeeringAsyncReplication bool
	ExternalHost            string
	ExternalHostProxy       string
	NodeID                  string
	TriggerCert             string
	TriggerKey              string
	TriggerCA               []string
	TriggerAsync            bool
}

// Fred is an instance of FReD.
type Fred struct {
	E *ExtHandler
	I *IntHandler
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

	s := newStoreService(config.Store, config.NaSe.GetNodeID())

	r := newReplicationService(s, config.Client, config.NaSe, config.PeeringAsyncReplication)

	t := newTriggerService(s, config.TriggerCert, config.TriggerKey, config.TriggerCA, config.TriggerAsync)

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
			items, err := config.Client.SendGetItem(addr, item.Keygroup, item.ID)
			if err != nil {
				log.Err(err).Msg("Was not able to get Items from node")
			}
			expiry, _ := config.NaSe.GetExpiry(item.Keygroup)

			for _, x := range items {
				_, err = s.update(x, expiry)
				if err != nil {
					log.Error().Msgf("Could not update missed item %s", x.ID)
				}
			}

		}
	} else {
		log.Trace().Msg("NodeStatus: No updates were missed by this node.")
	}

	return Fred{
		E: newExthandler(s, r, t, a, config.NaSe),
		I: newInthandler(s, r, t, config.NaSe),
	}
}
