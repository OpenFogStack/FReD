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

// RelayCreateKeygroup handles requests to the CreateKeygroup endpoint of the internal interface.
func (s *service) RelayCreateKeygroup(k keygroup.Keygroup) error {
	log.Debug().Msgf("RelayCreateKeygroup from replicationservice: in %v", k)
	kg := replication.Keygroup{
		Name: replication.KeygroupName(k.Name),
	}

	if s.n.ExistsKeygroup(kg)  {
		return nil
	}


	err := s.n.CreateKeygroup(kg)
	return err
}

// RelayDeleteKeygroup handles requests to the DeleteKeygroup endpoint of the internal interface.
func (s *service) RelayDeleteKeygroup(k keygroup.Keygroup) error {
	log.Debug().Msgf("RelayDeleteKeygroup from replicationservice: in %v", k)
	kg := replication.Keygroup{
		Name:    replication.KeygroupName(k.Name),
	}

	if !s.n.ExistsKeygroup(kg) {
		return errors.New("replicationservice: no such keygroup")
	}

	kg, err := s.n.GetKeygroup(kg)

	if err != nil {
		return err
	}

	for rn := range kg.Replica {
		node, err := s.n.GetNode(replication.Node{
			ID:   rn,
		})

		if err != nil {
			return err
		}

		if  err := s.c.SendDeleteKeygroup(node.IP, node.Port, k.Name); err != nil {
			return err
		}
	}

	if err := s.n.DeleteKeygroup(kg); err != nil {
		return err
	}

	return err
}

// RelayUpdate handles requests to the Update endpoint of the internal interface.
func (s *service) RelayUpdate(i data.Item) error {
	log.Debug().Msgf("RelayUpdate from replicationservice: in %v", i)
	kg := replication.Keygroup{
		Name:    replication.KeygroupName(i.Keygroup),
	}

	if !s.n.ExistsKeygroup(kg) {
		return errors.New("replicationservice: no such keygroup")
	}

	kg, err := s.n.GetKeygroup(kg)

	if err != nil {
		return err
	}

	for rn := range kg.Replica {
		node, err := s.n.GetNode(replication.Node{
			ID:   rn,
		})

		if err != nil {
			return err
		}

		if  err := s.c.SendUpdate(node.IP, node.Port, string(kg.Name), i.ID, i.Data); err != nil {
			return err
		}
	}

	return err
}

// RelayDelete handles requests to the Delete endpoint of the internal interface.
func (s *service) RelayDelete(i data.Item) error {
	log.Debug().Msgf("RelayDelete from replicationservice: in %v", i)
	kg := replication.Keygroup{
		Name:    replication.KeygroupName(i.Keygroup),
	}

	if !s.n.ExistsKeygroup(kg) {
		return errors.New("replicationservice: no such keygroup")
	}

	kg, err := s.n.GetKeygroup(kg)

	if err != nil {
		return err
	}

	for rn := range kg.Replica {
		node, err := s.n.GetNode(replication.Node{
			ID:   rn,
		})

		if err != nil {
			return err
		}

		if  err := s.c.SendDelete(node.IP, node.Port, string(kg.Name), i.ID); err != nil {
			return err
		}
	}

	return err
}


func (s *service) AddReplica(k keygroup.Keygroup, n replication.Node) error {
	log.Debug().Msgf("AddReplica from replicationservice: in kg=%v no=%v", k, n)
	kg := replication.Keygroup{
		Name: replication.KeygroupName(k.Name),
	}

	if !s.n.ExistsKeygroup(kg)  {
		if err := s.n.CreateKeygroup(kg); err != nil {
			return err
		}
		return nil
	}

	kg, err := s.n.GetKeygroup(kg)

	if err != nil {
		return err
	}

	err = s.n.AddReplica(kg, n)

	// TODO Tobias: "Da fehlt noch was"

	return err
}

func (s *service) RemoveReplica(k keygroup.Keygroup, n replication.Node) error {
	log.Debug().Msgf("RemoveReplica from replicationservice: in kg=%v no=%v", k, n)
	kg := replication.Keygroup{
		Name: replication.KeygroupName(k.Name),
	}

	if !s.n.ExistsKeygroup(kg)  {
		if err := s.n.CreateKeygroup(kg); err != nil {
			return err
		}
		return nil
	}

	kg, err := s.n.GetKeygroup(kg)

	if err != nil {
		return err
	}

	err = s.n.RemoveReplica(kg, n)

	return err
}

func (s *service) AddNode(n replication.Node) error {
	return s.n.CreateNode(n)
}

func (s *service) RemoveNode(n replication.Node) error {
	return s.n.DeleteNode(n)
}

func (s *service) GetNodes() ([]replication.Node, error) {
	return s.n.GetNodes()
}

func (s *service) GetReplica(k keygroup.Keygroup) ([]replication.Node, error) {
	log.Debug().Msgf("GetReplica from replicationservice: in %v", k)
	kg := replication.Keygroup{
		Name: replication.KeygroupName(k.Name),
	}

	if !s.n.ExistsKeygroup(kg)  {
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