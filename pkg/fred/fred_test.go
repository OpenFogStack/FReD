package fred_test

import (
	"net/url"
	"os"
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
	certBasePath = "../../nase/tls/"
	etcdDir      = ".default.etcd"
)

var f fred.Fred

func TestMain(m *testing.M) {
	nodeID := "X"

	log.Logger = log.Output(
		zerolog.ConsoleWriter{
			Out:     os.Stderr,
			NoColor: false,
		},
	)

	zerolog.SetGlobalLevel(zerolog.ErrorLevel)

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

	n, err := etcdnase.NewNameService(nodeID, []string{"127.0.0.1:6000"}, certBasePath+"client.crt", certBasePath+"client.key", certBasePath+"ca.crt")

	if err != nil {
		panic(err)
	}

	config := fred.Config{
		Store:             badgerdb.NewMemory(),
		Client:            peering.NewClient(),
		NaSe:              n,
		PeeringHost:       "127.0.0.1:8000",
		PeeringHostProxy:  "",
		ExternalHost:      "127.0.0.1:9000",
		ExternalHostProxy: "",
		NodeID:            nodeID,
		TriggerCert:       certBasePath + "client.crt",
		TriggerKey:        certBasePath + "client.key",
	}

	f = fred.New(&config)

	os.Exit(m.Run())
}

func testPut(t *testing.T, user, kg, id, value string) {
	err := f.E.HandleCreateKeygroup(user, fred.Keygroup{
		Name:    fred.KeygroupName(kg),
		Mutable: true,
		Expiry:  0,
	})

	if err != nil {
		log.Err(err).Msg(err.(*errors.Error).ErrorStack())
		t.Error(err)
	}

	err = f.E.HandleUpdate(user, fred.Item{
		Keygroup: fred.KeygroupName(kg),
		ID:       id,
		Val:      value,
	})

	if err != nil {
		log.Err(err).Msg(err.(*errors.Error).ErrorStack())
		t.Error(err)
	}

	i, err := f.E.HandleRead(user, fred.Item{
		Keygroup: fred.KeygroupName(kg),
		ID:       id,
	})

	if err != nil {
		log.Err(err).Msg(err.(*errors.Error).ErrorStack())
		t.Error(err)
	}

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

	if err != nil {
		log.Err(err).Msg(err.(*errors.Error).ErrorStack())
		t.Error(err)
	}

	err = f.E.HandleUpdate(user, fred.Item{
		Keygroup: fred.KeygroupName(kg),
		ID:       id,
		Val:      value,
	})

	if err != nil {
		log.Err(err).Msg(err.(*errors.Error).ErrorStack())
		t.Error(err)
	}

	err = f.E.HandleDelete(user, fred.Item{
		Keygroup: fred.KeygroupName(kg),
		ID:       id,
	})

	if err != nil {
		log.Err(err).Msg(err.(*errors.Error).ErrorStack())
		t.Error(err)
	}

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

	if err != nil {
		log.Err(err).Msg(err.(*errors.Error).ErrorStack())
		t.Error(err)
	}

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

	if err != nil {
		log.Err(err).Msg(err.(*errors.Error).ErrorStack())
		b.Error(err)
	}

	for i := 0; i < b.N; i++ {
		err = f.E.HandleUpdate(user, fred.Item{
			Keygroup: fred.KeygroupName(kg),
			ID:       id,
			Val:      value,
		})

		if err != nil {
			log.Err(err).Msg(err.(*errors.Error).ErrorStack())
			b.Error(err)
		}
	}

	err = f.E.HandleDeleteKeygroup(user, fred.Keygroup{
		Name: fred.KeygroupName(kg),
	})

	if err != nil {
		log.Err(err).Msg(err.(*errors.Error).ErrorStack())
		b.Error(err)
	}
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

	if err != nil {
		log.Err(err).Msg(err.(*errors.Error).ErrorStack())
		b.Error(err)
	}

	err = f.E.HandleUpdate(user, fred.Item{
		Keygroup: fred.KeygroupName(kg),
		ID:       id,
		Val:      value,
	})

	if err != nil {
		log.Err(err).Msg(err.(*errors.Error).ErrorStack())
		b.Error(err)
	}

	for i := 0; i < b.N; i++ {
		_, err := f.E.HandleRead(user, fred.Item{
			Keygroup: fred.KeygroupName(kg),
			ID:       id,
		})

		if err != nil {
			log.Err(err).Msg(err.(*errors.Error).ErrorStack())
			b.Error(err)
		}
	}

	err = f.E.HandleDeleteKeygroup(user, fred.Keygroup{
		Name: fred.KeygroupName(kg),
	})

	if err != nil {
		log.Err(err).Msg(err.(*errors.Error).ErrorStack())
		b.Error(err)
	}
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

	if err != nil {
		log.Err(err).Msg(err.(*errors.Error).ErrorStack())
		b.Error(err)
	}

	for i := 0; i < b.N; i++ {
		_, err = f.E.HandleAppend(user, fred.Item{
			Keygroup: fred.KeygroupName(kg),
			Val:      value,
		})

		if err != nil {
			log.Err(err).Msg(err.(*errors.Error).ErrorStack())
			b.Error(err)
		}
	}

	err = f.E.HandleDeleteKeygroup(user, fred.Keygroup{
		Name: fred.KeygroupName(kg),
	})

	if err != nil {
		log.Err(err).Msg(err.(*errors.Error).ErrorStack())
		b.Error(err)
	}
}
