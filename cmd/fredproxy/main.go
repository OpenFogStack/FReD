package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"git.tu-berlin.de/mcc-fred/fred/pkg/proxy"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	// we need:
	// - a port for the grpc client frontend
	// - a port for the grpc peering endpoint
	// - a list of addresses (we're simplifying this: you can't set different ports for your different machines)
	// - probably also the usual certificate stuff
	// - log level and handler as always
	clientPort := flag.Int("client-port", 0, "port to bind proxy to for client access")
	peeringPort := flag.Int("peering-port", 0, "port to bind to for peering access")
	machines := flag.String("machines", "", "list of machine addresses, comma-separated")
	loghandler := flag.String("log-handler", "dev", "dev=>pretty, prod=>json")
	loglevel := flag.String("log-level", "debug", "Log level, can be \"debug\", \"info\" ,\"warn\", \"error\", \"fatal\", \"panic\".")
	peeringCert := flag.String("peer-cert", "", "Certificate for peering connection.")
	peeringKey := flag.String("peer-key", "", "Key file for peering connection.")
	peeringCA := flag.String("peer-ca", "", "Certificate authority root certificate file for peering connections.")
	peeringSkipVerify := flag.Bool("peer-skip-verify", false, "Skip verification of peer certificate.")
	apiCert := flag.String("api-cert", "", "Certificate for API connection.")
	apiKey := flag.String("api-key", "", "Key file for API connection.")
	apiCA := flag.String("api-ca", "", "Certificate authority root certificate file for API connections.")
	apiSkipVerify := flag.Bool("api-skip-verify", false, "Skip verification of API certificate.")

	flag.Parse()

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

	// parse machines
	p := proxy.NewProxy(strings.Split(*machines, ","))

	pS, err := proxy.StartPeeringProxy(p, *peeringPort, *peeringCert, *peeringKey, *peeringCA, *peeringSkipVerify)
	if err != nil {
		log.Fatal().Err(err).Msg(err.Error())
	}

	aS, err := proxy.StartAPIProxy(p, *clientPort, *apiCert, *apiKey, *apiCA, *apiSkipVerify)
	if err != nil {
		log.Fatal().Err(err).Msg(err.Error())
	}

	pLis, err := net.Listen("tcp", fmt.Sprintf(":%d", *peeringPort))

	if err != nil {
		log.Fatal().Err(err).Msg("Failed to listen")
	}

	go func() {
		log.Debug().Msgf("Peering proxy starting to listen on :%d, proxying to %s", *peeringPort, *machines)
		err := pS.Serve(pLis)

		// if Serve returns without an error, we probably intentionally closed it
		if err != nil {
			log.Fatal().Msgf("PeeringProxy Server exited: %s", err.Error())
		}
	}()

	aLis, err := net.Listen("tcp", fmt.Sprintf(":%d", *clientPort))

	if err != nil {
		log.Fatal().Err(err).Msg("Failed to listen")
	}

	go func() {
		log.Debug().Msgf("API proxy starting to listen on :%d, proxying to %s", *clientPort, *machines)
		err := aS.Serve(aLis)

		// if Serve returns without an error, we probably intentionally closed it
		if err != nil {
			log.Fatal().Msgf("APIProxy Server exited: %s", err.Error())
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit,
		os.Interrupt,
		syscall.SIGTERM)

	<-quit
	log.Info().Msg("FReD Proxy Closing Now!")
	pS.GracefulStop()
	aS.GracefulStop()
}
