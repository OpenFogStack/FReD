package memoryrs

import (
	"errors"
	"sync"

	"github.com/rs/zerolog/log"

	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/commons"
	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/replication"
)

type node struct {
	addr replication.Address
	port int
}

// ReplicationStorage saves a set of all known nodes.
type ReplicationStorage struct {
	nodes     map[replication.ID]node
	kg        map[commons.KeygroupName]map[replication.ID]struct{}
	nodesLock sync.RWMutex
	kgLock    sync.RWMutex
	self      replication.Node
	needsSeed bool
}

// New creates a new NodeStorage.
func New(nodePort int) (rS *ReplicationStorage) {
	rS = &ReplicationStorage{
		nodes: make(map[replication.ID]node),
		kg:    make(map[commons.KeygroupName]map[replication.ID]struct{}),
		self: replication.Node{
			Port: nodePort,
		},
		needsSeed: true,
	}

	return
}

// CreateNode adds a node to the node storage in ReplicationStorage.
func (rS *ReplicationStorage) CreateNode(n replication.Node) error {
	if rS.needsSeed {
		rS.needsSeed = false
	}

	log.Debug().Msgf("CreateNode from memoryrs: in %#v", n)
	rS.nodesLock.RLock()
	_, ok := rS.nodes[n.ID]
	rS.nodesLock.RUnlock()

	if ok {
		return nil
	}

	rS.nodesLock.Lock()
	rS.nodes[n.ID] = node{
		addr: n.Addr,
		port: n.Port,
	}
	rS.nodesLock.Unlock()

	return nil
}

// DeleteNode removes a node from the node storage in ReplicationStorage.
func (rS *ReplicationStorage) DeleteNode(n replication.Node) error {
	log.Debug().Msgf("DeleteNode from memoryrs: in %#v", n)
	rS.nodesLock.RLock()
	_, ok := rS.nodes[n.ID]
	rS.nodesLock.RUnlock()

	if !ok {
		return errors.New("memoryrs: no such node")
	}

	rS.nodesLock.Lock()
	delete(rS.nodes, n.ID)
	rS.nodesLock.Unlock()

	return nil
}

// GetNode returns a node from the node storage in ReplicationStorage.
func (rS *ReplicationStorage) GetNode(n replication.Node) (replication.Node, error) {
	log.Debug().Msgf("GetNode from memoryrs: in %#v", n)
	rS.nodesLock.RLock()
	node, ok := rS.nodes[n.ID]
	rS.nodesLock.RUnlock()

	if !ok {
		return n, errors.New("memoryrs: no such node")
	}

	return replication.Node{
		ID:   n.ID,
		Addr: node.addr,
		Port: node.port,
	}, nil
}

// ExistsNode checks whether a node exists in the node storage in ReplicationStorage.
func (rS *ReplicationStorage) ExistsNode(n replication.Node) bool {
	rS.nodesLock.RLock()
	_, ok := rS.nodes[n.ID]
	rS.nodesLock.RUnlock()

	log.Debug().Msgf("ExistsNode from memoryrs: in %#v, out %t", n, ok)

	return ok
}

// CreateKeygroup creates a keygroup in the keygroup storage in ReplicationStorage.
func (rS *ReplicationStorage) CreateKeygroup(k replication.Keygroup) error {
	log.Debug().Msgf("CreateKeygroup from memoryrs: in %#v", k)
	rS.kgLock.RLock()
	_, ok := rS.kg[k.Name]
	rS.kgLock.RUnlock()

	if ok {
		return nil
	}

	rS.kgLock.Lock()
	rS.kg[k.Name] = make(map[replication.ID]struct{})
	rS.kgLock.Unlock()

	return nil
}

// DeleteKeygroup removes a keygroup from the keygroup storage in ReplicationStorage.
func (rS *ReplicationStorage) DeleteKeygroup(k replication.Keygroup) error {
	log.Debug().Msgf("DeleteKeygroup from memoryrs: in %#v", k)
	rS.kgLock.RLock()
	_, ok := rS.kg[k.Name]
	rS.kgLock.RUnlock()

	if !ok {
		return errors.New("memoryrs: no such node")
	}

	rS.kgLock.Lock()
	delete(rS.kg, k.Name)
	rS.kgLock.Unlock()

	return nil
}

// GetKeygroup returns a keygroup from the keygroup storage in ReplicationStorage.
func (rS *ReplicationStorage) GetKeygroup(k replication.Keygroup) (replication.Keygroup, error) {
	log.Debug().Msgf("GetKeygroup from memoryrs: in %#v", k)
	rS.kgLock.RLock()
	replicas, ok := rS.kg[k.Name]
	rS.kgLock.RUnlock()

	if !ok {
		return k, errors.New("memoryrs: no such node")
	}

	return replication.Keygroup{
		Name:    k.Name,
		Replica: replicas,
	}, nil
}

// ExistsKeygroup checks whether a keygroup exists in the keygroup storage in ReplicationStorage.
func (rS *ReplicationStorage) ExistsKeygroup(k replication.Keygroup) bool {
	rS.kgLock.RLock()
	_, ok := rS.kg[k.Name]
	rS.kgLock.RUnlock()

	log.Debug().Msgf("ExistsKeygroup from memoryrs: in %#v, out %t", k, ok)

	return ok
}

// AddReplica adds a replica node to the keygroup in the keygroup storage in ReplicationStorage.
func (rS *ReplicationStorage) AddReplica(k replication.Keygroup, n replication.Node) error {
	log.Debug().Msgf("AddReplica from memoryrs: in kg=%#v no=%#v", k, n)

	rS.kgLock.RLock()
	_, ok := rS.kg[k.Name]
	rS.kgLock.RUnlock()

	if !ok {
		return errors.New("memoryrs: no such keygroup")
	}

	rS.nodesLock.RLock()
	_, ok = rS.nodes[n.ID]
	rS.nodesLock.RUnlock()
	if !ok {
		return errors.New("memoryrs: no such node")
	}

	rS.kgLock.RLock()
	_, ok = rS.kg[k.Name][n.ID]
	rS.kgLock.RUnlock()

	if ok {
		return errors.New("memoryrs: node is already a replica node of keygroup")
	}

	rS.kgLock.Lock()
	rS.kg[k.Name][n.ID] = struct{}{}
	rS.kgLock.Unlock()

	return nil
}

// RemoveReplica removes a replica node from the keygroup in the keygroup storage in ReplicationStorage.
func (rS *ReplicationStorage) RemoveReplica(k replication.Keygroup, n replication.Node) error {
	log.Debug().Msgf("RemoveReplica from memoryrs: in kg=%#v no=%#v", k, n)
	rS.kgLock.RLock()
	_, ok := rS.kg[k.Name]
	rS.kgLock.RUnlock()

	if !ok {
		return errors.New("memoryrs: no such keygroup")
	}

	rS.kgLock.RLock()
	_, ok = rS.kg[k.Name][n.ID]
	rS.kgLock.RUnlock()

	if !ok {
		return errors.New("memoryrs: no such node")
	}

	rS.kgLock.Lock()
	delete(rS.kg[k.Name], n.ID)
	rS.kgLock.Unlock()

	return nil
}

// GetNodes returns all known replica nodes from the node storage in ReplicationStorage.
func (rS *ReplicationStorage) GetNodes() ([]replication.Node, error) {
	rS.nodesLock.RLock()
	defer rS.nodesLock.RUnlock()

	nodes := make([]replication.Node, len(rS.nodes))

	i := 0
	for id, node := range rS.nodes {
		nodes[i] = replication.Node{
			ID:   id,
			Addr: node.addr,
			Port: node.port,
		}

		i++
	}

	log.Debug().Msgf("GetNodes from memoryrs: found %d nodes", i)

	return nodes, nil
}

// GetReplica returns all known replica nodes for the given keygroup from the node storage in ReplicationStorage.
func (rS *ReplicationStorage) GetReplica(k replication.Keygroup) ([]replication.Node, error) {
	rS.kgLock.RLock()
	defer rS.kgLock.RUnlock()

	n, ok := rS.kg[k.Name]

	if !ok {
		return nil, errors.New("memoryrs: no such keygroup")
	}

	rS.nodesLock.RLock()
	defer rS.nodesLock.RUnlock()

	nodes := make([]replication.Node, len(n))

	i := 0
	for id := range n {
		node := rS.nodes[id]
		nodes[i] = replication.Node{
			ID:   id,
			Addr: node.addr,
			Port: node.port,
		}

		i++
	}
	log.Debug().Msgf("GetReplica from memoryrs: found %d nodes", i)

	return nodes, nil
}

// Seed seeds the node as a first node in the system and supplies it with information about itself.
func (rS *ReplicationStorage) Seed(n replication.Node) error {
	if !rS.needsSeed {
		return errors.New("memoryrs: node is already seeded or does not need seed")
	}

	rS.needsSeed = false

	rS.self = replication.Node{
		ID:   n.ID,
		Addr: n.Addr,
		Port: rS.self.Port,
	}

	return nil
}

// Unseed removes all data about self node from ReplicationStorage and sets "needsSeed" back to true, effectively removing the node from the system.
func (rS *ReplicationStorage) Unseed() error {
	if rS.needsSeed {
		return errors.New("memoryrs: node is already unseeded or needs seed")
	}

	rS.needsSeed = true

	rS.self = replication.Node{
		Port: rS.self.Port,
	}

	return nil

}

// GetSelf returns data about the self node from ReplicationStorage.
func (rS *ReplicationStorage) GetSelf() (replication.Node, error) {
	if rS.needsSeed {
		return replication.Node{}, errors.New("memoryrs: cannot return self, needs seed")
	}

	return rS.self, nil
}
