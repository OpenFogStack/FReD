package main

import (
	"flag"
	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/exthandler"
	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/memorykg"
	"log"

	"github.com/mmcloughlin/geohash"

	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/data"
	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/keygroup"
	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/leveldbsd"
	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/memorysd"
	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/webserver"
)

var addr = flag.String("addr", ":9001", "http service address")

var lat = *flag.Float64("lat", 52.514927933123914, "latitude of the server")
var lng = *flag.Float64("lng", 13.32676300345363, "longitude of the server")

var storageRuntime = *flag.String("storage-runtime", "leveldb", "storage runtime to use, e.g. leveldb or memory")
var dbPath = *flag.String("db-path", "./db", "path to use for database (only for databases that write to the file system, ignored otherwise)")

func main() {
	var is data.Service
	var ks keygroup.Service

	var i data.Store
	var k keygroup.Store

	switch storageRuntime {
	case "leveldb":
		i = leveldbsd.New(dbPath)
	case "memory":
		i = memorysd.New()
	default:
		panic("unknown storage backend")
	}

	k = memorykg.New()

	is = data.New(i)
	ks = keygroup.New(k, geohash.Encode(lat, lng))

	exthandler := exthandler.New(is, ks)
	//inthandler := inthandler.New(is, ks)

	log.Fatal(webserver.Setup(*addr, exthandler))
}
