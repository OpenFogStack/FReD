package replhandler

import (
	"github.com/rs/zerolog/log"

	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/data"
	errors "gitlab.tu-berlin.de/mcc-fred/fred/pkg/errors"
	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/keygroup"
	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/replication"
)

type service struct {
	n replication.Store
	c Client
}

// New creates a new handler for internal request (i.e. from peer nodes or the naming service).
func New(n replication.Store, c Client) replication.Service {
	return &service{
		n: n,
		c: c,
	}
}

// CreateKeygroup handles replication after requests to the CreateKeygroup endpoint of the internal interface.
func (s *service) CreateKeygroup(k keygroup.Keygroup) error {
	log.Debug().Msgf("CreateKeygroup from replservice: in %#v", k)
	kg := replication.Keygroup{
		Name: k.Name,
	}

	if s.n.ExistsKeygroup(kg) {
		return nil
	}

	err := s.n.CreateKeygroup(kg)
	return err
}

// DeleteKeygroup handles replication after requests to the DeleteKeygroup endpoint of the internal interface.
func (s *service) DeleteKeygroup(k keygroup.Keygroup) error {
	log.Debug().Msgf("RelayCreateKeygroup from replservice: in %#v", k)
	kg := replication.Keygroup{
		Name: k.Name,
	}

	err := s.n.DeleteKeygroup(kg)

	if err != nil {
		return err
	}
	return err
}

// RelayDeleteKeygroup handles replication after requests to the DeleteKeygroup endpoint of the external interface.
func (s *service) RelayDeleteKeygroup(k keygroup.Keygroup) error {
	log.Debug().Msgf("RelayDeleteKeygroup from replservice: in %#v", k)
	kg := replication.Keygroup{
		Name: k.Name,
	}

	if !s.n.ExistsKeygroup(kg) {
		log.Error().Msgf("RelayDeleteKeygroup from replservice: Keygroup does not exist: in %#v", k)
		return errors.New(errors.StatusNotFound, "replservice: no such keygroup")
	}

	kg, err := s.n.GetKeygroup(kg)

	if err != nil {
		return err
	}

	err = s.n.DeleteKeygroup(kg)

	if err != nil {
		return err
	}

	for rn := range kg.Replica {
		node, err := s.n.GetNode(replication.Node{
			ID: rn,
		})

		if err != nil {
			return err
		}

		log.Debug().Msgf("RelayDeleteKeygroup from replservice: sending %#v to %#v", k, node)
		if err := s.c.SendDeleteKeygroup(node.Addr, node.Port, k.Name); err != nil {
			return err
		}
	}

	return nil
}

// RelayUpdate handles replication after requests to the Update endpoint of the external interface.
func (s *service) RelayUpdate(i data.Item) error {
	log.Debug().Msgf("RelayUpdate from replservice: in %#v", i)
	kg := replication.Keygroup{
		Name: i.Keygroup,
	}

	if !s.n.ExistsKeygroup(kg) {
		log.Error().Msgf("RelayUpdate from replservice: Keygroup does not exist: in %#v", i)
		return errors.New(errors.StatusNotFound, "replservice: no such keygroup")
	}

	kg, err := s.n.GetKeygroup(kg)

	if err != nil {
		return err
	}

	log.Debug().Msgf("RelayUpdate sending to: in %#v", kg.Replica)

	for rn := range kg.Replica {
		node, err := s.n.GetNode(replication.Node{
			ID: rn,
		})

		if err != nil {
			return err
		}

		log.Debug().Msgf("RelayUpdate from replservice: sending %#v to %#v", i, node)
		if err := s.c.SendUpdate(node.Addr, node.Port, kg.Name, i.ID, i.Data); err != nil {
			return err
		}
	}

	return err
}

// RelayDelete handles replication after requests to the Delete endpoint of the external interface.
func (s *service) RelayDelete(i data.Item) error {
	log.Debug().Msgf("RelayDelete from replservice: in %#v", i)
	kg := replication.Keygroup{
		Name: i.Keygroup,
	}

	if !s.n.ExistsKeygroup(kg) {
		log.Error().Msgf("RelayDelete from replservice: Keygroup does not exist: in %#v", i)
		return errors.New(errors.StatusNotFound, "replservice: no such keygroup")
	}

	kg, err := s.n.GetKeygroup(kg)

	if err != nil {
		return err
	}

	for rn := range kg.Replica {
		node, err := s.n.GetNode(replication.Node{
			ID: rn,
		})

		if err != nil {
			return err
		}

		log.Debug().Msgf("RelayDelete from replservice: sending %#v to %#v", i, node)
		if err := s.c.SendDelete(node.Addr, node.Port, kg.Name, i.ID); err != nil {
			return err
		}
	}

	return err
}

// AddReplica handles replication after requests to the AddReplica endpoint. It relays this command if "relay" is set to "true".
func (s *service) AddReplica(k keygroup.Keygroup, n replication.Node, i []data.Item, relay bool) error {
	log.Debug().Msgf("AddReplica from replservice: in kg=%#v no=%#v", k, n)

	// first get the keygroup from the kgname in k
	kg := replication.Keygroup{
		Name: k.Name,
	}

	if !s.n.ExistsKeygroup(kg) {
		log.Error().Msgf("AddReplica from replservice: Keygroup does not exist: in %#v", k)
		return errors.New(errors.StatusNotFound, "replservice: no such keygroup")
	}

	kg, err := s.n.GetKeygroup(kg)

	if err != nil {
		return err
	}

	// if relay is set to true, we got the request from the external interface
	// and are responsible to bring the new replica up to speed
	// (-> send them past data, send them all other replicas, inform all other replicas)
	if relay {
		// let's get the information about this new replica first
		newNode, err := s.n.GetNode(n)

		if err != nil {
			return err
		}

		// tell this new node to create the keygroup they're now replicating
		log.Debug().Msgf("AddReplica from replservice: sending %#v to %#v", k, newNode)
		if err := s.c.SendCreateKeygroup(newNode.Addr, newNode.Port, kg.Name); err != nil {
			return err
		}

		// now tell this new node that we are also a replica node for that keygroup
		self, err := s.n.GetSelf()

		if err != nil {
			return err
		}

		log.Debug().Msgf("AddReplica from replservice: sending %#v to %#v", self, newNode)
		if err := s.c.SendAddReplica(newNode.Addr, newNode.Port, kg.Name, self); err != nil {
			return err
		}

		// now let's iterate over all other currently known replicas for this node (except ourselves)
		for rn := range kg.Replica {
			// get a replica node
			replNode, err := s.n.GetNode(replication.Node{
				ID: rn,
			})

			if err != nil {
				return err
			}

			// tell that replica node about the new node
			log.Debug().Msgf("AddReplica from replservice: sending %#v to %#v", newNode, replNode)
			if err := s.c.SendAddReplica(replNode.Addr, replNode.Port, kg.Name, newNode); err != nil {
				return err
			}

			// then tell the new node about that replica node
			log.Debug().Msgf("AddReplica from replservice: sending %#v to %#v", replNode, newNode)
			if err := s.c.SendAddReplica(newNode.Addr, newNode.Port, kg.Name, replNode); err != nil {
				return err
			}
		}

		// request came from external client interface, send past data as well
		for _, item := range i {
			// iterate over all data for that keygroup and send it to the new node
			// a batch might be better here
			log.Debug().Msgf("AddReplica from replservice: sending %#v to %#v", item, n)
			if err := s.c.SendUpdate(newNode.Addr, newNode.Port, kg.Name, item.ID, item.Data); err != nil {
				return err
			}
		}
	}

	// finally, in either case, save that the new node is now also a replica for the keygroup
	err = s.n.AddReplica(kg, n)

	if err != nil {
		return err
	}

	return nil
}

// RemoveReplica handles replication after requests to the RemoveReplica endpoint. It relays this command if "relay" is set to "true".
func (s *service) RemoveReplica(k keygroup.Keygroup, n replication.Node, relay bool) error {
	log.Debug().Msgf("RemoveReplica from replservice: in kg=%#v no=%#v", k, n)

	kg := replication.Keygroup{
		Name: k.Name,
	}

	if !s.n.ExistsKeygroup(kg) {
		log.Error().Msgf("RemoveReplica from replservice: Keygroup does not exist: in %#v", k)
		return errors.New(errors.StatusNotFound, "replservice: no such keygroup")
	}

	kg, err := s.n.GetKeygroup(kg)

	if err != nil {
		return err
	}

	err = s.n.RemoveReplica(kg, n)

	if err != nil {
		return err
	}

	if relay {
		node, err := s.n.GetNode(n)

		if err != nil {
			return err
		}

		log.Debug().Msgf("RemoveReplica from replservice: sending %#v to %#v", k, node)
		if err := s.c.SendDeleteKeygroup(node.Addr, node.Port, kg.Name); err != nil {
			return err
		}

		for rn := range kg.Replica {
			node, err := s.n.GetNode(replication.Node{
				ID: rn,
			})

			if err != nil {
				return err
			}

			log.Debug().Msgf("RemoveReplica from replservice: sending %#v to %#v", k, node)
			if err := s.c.SendRemoveReplica(node.Addr, node.Port, kg.Name, node); err != nil {
				log.Err(err).Msg("")
				return err
			}
		}
	}

	return nil
}

// AddNode handles replication after requests to the AddNode endpoint. It relays this command if "relay" is set to "true".
func (s *service) AddNode(n replication.Node, relay bool) error {
	if relay {
		nodes, err := s.n.GetNodes()

		if err != nil {
			return err
		}

		self, err := s.n.GetSelf()

		if err != nil {
			return err
		}

		if err := s.c.SendIntroduce(n.Addr, n.Port, self, n, nodes); err != nil {
			log.Err(err).Msg("")
			return err
		}

		for _, rn := range nodes {
			node, err := s.n.GetNode(replication.Node{
				ID: rn.ID,
			})

			if err != nil {
				log.Err(err).Msg("")
				return err
			}

			log.Debug().Msgf("AddNode from replservice: sending %#v to %#v", n, node)
			if err := s.c.SendAddNode(node.Addr, node.Port, n); err != nil {
				log.Err(err).Msg("")
				return err
			}
		}
	}

	// add the node afterwards to prevent it from being sent to itself
	if err := s.n.CreateNode(n); err != nil {
		log.Err(err).Msg("")
		return err
	}

	return nil
}

// RemoveNode handles replication after requests to the RemoveNode endpoint. It relays this command if "relay" is set to "true".
func (s *service) RemoveNode(n replication.Node, relay bool) error {

	if err := s.n.DeleteNode(n); err != nil {
		log.Err(err).Msg("")
		return err
	}

	if relay {
		nodes, err := s.n.GetNodes()

		if err != nil {
			return err
		}

		for _, rn := range nodes {
			node, err := s.n.GetNode(replication.Node{
				ID: rn.ID,
			})

			if err != nil {
				log.Err(err).Msg("")
				return err
			}

			log.Debug().Msgf("RemoveNode from replservice: sending %#v to %#v", n, node)
			if err := s.c.SendRemoveNode(node.Addr, node.Port, n); err != nil {
				log.Err(err).Msg("")
				return err
			}
		}
	}

	return nil
}

func (s *service) GetNode(n replication.Node) (replication.Node, error) {
	return s.n.GetNode(n)
}

// GetReplica returns a list of all known nodes.
func (s *service) GetNodes() ([]replication.Node, error) {
	return s.n.GetNodes()
}

// GetReplica returns a list of all replica nodes for a given keygroup.
func (s *service) GetReplica(k keygroup.Keygroup) ([]replication.Node, error) {
	log.Debug().Msgf("GetReplica from replservice: in %#v", k)
	kg := replication.Keygroup{
		Name: k.Name,
	}

	if !s.n.ExistsKeygroup(kg) {
		log.Error().Msgf("GetReplica from replservice: Keygroup does not exist: in %#v", k)
		return nil, errors.New(errors.StatusNotFound, "replservice: no such keygroup")
	}

	kg, err := s.n.GetKeygroup(kg)

	if err != nil {
		log.Err(err).Msgf("GetReplica from replservice: GetReplica did not work: in %#v", k)
		return nil, err
	}

	return s.n.GetReplica(kg)
}

func (s *service) Seed(n replication.Node) error {
	return s.n.Seed(n)
}

func (s *service) Unseed() error {
	return s.n.Unseed()
}