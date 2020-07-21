package fred

import (
	"sync"

	"github.com/rs/zerolog/log"
)

// Client is an interface to send replication messages across nodes.
type Client interface {
	SendCreateKeygroup(addr Address, port int, kgname KeygroupName) error
	SendDeleteKeygroup(addr Address, port int, kgname KeygroupName) error
	SendUpdate(addr Address, port int, kgname KeygroupName, id string, value string) error
	SendDelete(addr Address, port int, kgname KeygroupName, id string) error
	SendAddNode(addr Address, port int, node Node) error
	SendRemoveNode(addr Address, port int, node Node) error
	SendAddReplica(addr Address, port int, kgname KeygroupName, node Node) error
	SendRemoveReplica(addr Address, port int, kgname KeygroupName, node Node) error
	SendIntroduce(addr Address, port int, self Node, other Node, nodes []Node) error
	SendDetroduce(addr Address, port int) error
}

type replicationService struct {
	nodes     map[NodeID]Node
	kg        map[KeygroupName]map[NodeID]struct{}
	nodesLock sync.RWMutex
	kgLock    sync.RWMutex
	self      Node
	c         Client
}

// newReplicationService creates a new handler for internal request (i.e. from peer nodes or the naming service).
func newReplicationService(nodePort int, c Client) *replicationService {
	return &replicationService{
		nodes: make(map[NodeID]Node),
		kg:    make(map[KeygroupName]map[NodeID]struct{}),
		self: Node{
			Port: nodePort,
		},
		c: c,
	}
}

// CreateKeygroup handles replication after requests to the CreateKeygroup endpoint of the internal interface.
func (s *replicationService) CreateKeygroup(k Keygroup) error {
	log.Debug().Msgf("CreateKeygroup from replservice: in %#v", k)
	kg := Keygroup{
		Name: k.Name,
	}

	if s.existsKeygroup(kg) {
		return nil
	}

	log.Debug().Msgf("CreateKeygroup from memoryrs: in %#v", k)
	s.kgLock.RLock()
	_, ok := s.kg[kg.Name]
	s.kgLock.RUnlock()

	if ok {
		return nil
	}

	s.kgLock.Lock()
	s.kg[kg.Name] = make(map[NodeID]struct{})
	s.kgLock.Unlock()

	return nil
}

// DeleteKeygroup handles replication after requests to the DeleteKeygroup endpoint of the internal interface.
func (s *replicationService) DeleteKeygroup(k Keygroup) error {
	log.Debug().Msgf("RelayCreateKeygroup from replservice: in %#v", k)
	kg := Keygroup{
		Name: k.Name,
	}

	log.Debug().Msgf("DeleteKeygroup from memoryrs: in %#v", k)
	s.kgLock.RLock()
	_, ok := s.kg[kg.Name]
	s.kgLock.RUnlock()

	if !ok {
		return newError(StatusNotFound, "memoryrs: no such node")
	}

	s.kgLock.Lock()
	delete(s.kg, k.Name)
	s.kgLock.Unlock()

	return nil
}

// RelayDeleteKeygroup handles replication after requests to the DeleteKeygroup endpoint of the external interface.
func (s *replicationService) RelayDeleteKeygroup(k Keygroup) error {
	log.Debug().Msgf("RelayDeleteKeygroup from replservice: in %#v", k)
	kg := Keygroup{
		Name: k.Name,
	}

	if !s.existsKeygroup(kg) {
		log.Error().Msgf("RelayDeleteKeygroup from replservice: Keygroup does not exist: in %#v", k)
		return newError(StatusNotFound, "replservice: no such keygroup")
	}

	kg, err := s.getKeygroup(kg)

	if err != nil {
		return err
	}

	log.Debug().Msgf("DeleteKeygroup from memoryrs: in %#v", k)
	s.kgLock.RLock()
	_, ok := s.kg[kg.Name]
	s.kgLock.RUnlock()

	if !ok {
		return newError(StatusNotFound, "memoryrs: no such node")
	}

	s.kgLock.Lock()
	delete(s.kg, kg.Name)
	s.kgLock.Unlock()

	for rn := range kg.Replica {
		node, err := s.getNode(Node{
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
func (s *replicationService) RelayUpdate(i Item) error {
	log.Debug().Msgf("RelayUpdate from replservice: in %#v", i)
	kg := Keygroup{
		Name: i.Keygroup,
	}

	if !s.existsKeygroup(kg) {
		log.Error().Msgf("RelayUpdate from replservice: Keygroup does not exist: in %#v", i)
		return newError(StatusNotFound, "replservice: no such keygroup")
	}

	kg, err := s.getKeygroup(kg)

	if err != nil {
		return err
	}

	log.Debug().Msgf("RelayUpdate sending to: in %#v", kg.Replica)

	for rn := range kg.Replica {
		node, err := s.getNode(Node{
			ID: rn,
		})

		if err != nil {
			return err
		}

		log.Debug().Msgf("RelayUpdate from replservice: sending %#v to %#v", i, node)
		if err := s.c.SendUpdate(node.Addr, node.Port, kg.Name, i.ID, i.Val); err != nil {
			return err
		}
	}

	return err
}

// RelayDelete handles replication after requests to the Delete endpoint of the external interface.
func (s *replicationService) RelayDelete(i Item) error {
	log.Debug().Msgf("RelayDelete from replservice: in %#v", i)
	kg := Keygroup{
		Name: i.Keygroup,
	}

	if !s.existsKeygroup(kg) {
		log.Error().Msgf("RelayDelete from replservice: Keygroup does not exist: in %#v", i)
		return newError(StatusNotFound, "replservice: no such keygroup")
	}

	kg, err := s.getKeygroup(kg)

	if err != nil {
		return err
	}

	for rn := range kg.Replica {
		node, err := s.getNode(Node{
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
func (s *replicationService) AddReplica(k Keygroup, n Node, i []Item, relay bool) error {
	log.Debug().Msgf("AddReplica from replservice: in kg=%#v no=%#v", k, n)

	// first get the keygroup from the kgname in k
	kg := Keygroup{
		Name: k.Name,
	}

	if !s.existsKeygroup(kg) {
		log.Error().Msgf("AddReplica from replservice: Keygroup does not exist: in %#v", k)
		return newError(StatusNotFound, "replservice: no such keygroup")
	}

	kg, err := s.getKeygroup(kg)

	if err != nil {
		return err
	}

	// if relay is set to true, we got the request from the external interface
	// and are responsible to bring the new replica up to speed
	// (-> send them past val, send them all other replicas, inform all other replicas)
	if relay {
		// let's get the information about this new replica first
		newNode, err := s.getNode(n)

		if err != nil {
			return err
		}

		// tell this new node to create the keygroup they're now replicating
		log.Debug().Msgf("AddReplica from replservice: sending %#v to %#v", k, newNode)
		if err := s.c.SendCreateKeygroup(newNode.Addr, newNode.Port, kg.Name); err != nil {
			return err
		}

		// now tell this new node that we are also a replica node for that keygroup
		self, err := s.getSelf()

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
			replNode, err := s.getNode(Node{
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

		// request came from external client interface, send past val as well
		for _, item := range i {
			// iterate over all val for that keygroup and send it to the new node
			// a batch might be better here
			log.Debug().Msgf("AddReplica from replservice: sending %#v to %#v", item, n)
			if err := s.c.SendUpdate(newNode.Addr, newNode.Port, kg.Name, item.ID, item.Val); err != nil {
				return err
			}
		}
	}

	// finally, in either case, save that the new node is now also a replica for the keygroup
	log.Debug().Msgf("AddReplica from memoryrs: in kg=%#v no=%#v", k, n)

	s.kgLock.RLock()
	_, ok := s.kg[kg.Name]
	s.kgLock.RUnlock()

	if !ok {
		return newError(StatusNotFound, "memoryrs: no such keygroup")
	}

	s.nodesLock.RLock()
	_, ok = s.nodes[n.ID]
	s.nodesLock.RUnlock()
	if !ok {
		return newError(StatusNotFound, "memoryrs: no such node")
	}

	s.kgLock.RLock()
	_, ok = s.kg[kg.Name][n.ID]
	s.kgLock.RUnlock()

	if ok {
		return newError(StatusConflict, "memoryrs: node is already a replica node of keygroup")
	}

	s.kgLock.Lock()
	s.kg[k.Name][n.ID] = struct{}{}
	s.kgLock.Unlock()

	return nil
}

// RemoveReplica handles replication after requests to the RemoveReplica endpoint. It relays this command if "relay" is set to "true".
func (s *replicationService) RemoveReplica(k Keygroup, n Node, relay bool) error {
	log.Debug().Msgf("RemoveReplica from replservice: in kg=%#v no=%#v", k, n)

	kg := Keygroup{
		Name: k.Name,
	}

	if !s.existsKeygroup(kg) {
		log.Error().Msgf("RemoveReplica from replservice: Keygroup does not exist: in %#v", k)
		return newError(StatusNotFound, "replservice: no such keygroup")
	}

	kg, err := s.getKeygroup(kg)

	if err != nil {
		return err
	}

	log.Debug().Msgf("RemoveReplica from memoryrs: in kg=%#v no=%#v", k, n)
	s.kgLock.RLock()
	_, ok := s.kg[kg.Name]
	s.kgLock.RUnlock()

	if !ok {
		return newError(StatusNotFound, "memoryrs: no such keygroup")
	}

	s.kgLock.RLock()
	_, ok = s.kg[kg.Name][n.ID]
	s.kgLock.RUnlock()

	if !ok {
		return newError(StatusNotFound, "memoryrs: no such node")
	}

	s.kgLock.Lock()
	delete(s.kg[kg.Name], n.ID)
	s.kgLock.Unlock()

	if relay {
		node, err := s.getNode(n)

		if err != nil {
			return err
		}

		log.Debug().Msgf("RemoveReplica from replservice: sending %#v to %#v", k, node)
		if err := s.c.SendDeleteKeygroup(node.Addr, node.Port, kg.Name); err != nil {
			return err
		}

		for rn := range kg.Replica {
			node, err := s.getNode(Node{
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
func (s *replicationService) AddNode(n Node, relay bool) error {
	if relay {
		nodes, err := s.getNodes()

		if err != nil {
			return err
		}

		self, err := s.getSelf()

		if err != nil {
			return err
		}

		if err := s.c.SendIntroduce(n.Addr, n.Port, self, n, nodes); err != nil {
			log.Err(err).Msg("")
			return err
		}

		for _, rn := range nodes {
			node, err := s.getNode(Node{
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
	log.Debug().Msgf("CreateNode from memoryrs: in %#v", n)
	s.nodesLock.RLock()
	_, ok := s.nodes[n.ID]
	s.nodesLock.RUnlock()

	if ok {
		return nil
	}

	s.nodesLock.Lock()
	s.nodes[n.ID] = Node{
		Addr: n.Addr,
		Port: n.Port,
		ID:   n.ID,
	}
	s.nodesLock.Unlock()

	return nil
}

// RemoveNode handles replication after requests to the RemoveNode endpoint. It relays this command if "relay" is set to "true".
func (s *replicationService) RemoveNode(n Node, relay bool) error {

	log.Debug().Msgf("DeleteNode from memoryrs: in %#v", n)
	s.nodesLock.RLock()
	_, ok := s.nodes[n.ID]
	s.nodesLock.RUnlock()

	if !ok {
		return newError(StatusNotFound, "memoryrs: no such node")
	}

	s.nodesLock.Lock()
	delete(s.nodes, n.ID)
	s.nodesLock.Unlock()

	if relay {
		nodes, err := s.getNodes()

		if err != nil {
			return err
		}

		for _, rn := range nodes {
			node, err := s.getNode(Node{
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

func (s *replicationService) GetNode(n Node) (Node, error) {
	return s.getNode(n)
}

// GetReplica returns a list of all known nodes.
func (s *replicationService) GetNodes() ([]Node, error) {
	return s.getNodes()
}

// GetReplica returns a list of all replica nodes for a given keygroup.
func (s *replicationService) GetReplica(k Keygroup) ([]Node, error) {
	log.Debug().Msgf("GetReplica from replservice: in %#v", k)
	kg := Keygroup{
		Name: k.Name,
	}

	if !s.existsKeygroup(kg) {
		log.Error().Msgf("GetReplica from replservice: Keygroup does not exist: in %#v", k)
		return nil, newError(StatusNotFound, "replservice: no such keygroup")
	}

	kg, err := s.getKeygroup(kg)

	if err != nil {
		log.Err(err).Msgf("GetReplica from replservice: GetReplica did not work: in %#v", k)
		return nil, err
	}

	return s.getReplica(kg)
}

// getNode returns a node from the node storage in ReplicationStorage.
func (s *replicationService) getNode(n Node) (Node, error) {
	log.Debug().Msgf("GetNode from memoryrs: in %#v", n)
	s.nodesLock.RLock()
	node, ok := s.nodes[n.ID]
	s.nodesLock.RUnlock()

	if !ok {
		return n, newError(StatusNotFound, "memoryrs: no such node")
	}

	return Node{
		ID:   n.ID,
		Addr: node.Addr,
		Port: node.Port,
	}, nil
}

// getKeygroup returns a keygroup from the keygroup storage in ReplicationStorage.
func (s *replicationService) getKeygroup(k Keygroup) (Keygroup, error) {
	log.Debug().Msgf("GetKeygroup from memoryrs: in %#v", k)
	s.kgLock.RLock()
	replicas, ok := s.kg[k.Name]
	s.kgLock.RUnlock()

	if !ok {
		return k, newError(StatusNotFound, "memoryrs: no such node")
	}

	return Keygroup{
		Name:    k.Name,
		Replica: replicas,
	}, nil
}

// existsKeygroup checks whether a keygroup exists in the keygroup storage in ReplicationStorage.
func (s *replicationService) existsKeygroup(k Keygroup) bool {
	s.kgLock.RLock()
	_, ok := s.kg[k.Name]
	s.kgLock.RUnlock()

	log.Debug().Msgf("ExistsKeygroup from memoryrs: in %#v, out %t", k, ok)

	return ok
}

// getNodes returns all known replica nodes from the node storage in ReplicationStorage.
func (s *replicationService) getNodes() ([]Node, error) {
	s.nodesLock.RLock()
	defer s.nodesLock.RUnlock()

	nodes := make([]Node, len(s.nodes))

	i := 0
	for id, node := range s.nodes {
		nodes[i] = Node{
			ID:   id,
			Addr: node.Addr,
			Port: node.Port,
		}

		i++
	}

	log.Debug().Msgf("GetNodes from memoryrs: found %d nodes", i)

	return nodes, nil
}

// getReplica returns all known replica nodes for the given keygroup from the node storage in ReplicationStorage.
func (s *replicationService) getReplica(k Keygroup) ([]Node, error) {
	s.kgLock.RLock()
	defer s.kgLock.RUnlock()

	n, ok := s.kg[k.Name]

	if !ok {
		return nil, newError(StatusNotFound, "memoryrs: no such keygroup")
	}

	s.nodesLock.RLock()
	defer s.nodesLock.RUnlock()

	nodes := make([]Node, len(n))

	i := 0
	for id := range n {
		node := s.nodes[id]
		nodes[i] = Node{
			ID:   id,
			Addr: node.Addr,
			Port: node.Port,
		}

		i++
	}
	log.Debug().Msgf("GetReplica from memoryrs: found %d nodes", i)

	return nodes, nil
}

// getSelf returns val about the self node from ReplicationStorage.
func (s *replicationService) getSelf() (Node, error) {
	return s.self, nil
}
