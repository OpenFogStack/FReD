package main

// leave this in for cgo to work

import "C"

import (
	"os"

	"github.com/BurntSushi/toml"
	"github.com/alecthomas/kingpin"
	"github.com/mmcloughlin/geohash"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/badgerdb"
	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/fred"
	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/leveldb"
	storage "gitlab.tu-berlin.de/mcc-fred/fred/pkg/storageconnection"
	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/webserver"
	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/zmq"
)

type fredConfig struct {
	General struct {
		nodeID string `toml:nodeID`
	} `toml:general`
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
		Port int    `toml:"port"`
		Host string `toml:"host"`
	} `toml:"zmq"`
	Log struct {
		Level   string `toml:"level"`
		Handler string `toml:"handler"`
	} `toml:"log"`
	Remote struct {
		Host string `toml:"host"`
		Port int    `toml:"port"`
	} `toml:"remote"`
	Ldb struct {
		Path string `toml:"path"`
	} `toml:"leveldb"`
	NaSe struct {
		Host string `toml:"host"`
	} `toml:"nase"`
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
	wsPort            = kingpin.Flag("ws-port", "Port of webserver.").PlaceHolder("WS-PORT").Default("-1").Int() // Domain: [0,9999]
	wsSSL             = kingpin.Flag("use-tls", "Use TLS/SSL to serve over HTTPS. Works only if host argument is a FQDN.").PlaceHolder("USE-SSL").Bool()
	zmqPort           = kingpin.Flag("zmq-port", "Port of ZeroMQ.").PlaceHolder("ZMQ-PORT").Default("-1").Int() // Domain: [0,9999]
	zmqHost           = kingpin.Flag("zmq-host", "(Publicly reachable) address of this zmq server.").String()
	adaptor           = kingpin.Flag("adaptor", "Storage adaptor, can be \"leveldb\", \"remote\", \"badgerdb\", \"memory\".").Enum("leveldb", "remote", "badgerdb", "memory")
	logLevel          = kingpin.Flag("log-level", "Log level, can be \"debug\", \"info\" ,\"warn\", \"error\", \"fatal\", \"panic\".").Enum("debug", "info", "warn", "errors", "fatal", "panic")
	handler           = kingpin.Flag("handler", "Mode of log handler, can be \"dev\", \"prod\".").Enum("dev", "prod")
	remoteStorageHost = kingpin.Flag("remote-storage-host", "Host address of GRPC Server.").String()
	remoteStoragePort = kingpin.Flag("remote-storage-port", "Port of GRPC Server.").PlaceHolder("WS-PORT").Default("-1").Int()
	ldbPath           = kingpin.Flag("leveldb-path", "Path to the leveldb database").String()
	// TODO this should be a list of nodes. One node is enough, but if we want reliability we should accept multiple etcd nodes
	naseHost = kingpin.Flag("naseHost", "Host where the etcd-server runs").String()
	bdbPath  = kingpin.Flag("badgerdb-path", "Path to the badgerdb database").String()
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
	if *nodeID != "" {
		fc.General.nodeID = *nodeID
	}
	if *lat >= -90 && *lat <= 90 {
		fc.Location.Lat = *lat
	}
	if *lng >= -180 && *lng <= 180 {
		fc.Location.Lng = *lng
	}
	// If no NodeID is provided use lat/long as NodeID
	if fc.General.nodeID == "" {
		fc.General.nodeID = geohash.Encode(fc.Location.Lat, fc.Location.Lng)
	}
	if *wsHost != "" {
		fc.Server.Host = *wsHost
	}
	if *wsPort >= 0 {
		fc.Server.Port = *wsPort
	}
	if *wsSSL {
		fc.Server.UseTLS = *wsSSL
	}
	if *zmqHost != "" {
		fc.ZMQ.Host = *zmqHost
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
	if *remoteStorageHost != "" {
		fc.Remote.Host = *remoteStorageHost
	}
	if *remoteStoragePort >= 0 {
		fc.Remote.Port = *remoteStoragePort
	}
	if *ldbPath != "" {
		fc.Ldb.Path = *ldbPath
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
	log.Debug().Msgf("%v", fc)

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

	var store fred.Store

	switch fc.Storage.Adaptor {
	case "leveldb":
		// "%v": unly print field values. "%#v": also print field names
		log.Debug().Msgf("leveldb struct is: %#v", fc.Ldb)
		store = leveldb.New(fc.Ldb.Path)
		defer log.Err(store.Close()).Msg("error closing database")
	case "badgerdb":
		log.Debug().Msgf("badgerdb struct is: %#v", fc.Ldb)
		store = badgerdb.New(fc.Bdb.Path)
		defer log.Err(store.Close()).Msg("error closing database")
	case "memory":
		store = badgerdb.NewMemory()
		defer log.Err(store.Close()).Msg("error closing database")
	case "remote":
		store = storage.NewClient(fc.Remote.Host, fc.Remote.Port)
		defer log.Err(store.Close()).Msg("error closing database")
	default:
		log.Fatal().Msg("unknown storage backend")
	}

	c := zmq.NewClient()

	f := fred.New(&fred.Config{
		Store:     store,
		Client:    c,
		ZmqPort:   fc.ZMQ.Port,
		NodeID:    nodeID,
		NaSeHosts: []string{fc.NaSe.Host},
	})

	// Add more options here
	zmqH := zmq.New(f.I)

	zmqServer, err := zmq.Setup(fc.ZMQ.Port, nodeID, zmqH)

	if err != nil {
		panic("Cannot start zmqServer")
	}

	log.Fatal().Err(webserver.Setup(fc.Server.Host, fc.Server.Port, f.E, apiversion, fc.Server.UseTLS)).Msg("Webserver.Setup")

	// Shutdown
	zmqServer.Shutdown()
	c.Destroy()
}
