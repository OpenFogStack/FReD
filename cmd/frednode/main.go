package main

// leave this in for cgo to work

import "C"

import (
	"encoding/json"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/alecthomas/kingpin"
	"github.com/mmcloughlin/geohash"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/data"
	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/exthandler"
	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/inthandler"
	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/keygroup"
	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/leveldbsd"
	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/memorykg"
	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/memoryrs"
	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/memorysd"
	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/memoryzmq"
	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/replhandler"
	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/replication"
	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/webserver"
	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/zmqclient"
	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/zmqserver"
)

type fredConfig struct {
	Location struct {
		Lat float64 `toml:"lat"`
		Lng float64 `toml:"lng"`
	} `toml:"location"`
	Server struct {
		Host   string `toml:"host"`
		Port   int    `toml:"port"`
		UseTLS bool   `toml:"ssl"`
	} `toml:"webserver"`
	Storage struct {
		Adaptor string `toml:"adaptor"`
	} `toml:"storage"`
	ZMQ struct {
		Port int `toml:"port"`
	} `toml:"zmq"`
	Log struct {
		Level   string `toml:"level"`
		Handler string `toml:"handler"`
	} `toml:"log"`
}

const apiversion string = "/v0"

var (
	configPath = kingpin.Flag("config", "Path to .toml configuration file.").PlaceHolder("PATH").String()
	lat        = kingpin.Flag("lat", "Latitude of the node.").PlaceHolder("LATITUDE").Default("-200").Float64()   // Domain: [-90,90]
	lng        = kingpin.Flag("lng", "Longitude of the node.").PlaceHolder("LONGITUDE").Default("-200").Float64() // Domain: ]-180,180]
	wsHost     = kingpin.Flag("ws-host", "Host address of webserver.").String()
	wsPort     = kingpin.Flag("ws-port", "Port of webserver.").PlaceHolder("WS-PORT").Default("-1").Int() // Domain: [0,9999]
	wsSSL      = kingpin.Flag("use-tls", "Use TLS/SSL to serve over HTTPS. Works only if host argument is a FQDN.").PlaceHolder("USE-SSL").Bool()
	zmqPort    = kingpin.Flag("zmq-port", "Port of ZeroMQ.").PlaceHolder("ZMQ-PORT").Default("-1").Int() // Domain: [0,9999]
	adaptor    = kingpin.Flag("adaptor", "Storage adaptor, can be \"leveldb\", \"memory\".").Enum("leveldb", "memory")
	logLevel   = kingpin.Flag("log-level", "Log level, can be \"debug\", \"info\" ,\"warn\", \"error\", \"fatal\", \"panic\".").Enum("debug", "info", "warn", "errors", "fatal", "panic")
	handler    = kingpin.Flag("handler", "Mode of log handler, can be \"dev\", \"prod\".").Enum("dev", "prod")
)

func main() {
	kingpin.Version(apiversion)
	kingpin.HelpFlag.Short('h')
	kingpin.CommandLine.Help = "Fog Replicated Database"
	kingpin.Parse()

	var fc fredConfig
	if *configPath != "" {
		if _, err := toml.DecodeFile(*configPath, &fc); err != nil {
			log.Fatal().Err(err).Msg("Invalid configuration! Toml can not be decoded.")
		}
	}

	// replace with set cmd args, no real sanity checks
	// default value means unset -> don't replace
	// numbers have negative defaults outside their domain, simple domain checks are implemented
	// e.g. lat < -90 is ignored and toml is used (if available)
	if *lat >= -90 && *lat <= 90 {
		fc.Location.Lat = *lat
	}
	if *lng >= -180 && *lng <= 180 {
		fc.Location.Lng = *lng
	}
	if *wsHost != "" {
		fc.Server.Host = *wsHost
	}
	if *wsPort >= 0 {
		fc.Server.Port = *wsPort
	}
	if *wsSSL != false {
		fc.Server.UseTLS = *wsSSL
	}
	if *zmqPort >= 0 {
		fc.ZMQ.Port = *zmqPort
	}
	if *adaptor != "" {
		fc.Storage.Adaptor = *adaptor
	}
	if *logLevel != "" {
		fc.Log.Level = *logLevel
	}
	if *handler != "" {
		fc.Log.Handler = *handler
	}

	log.Debug().Msgf("Configuration: %s", (func() string {
		s, _ := json.MarshalIndent(fc, "", "    ")
		return string(s)
	})())

	// Setup Logging
	// In Dev the ConsoleWriter has nice colored output, but is not very fast.
	// In Prod the default handler is used. It writes json to stdout and is very fast.
	if fc.Log.Handler == "dev" {
		log.Logger = log.Output(
			zerolog.ConsoleWriter{
				Out:     os.Stderr,
				NoColor: false,
			},
		)
	} else if fc.Log.Handler != "prod" {
		log.Fatal().Msg("Log Handler has to be either dev or prod")
	}

	switch fc.Log.Level {
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "info":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "warn":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "errors":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	case "fatal":
		zerolog.SetGlobalLevel(zerolog.FatalLevel)
	case "panic":
		zerolog.SetGlobalLevel(zerolog.PanicLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
		log.Info().Msg("No Loglevel specified, using 'info'")
	}

	var nodeID = geohash.Encode(fc.Location.Lat, fc.Location.Lng)

	var is data.Service
	var ks keygroup.Service
	var rs replication.Service

	var i data.Store
	var k keygroup.Store
	var n replication.Store

	switch fc.Storage.Adaptor {
	case "leveldb":
		var ldbc struct {
			Config struct {
				Path string `toml:"path"`
			} `toml:"leveldb"`
		}

		if _, err := toml.DecodeFile(*configPath, &ldbc); err != nil {
			log.Fatal().Err(err).Msg("invalid leveldb configuration!")
		}

		// "%v": unly print field values. "%#v": also print field names
		log.Debug().Msgf("leveldb struct is: %#v", ldbc)

		i = leveldbsd.New(ldbc.Config.Path)
	case "memory":
		i = memorysd.New()
	default:
		log.Fatal().Msg("unknown storage backend")
	}

	// Add more options here
	k = memorykg.New()
	n = memoryrs.New(fc.ZMQ.Port)

	is = data.New(i)
	c := zmqclient.NewClient()

	ks = keygroup.New(k, nodeID)

	rs = replhandler.New(n, c)

	extH := exthandler.New(is, ks, rs)
	intH := inthandler.New(is, ks, rs)

	// Add more options here
	zmqH := memoryzmq.New(intH)

	zmqServer, err := zmqserver.Setup(fc.ZMQ.Port, nodeID, zmqH)

	if err != nil {
		panic("Cannot start zmqServer")
	}

	log.Fatal().Err(webserver.Setup(fc.Server.Host, fc.Server.Port, extH, apiversion, fc.Server.UseTLS)).Msg("Webserver.Setup")

	// Shutdown
	zmqServer.Shutdown()
	c.Destroy()
}
