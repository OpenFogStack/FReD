package replhandler

import (
	"errors"

	"github.com/rs/zerolog/log"

	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/commons"
	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/data"
	frederrors "gitlab.tu-berlin.de/mcc-fred/fred/pkg/errors"
	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/keygroup"
	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/nameservice"
	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/replication"
)

type service struct {
	// The ZMQ Client that communicates with other nodes
	zmqClient Client
	nase      nameservice.NameService
	dataStore data.Store
}

// New creates a new handler for internal request (i.e. from peer nodes or the naming service).
// The nameservice makes sure that the information is synced with the other nodes
func New(c Client, nase nameservice.NameService, dataStore data.Store) replication.Service {
	return &service{
		zmqClient: c,
		nase:      nase,
		dataStore: dataStore,
	}
}

// CreateKeygroup creates the keygroup with the NaSe and saves its existence locally
func (s *service) CreateKeygroup(k keygroup.Keygroup) error {
	log.Debug().Msgf("CreateKeygroup from replservice: in %#v", k)

	// Check if Keygroup already exists in NaSe
	exists, err := s.nase.ExistsKeygroup(k.Name)
	if err != nil {
		log.Err(err).Msg("Error checking whether Kg exists in NaSe")
		return err
	}
	if exists {
		log.Error().Msg("Cannot Create Keygroup since NaSe says the Kg already exists. Existing Keygroups can only be joined")
		// TODO s.localReplStore.ExistsKeygroup returns nil, but this returns an error. Whats better?
		return errors.New("keygroup already exists, cannot create it")
	}
	// Create Keygroup with nase, it returns an error
	err = s.nase.CreateKeygroup(k.Name)
	if err != nil {
		log.Err(err).Msg("Error creating Keygroup in NaSe")
		return err
	}

	// kg := replication.Keygroup{
	// 	Name: k.Name,
	// }
	// TODO save a local copy of the keygroup here
	return err
}

// DeleteKeygroup deletes the keygroup on the NaSe and deletes the local copy.
// It does not delete all the locally stored data from the keygroup.
// This gets called by every node that is in a keygroup if it is to be deleted (via RelayDeleteKeygroup, right below this method)
func (s *service) DeleteKeygroup(k keygroup.Keygroup) error {
	log.Debug().Msgf("DeleteKeygroup from replservice (does nothing): in %#v", k)

	// Deleting the Keygroup with the NaSe happens in RelayDeleteKeygroup, so this would just be double duty
	// err := s.nase.DeleteKeygroup(k.Name)

	// TODO delete the local copy of the keygroup here

	return nil
}

// RelayDeleteKeygroup deletes a keygroup locally and relays the deletion of a keygroup to all other nodes in this keygroup
// Calls DeleteKeygroup on every other node that is in this keygroup
// TODO is this really necessary with the nase? Maybe
func (s *service) RelayDeleteKeygroup(k keygroup.Keygroup) error {
	log.Debug().Msgf("RelayDeleteKeygroup from replservice: in %#v", k)

	exists, err := s.nase.ExistsKeygroup(k.Name)
	if err != nil {
		return err
	}
	if !exists {
		log.Error().Msgf("RelayUpdate from replservice: Keygroup does not exist according to NaSe: in %#v", k)
		return frederrors.New(frederrors.StatusNotFound, "replservice: no such keygroup according to NaSe")
	}

	// Inform all other nodes about the deletion so that they can delete their local copy
	ids, err := s.nase.GetKeygroupMembers(k.Name, true)
	if err != nil {
		log.Err(err).Msg("Cannot delete keygroup because the nase threw an error")
		return err
	}

	for _, id := range ids {
		addr, port, err := s.nase.GetNodeAdress(id)

		if err != nil {
			log.Err(err).Msg("Cannot Get node adress from NaSe")
			return err
		}

		log.Debug().Msgf("RelayDeleteKeygroup from replservice: sending %#v to %#v", k, addr)
		if err := s.zmqClient.SendDeleteKeygroup(addr, port, k.Name); err != nil {
			return err
		}
	}

	// Only now delete the keygroup with the nase
	s.nase.DeleteKeygroup(k.Name)

	return nil
}

// RelayUpdate handles replication after requests to the Update endpoint of the external interface.
// It sends the update to all other nodes by calling their Update method
func (s *service) RelayUpdate(i data.Item) error {
	log.Debug().Msgf("RelayUpdate from replservice: in %#v", i)

	exists, err := s.nase.ExistsKeygroup(i.Keygroup)
	if err != nil {
		return err
	}
	if !exists {
		log.Error().Msgf("RelayUpdate from replservice: Keygroup does not exist according to NaSe: in %#v", i)
		return frederrors.New(frederrors.StatusNotFound, "replservice: no such keygroup according to NaSe")
	}

	// inform all other nodes about the update, get a list of all nodes subscribed to this keygroup
	ids, err := s.nase.GetKeygroupMembers(i.Keygroup, true)
	if err != nil {
		log.Err(err).Msg("Cannot delete keygroup because the nase threw an error")
		return err
	}

	for _, id := range ids {
		addr, port, err := s.nase.GetNodeAdress(id)

		if err != nil {
			log.Err(err).Msg("Cannot Get node adress from NaSe")
			return err
		}

		log.Debug().Msgf("RelayUpdate from replservice: sending %#v to %#v", i, addr)
		if err := s.zmqClient.SendUpdate(addr, port, i.Keygroup, i.ID, i.Data); err != nil {
			return err
		}
	}

	return err
}

// RelayDelete handles replication after requests to the Delete endpoint of the external interface.
func (s *service) RelayDelete(i data.Item) error {
	log.Debug().Msgf("RelayDelete from replservice: in %#v", i)

	exists, err := s.nase.ExistsKeygroup(i.Keygroup)
	if err != nil {
		return err
	}
	if !exists {
		log.Error().Msgf("RelayDelete from replservice: Keygroup does not exist according to NaSe: in %#v", i)
		return frederrors.New(frederrors.StatusNotFound, "replservice: no such keygroup according to NaSe")
	}

	ids, err := s.nase.GetKeygroupMembers(i.Keygroup, true)

	for _, id := range ids {
		addr, port, err := s.nase.GetNodeAdress(id)

		if err != nil {
			log.Err(err).Msg("Cannot Get node adress from NaSe")
			return err
		}

		log.Debug().Msgf("RelayDelete from replservice: sending %#v to %#v", i, addr)
		if err := s.zmqClient.SendDelete(addr, port, i.Keygroup, i.ID); err != nil {
			return err
		}
	}

	return err
}

// AddReplica handles replication after requests to the AddReplica endpoint. It relays this command if "relay" is set to "true".
func (s *service) AddReplica(k keygroup.Keygroup, n replication.Node, i []data.Item, relay bool) error {
	log.Debug().Msgf("AddReplica from replservice: in kg=%#v no=%#v", k, n)

	// if relay is set to true, we got the request from the external interface
	// and are responsible to bring the new replica up to speed
	// (-> send them past data, send them all other replicas, inform all other replicas)
	if relay {

		exists, err := s.nase.ExistsKeygroup(k.Name)
		if err != nil {
			return err
		}
		if !exists {
			log.Error().Msgf("AddReplica from replservice: Keygroup does not exist according to NaSe: in %#v", i)
			return frederrors.New(frederrors.StatusNotFound, "replservice: no such keygroup according to NaSe")
		}

		// let's get the information about this new replica
		newNodeAddr, newNodePort, err := s.nase.GetNodeAdress(n.ID)
		newNode := &replication.Node{
			ID:   n.ID,
			Addr: newNodeAddr,
			Port: newNodePort,
		}

		if err != nil {
			return err
		}

		// Write the news into the NaSe first
		s.nase.JoinOtherNodeIntoKeygroup(k.Name, n.ID)

		// now let's iterate over all other currently known replicas for this node (except ourselves)
		// this also includes the newly added node, so it will receive a AddReplica that with itself as the new node.
		ids, err := s.nase.GetKeygroupMembers(k.Name, true)
		if err != nil {
			return err
		}
		for _, currID := range ids {
			// get a replica node
			replAddr, replPort, err := s.nase.GetNodeAdress(currID)
			replNode := &replication.Node{
				ID:   currID,
				Addr: replAddr,
				Port: replPort,
			}

			if err != nil {
				return err
			}

			// tell that replica node about the new node
			log.Debug().Msgf("AddReplica from replservice: sending %#v to %#v", newNode, replNode)
			if err := s.zmqClient.SendAddReplica(replAddr, replPort, k.Name, *newNode); err != nil {
				return err
			}

			// then tell the new node about that replica node
			log.Debug().Msgf("AddReplica from replservice: sending %#v to %#v", replNode, newNode)
			if err := s.zmqClient.SendAddReplica(newNode.Addr, newNode.Port, k.Name, *replNode); err != nil {
				return err
			}
		}

		// send all existing data to the new node
		// TODO so this one array contains all data items? Maybe not a good idea if there is a lot of data to be sent
		log.Debug().Msgf("AddReplica from replservice: About to send %d Elements to new node", len(i))
		for _, item := range i {
			// iterate over all data for that keygroup and send it to the new node
			// a batch might be better here
			log.Debug().Msgf("AddReplica from replservice: sending %#v to %#v", item, n)
			if err := s.zmqClient.SendUpdate(newNodeAddr, newNodePort, k.Name, item.ID, item.Data); err != nil {
				return err
			}
		}
	}
	// else {
	// 	// This request is not from the external interface, but from the internal one.
	// 	// TODO Nils the information received here should be cached locally
	// 	if string(n.ID) == s.nase.NodeID {
	// 		// TODO Nils this node just has been informed that it is now a loyal replicating member of a KG
	// 	} else {
	// 		// TODO Nils a third node has been informed that n.ID is a new member of a keygroup it is also replicating
	// 	}
	//
	// }
	return nil
}

// RemoveReplica handles replication after requests to the RemoveReplica endpoint
// If relay==true this call comes from the exthandler => relay it to other nodes.
// The other nodes will be called with relay=false
func (s *service) RemoveReplica(k keygroup.Keygroup, n replication.Node, relay bool) error {
	log.Debug().Msgf("RemoveReplica from replservice: in kg=%#v no=%#v", k, n)

	if relay {
		// This is the first removedNode to learn about it
		removedNodeAddr, removedNodePort, err := s.nase.GetNodeAdress(n.ID)
		if err != nil {
			return err
		}
		removedNode := &replication.Node{ID: n.ID,
			Addr: removedNodeAddr,
			Port: removedNodePort,
		}

		exists, err := s.nase.ExistsKeygroup(k.Name)
		if !exists {
			log.Error().Msgf("RemoveReplica from replservice: Keygroup does not exist according to NaSe: in %#v", k)
			return frederrors.New(frederrors.StatusNotFound, "replservice: no such keygroup according to NaSe")
		}
		if err != nil {
			return err
		}

		// First get all Replicas of this Keygroup to send the update to them
		kgMembers, err := s.nase.GetKeygroupMembers(k.Name, true)
		if err != nil {
			return err
		}
		// Now exit node from keygroup with nase.
		// Now one node (this node) will get the message that a node is removed from a keygroup
		// that it itself is not a member of ==> Delete the local copy (see after big if relay statement)
		err = s.nase.ExitOtherNodeFromKeygroup(k.Name, n.ID)
		if err != nil {
			return err
		}
		log.Debug().Msgf("RemoveReplica from replservice: sendingDeleteKeygroup %#v to %#v", k, removedNode)
		if err := s.zmqClient.SendDeleteKeygroup(removedNodeAddr, removedNodePort, k.Name); err != nil {
			return err
		}

		for _, idToInform := range kgMembers {
			nodeToInformAddr, nodeToInformPort, err := s.nase.GetNodeAdress(idToInform)

			if err != nil {
				return err
			}

			log.Debug().Msgf("RemoveReplica from replservice: sending RemoveReplica %#v to %#v", k, nodeToInformAddr)
			if err := s.zmqClient.SendRemoveReplica(nodeToInformAddr, nodeToInformPort, k.Name, replication.Node{ID: n.ID}); err != nil {
				log.Err(err).Msg("")
				return err
			}
		}
	}
	member, err := s.nase.IsKeygroupMember(s.nase.NodeID, k.Name)
	if !member {
		log.Debug().Msgf("RemoveReplica from Replservice: deleting local copy of this keygroup")
		s.dataStore.DeleteKeygroup(k.Name)
	}

	// TODO Nils delete the local copy that this node is in this keygroup

	return err
}

// GetNode returns the locally saved node with this ID
func (s *service) GetNode(n replication.Node) (replication.Node, error) {
	addr, port, err := s.nase.GetNodeAdress(n.ID)
	n = replication.Node{
		ID:   n.ID,
		Addr: addr,
		Port: port,
	}
	return n, err
}

// GetNodes returns a list of all known nodes.
func (s *service) GetNodes() ([]replication.Node, error) {
	return s.nase.GetAllNodes()
}

// GetReplica returns a list of all replica nodes for a given keygroup.
func (s *service) GetReplica(k keygroup.Keygroup) (nodes []replication.Node, err error) {
	log.Debug().Msgf("GetReplica from replservice: in %#v", k)

	exists, err := s.nase.ExistsKeygroup(k.Name)
	if !exists {
		log.Error().Msgf("RemoveReplica from replservice: Keygroup does not exist according to NaSe: in %#v", k)
		return nil, frederrors.New(frederrors.StatusNotFound, "replservice: no such keygroup according to NaSe")
	}
	if err != nil{
		return nil, err
	}

	ids, err := s.nase.GetKeygroupMembers(k.Name, true)
	for _, id := range ids {
		addr, port, err := s.nase.GetNodeAdress(id)
		if err != nil {
			return nil, err
		}
		newNode := &replication.Node{
			ID:   id,
			Addr: addr,
			Port: port,
		}
		// TODO is this the optimal solution?
		nodes = append(nodes, *newNode)
	}
	return
}

func (s *service) ExistsKeygroup(name commons.KeygroupName) (bool, error) {
	return s.nase.ExistsKeygroup(name)
}
