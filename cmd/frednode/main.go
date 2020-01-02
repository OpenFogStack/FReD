package main

import (
	"flag"
	"fmt"
	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/exthandler"
	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/memorykg"
	"log"

	"github.com/BurntSushi/toml"
	"github.com/mmcloughlin/geohash"

	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/data"
	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/keygroup"
	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/leveldbsd"
	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/memorysd"
	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/webserver"
)

var configPath = *flag.String("config", "", "path to .toml configuration file")

type fredConfig struct {
	location locationConfig
	server webserverConfig
	storage adaptorConfig
}

type locationConfig struct {
	Lat float64
	Lng float64
}

type webserverConfig struct {
	Host string
	Port int
}

type adaptorConfig struct {
	Adaptor string
}

func main() {
	if configPath == "" {
		log.Fatal("no configuration specified")
	}

	var fc fredConfig
	if _, err := toml.DecodeFile(configPath, &fc); err != nil {
		log.Fatalf("invalid configuration! error: %s", err)
	}

	var is data.Service
	var ks keygroup.Service

	var i data.Store
	var k keygroup.Store

	switch fc.storage.Adaptor {
	case "leveldb":
		type leveldbConfig struct {
			path string `toml:"leveldb.path"`
		}

		var ldbc leveldbConfig

		if _, err := toml.DecodeFile(configPath, &ldbc); err != nil {
			log.Fatalf("invalid leveldb configuration! error: %s", err)
		}

		i = leveldbsd.New(ldbc.path)
	case "memory":
		i = memorysd.New()
	default:
		panic("unknown storage backend")
	}

	k = memorykg.New()

	is = data.New(i)
	ks = keygroup.New(k, geohash.Encode(fc.location.Lat, fc.location.Lng))

	e := exthandler.New(is, ks)
	//inthandler := inthandler.New(is, ks)

	addr := fmt.Sprintf("%s:%d", fc.server.Host, fc.server.Port)
	log.Fatal(webserver.Setup(addr, e))
}
