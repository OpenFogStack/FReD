package main

import (
	"github.com/BurntSushi/toml"
	"github.com/alecthomas/kingpin"
	"github.com/go-errors/errors"
	"github.com/mmcloughlin/geohash"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/dynamoconnect"
	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/interconnection"
	"os"
	"os/signal"
	"syscall"

	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/badgerdb"
	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/fred"
	storage "gitlab.tu-berlin.de/mcc-fred/fred/pkg/storageconnection"
	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/webserver"
)

type fredConfig struct {
	General struct {
		nodeID string `toml:"nodeID"`
	} `toml:"general"`
	Location struct {
		Lat float64 `toml:"lat"`
		Lng float64 `toml:"lng"`
	} `toml:"location"`
	Server struct {
		Host   string `toml:"host"`
		UseTLS bool   `toml:"ssl"`
	} `toml:"webserver"`
	Storage struct {
		Adaptor string `toml:"adaptor"`
	} `toml:"storage"`
	Peering struct {
		Host string `toml:"host"`
	} `toml:"zmq"`
	Log struct {
		Level   string `toml:"level"`
		Handler string `toml:"handler"`
	} `toml:"log"`
	NaSe struct {
		Host string `toml:"host"`
	} `toml:"nase"`
	RemoteStore struct {
		Host string `toml:"host"`
	} `toml:"remotestore"`
	DynamoDB struct {
		Table      string `toml:"table"`
		Region     string `toml:"region"`
		PublicKey  string `toml:"publickey"`
		PrivateKey string `toml:"privatekey"`
	} `toml:"dynamodb"`
	Bdb struct {
		Path string `toml:"path"`
	} `toml:"badgerdb"`
}

const apiversion string = "/v0"

var (
	nodeID            = kingpin.Flag("nodeID", "Unique ID of this node. Will be calculated from lat/long if omitted").String()
	configPath        = kingpin.Flag("config", "Path to .toml configuration file.").PlaceHolder("PATH").String()
	lat               = kingpin.Flag("lat", "Latitude of the node.").PlaceHolder("LATITUDE").Default("-200").Float64()   // Domain: [-90,90]
	lng               = kingpin.Flag("lng", "Longitude of the node.").PlaceHolder("LONGITUDE").Default("-200").Float64() // Domain: ]-180,180]
	wsHost            = kingpin.Flag("ws-host", "Host address of webserver.").String()
	wsSSL             = kingpin.Flag("use-tls", "Use TLS/SSL to serve over HTTPS. Works only if host argument is a FQDN.").PlaceHolder("USE-SSL").Bool()
	zmqHost           = kingpin.Flag("zmq-host", "(Publicly reachable) address of this zmq server.").String()
	adaptor           = kingpin.Flag("adaptor", "Storage adaptor, can be \"remote\", \"badgerdb\", \"memory\", \"dynamo\".").Enum("remote", "badgerdb", "memory", "dynamo")
	logLevel          = kingpin.Flag("log-level", "Log level, can be \"debug\", \"info\" ,\"warn\", \"error\", \"fatal\", \"panic\".").Enum("debug", "info", "warn", "errors", "fatal", "panic")
	handler           = kingpin.Flag("handler", "Mode of log handler, can be \"dev\", \"prod\".").Enum("dev", "prod")
	remoteStorageHost = kingpin.Flag("remote-storage-host", "Host address of GRPC Server.").String()
	dynamoTable       = kingpin.Flag("dynamo-table", "AWS table for DynamoDB storage backend.").String()
	dynamoRegion      = kingpin.Flag("dynamo-region", "AWS region for DynamoDB storage backend.").String()

	// TODO this should be a list of nodes. One node is enough, but if we want reliability we should accept multiple etcd nodes
	naseHost = kingpin.Flag("nase-host", "Host where the etcd-server runs").String()
	bdbPath  = kingpin.Flag("badgerdb-path", "Path to the badgerdb database").String()
)

func main() {
	var err error

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
	if *nodeID != "" {
		fc.General.nodeID = *nodeID
	}
	// If no NodeID is provided use lat/long as NodeID
	if fc.General.nodeID == "" {
		fc.General.nodeID = geohash.Encode(fc.Location.Lat, fc.Location.Lng)
	}
	if *lat >= -90 && *lat <= 90 {
		fc.Location.Lat = *lat
	}
	if *lng >= -180 && *lng <= 180 {
		fc.Location.Lng = *lng
	}
	if *wsHost != "" {
		fc.Server.Host = *wsHost
	}
	if *wsSSL {
		fc.Server.UseTLS = *wsSSL
	}
	if *zmqHost != "" {
		fc.Peering.Host = *zmqHost
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
	if *remoteStorageHost != "" {
		fc.RemoteStore.Host = *remoteStorageHost
	}
	if *dynamoTable != "" {
		fc.DynamoDB.Table = *dynamoTable
	}
	if *dynamoRegion != "" {
		fc.DynamoDB.Region = *dynamoRegion
	}
	if *naseHost != "" {
		fc.NaSe.Host = *naseHost
	}
	if *bdbPath != "" {
		fc.Bdb.Path = *bdbPath
	}

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
		log.Fatal().Msg("Log ExtHandler has to be either dev or prod")
	}

	// Uncomment to print json config
	// log.Debug().Msgf("Configuration: %s", (func() string {
	// 	s, _ := json.MarshalIndent(fc, "", "    ")
	// 	return string(s)
	// })())
	log.Debug().Msg("Current configuration:")
	log.Debug().Msgf("%#v", fc)

	wsLogLevel := "release"

	switch fc.Log.Level {
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		wsLogLevel = "debug"
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

	var store fred.Store

	switch fc.Storage.Adaptor {
	case "badgerdb":
		log.Debug().Msgf("badgerdb struct is: %#v", fc.Bdb)
		store = badgerdb.New(fc.Bdb.Path)
	case "memory":
		store = badgerdb.NewMemory()
	case "remote":
		store = storage.NewClient(fc.RemoteStore.Host)
	case "dynamo":
		store, err = dynamoconnect.New(fc.DynamoDB.Table, fc.DynamoDB.Region)
		if err != nil {
			log.Fatal().Msgf("could not open a dynamo connection: %s", err.(*errors.Error).ErrorStack())
		}
	default:
		log.Fatal().Msg("unknown storage backend")
	}

	log.Debug().Msg("Starting Interconnection Client...")
	c := interconnection.NewClient()

	f := fred.New(&fred.Config{
		Store:       store,
		Client:      c,
		PeeringHost: fc.Peering.Host,
		NodeID:      fc.General.nodeID,
		NaSeHosts:   []string{fc.NaSe.Host},
	})

	log.Debug().Msg("Starting Interconnection Server...")
	s := interconnection.NewServer(fc.Peering.Host, f.I)

	log.Debug().Msg("Starting Web Server...")
	go log.Fatal().Err(webserver.Setup(fc.Server.Host, f.E, apiversion, fc.Server.UseTLS, wsLogLevel)).Msg("Webserver.Setup")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	func() {
		<-quit
		c.Destroy()
		log.Err(s.Close()).Msg("error closing interconnection server")
		log.Err(store.Close()).Msg("error closing database")
	}()
}
