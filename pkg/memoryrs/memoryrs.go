package memoryrs

import (
	"errors"
	"net"
	"sync"

	"github.com/rs/zerolog/log"

	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/commons"
	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/replication"
)

type node struct {
	ip   net.IP
	port int
}

// ReplicationStorage saves a set of all known nodes.
type ReplicationStorage struct {
	nodes     map[replication.ID]node
	kg        map[commons.KeygroupName]map[replication.ID]struct{}
	nodesLock sync.RWMutex
	kgLock    sync.RWMutex
}

// New creates a new NodeStorage.
func New() (rS *ReplicationStorage) {
	rS = &ReplicationStorage{
		nodes: make(map[replication.ID]node),
		kg:    make(map[commons.KeygroupName]map[replication.ID]struct{}),
	}

	return
}

// CreateNode adds a node to the node storage in ReplicationStorage.
func (rS *ReplicationStorage) CreateNode(n replication.Node) error {
	log.Debug().Msgf("CreateNode from memoryns: in %v", n)
	rS.nodesLock.RLock()
	_, ok := rS.nodes[n.ID]
	rS.nodesLock.RUnlock()

	if ok {
		return nil
	}

	rS.nodesLock.Lock()
	rS.nodes[n.ID] = node{
		ip:   n.IP,
		port: n.Port,
	}
	rS.nodesLock.Unlock()

	return nil
}

// DeleteNode removes a node from the node storage in ReplicationStorage.
func (rS *ReplicationStorage) DeleteNode(n replication.Node) error {
	log.Debug().Msgf("DeleteNode from memoryns: in %v", n)
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
	log.Debug().Msgf("GetNode from memoryns: in %v", n)
	rS.nodesLock.RLock()
	node, ok := rS.nodes[n.ID]
	rS.nodesLock.RUnlock()

	if !ok {
		return n, errors.New("memoryrs: no such node")
	}

	return replication.Node{
		ID:   n.ID,
		IP:   node.ip,
		Port: node.port,
	}, nil
}

// ExistsNode checks whether a node exists in the node storage in ReplicationStorage.
func (rS *ReplicationStorage) ExistsNode(n replication.Node) bool {
	rS.nodesLock.RLock()
	_, ok := rS.nodes[n.ID]
	rS.nodesLock.RUnlock()

	log.Debug().Msgf("ExistsNode from memoryns: in %v, out %t", n, ok)

	return ok
}

// CreateKeygroup creates a keygroup in the keygroup storage in ReplicationStorage.
func (rS *ReplicationStorage) CreateKeygroup(k replication.Keygroup) error {
	log.Debug().Msgf("CreateKeygroup from memoryns: in %v", k)
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
	log.Debug().Msgf("DeleteKeygroup from memoryns: in %v", k)
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
	log.Debug().Msgf("GetKeygroup from memoryns: in %v", k)
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

	log.Debug().Msgf("ExistsKeygroup from memoryns: in %v, out %t", k, ok)

	return ok
}

// AddReplica adds a replica node to the keygroup in the keygroup storage in ReplicationStorage.
func (rS *ReplicationStorage) AddReplica(k replication.Keygroup, n replication.Node) error {
	log.Debug().Msgf("AddReplica from memoryns: in kg=%v no=%v", k, n)
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
		return nil
	}

	rS.kgLock.Lock()
	rS.kg[k.Name][n.ID] = struct{}{}
	rS.kgLock.Unlock()

	return nil
}

// RemoveReplica removes a replica node from the keygroup in the keygroup storage in ReplicationStorage.
func (rS *ReplicationStorage) RemoveReplica(k replication.Keygroup, n replication.Node) error {
	log.Debug().Msgf("RemoveReplica from memoryns: in kg=%v no=%v", k, n)
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
			IP:   node.ip,
			Port: node.port,
		}

		i++
	}

	log.Debug().Msgf("GetNodes from memoryns: found %d nodes", i)

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
			IP:   node.ip,
			Port: node.port,
		}

		i++
	}
	log.Debug().Msgf("GetReplica from memoryns: found %d nodes", i)

	return nodes, nil
}
