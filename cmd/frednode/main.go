package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/BurntSushi/toml"
	"github.com/go-errors/errors"
	"github.com/mmcloughlin/geohash"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/api"
	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/dynamo"
	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/etcdnase"
	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/peering"
	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/storageclient"
	"gopkg.in/alecthomas/kingpin.v2"

	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/badgerdb"
	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/fred"
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
		Host  string `toml:"host"`
		Proxy string `toml:"proxy"`
		Cert  string `toml:"cert"`
		Key   string `toml:"key"`
		CA    string `toml:"ca"`
	} `toml:"server"`
	Storage struct {
		Adaptor string `toml:"adaptor"`
	} `toml:"storage"`
	Peering struct {
		Host      string `toml:"host"`
		HostProxy string `toml:"hostproxy"`
	} `toml:"peering"`
	Log struct {
		Level   string `toml:"level"`
		Handler string `toml:"handler"`
	} `toml:"log"`
	NaSe struct {
		Host string `toml:"host"`
		Cert string `toml:"cert"`
		Key  string `toml:"key"`
		CA   string `toml:"ca"`
	} `toml:"nase"`
	RemoteStore struct {
		Host string `toml:"host"`
		Cert string `toml:"cert"`
		Key  string `toml:"key"`
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
	Trigger struct {
		Cert string `toml:"cert"`
		Key  string `toml:"key"`
	} `toml:"trigger"`
}

var (
	// General configuration
	nodeID     = kingpin.Flag("nodeID", "Unique ID of this node. Will be calculated from lat/long if omitted").String()
	configPath = kingpin.Flag("config", "Path to .toml configuration file.").PlaceHolder("PATH").String()
	lat        = kingpin.Flag("lat", "Latitude of the node.").PlaceHolder("LATITUDE").Default("-200").Float64()   // Domain: [-90,90]
	lng        = kingpin.Flag("lng", "Longitude of the node.").PlaceHolder("LONGITUDE").Default("-200").Float64() // Domain: ]-180,180]

	// API server configuration
	grpcHost      = kingpin.Flag("host", "Host address of server for external connections.").String()
	grpcHostProxy = kingpin.Flag("host-proxy", "Publicly reachable host address of server for external connections (if behind a proxy).").String()
	grpcCert      = kingpin.Flag("cert", "Certificate for external connections.").String()
	grpcKey       = kingpin.Flag("key", "Key file for external connections.").String()
	grpcCA        = kingpin.Flag("ca-file", "Certificate authority root certificate file for external connections.").String()

	// peering configuration
	// this is the address that grpc will bind to (locally)
	peerHost = kingpin.Flag("peer-host", "local address of this peering server.").String()
	// this is the address that the node will advertise to nase
	peerHostProxy = kingpin.Flag("peer-host-proxy", "Publicly reachable address of this peering server (if behind a proxy)").String()

	// storage configuration
	adaptor = kingpin.Flag("adaptor", "Storage adaptor, can be \"remote\", \"badgerdb\", \"memory\", \"dynamo\".").Enum("remote", "badgerdb", "memory", "dynamo")

	remoteStorageHost = kingpin.Flag("remote-storage-host", "Host address of GRPC Server for storage connection.").String()
	remoteStorageCert = kingpin.Flag("remote-storage-cert", "Certificate for storage connection.").String()
	remoteStorageKey  = kingpin.Flag("remote-storage-key", "Key file for storage connection.").String()

	dynamoTable  = kingpin.Flag("dynamo-table", "AWS table for DynamoDB storage backend.").String()
	dynamoRegion = kingpin.Flag("dynamo-region", "AWS region for DynamoDB storage backend.").String()

	bdbPath = kingpin.Flag("badgerdb-path", "Path to the badgerdb database").String()

	// logging configuration
	logLevel = kingpin.Flag("log-level", "Log level, can be \"debug\", \"info\" ,\"warn\", \"error\", \"fatal\", \"panic\".").Enum("debug", "info", "warn", "errors", "fatal", "panic", "")
	handler  = kingpin.Flag("handler", "Mode of log handler, can be \"dev\", \"prod\".").Enum("dev", "prod")

	// Nameservice configuration
	// TODO this should be a list of nodes. One node is enough, but if we want reliability we should accept multiple etcd nodes
	naseHost = kingpin.Flag("nase-host", "Host where the etcd-server runs").String()
	naseCert = kingpin.Flag("nase-cert", "Certificate file to authenticate against etcd").String()
	naseKey  = kingpin.Flag("nase-key", "Key file to authenticate against etcd").String()
	naseCA   = kingpin.Flag("nase-ca", "CA certificate file to authenticate against etcd").String()

	// trigger node tls configuration
	triggerCert = kingpin.Flag("trigger-cert", "Certificate for trigger node connection.").String()
	triggerKey  = kingpin.Flag("trigger-key", "Key file for trigger node connection.").String()
)

func main() {
	var err error

	kingpin.Version("v0")
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
	if *grpcHost != "" {
		fc.Server.Host = *grpcHost
	}
	if *grpcHostProxy != "" {
		fc.Server.Proxy = *grpcHostProxy
	}
	if *grpcCert != "" {
		fc.Server.Cert = *grpcCert
	}
	if *grpcKey != "" {
		fc.Server.Key = *grpcKey
	}
	if *grpcCA != "" {
		fc.Server.CA = *grpcCA
	}
	if *peerHost != "" {
		fc.Peering.Host = *peerHost
	}
	if *peerHostProxy != "" {
		fc.Peering.HostProxy = *peerHostProxy
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
	if *remoteStorageCert != "" {
		fc.RemoteStore.Cert = *remoteStorageCert
	}
	if *remoteStorageKey != "" {
		fc.RemoteStore.Key = *remoteStorageKey
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
	if *naseCert != "" {
		fc.NaSe.Cert = *naseCert
	}
	if *naseKey != "" {
		fc.NaSe.Key = *naseKey
	}
	if *naseCA != "" {
		fc.NaSe.CA = *naseCA
	}
	if *bdbPath != "" {
		fc.Bdb.Path = *bdbPath
	}
	if *triggerCert != "" {
		fc.Trigger.Cert = *triggerCert
	}
	if *triggerKey != "" {
		fc.Trigger.Key = *triggerKey
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
	// log.Debug().Msgf("Configuration: %is", (func() string {
	// 	is, _ := json.MarshalIndent(fc, "", "    ")
	// 	return string(is)
	// })())
	log.Debug().Msg("Current configuration:")
	log.Debug().Msgf("%#v", fc)

	switch fc.Log.Level {
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

	var store fred.Store

	switch fc.Storage.Adaptor {
	case "badgerdb":
		log.Debug().Msgf("badgerdb struct is: %#v", fc.Bdb)
		store = badgerdb.New(fc.Bdb.Path)
	case "memory":
		store = badgerdb.NewMemory()
	case "remote":
		store = storageclient.NewClient(fc.RemoteStore.Host, fc.RemoteStore.Cert, fc.RemoteStore.Key)
	case "dynamo":
		store, err = dynamo.New(fc.DynamoDB.Table, fc.DynamoDB.Region)
		if err != nil {
			log.Fatal().Msgf("could not open a dynamo connection: %s", err.(*errors.Error).ErrorStack())
		}
	default:
		log.Fatal().Msg("unknown storage backend")
	}

	log.Debug().Msg("Starting Interconnection Client...")
	c := peering.NewClient()

	log.Debug().Msg("Starting NaSe Client...")
	n, err := etcdnase.NewNameService(fc.General.nodeID, []string{fc.NaSe.Host}, fc.NaSe.Cert, fc.NaSe.Key, fc.NaSe.CA)

	if err != nil {
		log.Err(err).Msg(err.(*errors.Error).ErrorStack())
		panic(err)
	}

	f := fred.New(&fred.Config{
		Store:            store,
		Client:           c,
		NaSe:             n,
		PeeringHost:      fc.Peering.Host,
		PeeringHostProxy: fc.Peering.HostProxy,
		TriggerCert:      fc.Trigger.Cert,
		TriggerKey:       fc.Trigger.Key,
	})

	log.Debug().Msg("Starting Interconnection Server...")
	is := peering.NewServer(fc.Peering.Host, f.I)

	log.Debug().Msg("Starting GRPC Server for Client (==Externalconnection)...")
	isProxied := fc.Server.Proxy != "" && fc.Server.Host != fc.Server.Proxy
	es := api.NewServer(fc.Server.Host, f.E, fc.Server.Cert, fc.Server.Key, fc.Server.CA, isProxied, fc.Server.Proxy)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	func() {
		<-quit
		c.Destroy()
		log.Err(is.Close()).Msg("error closing peering server")
		log.Err(es.Close()).Msg("error closing api server")
		log.Err(store.Close()).Msg("error closing database")
	}()
}
