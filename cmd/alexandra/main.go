package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"git.tu-berlin.de/mcc-fred/fred/pkg/alexandra"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type config struct {
	lightHouse    string
	caCert        string
	alexandraCert string
	alexandraKey  string
	nodesCert     string
	nodesKey      string
	loglevel      string
	logHandler    string
	isProxied     bool
	proxyHost     string
	address       string
}

func parseArgs() (c config) {
	flag.StringVar(&(c.lightHouse), "lighthouse", "", "ip of the first fred node to connect to")
	flag.StringVar(&(c.caCert), "ca-cert", "", "certificate of the ca")
	flag.StringVar(&(c.alexandraCert), "alexandra-cert", "", "Certificate to check clients")
	flag.StringVar(&(c.alexandraKey), "alexandra-key", "", "key to show clients")
	flag.StringVar(&(c.nodesCert), "clients-cert", "", "Certificate to check fred nodes")
	flag.StringVar(&(c.nodesKey), "clients-key", "", "key to show fred nodes")
	flag.StringVar(&(c.loglevel), "log-level", "debug", "Log level, can be \"debug\", \"info\" ,\"warn\", \"error\", \"fatal\", \"panic\".")
	flag.StringVar(&(c.logHandler), "log-handler", "dev", "dev or prod")
	flag.BoolVar(&(c.isProxied), "is-proxy", false, "Is this behind a proxy?")
	flag.StringVar(&(c.proxyHost), "proxy-host", "", "Proxy host if this is proxied")
	flag.StringVar(&(c.address), "address", "172.26.4.1:10000", "where to start the server")
	flag.Parse()
	return
}

func main() {
	c := parseArgs()

	log.Info().Msgf("%#v", c)

	// Setup Logging as always
	if c.logHandler == "dev" {
		log.Logger = log.Output(
			zerolog.ConsoleWriter{
				Out:     os.Stderr,
				NoColor: false,
			},
		)
	} else if c.logHandler != "prod" {
		log.Fatal().Msg("Log Handler has to be either dev or prod")
	}

	switch c.loglevel {
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

	// Setup alexandra
	server := alexandra.NewServer(c.address, c.caCert, c.alexandraCert, c.alexandraKey, c.nodesCert, c.nodesKey, c.lightHouse, c.isProxied, c.proxyHost)

	// Quitting stuff
	quit := make(chan os.Signal, 1)
	signal.Notify(quit,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	go func() {
		<-quit
		server.Stop()
	}()

	server.ServeBlocking()
}
