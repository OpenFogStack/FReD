package main

import (
	"crypto/tls"
	"flag"
	"net"
	"os"

	"git.tu-berlin.de/mcc-fred/fred/pkg/grpcutil"
	storage2 "git.tu-berlin.de/mcc-fred/fred/proto/storage"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"google.golang.org/grpc"

	"git.tu-berlin.de/mcc-fred/fred/pkg/badgerdb"
	"git.tu-berlin.de/mcc-fred/fred/pkg/fred"
	storage "git.tu-berlin.de/mcc-fred/fred/pkg/storageserver"
)

func main() {
	path := flag.String("path", "./db", "Path for badgerdb")
	host := flag.String("port", ":1337", "Host for the server to listen to")
	loghandler := flag.String("log-handler", "dev", "dev=>pretty, prod=>json")
	loglevel := flag.String("log-level", "debug", "Log level, can be \"debug\", \"info\" ,\"warn\", \"error\", \"fatal\", \"panic\".")

	cert := flag.String("cert", "", "certificate file for grpc server")
	key := flag.String("key", "", "key file for grpc server")
	ca := flag.String("ca-file", "", "CA root for grpc server")

	flag.Parse()
	lis, err := net.Listen("tcp", *host)
	if err != nil {
		log.Fatal().Msgf("failed to listen: %v", err)
	}

	// Setup Logging
	// In Dev the ConsoleWriter has nice colored output, but is not very fast.
	// In Prod the default handler is used. It writes json to stdout and is very fast.
	if *loghandler == "dev" {
		log.Logger = log.Output(
			zerolog.ConsoleWriter{
				Out:     os.Stderr,
				NoColor: false,
			},
		)

		zerolog.DisableSampling(true)
	} else if *loghandler != "prod" {
		log.Fatal().Msg("Log Handler has to be either dev or prod")
	}

	switch *loglevel {
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
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		log.Info().Msg("No Loglevel specified, using 'debug'")
	}

	creds, _, err := grpcutil.GetCredsFromConfig(*cert, *key, []string{*ca}, false, &tls.Config{ClientAuth: tls.RequireAndVerifyClientCert})

	if err != nil {
		log.Fatal().Err(err).Msg("Error getting credentials")
	}

	var store fred.Store = badgerdb.New(*path)
	grpcServer := grpc.NewServer(grpc.Creds(creds))
	storage2.RegisterDatabaseServer(grpcServer, storage.NewStorageServer(&store))
	log.Debug().Msgf("Server is listening on port %s", *host)
	log.Fatal().Err(grpcServer.Serve(lis))
	log.Err(store.Close()).Msg("error closing database")
	log.Debug().Msg("Server is done.")
}
