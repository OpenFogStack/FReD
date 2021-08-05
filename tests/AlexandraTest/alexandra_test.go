package alexandratest

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	composePath  = "../runner/"
	certBasePath = "../runner/certificates/"
)

var c *alexandraClient

func TestMain(m *testing.M) {
	log.Logger = log.Output(
		zerolog.ConsoleWriter{
			Out:     os.Stderr,
			NoColor: false,
		},
	)

	dc := testcontainers.NewLocalDockerCompose([]string{
		composePath + "fredwork.yml",
		composePath + "etcd.yml",
		composePath + "nodeA.yml",
		composePath + "nodeB.yml",
		composePath + "nodeC.yml",
		composePath + "trigger.yml",
		composePath + "alexandra.yml",
	}, "3nodetest")

	execError := dc.Down()

	if execError.Error != nil {
		panic(fmt.Sprintf("Failed when running: %v", execError.Command))
	}

	execError = dc.WithCommand([]string{"build"}).Invoke()

	if execError.Error != nil {
		panic(fmt.Sprintf("Failed when running: %v", execError))
	}

	dc.WithExposedService("alexandra", 10000, wait.NewHostPortStrategy("10000"))

	execError = dc.WithCommand([]string{"up", "-d", "--force-recreate", "--renew-anon-volumes", "--remove-orphans"}).Invoke()

	if execError.Error != nil {
		panic(fmt.Sprintf("Failed when running: %v", execError))
	}

	c = newAlexandraClient("127.0.0.1:10000", certBasePath+"alexandraTester.crt", certBasePath+"alexandraTester.key", certBasePath+"ca.crt")

	time.Sleep(20 * time.Second)

	stat := m.Run()

	execError = dc.Down()

	if execError.Error != nil {
		panic(fmt.Sprintf("Failed when running: %v", execError.Command))
	}

	os.Exit(stat)
}

func TestBasic(t *testing.T) {
	kg := "alexandraTest"
	id := "id"
	val := "data"

	// Create a keygroup
	err := c.createKeygroup("nodeB", kg, true, 10_000)
	assert.NoError(t, err)

	// Put a value into it.
	err = c.update(kg, id, val)
	assert.NoError(t, err)

	// read this value from anywhere
	read, err := c.read(kg, id, 500)
	assert.NoError(t, err)
	assert.Equal(t, val, read)

	// Add the other nodes to the keygroup
	err = c.addKeygroupReplica(kg, "nodeA", 600)
	assert.NoError(t, err)

	err = c.addKeygroupReplica(kg, "nodeC", 600)
	assert.NoError(t, err)

	// read from somewhere else
	read, err = c.read(kg, id, 500)
	assert.NoError(t, err)
	assert.Equal(t, val, read)
}
