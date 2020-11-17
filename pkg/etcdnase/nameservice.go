package etcdnase

import (
	"fmt"
	"time"

	"github.com/go-errors/errors"
	"github.com/rs/zerolog/log"
	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/fred"
	"go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/pkg/transport"
)

const (
	fmtKgNodeString         = "kg|%s|node|%s"
	fmtKgStatusString       = "kg|%s|status"
	fmtKgMutableString      = "kg|%s|mutable"
	fmtKgExpiryString       = "kg|%s|expiry|node|%s"
	fmtNodeAdressString     = "node|%s|address"
	fmtUserPermissionPrefix = "user|%s|kg"
	fmtUserPermissionString = "user|%s|kg|%s|method|%s"
	fmtFailedNodeKgString   = "failnode|%s|kg|%s|%s" // Node, Keygroup, ID
	fmtFailedNodePrefix     = "failnode|%s|"
	nodePrefixString        = "node|"
	sep                     = "|"
	timeout                 = 5 * time.Second
)

// NameService is the interface to the etcd server that serves as NaSe
// It is used by the replservice to sync updates to keygroups with other nodes and thereby makes sure that ReplicationStorage always has up to date information
type NameService struct {
	cli    *clientv3.Client
	NodeID string
}

// NewNameService creates a new NameService
func NewNameService(nodeID string, endpoints []string, certfFile string, keyFile string, caFile string) (*NameService, error) {
	tlsInfo := transport.TLSInfo{
		CertFile:      certfFile,
		KeyFile:       keyFile,
		TrustedCAFile: caFile,
	}

	tlsConfig, err := tlsInfo.ClientConfig()

	if err != nil {
		return nil, errors.Errorf("Error configuring certificates for the etcd client: %v", err)
	}

	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: 5 * time.Second,
		TLS:         tlsConfig,
	})

	if err != nil {
		// Deadline Exceeded
		return nil, errors.Errorf("Error starting the etcd client: %v", err)
	}

	return &NameService{
		cli: cli, NodeID: nodeID,
	}, nil
}

// RegisterSelf stores information about this node
func (n *NameService) RegisterSelf(host string) error {
	key := fmt.Sprintf(fmtNodeAdressString, n.NodeID)
	log.Debug().Msgf("NaSe: registering self as %s // %s", key, host)

	return n.put(key, host)
}

// GetNodeID returns the ID of this node.
func (n *NameService) GetNodeID() fred.NodeID {
	return fred.NodeID(n.NodeID)
}
