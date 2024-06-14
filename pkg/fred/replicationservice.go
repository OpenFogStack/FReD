package fred

import (
	"sync"
	"time"

	"git.tu-berlin.de/mcc-fred/vclock"
	"github.com/go-errors/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Client is an interface to send replication messages across nodes.
type Client interface {
	SendCreateKeygroup(host string, kgname KeygroupName, expiry int) error
	SendDeleteKeygroup(host string, kgname KeygroupName) error
	SendUpdate(host string, kgname KeygroupName, id string, value string, tombstoned bool, vvector vclock.VClock) error
	SendAppend(host string, kgname KeygroupName, id string, value string) error
	StreamItems(host string, kgname KeygroupName, append bool, items []Item) error
	SendGetItem(host string, kgname KeygroupName, id string) ([]Item, error)
	SendGetAllItems(host string, kgname KeygroupName) ([]Item, error)
}

type replicationService struct {
	c     Client
	s     *storeService
	n     NameService
	async bool
}

// newReplicationService creates a new handler for internal request (i.e. from peer nodes or the naming service).
// The nameservice makes sure that the information is synced with the other nodes.
// The experimental async flag determines whether the replication service should send the replication messages asynchronously.
// This only applies to data item updates (update, delete, append), not keygroup modification.
func newReplicationService(s *storeService, c Client, n NameService, async bool) *replicationService {
	service := &replicationService{
		s:     s,
		c:     c,
		n:     n,
		async: async,
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
	log.Trace().Msgf("CreateKeygroup from replservice: in keygroup=%+v", k)

	// Check if Keygroup already exists in NaSe
	exists, err := s.n.ExistsKeygroup(k.Name)
	if err != nil {
		log.Err(err).Msg("Error checking whether Kg exists in NaSe")
		return err
	}

	if exists {
		log.Error().Msg("Cannot Create Keygroup since NaSe says the Kg already exists. Existing Keygroups can only be joined")
		err = errors.Errorf("keygroup already exists, cannot create it: %+v", k)
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
	log.Trace().Msgf("DeleteKeygroup from replservice (does nothing): in %+v", k)

	// Deleting the Keygroup with the NaSe happens in RelayDeleteKeygroup, so this would just be double duty
	// err := s.nase.DeleteKeygroup(k.Name)

	// TODO delete the local copy of the keygroup here

	return nil
}

// relayDeleteKeygroup deletes a keygroup locally and relays the deletion of a keygroup to all other nodes in this keygroup
// Calls DeleteKeygroup on every other node that is in this keygroup
// TODO is this really necessary with the nase? Maybe
func (s *replicationService) relayDeleteKeygroup(k Keygroup) error {
	log.Trace().Msgf("RelayDeleteKeygroup from replservice: in %+v", k)

	exists, err := s.n.ExistsKeygroup(k.Name)
	if err != nil {
		return err
	}
	if !exists {
		err = errors.Errorf("no such keygroup according to NaSe: %+v", k)
		return err
	}

	// Inform all other nodes about the deletion so that they can delete their local copy
	ids, err := s.n.GetKeygroupMembers(k.Name, true)
	if err != nil {
		log.Err(err).Msg("Cannot delete keygroup because the nase threw an error")
		return err
	}

	addrs := make(map[NodeID]string)

	for id := range ids {
		addr, err := s.n.GetNodeAddress(id)

		if err != nil {
			log.Err(err).Msg("Cannot Get node address from NaSe")
			return err
		}

		addrs[id] = addr
	}

	// wait group to wait for all the other nodes to delete the keygroup
	var wg sync.WaitGroup

	for id, addr := range addrs {
		wg.Add(1)
		go func(addr string, _ NodeID) {
			defer wg.Done()
			log.Trace().Msgf("RelayDeleteKeygroup from replservice: sending %+v to %+v", k, addr)
			err = s.c.SendDeleteKeygroup(addr, k.Name)

			if err != nil {
				log.Trace().Msg(err.Error())
			}
		}(addr, id)
	}

	wg.Wait()

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
	log.Trace().Msgf("RelayUpdate from replservice: in %+v", i)

	// inform all other nodes about the update, get a list of all nodes subscribed to this keygroup
	ids, err := s.n.GetKeygroupMembers(i.Keygroup, true)
	if err != nil {
		log.Err(err).Msg("Cannot relayUpdate because the nase threw an error")
		return err
	}

	addrs := make(map[NodeID]string)

	for id := range ids {
		addr, err := s.n.GetNodeAddress(id)

		if err != nil {
			log.Err(err).Msg("Cannot Get node address from NaSe")
			return err
		}

		addrs[id] = addr
	}

	// wait group to wait for all the other nodes to relay update the keygroup
	var wg sync.WaitGroup

	for id, addr := range addrs {
		wg.Add(1)
		go func(addr string, id NodeID) {
			defer wg.Done()

			log.Trace().Msgf("RelayUpdate from replservice: sending %+v to %+v", i, addr)
			if err := s.c.SendUpdate(addr, i.Keygroup, i.ID, i.Val, i.Tombstoned, i.Version); err != nil {
				err = s.reportNodeFail(id, i.Keygroup, i.ID)

				if err != nil {
					log.Trace().Msg(err.Error())
				}
			}
		}(addr, id)
	}

	if !s.async {
		wg.Wait()
	}

	return nil
}

// relayAppend handles replication after requests to the Append endpoint of the client interface.
// It sends the append to all other nodes by calling their Append method
func (s *replicationService) relayAppend(i Item) error {
	log.Trace().Msgf("relayAppend from replservice: in %+v", i)

	exists, err := s.n.ExistsKeygroup(i.Keygroup)
	if err != nil {
		return err
	}
	if !exists {
		err = errors.Errorf("no such keygroup according to NaSe: %+v", i.Keygroup)
		return err
	}

	// inform all other nodes about the append, get a list of all nodes subscribed to this keygroup
	ids, err := s.n.GetKeygroupMembers(i.Keygroup, true)
	if err != nil {
		log.Err(err).Msg("Cannot relayAppend because the nase threw an error")
		return err
	}

	addrs := make(map[NodeID]string)

	for id := range ids {
		addr, err := s.n.GetNodeAddress(id)

		if err != nil {
			log.Err(err).Msg("Cannot Get node address from NaSe")
			return err
		}

		addrs[id] = addr
	}

	// wait group to wait for all the other nodes to relay append the keygroup
	var wg sync.WaitGroup

	for id, addr := range addrs {
		wg.Add(1)
		go func(addr string, id NodeID) {
			defer wg.Done()

			log.Trace().Msgf("relayAppend from replservice: sending %+v to %+v", i, addr)
			if err := s.c.SendAppend(addr, i.Keygroup, i.ID, i.Val); err != nil {
				err = s.reportNodeFail(id, i.Keygroup, i.ID)

				if err != nil {
					log.Trace().Msg(err.Error())
				}
			}
		}(addr, id)
	}

	if !s.async {
		wg.Wait()
	}

	return nil
}

// addReplica handles replication after requests to the AddReplica endpoint.
func (s *replicationService) addReplica(k Keygroup, n Node) error {
	log.Debug().Msgf("AddReplica from replservice: in kg=%+v no=%+v", k, n)

	// we got the request from the client interface
	// and are responsible to bring the new replica up to speed
	// (-> send them past data, send them all other replicas, inform all other replicas)
	// HOWEVER: if we are the new node, request that data from somewhere else
	// this method is not called internally as the naming service handles that

	exists, err := s.n.ExistsKeygroup(k.Name)
	if err != nil {
		return err
	}
	if !exists {
		err = errors.Errorf("no such keygroup according to NaSe: %+v", k)
		return err
	}

	// let's get the information about this new replica
	newNodeAddr, err := s.n.GetNodeAddress(n.ID)

	if err != nil {
		log.Err(err).Msg(err.(*errors.Error).ErrorStack())
		return err
	}

	log.Debug().Msgf("AddReplica from replservice: new node address is %+v", newNodeAddr)

	// if we are not the new node and we are also not a replica for that keygroup, we have nothing to do with that
	// request: abort
	if n.ID != s.n.GetNodeID() && !s.s.existsKeygroup(k.Name) {
		return errors.Errorf("keygroup %+v is unknown and we are not the new node %+v", k, n)
	}

	// Write the news into the NaSe first
	err = s.n.JoinNodeIntoKeygroup(k.Name, n.ID, k.Expiry)

	if err != nil {
		log.Err(err).Msg(err.(*errors.Error).ErrorStack())
		return err
	}

	log.Debug().Msgf("AddReplica from replservice: joined node %+v into keygroup %+v", n, k)

	// let's tell this new node that it should create a local copy of this keygroup
	err = s.c.SendCreateKeygroup(newNodeAddr, k.Name, k.Expiry)
	if err != nil {
		log.Err(err).Msg(err.(*errors.Error).ErrorStack())
		return err
	}

	// send all existing data to the new node
	// there are three basic possibilities:
	// we are a replica for this keygroup and we need to send someone else the data
	// we are not a replica for this keygroup and we need to get the data from somewhere
	// we are not a replica but we are also not the new replica
	// in the last case we just have nothing to do with this and error out
	// TODO so this one array contains all data items? Maybe not a good idea if there is a lot of data to be sent
	var i []Item

	readStart := time.Now()

	if n.ID != s.n.GetNodeID() {
		log.Debug().Msgf("AddReplica from replservice: we are not the new node (self %s != new node %s), so we need to send data to the new node", s.n.GetNodeID(), n.ID)
		// we are adding a new node and we have all the data: send our data
		i, err = s.s.readAll(k.Name)

		if err != nil {
			log.Err(err).Msg(err.(*errors.Error).ErrorStack())
			return errors.Errorf("error adding replica of keygroup %s to %s", k.Name, newNodeAddr)
		}

		log.Debug().Msgf("AddReplica from replservice: About to send %d elements from keygroup %s to new node %s", len(i), k.Name, newNodeAddr)
	} else {
		// oh no! We are the new node and have no data locally, let's request it from somewhere
		// take a victim with a higher expiry to request data from (but not ourselves!)
		_, addr := s.n.GetNodeWithBiggerExpiry(k.Name)

		if addr == "" {
			log.Error().Msgf("AddReplica: Cannot find node to get data for keygroup %s from, so this keygroup is empty", k.Name)
			// TODO is this an error? because there is nothing you can do (except tell the user to not fuck up their replica placement)
			return nil
		}

		log.Debug().Msgf("got node %s as replication source for keygroup %s", addr, k.Name)

		i, err = s.c.SendGetAllItems(addr, k.Name)

		if err != nil {
			log.Err(err).Msg(err.(*errors.Error).ErrorStack())
			return errors.Errorf("error adding replica of keygroup %s to %s", k.Name, newNodeAddr)
		}

		log.Debug().Msgf("AddReplica got %d items to replicate keygroup %s from %s to %s", len(i), k.Name, addr, newNodeAddr)
	}

	log.Debug().Msgf("AddReplica from replservice: About to send %d elements for keygroup %s to new node %s (took %s to get then all)", len(i), k.Name, newNodeAddr, time.Since(readStart))
	log.Debug().Msgf("Replicating %d items from keygroup %s to %s", len(i), k.Name, newNodeAddr)

	mutable, err := s.n.IsMutable(k.Name)

	if err != nil {
		return err
	}

	// an expensive operation
	if zerolog.GlobalLevel() == zerolog.DebugLevel {
		totalSize := 0
		for _, item := range i {
			totalSize += len(item.ID) + len(item.Val)
		}
		log.Debug().Msgf("Total transfer size is %d characters for keygroup %s to node %s", totalSize, k.Name, newNodeAddr)
	}

	writeStart := time.Now()
	err = s.c.StreamItems(newNodeAddr, k.Name, !mutable, i)

	if err != nil {
		return err
	}

	log.Debug().Msgf("Sent %d items from keygroup %s to %s (took %s)", len(i), k.Name, newNodeAddr, time.Since(writeStart))

	return nil

}

// removeReplica handles replication after requests to the RemoveReplica endpoint
func (s *replicationService) removeReplica(k Keygroup, n Node) error {
	log.Trace().Msgf("RemoveReplica from replservice: in kg=%+v no=%+v", k, n)

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
		err = errors.Errorf("no such keygroup according to NaSe: %+v", k)
		return err
	}
	if err != nil {
		return err
	}

	members, err := s.n.GetKeygroupMembers(k.Name, false)
	if err != nil {
		return err
	}
	// Check if to-be-deleted node is in the keygorup
	if _, ok := members[n.ID]; !ok {
		return errors.Errorf("can not remove node from keygroup it is not a member of")
	}
	// Check if the node is the last member of the keygroup (so the only one holding data)
	if len(members) <= 1 {
		log.Error().Msgf("trying to exit the only node left from the keygroup. please delete keygroup instead")
		return errors.Errorf("can not exit last node from keygroup %s. maybe you want to delete the keygroup instead", k.Name)
	}

	// let's tell this new node that it should delete the local copy of this keygroup
	err = s.c.SendDeleteKeygroup(removedNode.Host, k.Name)
	if err != nil {
		log.Err(err).Msg(err.(*errors.Error).ErrorStack())
		return err
	}

	// Now exit node from keygroup with nase.
	// Now one node (this node) will get the message that a node is removed from a keygroup
	// that it itself is not a member of ==> Delete the local copy (see after big if relay statement)
	err = s.n.ExitOtherNodeFromKeygroup(k.Name, n.ID)
	if err != nil {
		return err
	}

	return nil
}

// getNode returns the locally saved node with this ID
func (s *replicationService) getNodeExternal(n Node) (Node, error) {
	addr, err := s.n.GetNodeAddressExternal(n.ID)
	n = Node{
		ID:   n.ID,
		Host: addr,
	}

	return n, err
}

// getNodesExternalAdress returns a list of all known nodes with the address they can be reached at externally.
func (s *replicationService) getNodesExternalAdress() ([]Node, error) {
	return s.n.GetAllNodesExternal()
}

// getReplicaExternal returns a list of all replica nodes for a given keygroup.
func (s *replicationService) getReplicaExternal(k Keygroup) (nodes []Node, expiries map[NodeID]int, err error) {
	log.Trace().Msgf("GetReplicaExternal from replservice: in %+v", k)

	exists, err := s.n.ExistsKeygroup(k.Name)
	if !exists {
		log.Err(err).Msg(err.(*errors.Error).ErrorStack())
		return nil, nil, errors.Errorf("no such keygroup according to NaSe")
	}
	if err != nil {
		return nil, nil, err
	}

	expiries = make(map[NodeID]int)

	ids, err := s.n.GetKeygroupMembers(k.Name, false)
	log.Trace().Msgf("...got Nodes: %+v", ids)
	for id, expiry := range ids {
		addr, err := s.n.GetNodeAddressExternal(id)

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
