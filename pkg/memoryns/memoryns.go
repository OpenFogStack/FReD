package memoryns

import (
	"errors"
	"net"
	"sync"

	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/replication"
)

type node struct {
	ip net.IP
	port int
}

type nodes struct {
	nodes map[replication.ID]node
	sync.RWMutex
}

type kg struct {
	kg map[replication.KeygroupName]map[replication.ID]struct{}
	sync.RWMutex
}

// ReplicationStorage saves a set of all known nodes.
type ReplicationStorage struct {
	nodes nodes
	kg kg
	sync.RWMutex
}

// New creates a new NodeStorage.
func New() (rS *ReplicationStorage) {
	rS = &ReplicationStorage{
		nodes: nodes{
			nodes: make(map[replication.ID]node),
		},
		kg: kg{
				kg: make(map[replication.KeygroupName]map[replication.ID]struct{}),
		},
	}

	return
}

// CreateNode adds a node to the node storage in ReplicationStorage.
func (nS *ReplicationStorage) CreateNode(n replication.Node) error {
	nS.nodes.RLock()
	_, ok := nS.nodes.nodes[n.ID]
	nS.nodes.RUnlock()

	if ok {
		return nil
	}

	nS.nodes.Lock()
	nS.nodes.nodes[n.ID] = node{
		ip: n.IP,
		port: n.Port,
	}
	nS.nodes.Unlock()

	return nil
}

// DeleteNode removes a node from the node storage in ReplicationStorage.
func (nS *ReplicationStorage) DeleteNode(n replication.Node) error {
	nS.nodes.RLock()
	_, ok := nS.nodes.nodes[n.ID]
	nS.nodes.RUnlock()

	if !ok {
		return errors.New("memoryns: no such node")
	}

	nS.nodes.Lock()
	delete(nS.nodes.nodes, n.ID)
	nS.nodes.Unlock()

	return nil
}

// GetNode returns a node from the node storage in ReplicationStorage.
func (nS *ReplicationStorage) GetNode(n replication.Node) (replication.Node, error) {
	nS.nodes.RLock()
	node, ok := nS.nodes.nodes[n.ID]
	nS.nodes.RUnlock()

	if !ok {
		return n, errors.New("memoryns: no such node")
	}

	return replication.Node{
		ID:   n.ID,
		IP:   node.ip,
		Port: node.port,
	}, nil
}

// ExistsNode checks whether a node exists in the node storage in ReplicationStorage.
func (nS *ReplicationStorage) ExistsNode(n replication.Node) bool {
	nS.nodes.RLock()
	_, ok := nS.nodes.nodes[n.ID]
	nS.nodes.RUnlock()

	return ok
}

// CreateKeygroup creates a keygroup in the keygroup storage in ReplicationStorage.
func (nS *ReplicationStorage) CreateKeygroup(k replication.Keygroup) error {
	nS.kg.RLock()
	_, ok := nS.kg.kg[k.Name]
	nS.kg.RUnlock()

	if ok {
		return nil
	}

	nS.kg.Lock()
	nS.kg.kg[k.Name] = make(map[replication.ID]struct{})
	nS.kg.Unlock()

	return nil
}

// DeleteKeygroup removes a keygroup from the keygroup storage in ReplicationStorage.
func (nS *ReplicationStorage) DeleteKeygroup(k replication.Keygroup) error {
	nS.kg.RLock()
	_, ok := nS.kg.kg[k.Name]
	nS.kg.RUnlock()

	if !ok {
		return errors.New("memoryns: no such node")
	}

	nS.kg.Lock()
	delete(nS.kg.kg, k.Name)
	nS.kg.Unlock()

	return nil
}

// GetKeygroup returns a keygroup from the keygroup storage in ReplicationStorage.
func (nS *ReplicationStorage) GetKeygroup(k replication.Keygroup) (replication.Keygroup, error) {
	nS.kg.RLock()
	replicas, ok := nS.kg.kg[k.Name]
	nS.kg.RUnlock()

	if !ok {
		return k, errors.New("memoryns: no such node")
	}

	return replication.Keygroup{
		Name:    k.Name,
		Replica: replicas,
	}, nil
}

// ExistsKeygroup checks whether a keygroup exists in the keygroup storage in ReplicationStorage.
func (nS *ReplicationStorage) ExistsKeygroup(k replication.Keygroup) bool {
	nS.kg.RLock()
	_, ok := nS.kg.kg[k.Name]
	nS.kg.RUnlock()

	return ok
}

// AddReplica adds a replica node to the keygroup in the keygroup storage in ReplicationStorage.
func (nS *ReplicationStorage) AddReplica(k replication.Keygroup, n replication.Node) error {
	nS.kg.RLock()
	_, ok := nS.kg.kg[k.Name]
	nS.kg.RUnlock()

	if ok {
		return errors.New("memoryns: no such keygroup")
	}

	nS.kg.RLock()
	_, ok = nS.kg.kg[k.Name][n.ID]
	nS.kg.RUnlock()

	if ok {
		return nil
	}

	nS.kg.Lock()
	nS.kg.kg[k.Name][n.ID] = struct{}{}
	nS.kg.Unlock()

	return nil
}

// RemoveReplica removes a replica node from the keygroup in the keygroup storage in ReplicationStorage.
func (nS *ReplicationStorage) RemoveReplica(k replication.Keygroup, n replication.Node) error {
	nS.kg.RLock()
	_, ok := nS.kg.kg[k.Name]
	nS.kg.RUnlock()

	if ok {
		return errors.New("memoryns: no such keygroup")
	}

	nS.kg.RLock()
	_, ok = nS.kg.kg[k.Name][n.ID]
	nS.kg.RUnlock()

	if ok {
		return errors.New("memoryns: no such node")
	}

	nS.kg.Lock()
	delete(nS.kg.kg[k.Name], n.ID)
	nS.kg.Unlock()

	return nil
}