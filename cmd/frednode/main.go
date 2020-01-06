package main

// leave this in for cgo to work

import "C"

import (
	"flag"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/BurntSushi/toml"
	"github.com/mmcloughlin/geohash"

	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/data"
	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/exthandler"
	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/inthandler"
	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/keygroup"
	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/leveldbsd"
	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/memorykg"
	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/memorysd"
	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/memoryzmq"
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
		Level string `toml:level`
	} `toml:log`
}

func main() {
	configPath := flag.String("config", "", "path to .toml configuration file")

	flag.Parse()

	if *configPath == "" {
		log.Fatal().Msg("no configuration specified")
	}

	var fc fredConfig
	if _, err := toml.DecodeFile(*configPath, &fc); err != nil {
		log.Fatal().Msgf("invalid configuration! error: %s", err)
	}

	// Setup Logging
	// This writer has nice colored output, but is not very fast.
	// In Prod another writer should be used. See Readme of zerolog
	log.Logger = log.Output(
		zerolog.ConsoleWriter{
			Out:     os.Stderr,
			NoColor: false,
		},
	)

	if gin.IsDebugging() {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		switch fc.Log.Level {
		case "debug": zerolog.SetGlobalLevel(zerolog.DebugLevel)
		case "info": zerolog.SetGlobalLevel(zerolog.InfoLevel)
		case "warn": zerolog.SetGlobalLevel(zerolog.WarnLevel)
		case "error": zerolog.SetGlobalLevel(zerolog.ErrorLevel)
		case "fatal": zerolog.SetGlobalLevel(zerolog.FatalLevel)
		case "panic": zerolog.SetGlobalLevel(zerolog.PanicLevel)
		default:
			zerolog.SetGlobalLevel(zerolog.InfoLevel)
			log.Info().Msg("No Loglevel specified, using 'info'")
		}
	}

	var nodeID = geohash.Encode(fc.Location.Lat, fc.Location.Lng)

	var is data.Service
	var ks keygroup.Service

	var i data.Store
	var k keygroup.Store

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
		log.Fatal().Msgf("unknown storage backend")
	}

	// Add more options here
	k = memorykg.New()

	is = data.New(i)

	ks = keygroup.New(k, nodeID)

	addr := fmt.Sprintf("%s:%d", fc.Server.Host, fc.Server.Port)

	extH := exthandler.New(is, ks)
	intH := inthandler.New(is, ks)

	// Add more options here
	zmqH := memoryzmq.New(intH)

	zmqServer, err := zmqserver.Setup(fc.ZMQ.Port, nodeID, zmqH)

	if err != nil {
		panic("Cannot start zmqServer")
	}

	zmqClient := zmqclient.NewClient()

	zmqClient.SendCreateKeygroup("localhost", fc.ZMQ.Port, "Hello")
	zmqClient.SendCreateKeygroup("localhost", fc.ZMQ.Port, "World")

	log.Fatal().Err(webserver.Setup(addr, extH)).Msg("Websever.Setup")

	// Shutdown
	zmqServer.Shutdown()
	zmqClient.Destroy()
}
