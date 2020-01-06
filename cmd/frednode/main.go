package main

// leave this in for cgo to work

import "C"

import (
	"flag"
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/gin-gonic/gin"
	"github.com/mmcloughlin/geohash"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/data"
	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/exthandler"
	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/inthandler"
	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/keygroup"
	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/leveldbsd"
	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/memorykg"
	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/memoryns"
	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/memorysd"
	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/memoryzmq"
	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/replication"
	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/replicationhandler"
	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/webserver"
	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/zmqclient"
	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/zmqserver"
)

type fredConfig struct {
	Location struct {
		Lat float64 `toml:"lat"`
		Lng float64 `toml:"lng"`
	} `toml:"location"`
	Server struct {
		Host string `toml:"host"`
		Port int    `toml:"port"`
	} `toml:"webserver"`
	Storage struct {
		Adaptor string `toml:"adaptor"`
	} `toml:"storage"`
	ZMQ struct {
		Port int `toml:"port"`
	} `toml:"zmq"`
	Log struct {
		Level   string `toml:"level"`
		Handler string `toml:"handler"`
	} `toml:"log"`
}

func main() {
	configPath := flag.String("config", "", "path to .toml configuration file")

	flag.Parse()

	if *configPath == "" {
		log.Fatal().Msg("no configuration specified")
	}

	var fc fredConfig
	if _, err := toml.DecodeFile(*configPath, &fc); err != nil {
		log.Fatal().Err(err).Msg("invalid configuration!")
	}

	// Setup Logging
	// This writer has nice colored output, but is not very fast.
	// In Prod another writer should be used. See Readme of zerolog

	if fc.Log.Handler == "dev" {
		log.Logger = log.Output(
			zerolog.ConsoleWriter{
				Out:     os.Stderr,
				NoColor: false,
			},
		)
	} else if fc.Log.Handler != "prod" {
		log.Fatal().Msg("Log Handler has to be either dev or prod")
	}

	if gin.IsDebugging() {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		switch fc.Log.Level {
		case "debug":
			zerolog.SetGlobalLevel(zerolog.DebugLevel)
		case "info":
			zerolog.SetGlobalLevel(zerolog.InfoLevel)
		case "warn":
			zerolog.SetGlobalLevel(zerolog.WarnLevel)
		case "error":
			zerolog.SetGlobalLevel(zerolog.ErrorLevel)
		case "fatal":
			zerolog.SetGlobalLevel(zerolog.FatalLevel)
		case "panic":
			zerolog.SetGlobalLevel(zerolog.PanicLevel)
		default:
			zerolog.SetGlobalLevel(zerolog.InfoLevel)
			log.Info().Msg("No Loglevel specified, using 'info'")
		}
	}

	var nodeID = geohash.Encode(fc.Location.Lat, fc.Location.Lng)

	var is data.Service
	var ks keygroup.Service
	var rs replication.Service

	var i data.Store
	var k keygroup.Store
	var n replication.Store

	switch fc.Storage.Adaptor {
	case "leveldb":
		var ldbc struct {
			Config struct {
				Path string `toml:"path"`
			} `toml:"leveldb"`
		}

		if _, err := toml.DecodeFile(*configPath, &ldbc); err != nil {
			log.Fatal().Err(err).Msg("invalid leveldb configuration!")
		}

		log.Print(ldbc)

		i = leveldbsd.New(ldbc.Config.Path)
	case "memory":
		i = memorysd.New()
	default:
		log.Fatal().Msg("unknown storage backend")
	}

	// Add more options here
	k = memorykg.New()
	n = memoryns.New()

	is = data.New(i)
	c := zmqclient.NewClient()

	ks = keygroup.New(k, nodeID)

	rs = replicationhandler.New(n, c)

	addr := fmt.Sprintf("%s:%d", fc.Server.Host, fc.Server.Port)

	extH := exthandler.New(is, ks, rs)
	intH := inthandler.New(is, ks)

	// Add more options here
	zmqH := memoryzmq.New(intH)

	zmqServer, err := zmqserver.Setup(fc.ZMQ.Port, nodeID, zmqH)

	if err != nil {
		panic("Cannot start zmqServer")
	}

	log.Fatal().Err(webserver.Setup(addr, extH)).Msg("Websever.Setup")

	// Shutdown
	zmqServer.Shutdown()
	c.Destroy()
}
