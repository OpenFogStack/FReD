package main

import (
	"flag"
	"fmt"
	"log"

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
}

func main() {
	configPath := flag.String("config", "", "path to .toml configuration file")

	flag.Parse()

	if *configPath == "" {
		log.Fatal("no configuration specified")
	}

	var fc fredConfig
	if _, err := toml.DecodeFile(*configPath, &fc); err != nil {
		log.Fatalf("invalid configuration! error: %s", err)
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
			log.Fatalf("invalid leveldb configuration! error: %s", err)
		}

		log.Print(ldbc)

		i = leveldbsd.New(ldbc.Config.Path)
	case "memory":
		i = memorysd.New()
	default:
		log.Fatalf("unknown storage backend")
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

	go zmqserver.Setup(fc.ZMQ.Port, nodeID, zmqH)

	log.Fatal(webserver.Setup(addr, extH))
}
