package main

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"io/ioutil"
	"net"
	"os"

	storage2 "gitlab.tu-berlin.de/mcc-fred/fred/proto/storage"
	"google.golang.org/grpc/credentials"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"google.golang.org/grpc"

	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/badgerdb"
	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/fred"
	storage "gitlab.tu-berlin.de/mcc-fred/fred/pkg/storageserver"
)

func main() {
	path := flag.String("path", "./db", "Path for badgerdb")
	host := flag.String("port", ":1337", "Host for the server to listen to")
	loglevel := flag.String("loglevel", "dev", "dev=>pretty, prod=>json")

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

	// Load server's certificate and private key
	serverCert, err := tls.LoadX509KeyPair(*cert, *key)

	if err != nil {
		log.Fatal().Msgf("could not load key pair: %v", err)
	}

	// Create a new cert pool and add our own CA certificate
	rootCAs := x509.NewCertPool()

	loaded, err := ioutil.ReadFile(*ca)

	if err != nil {
		log.Fatal().Msgf("unexpected missing certfile: %v", err)
	}

	rootCAs.AppendCertsFromPEM(loaded)

	// Create the credentials and return it
	config := &tls.Config{
		Certificates: []tls.Certificate{serverCert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    rootCAs,
		MinVersion:   tls.VersionTLS12,
	}

	var store fred.Store = badgerdb.New(*path)
	grpcServer := grpc.NewServer(grpc.Creds(credentials.NewTLS(config)))
	storage2.RegisterDatabaseServer(grpcServer, storage.NewStorageServer(&store))
	log.Debug().Msgf("Server is listening on port %s", *host)
	log.Fatal().Err(grpcServer.Serve(lis))
	log.Err(store.Close()).Msg("error closing database")
	log.Debug().Msg("Server is done.")
}
