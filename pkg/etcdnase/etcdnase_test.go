//go:build !race
// +build !race

package etcdnase

import (
	"fmt"
	"math/rand"
	"net/url"
	"os"
	"strconv"
	"testing"
	"time"

	"git.tu-berlin.de/mcc-fred/fred/pkg/fred"
	"github.com/go-errors/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
	"go.etcd.io/etcd/client/pkg/v3/transport"
	"go.etcd.io/etcd/server/v3/embed"
)

const (
	certBasePath = "../../tests/certificates/"
	etcdDir      = ".default.etcd"
	nodeID       = "A"
	host         = "localhost:8000"
	extHost      = "localhost:9000"
)

var n *NameService

func TestMain(m *testing.M) {
	log.Logger = log.Output(
		zerolog.ConsoleWriter{
			Out:     os.Stderr,
			NoColor: false,
		},
	)

	zerolog.SetGlobalLevel(zerolog.FatalLevel)
	// zerolog.SetGlobalLevel(zerolog.DebugLevel)

	fInfo, err := os.Stat(etcdDir)

	if err == nil {
		if !fInfo.IsDir() {
			panic(errors.Errorf("%s is not a directory", etcdDir))
		}

		err = os.RemoveAll(etcdDir)
		if err != nil {
			panic(err)
		}
	}

	cfg := embed.NewConfig()
	cfg.Dir = etcdDir
	cURL, _ := url.Parse("https://127.0.0.1:6000")
	pURL, _ := url.Parse("http://127.0.0.1:6001")
	cfg.LCUrls = []url.URL{*cURL}
	cfg.ACUrls = []url.URL{*cURL}
	cfg.LPUrls = []url.URL{*pURL}
	cfg.APUrls = []url.URL{*pURL}
	cfg.ForceNewCluster = true

	cfg.ClientTLSInfo = transport.TLSInfo{
		CertFile:       certBasePath + "etcd.crt",
		KeyFile:        certBasePath + "etcd.key",
		TrustedCAFile:  certBasePath + "ca.crt",
		ClientCertAuth: true,
	}

	cfg.LogLevel = "error"

	cfg.InitialCluster = cfg.InitialClusterFromName(cfg.Name)

	e, err := embed.StartEtcd(cfg)

	if err != nil {
		panic(err)
	}

	<-e.Server.ReadyNotify()

	n, err = NewNameService(nodeID, []string{"127.0.0.1:6000"}, certBasePath+"nodeA.crt", certBasePath+"nodeA.key", certBasePath+"ca.crt", false, true)

	if err != nil {
		panic(err)
	}

	err = n.RegisterSelf(host, extHost)

	if err != nil {
		panic(err)
	}

	stat := m.Run()

	e.Close()

	fInfo, err = os.Stat(etcdDir)

	if err == nil {
		if !fInfo.IsDir() {
			panic(errors.Errorf("%s is not a directory", etcdDir))
		}

		err = os.RemoveAll(etcdDir)
		if err != nil {
			panic(err)
		}
	}

	os.Exit(stat)
}

func TestNodeID(t *testing.T) {
	id := n.GetNodeID()
	assert.Equal(t, nodeID, string(id))
}

func TestFailure(t *testing.T) {
	items := n.RequestNodeStatus(nodeID)

	assert.Len(t, items, 0)
}

func TestNodes(t *testing.T) {
	addr, err := n.GetNodeAddress(nodeID)

	assert.NoError(t, err)
	assert.Equal(t, host, addr)

	addr, err = n.GetNodeAddressExternal(nodeID)

	assert.NoError(t, err)
	assert.Equal(t, extHost, addr)

	nodes, err := n.GetAllNodes()

	assert.NoError(t, err)
	assert.Len(t, nodes, 1)
	assert.Equal(t, nodes[0].ID, fred.NodeID(nodeID))
	assert.Equal(t, nodes[0].Host, host)

	nodes, err = n.GetAllNodesExternal()

	assert.NoError(t, err)
	assert.Len(t, nodes, 1)
	assert.Equal(t, nodes[0].ID, fred.NodeID(nodeID))
	assert.Equal(t, nodes[0].Host, extHost)
}

func TestNonexistentKeygroups(t *testing.T) {
	kg := fred.KeygroupName("kg-noexist")

	exists, err := n.ExistsKeygroup(kg)

	assert.NoError(t, err)
	assert.False(t, exists)

	err = n.DeleteKeygroup(kg)

	assert.Error(t, err)

	_, err = n.IsMutable(kg)

	assert.Error(t, err)

	_, err = n.GetExpiry(kg)

	assert.Error(t, err)

	_, err = n.GetKeygroupMembers(kg, false)

	assert.Error(t, err)
}

// now create the keygroup
func TestKeygroups(t *testing.T) {
	kg := fred.KeygroupName("kg")
	mutable := true
	expiry := 10
	err := n.CreateKeygroup(kg, mutable, expiry)

	assert.NoError(t, err)

	m, err := n.IsMutable(kg)

	assert.NoError(t, err)
	assert.Equal(t, mutable, m)

	exp, err := n.GetExpiry(kg)

	assert.NoError(t, err)
	assert.Equal(t, expiry, exp)
}

func TestKeygroupMembers(t *testing.T) {
	kg := fred.KeygroupName("kg-members")
	mutable := true
	expiry := 124
	badNode := fred.NodeID("B")

	err := n.CreateKeygroup(kg, mutable, expiry)

	assert.NoError(t, err)

	members, err := n.GetKeygroupMembers(kg, false)

	assert.NoError(t, err)
	assert.Len(t, members, 1)
	assert.Contains(t, members, fred.NodeID(nodeID))
	assert.Equal(t, expiry, members[nodeID])

	err = n.JoinNodeIntoKeygroup(kg, badNode, 0)

	assert.Error(t, err)

	err = n.ExitOtherNodeFromKeygroup(kg, badNode)

	assert.Error(t, err)
}

func TestPermissions(t *testing.T) {
	user := "user"
	kg := fred.KeygroupName("kg-permissions")

	_, err := n.GetUserPermissions(user, kg)

	assert.Error(t, err)

	err = n.RevokeUserPermissions(user, fred.Update, kg)

	assert.Error(t, err)

	err = n.CreateKeygroup(kg, true, 0)

	assert.NoError(t, err)

	perm, err := n.GetUserPermissions(user, kg)

	assert.NoError(t, err)
	assert.Len(t, perm, 0)

	// expect no error even if user doesn't have permissions
	err = n.RevokeUserPermissions(user, fred.Update, kg)

	assert.NoError(t, err)

	err = n.AddUserPermissions(user, fred.Read, kg)

	assert.NoError(t, err)

	perm, err = n.GetUserPermissions(user, kg)

	assert.NoError(t, err)
	assert.Len(t, perm, 1)
	assert.Contains(t, perm, fred.Read)

	err = n.AddUserPermissions(user, fred.Update, kg)

	assert.NoError(t, err)

	err = n.RevokeUserPermissions(user, fred.Read, kg)

	assert.NoError(t, err)

	perm, err = n.GetUserPermissions(user, kg)

	assert.NoError(t, err)
	assert.Len(t, perm, 1)
	assert.Contains(t, perm, fred.Update)
}

func TestCache(t *testing.T) {
	altNodeID := fred.NodeID("B")
	altNodeHost := "localhost:8001"
	altNodeExtHost := "localhost:8002"

	kg := fred.KeygroupName("kg-caching")
	user := "user2"

	n2, err := NewNameService(string(altNodeID), []string{"127.0.0.1:6000"}, certBasePath+"nodeB.crt", certBasePath+"nodeB.key", certBasePath+"ca.crt", false, true)

	assert.NoError(t, err)

	err = n2.RegisterSelf(altNodeHost, altNodeExtHost)

	assert.NoError(t, err)

	err = n.CreateKeygroup(kg, true, 0)

	assert.NoError(t, err)

	err = n.JoinNodeIntoKeygroup(kg, altNodeID, 0)

	assert.NoError(t, err)

	err = n.AddUserPermissions(user, fred.Read, kg)

	assert.NoError(t, err)

	err = n2.RevokeUserPermissions(user, fred.Read, kg)

	assert.NoError(t, err)

	perm, err := n.GetUserPermissions(user, kg)

	assert.NoError(t, err)
	assert.Len(t, perm, 0)

	err = n.AddUserPermissions(user, fred.Update, kg)

	assert.NoError(t, err)

	perm, err = n.GetUserPermissions(user, kg)

	assert.NoError(t, err)
	assert.Len(t, perm, 1)
	assert.Contains(t, perm, fred.Update)

	err = n2.RevokeUserPermissions(user, fred.Update, kg)

	assert.NoError(t, err)

	err = n2.AddUserPermissions(user, fred.RemoveUser, kg)

	assert.NoError(t, err)

	time.Sleep(1 * time.Millisecond)

	perm, err = n.GetUserPermissions(user, kg)

	assert.NoError(t, err)
	assert.Len(t, perm, 1)
	assert.Contains(t, perm, fred.RemoveUser)
}

// GORACE="halt_on_error=1" go test -run=CacheRace -race
// currently this triggers a race condition, hence this suite is excluded from race testing
func TestCacheRace(t *testing.T) {
	concurrent := 10
	ops := 100
	kg := fred.KeygroupName("kg-cache-race")
	user := "user"

	err := n.CreateKeygroup(kg, true, 0)

	assert.NoError(t, err)

	nases := make([]*NameService, concurrent)

	for i := 0; i < concurrent; i++ {
		id := "node" + strconv.Itoa(i)
		x, err := NewNameService(id, []string{"127.0.0.1:6000"}, certBasePath+"nodeA.crt", certBasePath+"nodeA.key", certBasePath+"ca.crt", false, true)
		assert.NoError(t, err)

		err = x.RegisterSelf(fmt.Sprintf("localhost:10%d01", i), fmt.Sprintf("localhost:10%d02", i))

		assert.NoError(t, err)

		err = n.JoinNodeIntoKeygroup(kg, fred.NodeID(id), 0)

		assert.NoError(t, err)

		nases[i] = x
	}

	done := make(chan struct{})

	for i := 0; i < concurrent; i++ {
		for j := 0; j < concurrent; j++ {
			go func(nase *NameService) {
				for k := 0; k < ops; k++ {
					op := rand.Intn(concurrent)
					p := []fred.Method{
						fred.CreateKeygroup,
						fred.DeleteKeygroup,
						fred.Read,
						fred.Update,
						fred.Delete,
						fred.AddReplica,
						fred.GetReplica,
						fred.RemoveReplica,
						fred.GetAllReplica,
						fred.GetTrigger,
						fred.AddTrigger,
						fred.RemoveTrigger,
						fred.AddUser,
						fred.RemoveUser,
					}

					switch op {
					case 0:
						m := p[rand.Intn(len(p))]
						log.Debug().Msgf("%s adding %s", nase.NodeID, m)
						err = nase.AddUserPermissions(user, m, kg)
						assert.NoError(t, err)

					case 1:
						m := p[rand.Intn(len(p))]
						log.Debug().Msgf("%s removing %s", nase.NodeID, m)
						err = nase.RevokeUserPermissions(user, m, kg)
						assert.NoError(t, err)

					default:
						log.Debug().Msgf("%s getting", nase.NodeID)
						_, err = nase.GetUserPermissions(user, kg)
						assert.NoError(t, err)
					}
				}
				done <- struct{}{}
			}(nases[i])
		}
	}

	for i := 0; i < concurrent*concurrent; i++ {
		<-done
	}
}

func BenchmarkGet(b *testing.B) {
	key := "key-get"
	val := "val"
	err := n.put(key, val)
	assert.NoError(b, err)
	for i := 0; i < b.N; i++ {
		_, _ = n.getExact(key)
	}
}

func BenchmarkGetPrefix(b *testing.B) {
	key := "key-get-nocache|prefix"
	prefix := "key-get-nocache|"
	val := "val"
	err := n.put(key, val, prefix)
	assert.NoError(b, err)
	for i := 0; i < b.N; i++ {
		_, _ = n.getPrefix(prefix)
	}
}

func BenchmarkPut(b *testing.B) {
	key := "key-put"
	val := "val"
	for i := 0; i < b.N; i++ {
		_ = n.put(key, val)
	}
}

func BenchmarkGetNoCache(b *testing.B) {
	n, err := NewNameService(nodeID, []string{"127.0.0.1:6000"}, certBasePath+"nodeA.crt", certBasePath+"nodeA.key", certBasePath+"ca.crt", false, false)
	assert.NoError(b, err)
	b.ResetTimer()

	key := "key-get-nocache"
	val := "val"
	err = n.put(key, val)
	assert.NoError(b, err)
	for i := 0; i < b.N; i++ {
		_, _ = n.getExact(key)
	}
}

func BenchmarkGetPrefixNoCache(b *testing.B) {
	n, err := NewNameService(nodeID, []string{"127.0.0.1:6000"}, certBasePath+"nodeA.crt", certBasePath+"nodeA.key", certBasePath+"ca.crt", false, false)
	assert.NoError(b, err)
	b.ResetTimer()

	key := "key-get-nocache|prefix"
	prefix := "key-get-nocache|"
	val := "val"
	err = n.put(key, val, prefix)
	assert.NoError(b, err)
	for i := 0; i < b.N; i++ {
		_, _ = n.getPrefix(prefix)
	}
}

func BenchmarkPutNoCache(b *testing.B) {
	n, err := NewNameService(nodeID, []string{"127.0.0.1:6000"}, certBasePath+"nodeA.crt", certBasePath+"nodeA.key", certBasePath+"ca.crt", false, false)
	assert.NoError(b, err)
	b.ResetTimer()

	key := "key-put-nocache"
	val := "val"
	for i := 0; i < b.N; i++ {
		_ = n.put(key, val)
	}
}
