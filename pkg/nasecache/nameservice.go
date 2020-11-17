package nasecache

import (
	"time"

	"github.com/allegro/bigcache"
	"github.com/go-errors/errors"
	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/fred"
)

const (
	sep = "|"
)

// NameServiceCache functions like NameService, but inserts local caching
type NameServiceCache struct {
	regularNase fred.NameService
	cache       *bigcache.BigCache
}

// NewNameServiceCache creates a new NameServiceCache
func NewNameServiceCache(regularNase fred.NameService) (*NameServiceCache, error) {
	cache, err := bigcache.NewBigCache(bigcache.DefaultConfig(40 * time.Second))

	if err != nil {
		return nil, errors.Errorf("Error initializing the cache")
	}

	return &NameServiceCache{
		regularNase: regularNase,
		cache:       cache,
	}, nil
}

// manage information about this node

// GetNodeID returns the ID of this node.
func (n *NameServiceCache) GetNodeID() fred.NodeID {
	return n.regularNase.GetNodeID()
}

// RegisterSelf stores information about this node
func (n *NameServiceCache) RegisterSelf(host string) error {
	return n.regularNase.RegisterSelf(host)
}

// ReportFailedNode report nonreachable node
func (n *NameServiceCache) ReportFailedNode(nodeID fred.NodeID, kg fred.KeygroupName, id string) error {
	return n.regularNase.ReportFailedNode(nodeID, kg, id)
}

// RequestNodeStatus Request whether this node has missed items
func (n *NameServiceCache) RequestNodeStatus(nodeID fred.NodeID) []fred.Item {
	return n.regularNase.RequestNodeStatus(nodeID)
}

// GetNodeWithBiggerExpiry get a node that has a bigger expiry of this keygroup than local
func (n *NameServiceCache) GetNodeWithBiggerExpiry(kg fred.KeygroupName) (nodeID fred.NodeID, addr string) {
	return n.regularNase.GetNodeWithBiggerExpiry(kg)
}