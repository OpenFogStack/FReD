package fred

import (
	"github.com/go-errors/errors"
	"github.com/rs/zerolog/log"
)

// Client is an interface to send replication messages across nodes.
type Client interface {
	SendCreateKeygroup(host string, kgname KeygroupName, expiry int) error
	SendDeleteKeygroup(host string, kgname KeygroupName) error
	SendUpdate(host string, kgname KeygroupName, id string, value string) error
	SendDelete(host string, kgname KeygroupName, id string) error
	SendAddReplica(host string, kgname KeygroupName, node Node, expiry int) error
	SendRemoveReplica(host string, kgname KeygroupName, node Node) error
	SendGetItem(host string, kgname KeygroupName, id string) (Item, error)
	SendGetAllItems(host string, kgname KeygroupName) ([]Item, error)
}

type replicationService struct {
	c Client
	s *storeService
	n NameService
}

// newReplicationService creates a new handler for internal request (i.e. from peer nodes or the naming service).
// The nameservice makes sure that the information is synced with the other nodes
func newReplicationService(s *storeService, c Client, n NameService) *replicationService {
	service := &replicationService{
		s: s,
		c: c,
		n: n,
	}

	return service
}

// reportNodeFail report that this node was not able to receive the following kg/id
// If this returns an error the NaSe is also not reachable => this node is likely down
func (s *replicationService) reportNodeFail(nodeID NodeID, kg KeygroupName, id string) error {
	// This node is not able to reach node id
	// If it can reach the NaSe then maybe node id ist not available at the moment
	// So write that into the NaSe!
	log.Warn().Msg("reportNodeFail from replservice: could not reach node, reporting it as offline in the name service")
	repErr := s.n.ReportFailedNode(nodeID, kg, id)
	if repErr != nil {
		log.Error().Msg("reportNodeFail from replservice: this node is probably offline since it cant reach the NameService or the node...")
		return repErr
	}
	// Don't return an error here since the failed node is not stopping us from working...
	return nil
}

// createKeygroup creates the keygroup with the NaSe and saves its existence locally
func (s *replicationService) createKeygroup(k Keygroup) error {
	log.Debug().Msgf("CreateKeygroup from replservice: in keygroup=%#v", k)

	// Check if Keygroup already exists in NaSe
	exists, err := s.n.ExistsKeygroup(k.Name)
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
	err = s.n.CreateKeygroup(k.Name, k.Mutable, k.Expiry)
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

	exists, err := s.n.ExistsKeygroup(k.Name)
	if err != nil {
		return err
	}
	if !exists {
		err = errors.Errorf("no such keygroup according to NaSe: %#v", k)
		return err
	}

	// Inform all other nodes about the deletion so that they can delete their local copy
	ids, err := s.n.GetKeygroupMembers(k.Name, true)
	if err != nil {
		log.Err(err).Msg("Cannot delete keygroup because the nase threw an error")
		return err
	}

	for id := range ids {
		addr, err := s.n.GetNodeAddress(id)

		if err != nil {
			log.Err(err).Msg("Cannot Get node adress from NaSe")
			return err
		}

		log.Debug().Msgf("RelayDeleteKeygroup from replservice: sending %#v to %#v", k, addr)
		s.c.SendDeleteKeygroup(addr, k.Name)
		// No error handling because:
		// There is no need to report a node fail here since the node always gets this information from NaSe
		// It is also not necessary to fail here because of this
	}

	// Only now delete the keygroup with the nase
	err = s.n.DeleteKeygroup(k.Name)

	if err != nil {
		log.Err(err).Msg(err.(*errors.Error).ErrorStack())
		return err
	}

	return nil
}

// relayUpdate handles replication after requests to the Update endpoint of the client interface.
// It sends the update to all other nodes by calling their Update method
func (s *replicationService) relayUpdate(i Item) error {
	log.Debug().Msgf("RelayUpdate from replservice: in %#v", i)

	exists, err := s.n.ExistsKeygroup(i.Keygroup)
	if err != nil {
		return err
	}
	if !exists {
		err = errors.Errorf("no such keygroup according to NaSe: %#v", i.Keygroup)
		return err
	}

	// inform all other nodes about the update, get a list of all nodes subscribed to this keygroup
	ids, err := s.n.GetKeygroupMembers(i.Keygroup, true)
	if err != nil {
		log.Err(err).Msg("Cannot delete keygroup because the nase threw an error")
		return err
	}

	for id := range ids {
		addr, err := s.n.GetNodeAddress(id)

		if err != nil {
			log.Err(err).Msg("Cannot Get node adress from NaSe")
			return err
		}

		log.Debug().Msgf("RelayUpdate from replservice: sending %#v to %#v", i, addr)
		if err := s.c.SendUpdate(addr, i.Keygroup, i.ID, i.Val); err != nil {
			s.reportNodeFail(id, i.Keygroup, i.ID)
		}
	}

	return nil
}

// relayDelete handles replication after requests to the Delete endpoint of the client interface.
func (s *replicationService) relayDelete(i Item) error {
	log.Debug().Msgf("RelayDelete from replservice: in %#v", i)

	exists, err := s.n.ExistsKeygroup(i.Keygroup)
	if err != nil {
		return err
	}
	if !exists {
		err = errors.Errorf("no such keygroup according to NaSe: %#v", i.Keygroup)
		return err
	}

	ids, err := s.n.GetKeygroupMembers(i.Keygroup, true)

	for id := range ids {
		addr, err := s.n.GetNodeAddress(id)

		if err != nil {
			log.Err(err).Msg("Cannot Get node adress from NaSe")
			return err
		}

		log.Debug().Msgf("RelayDelete from replservice: sending %#v to %#v", i, addr)
		if err := s.c.SendDelete(addr, i.Keygroup, i.ID); err != nil {
			s.reportNodeFail(id, i.Keygroup, i.ID)
		}
	}

	return err
}

// addReplica handles replication after requests to the AddReplica endpoint. It relays this command if "relay" is set to "true".
func (s *replicationService) addReplica(k Keygroup, n Node, relay bool) error {
	log.Debug().Msgf("AddReplica from replservice: in kg=%#v no=%#v", k, n)

	// if relay is set to true, we got the request from the client interface
	// and are responsible to bring the new replica up to speed
	// (-> send them past data, send them all other replicas, inform all other replicas)
	// HOWEVER: if we are the new node, request that data from somewhere else
	if relay {

		exists, err := s.n.ExistsKeygroup(k.Name)
		if err != nil {
			return err
		}
		if !exists {
			err = errors.Errorf("no such keygroup according to NaSe: %#v", k)
			return err
		}

		// let's get the information about this new replica
		newNodeAddr, err := s.n.GetNodeAddress(n.ID)

		if err != nil {
			log.Err(err).Msg(err.(*errors.Error).ErrorStack())
			return err
		}

		newNode := &Node{
			ID:   n.ID,
			Host: newNodeAddr,
		}

		// Write the news into the NaSe first
		err = s.n.JoinNodeIntoKeygroup(k.Name, n.ID, k.Expiry)

		if err != nil {
			log.Err(err).Msg(err.(*errors.Error).ErrorStack())
			return err
		}

		// let's tell this new node that it should create a local copy of this keygroup
		err = s.c.SendCreateKeygroup(newNodeAddr, k.Name, k.Expiry)
		if err != nil {
			log.Err(err).Msg(err.(*errors.Error).ErrorStack())
			return err
		}

		// now let's iterate over all other currently known replicas for this node (except ourselves)
		// this also includes the newly added node, so it will receive a AddReplica with itself as the new node.
		ids, err := s.n.GetKeygroupMembers(k.Name, true)
		if err != nil {
			return err
		}

		for currID := range ids {
			// get a replica node
			replAddr, err := s.n.GetNodeAddress(currID)
			replNode := &Node{
				ID:   NodeID(currID),
				Host: replAddr,
			}

			if err != nil {
				return err
			}

			// tell that replica node about the new node
			log.Debug().Msgf("AddReplica from replservice: sending %#v to %#v", newNode, replNode)
			if err := s.c.SendAddReplica(replAddr, k.Name, *newNode, k.Expiry); err != nil {
				return err
			}

			// then tell the new node about that replica node
			log.Debug().Msgf("AddReplica from replservice: sending %#v to %#v", replNode, newNode)
			if err := s.c.SendAddReplica(newNode.Host, k.Name, *replNode, k.Expiry); err != nil {
				return err
			}
		}

		// send all existing data to the new node
		// TODO so this one array contains all data items? Maybe not a good idea if there is a lot of data to be sent
		var i []Item

		if n.ID != NodeID(s.n.GetNodeID()) {
			// we are adding a new node and we have all the data: send our data
			i, err = s.s.readAll(k.Name)

			if err != nil {
				log.Err(err).Msg(err.(*errors.Error).ErrorStack())
				return errors.Errorf("error adding replica")
			}
		} else {
			// oh no! We are the new node and have no data locally, let's request it from somewhere
			// take a victim with a higher expiry to request data from (but not ourselves!)
			_, addr := s.n.GetNodeWithBiggerExpiry(k.Name)

			if err != nil {
				return err
			}

			i, err = s.c.SendGetAllItems(addr, k.Name)

			if err != nil {
				log.Err(err).Msg(err.(*errors.Error).ErrorStack())
				return errors.Errorf("error adding replica")
			}
		}

		log.Debug().Msgf("AddReplica from replservice: About to send %d Elements to new node", len(i))
		for _, item := range i {
			// iterate over all data for that keygroup and send it to the new node
			// a batch might be better here
			log.Debug().Msgf("AddReplica from replservice: sending %#v to %#v", item, n)
			if err := s.c.SendUpdate(newNodeAddr, k.Name, item.ID, item.Val); err != nil {
				s.reportNodeFail(n.ID, k.Name, item.ID)
			}
		}

		return nil

	}
	// else {
	// 	// This request is not from the client interface, but from the internal one.
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
		removedNodeAddr, err := s.n.GetNodeAddress(n.ID)
		if err != nil {
			return err
		}
		removedNode := &Node{
			ID:   n.ID,
			Host: removedNodeAddr,
		}

		exists, err := s.n.ExistsKeygroup(k.Name)
		if !exists {
			err = errors.Errorf("no such keygroup according to NaSe: %#v", k)
			return err
		}

		if err != nil {
			return err
		}

		// let's tell this new node that it should delete the local copy of this keygroup
		err = s.c.SendDeleteKeygroup(removedNode.Host, k.Name)
		if err != nil {
			log.Err(err).Msg(err.(*errors.Error).ErrorStack())
			return err
		}

		// First get all Replicas of this Keygroup to send the update to them
		kgMembers, err := s.n.GetKeygroupMembers(k.Name, true)
		if err != nil {
			return err
		}
		// Now exit node from keygroup with nase.
		// Now one node (this node) will get the message that a node is removed from a keygroup
		// that it itself is not a member of ==> Delete the local copy (see after big if relay statement)
		err = s.n.ExitOtherNodeFromKeygroup(k.Name, n.ID)
		if err != nil {
			return err
		}

		for idToInform := range kgMembers {
			nodeToInformAddr, err := s.n.GetNodeAddress(idToInform)

			if err != nil {
				return err
			}

			log.Debug().Msgf("RemoveReplica from replservice: sending RemoveReplica %#v to %#v", k, nodeToInformAddr)
			if err := s.c.SendRemoveReplica(nodeToInformAddr, k.Name, Node{ID: n.ID}); err != nil {
				log.Err(err).Msg("Sending Remove Replica has failed, continuing anyways")
			}
		}
	}

	// TODO Nils delete the local cache copy that this node is in this keygroup

	return nil
}

// getNode returns the locally saved node with this ID
func (s *replicationService) getNode(n Node) (Node, error) {
	addr, err := s.n.GetNodeAddress(n.ID)

	n = Node{
		ID:   n.ID,
		Host: addr,
	}

	return n, err
}

// getNodes returns a list of all known nodes.
func (s *replicationService) getNodes() ([]Node, error) {
	return s.n.GetAllNodes()
}

// getReplica returns a list of all replica nodes for a given keygroup.
func (s *replicationService) getReplica(k Keygroup) (nodes []Node, expiries map[NodeID]int, err error) {
	log.Debug().Msgf("GetReplica from replservice: in %#v", k)

	exists, err := s.n.ExistsKeygroup(k.Name)
	if !exists {
		log.Err(err).Msg(err.(*errors.Error).ErrorStack())
		return nil, nil, errors.Errorf("no such keygroup according to NaSe")
	}
	if err != nil {
		return nil, nil, err
	}

	expiries = make(map[NodeID]int)

	ids, err := s.n.GetKeygroupMembers(k.Name, true)
	for id, expiry := range ids {
		addr, err := s.n.GetNodeAddress(id)

		if err != nil {
			return nil, nil, err
		}

		newNode := &Node{
			ID:   id,
			Host: addr,
		}
		// TODO is this the optimal solution?
		nodes = append(nodes, *newNode)
		expiries[id] = expiry
	}
	return
}
