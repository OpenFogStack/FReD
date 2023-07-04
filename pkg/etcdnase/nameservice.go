package etcdnase

import (
	"fmt"
	"time"

	"git.tu-berlin.de/mcc-fred/fred/pkg/fred"
	"git.tu-berlin.de/mcc-fred/fred/pkg/grpcutil"
	"github.com/dgraph-io/ristretto"
	"github.com/go-errors/errors"
	"github.com/rs/zerolog/log"
	"go.etcd.io/etcd/client/pkg/v3/transport"
	"go.etcd.io/etcd/client/v3"
)

const (
	fmtKgNodeStringPrefix         = "kg|%s|node|"
	fmtKgStatusString             = "kg|%s|status"
	fmtKgMutableString            = "kg|%s|mutable"
	fmtKgExpiryStringPrefix       = "kg|%s|expiry|node|"
	fmtNodeAdressString           = "node|%s|address"
	fmtNodeExternalAdressString   = "node|%s|extaddress"
	fmtUserPermissionStringPrefix = "user|%s|kg|%s|method|"
	fmtFailedNodeKgStringPrefix   = "failnode|%s|kg|%s|" // Node, Keygroup, ID
	fmtFailedNodePrefix           = "failnode|%s|"
	nodePrefixString              = "node|"
	sep                           = "|"
	timeout                       = 5 * time.Second
)

// NameService is the interface to the etcd server that serves as NaSe
// It is used by the replservice to sync updates to keygroups with other nodes and thereby makes sure that ReplicationStorage always has up to date information
type NameService struct {
	cli     *clientv3.Client
	watcher clientv3.Watcher
	local   *ristretto.Cache
	cached  bool
	NodeID  string
}

// NewNameService creates a new NameService
func NewNameService(nodeID string, endpoints []string, certFile string, keyFile string, caFile string, cached bool) (*NameService, error) {

	_, _, err := grpcutil.GetCreds(certFile, keyFile, []string{caFile}, false)

	if err != nil {
		return nil, errors.Errorf("Error configuring certificates for the etcd client: %v", err)
	}

	tlsInfo := transport.TLSInfo{
		CertFile:      certFile,
		KeyFile:       keyFile,
		TrustedCAFile: caFile,
	}

	tlsConfig, err := tlsInfo.ClientConfig()

	if err != nil {
		return nil, errors.Errorf("Error configuring certificates for the etcd client: %v", err)
	}

	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: timeout,
		TLS:         tlsConfig,
	})

	if err != nil {
		// Deadline Exceeded
		return nil, errors.Errorf("Error starting the etcd client: %v", err)
	}

	var cache *ristretto.Cache
	var watcher clientv3.Watcher
	if cached {
		cache, err = ristretto.NewCache(&ristretto.Config{
			NumCounters: 1e7,     // number of keys to track frequency of (10M).
			MaxCost:     1 << 30, // maximum cost of cache (1GB).
			BufferItems: 64,      // number of keys per Get buffer.
		})

		if err != nil {
			return nil, errors.Errorf("error creating a local ristretto cache: %s", err.Error())
		}
		watcher = clientv3.NewWatcher(cli)
	}

	return &NameService{
		cli:     cli,
		watcher: watcher,
		local:   cache,
		NodeID:  nodeID,
		cached:  cached,
	}, nil
}

// RegisterSelf stores information about this node
func (n *NameService) RegisterSelf(host string, externalHost string) error {
	key := fmt.Sprintf(fmtNodeAdressString, n.NodeID)
	log.Debug().Msgf("NaSe: registering self as %s // %s", key, host)

	err := n.put(key, host)

	if err != nil {
		return err
	}

	key = fmt.Sprintf(fmtNodeExternalAdressString, n.NodeID)
	log.Debug().Msgf("NaSe: registering external address as %s // %s", key, externalHost)
	return n.put(key, externalHost)
}

// GetNodeID returns the ID of this node.
func (n *NameService) GetNodeID() fred.NodeID {
	return fred.NodeID(n.NodeID)
}
