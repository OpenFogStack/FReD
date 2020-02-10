package main

import (
	"flag"
	"fmt"
	"net"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"google.golang.org/grpc"

	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/data"
	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/leveldbsd"
	storage "gitlab.tu-berlin.de/mcc-fred/fred/pkg/storageconnection"
)

func main() {
	path := flag.String("path", "./db", "Path for leveldb")
	port := flag.Int("port", 1337, "Port for the server to listen to")
	loglevel := flag.String("loglevel", "dev", "dev=>pretty, prod=>json")
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatal().Msgf("failed to listen: %v", err)
	}

	// Setup Logging
	// In Dev the ConsoleWriter has nice colored output, but is not very fast.
	// In Prod the default handler is used. It writes json to stdout and is very fast.
	if *loglevel == "dev" {
		log.Logger = log.Output(
			zerolog.ConsoleWriter{
				Out:     os.Stderr,
				NoColor: false,
			},
		)
	} else if *loglevel != "prod" {
		log.Fatal().Msg("Log Handler has to be either dev or prod")
	}

	var store data.Store = leveldbsd.New(*path)
	grpcServer := grpc.NewServer()
	storage.RegisterDatabaseServer(grpcServer, storage.NewStorageServer(&store))
	log.Debug().Msgf("Server is listening on port %d", *port)
	grpcServer.Serve(lis)
	log.Debug().Msg("Server is done.")
}
