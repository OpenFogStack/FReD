package fred_test

import (
	"net/url"
	"os"
	"strconv"
	"testing"

	"git.tu-berlin.de/mcc-fred/fred/pkg/badgerdb"
	"git.tu-berlin.de/mcc-fred/fred/pkg/etcdnase"
	"git.tu-berlin.de/mcc-fred/fred/pkg/fred"
	"git.tu-berlin.de/mcc-fred/fred/pkg/peering"
	"github.com/go-errors/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
	"go.etcd.io/etcd/client/pkg/v3/transport"
	"go.etcd.io/etcd/server/v3/embed"
)

const (
	certBasePath = "../../tests/runner/certificates/"
	etcdDir      = ".default.etcd"
	nodeID       = fred.NodeID("X")
)

var f fred.Fred

func TestMain(m *testing.M) {
	log.Logger = log.Output(
		zerolog.ConsoleWriter{
			Out:     os.Stderr,
			NoColor: false,
		},
	)

	zerolog.SetGlobalLevel(zerolog.FatalLevel)

	fInfo, err := os.Stat(etcdDir)

	if err == nil {
		if !fInfo.IsDir() {
			panic(errors.Errorf("%s is not a directory!", etcdDir))
		}

		err = os.RemoveAll(etcdDir)
		if err != nil {
			panic(err)
		}
	}

	cfg := embed.NewConfig()
	cfg.Dir = etcdDir
	u, _ := url.Parse("https://127.0.0.1:6000")
	cfg.LCUrls = []url.URL{*u}
	cfg.ACUrls = []url.URL{*u}
	cfg.ForceNewCluster = true

	cfg.ClientTLSInfo = transport.TLSInfo{
		CertFile:       certBasePath + "etcd.crt",
		KeyFile:        certBasePath + "etcd.key",
		TrustedCAFile:  certBasePath + "ca.crt",
		ClientCertAuth: true,
	}

	cfg.LogLevel = "error"

	e, err := embed.StartEtcd(cfg)

	if err != nil {
		panic(err)
	}

	<-e.Server.ReadyNotify()

	n, err := etcdnase.NewNameService(string(nodeID), []string{"127.0.0.1:6000"}, certBasePath+"nodeA.crt", certBasePath+"nodeA.key", certBasePath+"ca.crt", true)

	if err != nil {
		panic(err)
	}

	config := fred.Config{
		Store:             badgerdb.NewMemory(),
		Client:            peering.NewClient(certBasePath+"nodeA.crt", certBasePath+"nodeA.key", certBasePath+"ca.crt"),
		NaSe:              n,
		PeeringHost:       "127.0.0.1:8000",
		PeeringHostProxy:  "",
		ExternalHost:      "127.0.0.1:9000",
		ExternalHostProxy: "",
		NodeID:            string(nodeID),
		TriggerCert:       certBasePath + "nodeA.crt",
		TriggerKey:        certBasePath + "nodeA.key",
		TriggerCA:         []string{certBasePath + "ca.crt"},
	}

	f = fred.New(&config)

	stat := m.Run()

	fInfo, err = os.Stat(etcdDir)

	if err == nil {
		if !fInfo.IsDir() {
			panic(errors.Errorf("%s is not a directory!", etcdDir))
		}

		err = os.RemoveAll(etcdDir)
		if err != nil {
			panic(err)
		}
	}

	os.Exit(stat)
}

func testPut(t *testing.T, user, kg, id, value string) {
	err := f.E.HandleCreateKeygroup(user, fred.Keygroup{
		Name:    fred.KeygroupName(kg),
		Mutable: true,
		Expiry:  0,
	})

	assert.NoError(t, err)

	err = f.E.HandleUpdate(user, fred.Item{
		Keygroup: fred.KeygroupName(kg),
		ID:       id,
		Val:      value,
	})

	assert.NoError(t, err)

	i, err := f.E.HandleRead(user, fred.Item{
		Keygroup: fred.KeygroupName(kg),
		ID:       id,
	})

	assert.NoError(t, err)

	assert.Equal(t, kg, string(i.Keygroup))
	assert.Equal(t, id, i.ID)
	assert.Equal(t, value, i.Val)
}

func TestPut(t *testing.T) {
	testPut(t, "user", "testkg", "id", "value")
	testPut(t, "user", "testkg2", "id", "value")
	testPut(t, "user", "testkg3", "id", "value")
	testPut(t, "user", "testkg4", "id", "value")
	testPut(t, "user", "testkg5", "id", "value")
}

func testDelete(t *testing.T, user, kg, id, value string) {
	err := f.E.HandleCreateKeygroup(user, fred.Keygroup{
		Name:    fred.KeygroupName(kg),
		Mutable: true,
		Expiry:  0,
	})

	assert.NoError(t, err)

	err = f.E.HandleUpdate(user, fred.Item{
		Keygroup: fred.KeygroupName(kg),
		ID:       id,
		Val:      value,
	})

	assert.NoError(t, err)

	err = f.E.HandleDelete(user, fred.Item{
		Keygroup: fred.KeygroupName(kg),
		ID:       id,
	})

	assert.NoError(t, err)

	_, err = f.E.HandleRead(user, fred.Item{
		Keygroup: fred.KeygroupName(kg),
		ID:       id,
	})

	assert.Error(t, err)
}

func TestDelete(t *testing.T) {
	testDelete(t, "user", "testdelete", "id", "value")
	testDelete(t, "user", "testdelete2", "id", "value")
	testDelete(t, "user", "testdelete3", "id", "value")
	testDelete(t, "user", "testdelete4", "id", "value")
	testDelete(t, "user", "testdelete5", "id", "value")
}

func testMisformedKVInput(t *testing.T, user, kg, id, value string) {
	err := f.E.HandleCreateKeygroup(user, fred.Keygroup{
		Name:    fred.KeygroupName(kg),
		Mutable: true,
		Expiry:  0,
	})

	assert.NoError(t, err)

	err = f.E.HandleUpdate(user, fred.Item{
		Keygroup: fred.KeygroupName(kg),
		ID:       id,
		Val:      value,
	})

	assert.Error(t, err)
}

func TestMisformedKVInput(t *testing.T) {
	testMisformedKVInput(t, "user", "misformed", "id(", "value")
	testMisformedKVInput(t, "user", "misformed2", "id(", "value=")
}

func testMisformedKeygroupInput(t *testing.T, user, kg, id, value string) {
	err := f.E.HandleCreateKeygroup(user, fred.Keygroup{
		Name:    fred.KeygroupName(kg),
		Mutable: true,
		Expiry:  0,
	})

	assert.Error(t, err)
}

func TestMisformedKeygroupInput(t *testing.T) {
	// TODO: allow this?
	// testMisformedKeygroupInput(t, "user920194\n", "misformed2", "id", "value")
	testMisformedKeygroupInput(t, "user", "misf%?ormed", "id", "value")
	testMisformedKeygroupInput(t, "use|r", "misf|ormed", "id|", "val|ue")
}

func testScan(t *testing.T, user string, kg string, mutable bool, updates int, scanStart int, scanRange int) {

	// 1. create a keygroup
	err := f.E.HandleCreateKeygroup(user, fred.Keygroup{
		Name:    fred.KeygroupName(kg),
		Mutable: mutable,
		Expiry:  0,
	})

	assert.NoError(t, err)

	defer func() {
		err := f.E.HandleDeleteKeygroup(user, fred.Keygroup{
			Name: fred.KeygroupName(kg),
		})

		assert.NoError(t, err)
	}()

	// 2. put in a bunch of items
	ids := make([]string, updates)
	vals := make([]string, updates)

	for i := 0; i < updates; i++ {
		vals[i] = "val" + strconv.Itoa(i)
		if mutable {
			ids[i] = "id" + strconv.Itoa(i)
			err := f.E.HandleUpdate(user, fred.Item{
				Keygroup: fred.KeygroupName(kg),
				ID:       ids[i],
				Val:      vals[i],
			})
			assert.NoError(t, err)
		} else {
			item, err := f.E.HandleAppend(user, fred.Item{
				Keygroup: fred.KeygroupName(kg),
				Val:      vals[i],
			})

			if !assert.NoError(t, err) {
				continue
			}

			ids[i] = item.ID
		}
	}

	// 3. do a scan read
	// we expect [scanRange] amount of items, starting with [scanStart]
	var startKey string
	if mutable {
		startKey = "id" + strconv.Itoa(scanStart)
	} else {
		startKey = strconv.Itoa(scanStart)
	}

	items, err := f.E.HandleScan(user, fred.Item{
		Keygroup: fred.KeygroupName(kg),
		ID:       startKey,
	}, uint64(scanRange))

	// if scanRange == 0 we expect an empty array
	if scanRange == 0 {
		assert.Len(t, items, 0, "scanRange %d == 0", scanRange)
		return
	}

	// if startkey (scanStart not in [0, updates[) is not found, we also expect an error
	if scanStart < 0 || scanStart >= updates {
		assert.Error(t, err, "scanStart %d not in [0, updates[", scanStart)
		return
	}

	// if scanRange < updates, we expect exactly min(updates - scanStart, scanRange) items
	// in this case we obviously want all the values to be correct as well!

	assert.NoError(t, err)

	expected := updates - scanStart
	if scanRange < expected {
		expected = scanRange
	}

	assert.Len(t, items, expected, "expect exactly updates - scanStart items")

	// we make no assumptions about the ordering of key-value pairs on a scan read

	res := make(map[string]string)

	for _, i := range items {
		res[i.ID] = i.Val
	}

	for i := 0; i < updates; i++ {
		if i < scanStart || i >= scanStart+expected {
			continue
		}
		assert.Contains(t, res, ids[i])
		assert.Equal(t, vals[i], res[ids[i]])
	}
}

func TestScan(t *testing.T) {
	testScan(t, "user", "scankeygroup1", true, 10, 3, 5)
	testScan(t, "user", "scankeygroup2", false, 10, 3, 5)
	testScan(t, "user", "scankeygroup3", true, 100, 0, 100)
	testScan(t, "user", "scankeygroup4", false, 100, 0, 100)
	testScan(t, "user", "scankeygroup5", true, 10, 11, 5)
	testScan(t, "user", "scankeygroup6", true, 10, 0, 11)
	testScan(t, "user", "scankeygroup7", false, 10, 11, 5)
	testScan(t, "user", "scankeygroup8", false, 10, 0, 11)
}

func TestPermissions(t *testing.T) {
	user1 := "user1"
	user2 := "user2"
	kg := "permissiontest"

	err := f.E.HandleCreateKeygroup(user1, fred.Keygroup{
		Name:    fred.KeygroupName(kg),
		Mutable: true,
		Expiry:  0,
	})

	assert.NoError(t, err)

	err = f.E.HandleUpdate(user1, fred.Item{
		Keygroup: fred.KeygroupName(kg),
		ID:       "id",
		Val:      "value",
	})

	assert.NoError(t, err)

	err = f.E.HandleUpdate(user2, fred.Item{
		Keygroup: fred.KeygroupName(kg),
		ID:       "id2",
		Val:      "value2",
	})

	assert.Error(t, err)

	_, err = f.E.HandleRead(user2, fred.Item{
		Keygroup: fred.KeygroupName(kg),
		ID:       "id",
	})

	assert.Error(t, err)

	err = f.E.HandleAddUser(user1, user2, fred.Keygroup{Name: fred.KeygroupName(kg)}, fred.ConfigureKeygroups)

	assert.NoError(t, err)

	_, err = f.E.HandleRead(user2, fred.Item{
		Keygroup: fred.KeygroupName(kg),
		ID:       "id",
	})

	assert.Error(t, err)

	err = f.E.HandleAddUser(user1, user2, fred.Keygroup{Name: fred.KeygroupName(kg)}, fred.ReadKeygroup)

	assert.NoError(t, err)

	i, err := f.E.HandleRead(user2, fred.Item{
		Keygroup: fred.KeygroupName(kg),
		ID:       "id",
	})

	assert.NoError(t, err)

	assert.Equal(t, kg, string(i.Keygroup))
	assert.Equal(t, "id", i.ID)
	assert.Equal(t, "value", i.Val)

	err = f.E.HandleUpdate(user2, fred.Item{
		Keygroup: fred.KeygroupName(kg),
		ID:       "id2",
		Val:      "value2",
	})

	assert.Error(t, err)

	err = f.E.HandleRemoveUser(user2, user1, fred.Keygroup{Name: fred.KeygroupName(kg)}, fred.ReadKeygroup)

	assert.NoError(t, err)

	i, err = f.E.HandleRead(user1, fred.Item{
		Keygroup: fred.KeygroupName(kg),
		ID:       "id",
	})

	assert.Error(t, err)
}

func TestAllReplicas(t *testing.T) {
	nodes, err := f.E.HandleGetAllReplica("user")

	assert.NoError(t, err)

	assert.Len(t, nodes, 1)
	assert.Equal(t, nodeID, nodes[0].ID)
	assert.Equal(t, "127.0.0.1:9000", nodes[0].Host)
}

func TestSingleReplica(t *testing.T) {
	node, err := f.E.HandleGetReplica("user", fred.Node{ID: nodeID})

	assert.NoError(t, err)
	assert.Equal(t, nodeID, node.ID)
	assert.Equal(t, "127.0.0.1:9000", node.Host)

	_, err = f.E.HandleGetReplica("user", fred.Node{ID: "Y"})

	assert.Error(t, err)
}

func TestKeygroupReplicas(t *testing.T) {
	user := "user1"
	var kg fred.KeygroupName = "replicakeygroup"

	err := f.E.HandleCreateKeygroup(user, fred.Keygroup{
		Name:    kg,
		Mutable: true,
		Expiry:  0,
	})

	assert.NoError(t, err)

	nodes, expiries, err := f.E.HandleGetKeygroupReplica(user, fred.Keygroup{Name: kg})

	assert.NoError(t, err)
	assert.Len(t, nodes, 1)
	assert.Equal(t, nodeID, nodes[0].ID)
	assert.Equal(t, "127.0.0.1:9000", nodes[0].Host)
	assert.Len(t, expiries, 1)
	assert.Equal(t, expiries[nodes[0].ID], 0)
}

func TestTriggers(t *testing.T) {
	user := "user1"
	var kg fred.KeygroupName = "triggerkeygroup"
	tid := "trigger1"
	thost := "5.5.5.5:9000"

	err := f.E.HandleAddTrigger(user, fred.Keygroup{
		Name: kg,
	}, fred.Trigger{
		ID:   tid,
		Host: thost,
	})

	assert.Error(t, err)

	err = f.E.HandleCreateKeygroup(user, fred.Keygroup{
		Name:    kg,
		Mutable: true,
		Expiry:  0,
	})

	assert.NoError(t, err)

	err = f.E.HandleAddTrigger(user, fred.Keygroup{
		Name: kg,
	}, fred.Trigger{
		ID:   tid,
		Host: thost,
	})

	assert.NoError(t, err)

	triggers, err := f.E.HandleGetKeygroupTriggers(user, fred.Keygroup{
		Name: kg,
	})

	assert.NoError(t, err)
	assert.Len(t, triggers, 1)
	assert.Equal(t, triggers[0].ID, tid)
	assert.Equal(t, triggers[0].Host, thost)

	err = f.E.HandleRemoveTrigger(user, fred.Keygroup{
		Name: kg,
	}, fred.Trigger{
		ID: tid,
	})

	assert.NoError(t, err)

	triggers, err = f.E.HandleGetKeygroupTriggers(user, fred.Keygroup{
		Name: kg,
	})

	assert.NoError(t, err)
	assert.Len(t, triggers, 0)

}

func TestInternalPut(t *testing.T) {
	kg := fred.KeygroupName("kginternalput")
	id := "item1"
	val := "value1"

	err := f.E.HandleCreateKeygroup("user", fred.Keygroup{
		Name:    kg,
		Mutable: true,
		Expiry:  0,
	})

	assert.NoError(t, err)

	err = f.I.HandleUpdate(fred.Item{
		Keygroup: kg,
		ID:       id,
		Val:      val,
	})

	assert.NoError(t, err)

	i, err := f.I.HandleGet(fred.Item{
		Keygroup: kg,
		ID:       id,
	})

	assert.NoError(t, err)

	assert.Equal(t, kg, i.Keygroup)
	assert.Equal(t, id, i.ID)
	assert.Equal(t, val, i.Val)

	items, err := f.I.HandleGetAllItems(fred.Keygroup{
		Name: kg,
	})

	assert.NoError(t, err)
	assert.Len(t, items, 1)
	assert.Equal(t, kg, items[0].Keygroup)
	assert.Equal(t, id, items[0].ID)
	assert.Equal(t, val, items[0].Val)
}

func TestInternalDelete(t *testing.T) {
	kg := fred.KeygroupName("kginternaldelete")
	id := "item1"
	val := "value1"

	err := f.E.HandleCreateKeygroup("user", fred.Keygroup{
		Name:    kg,
		Mutable: true,
		Expiry:  0,
	})

	assert.NoError(t, err)

	err = f.E.HandleUpdate("user", fred.Item{
		Keygroup: kg,
		ID:       id,
		Val:      val,
	})

	assert.NoError(t, err)

	items, err := f.I.HandleGetAllItems(fred.Keygroup{
		Name: kg,
	})

	assert.NoError(t, err)
	assert.Len(t, items, 1)

	err = f.I.HandleDelete(fred.Item{
		Keygroup: kg,
		ID:       id,
	})

	assert.NoError(t, err)

	items, err = f.I.HandleGetAllItems(fred.Keygroup{
		Name: kg,
	})

	assert.NoError(t, err)
	assert.Len(t, items, 0)
}

func TestInternalAppend(t *testing.T) {
	kg := fred.KeygroupName("kginternappend")
	id := "0"
	val := "value1"

	err := f.E.HandleCreateKeygroup("user", fred.Keygroup{
		Name:    kg,
		Mutable: false,
		Expiry:  0,
	})

	assert.NoError(t, err)

	err = f.I.HandleAppend(fred.Item{
		Keygroup: kg,
		ID:       id,
		Val:      val,
	})

	assert.NoError(t, err)

	i, err := f.I.HandleGet(fred.Item{
		Keygroup: kg,
		ID:       id,
	})

	assert.NoError(t, err)

	assert.Equal(t, kg, i.Keygroup)
	assert.Equal(t, id, i.ID)
	assert.Equal(t, val, i.Val)

	items, err := f.I.HandleGetAllItems(fred.Keygroup{
		Name: kg,
	})

	assert.NoError(t, err)
	assert.Len(t, items, 1)
	assert.Equal(t, kg, items[0].Keygroup)
	assert.Equal(t, id, items[0].ID)
	assert.Equal(t, val, items[0].Val)

	i, err = f.E.HandleAppend("user", fred.Item{
		Keygroup: kg,
		Val:      "value2",
	})

	assert.NoError(t, err)
	assert.Equal(t, kg, i.Keygroup)
	assert.Equal(t, "1", i.ID)
	assert.Equal(t, "value2", i.Val)
}

func BenchmarkPut(b *testing.B) {
	user := "user"
	kg := "benchmarkPut"
	id := "benchmarkItem"
	value := "benchmarkVal"

	err := f.E.HandleCreateKeygroup(user, fred.Keygroup{
		Name:    fred.KeygroupName(kg),
		Mutable: true,
		Expiry:  0,
	})

	assert.NoError(b, err)

	for i := 0; i < b.N; i++ {
		err = f.E.HandleUpdate(user, fred.Item{
			Keygroup: fred.KeygroupName(kg),
			ID:       id,
			Val:      value,
		})

		assert.NoError(b, err)
	}

	err = f.E.HandleDeleteKeygroup(user, fred.Keygroup{
		Name: fred.KeygroupName(kg),
	})

	assert.NoError(b, err)
}

func BenchmarkGet(b *testing.B) {
	user := "user"
	kg := "benchmarkGet"
	id := "benchmarkItem"
	value := "benchmarkVal"

	err := f.E.HandleCreateKeygroup(user, fred.Keygroup{
		Name:    fred.KeygroupName(kg),
		Mutable: true,
		Expiry:  0,
	})

	assert.NoError(b, err)

	err = f.E.HandleUpdate(user, fred.Item{
		Keygroup: fred.KeygroupName(kg),
		ID:       id,
		Val:      value,
	})

	assert.NoError(b, err)

	for i := 0; i < b.N; i++ {
		_, err := f.E.HandleRead(user, fred.Item{
			Keygroup: fred.KeygroupName(kg),
			ID:       id,
		})

		assert.NoError(b, err)
	}

	err = f.E.HandleDeleteKeygroup(user, fred.Keygroup{
		Name: fred.KeygroupName(kg),
	})

	assert.NoError(b, err)
}

func BenchmarkAppend(b *testing.B) {
	user := "user"
	kg := "benchmarkAppend"
	value := "benchmarkVal"

	err := f.E.HandleCreateKeygroup(user, fred.Keygroup{
		Name:    fred.KeygroupName(kg),
		Mutable: false,
		Expiry:  0,
	})

	assert.NoError(b, err)

	for i := 0; i < b.N; i++ {
		_, err = f.E.HandleAppend(user, fred.Item{
			Keygroup: fred.KeygroupName(kg),
			Val:      value,
		})

		assert.NoError(b, err)
	}

	err = f.E.HandleDeleteKeygroup(user, fred.Keygroup{
		Name: fred.KeygroupName(kg),
	})

	assert.NoError(b, err)
}
