package main

import (
	"crypto/tls"
	"flag"
	"net"
	"os"
	"os/signal"
	"syscall"

	"git.tu-berlin.de/mcc-fred/fred/pkg/alexandra"
	"git.tu-berlin.de/mcc-fred/fred/pkg/grpcutil"
	"git.tu-berlin.de/mcc-fred/fred/proto/middleware"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

type config struct {
	lightHouse          string
	caCert              string
	alexandraCert       string
	alexandraKey        string
	alexandraSkipVerify bool
	alexandraInsecure   bool
	nodesCert           string
	nodesKey            string
	nodesSkipVerify     bool
	loglevel            string
	logHandler          string
	isProxied           bool
	proxyHost           string
	address             string
	experimental        bool
}

func parseArgs() (c config) {
	flag.StringVar(&(c.lightHouse), "lighthouse", "", "ip of the first fred node to connect to")
	flag.StringVar(&(c.caCert), "ca-cert", "", "certificate of the ca")
	flag.StringVar(&(c.alexandraCert), "alexandra-cert", "", "Certificate to check clients")
	flag.StringVar(&(c.alexandraKey), "alexandra-key", "", "key to show clients")
	flag.BoolVar(&(c.alexandraSkipVerify), "alexandra-skip-verify", false, "skip tls auth for clients")
	flag.BoolVar(&(c.alexandraInsecure), "alexandra-insecure", false, "disable tls for clients")
	flag.StringVar(&(c.nodesCert), "clients-cert", "", "Certificate to check fred nodes")
	flag.StringVar(&(c.nodesKey), "clients-key", "", "key to show fred nodes")
	flag.BoolVar(&(c.nodesSkipVerify), "clients-skip-verify", false, "skip tls auth for fred nodes")
	flag.StringVar(&(c.loglevel), "log-level", "debug", "Log level, can be \"debug\", \"info\" ,\"warn\", \"error\", \"fatal\", \"panic\".")
	flag.StringVar(&(c.logHandler), "log-handler", "dev", "dev or prod")
	flag.BoolVar(&(c.isProxied), "is-proxy", false, "Is this behind a proxy?")
	flag.StringVar(&(c.proxyHost), "proxy-host", "", "Proxy host if this is proxied")
	flag.StringVar(&(c.address), "address", "", "where to start the server")
	flag.BoolVar(&(c.experimental), "experimental", false, "enable experimental features")
	flag.Parse()
	return
}

func main() {
	c := parseArgs()

	// Setup Logging as always
	if c.logHandler == "dev" {
		log.Logger = log.Output(
			zerolog.ConsoleWriter{
				Out:     os.Stderr,
				NoColor: false,
			},
		)

		zerolog.DisableSampling(true)

		log.Info().Msgf("Current configuration:\n%+v", c)
	} else if c.logHandler != "prod" {
		log.Fatal().Msg("Log Handler has to be either dev or prod")
	}

	switch c.loglevel {
	case "trace":
		zerolog.SetGlobalLevel(zerolog.TraceLevel)
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

	if c.alexandraCert == "" {
		log.Fatal().Msg("alexandra server: no certificate file given")
	}

	if c.alexandraKey == "" {
		log.Fatal().Msg("alexandra server: no key file given")
	}

	if c.caCert == "" {
		log.Fatal().Msg("alexandra server: no root certificate file given")
	}

	// Setup alexandra
	m := alexandra.NewMiddleware(c.nodesCert, c.nodesKey, c.caCert, c.lightHouse, c.isProxied, c.proxyHost, c.experimental, c.nodesSkipVerify)

	// Setup TLS
	creds, _, err := grpcutil.GetCredsFromConfig(c.alexandraCert, c.alexandraKey, []string{c.caCert}, c.alexandraInsecure, c.alexandraSkipVerify, &tls.Config{ClientAuth: tls.RequireAndVerifyClientCert})

	if err != nil {
		log.Fatal().Err(err).Msg("Failed to get credentials")
		return
	}

	lis, err := net.Listen("tcp", c.address)

	if err != nil {
		log.Fatal().Err(err).Msg("Failed to listen")
		return
	}

	s := grpc.NewServer(grpc.Creds(creds))

	middleware.RegisterMiddlewareServer(s, m)

	log.Debug().Msgf("Alexandra Server is listening on %s", c.address)

	// Quitting stuff
	quit := make(chan os.Signal, 1)
	signal.Notify(quit,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	go func() {
		<-quit
		s.Stop()
	}()

	err = s.Serve(lis)

	if err != nil {
		log.Fatal().Msgf("Alexandra Server exited with error: %s", err.Error())
	}

}
