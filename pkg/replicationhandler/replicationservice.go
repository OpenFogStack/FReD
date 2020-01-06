package replicationhandler

import (
	"errors"
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

	for rn := range kg.Replica {
		node, err := s.n.GetNode(replication.Node{
			ID:   rn,
		})

		if err != nil {
			return err
		}

		if  err := s.c.SendCreateKeygroup(node.IP, node.Port, string(kg.Name)); err != nil {
			return err
		}
	}

	return err
}

// RelayDeleteKeygroup handles requests to the DeleteKeygroup endpoint of the internal interface.
func (s *service) RelayDeleteKeygroup(k keygroup.Keygroup) error {
	kg := replication.Keygroup{
		Name:    replication.KeygroupName(k.Name),
	}

	if s.n.ExistsKeygroup(kg) {
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
	kg := replication.Keygroup{
		Name:    replication.KeygroupName(i.Keygroup),
	}

	if s.n.ExistsKeygroup(kg) {
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
	kg := replication.Keygroup{
		Name:    replication.KeygroupName(i.Keygroup),
	}

	if s.n.ExistsKeygroup(kg) {
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

	return err
}

func (s *service) RemoveReplica(k keygroup.Keygroup, n replication.Node) error {
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