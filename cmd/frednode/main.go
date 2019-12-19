package main

import (
	"flag"
	"log"

	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/app"
	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/webserver"
)

var addr = flag.String("addr", ":9001", "http service address")

var lat = *flag.Float64("lat", 52.514927933123914, "latitude of the server")
var lng = *flag.Float64("lng", 13.32676300345363, "longitude of the server")

var storageRuntime = *flag.String("storage-runtime", "leveldb", "storage runtime to use, e.g. leveldb or memory")
var dbPath = *flag.String("db-path", "./db", "path to use for database (only for databases that write to the file system, ignored otherwise)")

func main() {
	a, err := app.New(lat, lng, storageRuntime, dbPath)

	if err != nil {
		panic(err)
	}

	log.Fatal(webserver.SetupRouter(*addr, a))
}
