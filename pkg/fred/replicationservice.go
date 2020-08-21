package fred

import (
	"github.com/go-errors/errors"
	"github.com/rs/zerolog/log"
)

// Client is an interface to send replication messages across nodes.
type Client interface {
	SendCreateKeygroup(addr Address, port int, kgname KeygroupName) error
	SendDeleteKeygroup(addr Address, port int, kgname KeygroupName) error
	SendUpdate(addr Address, port int, kgname KeygroupName, id string, value string) error
	SendDelete(addr Address, port int, kgname KeygroupName, id string) error
	SendAddReplica(addr Address, port int, kgname KeygroupName, node Node) error
	SendRemoveReplica(addr Address, port int, kgname KeygroupName, node Node) error
}

type replicationService struct {
	c Client
	s *storeService
	n *nameService
}

// newReplicationService creates a new handler for internal request (i.e. from peer nodes or the naming service).
// The nameservice makes sure that the information is synced with the other nodes
func newReplicationService(s *storeService, c Client, n *nameService) *replicationService {
	return &replicationService{
		s: s,
		c: c,
		n: n,
	}
}

// createKeygroup creates the keygroup with the NaSe and saves its existence locally
func (s *replicationService) createKeygroup(k Keygroup) error {
	log.Debug().Msgf("CreateKeygroup from replservice: in keygroup=%#v", k)

	// Check if Keygroup already exists in NaSe
	exists, err := s.n.existsKeygroup(k.Name)
	if err != nil {
		log.Err(err).Msg("Error checking whether Kg exists in NaSe")
		return err
	}

	if exists {
		log.Error().Msg("Cannot Create Keygroup since NaSe says the Kg already exists. Existing Keygroups can only be joined")
		// TODO s.localReplStore.ExistsKeygroup returns nil, but this returns an error. Whats better?
		err = errors.Errorf("keygroup already exists, cannot create it: %#v", k)
		return err
	}

	// Create Keygroup with nase, it returns an error
	err = s.n.createKeygroup(k.Name)
	if err != nil {
		log.Err(err).Msg("Error creating Keygroup in NaSe")
		return err
	}

	return err
}

// deleteKeygroup deletes the keygroup on the NaSe and deletes the local copy.
// It does not delete all the locally stored data from the keygroup.
// This gets called by every node that is in a keygroup if it is to be deleted (via RelayDeleteKeygroup, right below this method)
func (s *replicationService) deleteKeygroup(k Keygroup) error {
	log.Debug().Msgf("DeleteKeygroup from replservice (does nothing): in %#v", k)

	// Deleting the Keygroup with the NaSe happens in RelayDeleteKeygroup, so this would just be double duty
	// err := s.nase.DeleteKeygroup(k.Name)

	// TODO delete the local copy of the keygroup here

	return nil
}

// relayDeleteKeygroup deletes a keygroup locally and relays the deletion of a keygroup to all other nodes in this keygroup
// Calls DeleteKeygroup on every other node that is in this keygroup
// TODO is this really necessary with the nase? Maybe
func (s *replicationService) relayDeleteKeygroup(k Keygroup) error {
	log.Debug().Msgf("RelayDeleteKeygroup from replservice: in %#v", k)

	exists, err := s.n.existsKeygroup(k.Name)
	if err != nil {
		return err
	}
	if !exists {
		err = errors.Errorf("no such keygroup according to NaSe: %#v", k)
		return err
	}

	// Inform all other nodes about the deletion so that they can delete their local copy
	ids, err := s.n.getKeygroupMembers(k.Name, true)
	if err != nil {
		log.Err(err).Msg("Cannot delete keygroup because the nase threw an error")
		return err
	}

	for _, id := range ids {
		addr, port, err := s.n.getNodeAddress(id)

		if err != nil {
			log.Err(err).Msg("Cannot Get node adress from NaSe")
			return err
		}

		log.Debug().Msgf("RelayDeleteKeygroup from replservice: sending %#v to %#v", k, addr)
		if err := s.c.SendDeleteKeygroup(addr, port, k.Name); err != nil {
			return err
		}
	}

	// Only now delete the keygroup with the nase
	err = s.n.deleteKeygroup(k.Name)

	if err != nil {
		log.Err(err).Msg(err.(*errors.Error).ErrorStack())
		return err
	}

	return nil
}

// relayUpdate handles replication after requests to the Update endpoint of the external interface.
// It sends the update to all other nodes by calling their Update method
func (s *replicationService) relayUpdate(i Item) error {
	log.Debug().Msgf("RelayUpdate from replservice: in %#v", i)

	exists, err := s.n.existsKeygroup(i.Keygroup)
	if err != nil {
		return err
	}
	if !exists {
		err = errors.Errorf("no such keygroup according to NaSe: %#v", i.Keygroup)
		return err
	}

	// inform all other nodes about the update, get a list of all nodes subscribed to this keygroup
	ids, err := s.n.getKeygroupMembers(i.Keygroup, true)
	if err != nil {
		log.Err(err).Msg("Cannot delete keygroup because the nase threw an error")
		return err
	}

	for _, id := range ids {
		addr, port, err := s.n.getNodeAddress(id)

		if err != nil {
			log.Err(err).Msg("Cannot Get node adress from NaSe")
			return err
		}

		log.Debug().Msgf("RelayUpdate from replservice: sending %#v to %#v", i, addr)
		if err := s.c.SendUpdate(addr, port, i.Keygroup, i.ID, i.Val); err != nil {
			return err
		}
	}

	return err
}

// relayDelete handles replication after requests to the Delete endpoint of the external interface.
func (s *replicationService) relayDelete(i Item) error {
	log.Debug().Msgf("RelayDelete from replservice: in %#v", i)

	exists, err := s.n.existsKeygroup(i.Keygroup)
	if err != nil {
		return err
	}
	if !exists {
		err = errors.Errorf("no such keygroup according to NaSe: %#v", i.Keygroup)
		return err
	}

	ids, err := s.n.getKeygroupMembers(i.Keygroup, true)

	for _, id := range ids {
		addr, port, err := s.n.getNodeAddress(id)

		if err != nil {
			log.Err(err).Msg("Cannot Get node adress from NaSe")
			return err
		}

		log.Debug().Msgf("RelayDelete from replservice: sending %#v to %#v", i, addr)
		if err := s.c.SendDelete(addr, port, i.Keygroup, i.ID); err != nil {
			return err
		}
	}

	return err
}

// addReplica handles replication after requests to the AddReplica endpoint. It relays this command if "relay" is set to "true".
func (s *replicationService) addReplica(k Keygroup, n Node, relay bool) error {
	log.Debug().Msgf("AddReplica from replservice: in kg=%#v no=%#v", k, n)

	// if relay is set to true, we got the request from the external interface
	// and are responsible to bring the new replica up to speed
	// (-> send them past data, send them all other replicas, inform all other replicas)
	if relay {

		exists, err := s.n.existsKeygroup(k.Name)
		if err != nil {
			return err
		}
		if !exists {
			err = errors.Errorf("no such keygroup according to NaSe: %#v", k)
			return err
		}

		// let's get the information about this new replica
		newNodeAddr, newNodePort, err := s.n.getNodeAddress(n.ID)
		newNode := &Node{
			ID:   n.ID,
			Addr: newNodeAddr,
			Port: newNodePort,
		}

		if err != nil {
			return err
		}

		// Write the news into the NaSe first
		err = s.n.joinOtherNodeIntoKeygroup(k.Name, n.ID)

		if err != nil {
			log.Err(err).Msg(err.(*errors.Error).ErrorStack())
			return err
		}

		// let's tell this new node that it should create a local copy of this keygroup
		err = s.c.SendCreateKeygroup(newNode.Addr, newNode.Port, k.Name)
		if err != nil {
			log.Err(err).Msg(err.(*errors.Error).ErrorStack())
			return err
		}

		// now let's iterate over all other currently known replicas for this node (except ourselves)
		// this also includes the newly added node, so it will receive a AddReplica that with itself as the new node.
		ids, err := s.n.getKeygroupMembers(k.Name, true)
		if err != nil {
			return err
		}
		for _, currID := range ids {
			// get a replica node
			replAddr, replPort, err := s.n.getNodeAddress(currID)
			replNode := &Node{
				ID:   currID,
				Addr: replAddr,
				Port: replPort,
			}

			if err != nil {
				return err
			}

			// tell that replica node about the new node
			log.Debug().Msgf("AddReplica from replservice: sending %#v to %#v", newNode, replNode)
			if err := s.c.SendAddReplica(replAddr, replPort, k.Name, *newNode); err != nil {
				return err
			}

			// then tell the new node about that replica node
			log.Debug().Msgf("AddReplica from replservice: sending %#v to %#v", replNode, newNode)
			if err := s.c.SendAddReplica(newNode.Addr, newNode.Port, k.Name, *replNode); err != nil {
				return err
			}
		}

		// send all existing data to the new node
		// TODO so this one array contains all data items? Maybe not a good idea if there is a lot of data to be sent
		i, err := s.s.readAll(k.Name)

		if err != nil {
			log.Err(err).Msg(err.(*errors.Error).ErrorStack())
			return errors.Errorf("error adding replica")
		}

		log.Debug().Msgf("AddReplica from replservice: About to send %d Elements to new node", len(i))
		for _, item := range i {
			// iterate over all data for that keygroup and send it to the new node
			// a batch might be better here
			log.Debug().Msgf("AddReplica from replservice: sending %#v to %#v", item, n)
			if err := s.c.SendUpdate(newNodeAddr, newNodePort, k.Name, item.ID, item.Val); err != nil {
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

// removeReplica handles replication after requests to the RemoveReplica endpoint
// If relay==true this call comes from the exthandler => relay it to other nodes.
// The other nodes will be called with relay=false
func (s *replicationService) removeReplica(k Keygroup, n Node, relay bool) error {
	log.Debug().Msgf("RemoveReplica from replservice: in kg=%#v no=%#v", k, n)

	if relay {
		// This is the first removedNode to learn about it
		removedNodeAddr, removedNodePort, err := s.n.getNodeAddress(n.ID)
		if err != nil {
			return err
		}
		removedNode := &Node{
			ID:   n.ID,
			Addr: removedNodeAddr,
			Port: removedNodePort,
		}

		exists, err := s.n.existsKeygroup(k.Name)
		if !exists {
			err = errors.Errorf("no such keygroup according to NaSe: %#v", k)
			return err
		}

		if err != nil {
			return err
		}

		// let's tell this new node that it should delete the local copy of this keygroup
		err = s.c.SendDeleteKeygroup(removedNode.Addr, removedNode.Port, k.Name)
		if err != nil {
			log.Err(err).Msg(err.(*errors.Error).ErrorStack())
			return err
		}

		// First get all Replicas of this Keygroup to send the update to them
		kgMembers, err := s.n.getKeygroupMembers(k.Name, true)
		if err != nil {
			return err
		}
		// Now exit node from keygroup with nase.
		// Now one node (this node) will get the message that a node is removed from a keygroup
		// that it itself is not a member of ==> Delete the local copy (see after big if relay statement)
		err = s.n.exitOtherNodeFromKeygroup(k.Name, n.ID)
		if err != nil {
			return err
		}
		log.Debug().Msgf("RemoveReplica from replservice: sendingDeleteKeygroup %#v to %#v", k, removedNode)
		if err := s.c.SendDeleteKeygroup(removedNodeAddr, removedNodePort, k.Name); err != nil {
			return err
		}

		for _, idToInform := range kgMembers {
			nodeToInformAddr, nodeToInformPort, err := s.n.getNodeAddress(idToInform)

			if err != nil {
				return err
			}

			log.Debug().Msgf("RemoveReplica from replservice: sending RemoveReplica %#v to %#v", k, nodeToInformAddr)
			if err := s.c.SendRemoveReplica(nodeToInformAddr, nodeToInformPort, k.Name, Node{ID: n.ID}); err != nil {
				log.Err(err).Msg("")
				return err
			}
		}
	}

	// TODO Nils delete the local cache copy that this node is in this keygroup

	return nil
}

// getNode returns the locally saved node with this ID
func (s *replicationService) getNode(n Node) (Node, error) {
	addr, port, err := s.n.getNodeAddress(n.ID)
	n = Node{
		ID:   n.ID,
		Addr: addr,
		Port: port,
	}
	return n, err
}

// getNodes returns a list of all known nodes.
func (s *replicationService) getNodes() ([]Node, error) {
	return s.n.getAllNodes()
}

// getReplica returns a list of all replica nodes for a given keygroup.
func (s *replicationService) getReplica(k Keygroup) (nodes []Node, err error) {
	log.Debug().Msgf("GetReplica from replservice: in %#v", k)

	exists, err := s.n.existsKeygroup(k.Name)
	if !exists {
		log.Err(err).Msg(err.(*errors.Error).ErrorStack())
		return nil, errors.Errorf("no such keygroup according to NaSe")
	}
	if err != nil {
		return nil, err
	}

	ids, err := s.n.getKeygroupMembers(k.Name, true)
	for _, id := range ids {
		addr, port, err := s.n.getNodeAddress(id)
		if err != nil {
			return nil, err
		}
		newNode := &Node{
			ID:   id,
			Addr: addr,
			Port: port,
		}
		// TODO is this the optimal solution?
		nodes = append(nodes, *newNode)
	}
	return
}