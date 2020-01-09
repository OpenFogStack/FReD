package replicationhandler

import (
	"errors"

	"github.com/rs/zerolog/log"

	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/data"
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

// RelayCreateKeygroup handles replication after requests to the CreateKeygroup endpoint of the external interface.
func (s *service) RelayCreateKeygroup(k keygroup.Keygroup) error {
	log.Debug().Msgf("RelayCreateKeygroup from replicationservice: in %v", k)
	kg := replication.Keygroup{
		Name: k.Name,
	}

	if s.n.ExistsKeygroup(kg) {
		return nil
	}

	err := s.n.CreateKeygroup(kg)
	return err
}

// RelayDeleteKeygroup handles replication after requests to the DeleteKeygroup endpoint of the external interface.
func (s *service) RelayDeleteKeygroup(k keygroup.Keygroup) error {
	log.Debug().Msgf("RelayDeleteKeygroup from replicationservice: in %v", k)
	kg := replication.Keygroup{
		Name: k.Name,
	}

	if !s.n.ExistsKeygroup(kg) {
		log.Error().Msgf("RelayDeleteKeygroup from replicationservice: Keygroup does not exist: in %v", k)
		return errors.New("replicationservice: no such keygroup")
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

		log.Debug().Msgf("RelayDeleteKeygroup from replicationservice: sending %v to %v", k, node)
		if err := s.c.SendDeleteKeygroup(node.IP, node.Port, k.Name); err != nil {
			return err
		}
	}



	return nil
}

// RelayUpdate handles replication after requests to the Update endpoint of the external interface.
func (s *service) RelayUpdate(i data.Item) error {
	log.Debug().Msgf("RelayUpdate from replicationservice: in %v", i)
	kg := replication.Keygroup{
		Name: i.Keygroup,
	}

	if !s.n.ExistsKeygroup(kg) {
		log.Error().Msgf("RelayUpdate from replicationservice: Keygroup does not exist: in %v", i)
		return errors.New("replicationservice: no such keygroup")
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

		log.Debug().Msgf("RelayUpdate from replicationservice: sending %v to %v", i, node)
		if err := s.c.SendUpdate(node.IP, node.Port, kg.Name, i.ID, i.Data); err != nil {
			return err
		}
	}

	return err
}

// RelayDelete handles replication after requests to the Delete endpoint of the external interface.
func (s *service) RelayDelete(i data.Item) error {
	log.Debug().Msgf("RelayDelete from replicationservice: in %v", i)
	kg := replication.Keygroup{
		Name: i.Keygroup,
	}

	if !s.n.ExistsKeygroup(kg) {
		log.Error().Msgf("RelayDelete from replicationservice: Keygroup does not exist: in %v", i)
		return errors.New("replicationservice: no such keygroup")
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

		log.Debug().Msgf("RelayDelete from replicationservice: sending %v to %v", i, node)
		if err := s.c.SendDelete(node.IP, node.Port, kg.Name, i.ID); err != nil {
			return err
		}
	}

	return err
}

// AddReplica handles replication after requests to the AddReplica endpoint. It relays this command if "relay" is set to "true".
func (s *service) AddReplica(k keygroup.Keygroup, n replication.Node, relay bool) error {
	log.Debug().Msgf("AddReplica from replicationservice: in kg=%v no=%v", k, n)
	kg := replication.Keygroup{
		Name: k.Name,
	}

	if !s.n.ExistsKeygroup(kg) {
		log.Error().Msgf("AddReplica from replicationservice: Keygroup does not exist: in %v", k)
		return errors.New("replicationservice: no such keygroup")
	}

	kg, err := s.n.GetKeygroup(kg)

	if err != nil {
		return err}

	err = s.n.AddReplica(kg, n)
	if err != nil {
		return err
	}


	if relay {
		node, err := s.n.GetNode(n)

		if err != nil {
			return err
		}

		log.Debug().Msgf("AddReplica from replicationservice: sending %v to %v", k, node)
		if err := s.c.SendCreateKeygroup(node.IP, node.Port, kg.Name); err != nil {
			return err
		}

		for rn := range kg.Replica {
			node, err := s.n.GetNode(replication.Node{
				ID: rn,
			})

			if err != nil {
				return err
			}

			log.Debug().Msgf("AddReplica from replicationservice: sending %v to %v", k, node)
			if err := s.c.SendAddReplica(node.IP, node.Port, kg.Name, node.ID); err != nil {
				return err
			}
		}
	}

	return nil
}

// RemoveReplica handles replication after requests to the RemoveReplica endpoint. It relays this command if "relay" is set to "true".
func (s *service) RemoveReplica(k keygroup.Keygroup, n replication.Node, relay bool) error {
	log.Debug().Msgf("RemoveReplica from replicationservice: in kg=%v no=%v", k, n)

	kg := replication.Keygroup{
		Name: k.Name,
	}

	if !s.n.ExistsKeygroup(kg) {
		log.Error().Msgf("RemoveReplica from replicationservice: Keygroup does not exist: in %v", k)
		return errors.New("replicationservice: no such keygroup")
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

		log.Debug().Msgf("RemoveReplica from replicationservice: sending %v to %v", k, node)
		if err := s.c.SendDeleteKeygroup(node.IP, node.Port, kg.Name); err != nil {
			return err
		}

		for rn := range kg.Replica {
			node, err := s.n.GetNode(replication.Node{
				ID: rn,
			})

			if err != nil {
				return err
			}

			log.Debug().Msgf("RemoveReplica from replicationservice: sending %v to %v", k, node)
			if err := s.c.SendRemoveReplica(node.IP, node.Port, kg.Name, node.ID); err != nil {
				log.Err(err).Msg("")
				return err
			}
		}
	}

	return nil
}

// AddNode handles replication after requests to the AddNode endpoint. It relays this command if "relay" is set to "true".
func (s *service) AddNode(n replication.Node, relay bool) error {

	if err := s.n.CreateNode(n); err != nil {
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

			log.Debug().Msgf("AddNode from replicationservice: sending %v to %v", n, node)
			if err := s.c.SendAddNode(node.IP, node.Port, n.ID, n.IP, n.Port); err != nil {
				log.Err(err).Msg("")
				return err
			}
		}
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

			log.Debug().Msgf("RemoveNode from replicationservice: sending %v to %v", n, node)
			if err := s.c.SendRemoveNode(node.IP, node.Port, n.ID); err != nil {
				log.Err(err).Msg("")
				return err
			}
		}
	}

	return nil
}

// GetReplica returns a list of all known nodes.
func (s *service) GetNodes() ([]replication.Node, error) {
	return s.n.GetNodes()
}

// GetReplica returns a list of all replica nodes for a given keygroup.
func (s *service) GetReplica(k keygroup.Keygroup) ([]replication.Node, error) {
	log.Debug().Msgf("GetReplica from replicationservice: in %v", k)
	kg := replication.Keygroup{
		Name: k.Name,
	}

	if !s.n.ExistsKeygroup(kg) {
		log.Error().Msgf("GetReplica from replicationservice: Keygroup does not exist: in %v", k)
		return nil, errors.New("replicationservice: no such keygroup")
	}

	kg, err := s.n.GetKeygroup(kg)

	if err != nil {
		log.Err(err).Msgf("GetReplica from replicationservice: GetReplica did not work: in %v", k)
		return nil, err
	}

	return s.n.GetReplica(kg)
}
