package main

import (
	"flag"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"gitlab.tu-berlin.de/mcc-fred/fred/tests/3NodeTest/pkg"
)

func main() {
	// Logging Setup
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	log.Logger = log.Output(
		zerolog.ConsoleWriter{
			Out:     os.Stderr,
			NoColor: false,
		},
	)

	// Parse Flags
	nodeAurl := flag.String("nodeA", "http://localhost:9001/v0/", "ip:port/apiVersion/ of nodeA")
	nodeBurl := flag.String("nodeB", "http://localhost:9002/v0/", "ip:port/apiVersion/ of nodeB")
	nodeCurl := flag.String("nodeC", "http://localhost:9003/v0/", "ip:port/apiVersion/ of nodeC")
	flag.Parse()
	log.Debug().Str("nodeAurl", *nodeAurl).Str("nodeBurl", *nodeBurl).Str("nodeCurl", *nodeCurl).Msg("Connecting to nodes")

	nodeA := pkg.NewNode(*nodeAurl)
	log.Debug().Msg("Creating a Keygroup...")
	nodeA.CreateKeygroup("testing", 200)
	log.Debug().Msg("Deleting the Keygroup...")
	nodeA.DeleteKeygroup("testing", 200)
	log.Debug().Msg("Creating another Keygroup...")
	nodeA.CreateKeygroup("lol", 200)
	log.Debug().Msg("Creating a Keygroup with a difficult name...")
	nodeA.CreateKeygroup("🐧", 200)
}
