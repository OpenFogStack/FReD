package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"git.tu-berlin.de/mcc-fred/fred/pkg/api"
	"git.tu-berlin.de/mcc-fred/fred/pkg/dynamo"
	"git.tu-berlin.de/mcc-fred/fred/pkg/etcdnase"
	"git.tu-berlin.de/mcc-fred/fred/pkg/nasecache"
	"git.tu-berlin.de/mcc-fred/fred/pkg/peering"
	"git.tu-berlin.de/mcc-fred/fred/pkg/storageclient"
	"github.com/caarlos0/env/v6"
	"github.com/go-errors/errors"
	"github.com/mmcloughlin/geohash"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"git.tu-berlin.de/mcc-fred/fred/pkg/badgerdb"
	"git.tu-berlin.de/mcc-fred/fred/pkg/fred"
)

type fredConfig struct {
	General struct {
		nodeID string `env:"NODEID"`
	}
	Location struct {
		Lat float64 `env:"LAT"`
		Lng float64 `env:"LONG"`
	}
	Server struct {
		Host  string `env:"HOST"`
		Proxy string `env:"PROXY"`
		Cert  string `env:"CERT"`
		Key   string `env:"KEY"`
		CA    string `env:"CA_FILE"`
	}
	Storage struct {
		Adaptor string `env:"STORAGE_ADAPTOR"`
	}
	Peering struct {
		Host      string `env:"PEERING_HOST"`
		HostProxy string `env:"PEERING_PROXY"`
	}
	Log struct {
		Level   string `env:"LOG_LEVEL"`
		Handler string `env:"LOG_HANDLER"`
	}
	NaSe struct {
		Host   string `env:"NASE_HOST"`
		Cert   string `env:"NASE_CERT"`
		Key    string `env:"NASE_KEY"`
		CA     string `env:"NASE_CA"`
		Cached bool   `env:"NASE_CACHED"`
	}
	RemoteStore struct {
		Host string `env:"REMOTE_STORAGE_HOST"`
		Cert string `env:"REMOTE_STORAGE_CERT"`
		Key  string `env:"REMOTE_STORAGE_KEY"`
	}
	DynamoDB struct {
		Table      string `env:"DYNAMODB_TABLE"`
		Region     string `env:"DYNAMODB_REGION"`
		PublicKey  string `env:"DYNAMODB_PUBLIC_KEY"`
		PrivateKey string `env:"DYNAMODB_PRIVATE_KEY"`
	}
	Bdb struct {
		Path string `env:"BADGERDB_PATH"`
	}
	Trigger struct {
		Cert string `env:"TRIGGER_CERT"`
		Key  string `env:"TRIGGER_KEY"`
	}
}

func parseArgs() (fc fredConfig) {

	// General configuration
	flag.StringVar(&(fc.General.nodeID), "nodeID", "", "Unique ID of this node. Will be calculated from lat/long if omitted. (Env: NODEID)")

	// Location configuration
	flag.Float64Var(&(fc.Location.Lat), "lat", 0, "Latitude of the node. (Env: LAT)")   // Domain: [-90,90]
	flag.Float64Var(&(fc.Location.Lng), "lng", 0, "Longitude of the node. (Env: LONG)") // Domain: ]-180,180]

	// API server configuration
	flag.StringVar(&(fc.Server.Host), "host", "", "Host address of server for external connections. (Env: HOST)")
	flag.StringVar(&(fc.Server.Proxy), "host-proxy", "", "Publicly reachable host address of server for external connections (if behind a proxy). (Env: PROXY)")
	flag.StringVar(&(fc.Server.Cert), "cert", "", "Certificate for external connections. (Env: CERT)")
	flag.StringVar(&(fc.Server.Key), "key", "", "Key file for external connections. (Env: KEY)")
	flag.StringVar(&(fc.Server.CA), "ca-file", "", "Certificate authority root certificate file for external connections. (Env: CA_FILE)")

	// peering configuration
	// this is the address that grpc will bind to (locally)
	flag.StringVar(&(fc.Peering.Host), "peer-host", "", "local address of this peering server. (Env: PEERING_HOST)")
	// this is the address that the node will advertise to nase
	flag.StringVar(&(fc.Peering.HostProxy), "peer-host-proxy", "", "Publicly reachable address of this peering server (if behind a proxy). (Env: PEERING_PROXY)")

	// storage configuration
	flag.StringVar(&(fc.Storage.Adaptor), "adaptor", "", "Storage adaptor, can be \"remote\", \"badgerdb\", \"memory\", \"dynamo\". (Env: STORAGE_ADAPTOR)")

	flag.StringVar(&(fc.RemoteStore.Host), "remote-storage-host", "", "Host address of GRPC Server for storage connection. (Env: REMOTE_STORAGE_HOST)")
	flag.StringVar(&(fc.RemoteStore.Cert), "remote-storage-cert", "", "Certificate for storage connection. (Env: REMOTE_STORAGE_CERT)")
	flag.StringVar(&(fc.RemoteStore.Key), "remote-storage-key", "", "Key file for storage connection. (Env: REMOTE_STORAGE_KEY)")

	flag.StringVar(&(fc.DynamoDB.Table), "dynamo-table", "", "AWS table for DynamoDB storage backend. (Env: DYNAMODB_TABLE)")
	flag.StringVar(&(fc.DynamoDB.Region), "dynamo-region", "", "AWS region for DynamoDB storage backend. (Env: DYNAMODB_REGION)")
	flag.StringVar(&(fc.DynamoDB.PublicKey), "dynamo-public-key", "", "AWS public key for DynamoDB storage backend. (Env: DYNAMODB_PUBLIC_KEY)")
	flag.StringVar(&(fc.DynamoDB.PrivateKey), "dynamo-private-key", "", "AWS private key for DynamoDB storage backend. (Env: DYNAMODB_PRIVATE_KEY)")

	flag.StringVar(&(fc.Bdb.Path), "badgerdb-path", "", "Path to the BadgerDB database. (Env: BADGERDB_PATH)")

	// logging configuration
	flag.StringVar(&(fc.Log.Level), "log-level", "debug", "Log level, can be \"debug\", \"info\" ,\"warn\", \"error\", \"fatal\", \"panic\". (Env: LOG_LEVEL)")
	flag.StringVar(&(fc.Log.Handler), "handler", "dev", "Mode of log handler, can be \"dev\", \"prod\". (Env: LOG_HANDLER)")

	// Nameservice configuration
	// TODO this should be a list of nodes. One node is enough, but if we want reliability we should accept multiple etcd nodes
	flag.StringVar(&(fc.NaSe.Host), "nase-host", "", "Host where the etcd server runs. (Env: NASE_HOST)")
	flag.StringVar(&(fc.NaSe.Cert), "nase-cert", "", "Certificate file to authenticate against etcd. (Env: NASE_CERT)")
	flag.StringVar(&(fc.NaSe.Key), "nase-key", "", "Key file to authenticate against etcd. (Env: NASE_KEY)")
	flag.StringVar(&(fc.NaSe.CA), "nase-ca", "", "CA certificate file to authenticate against etcd. (Env: NASE_CA)")
	flag.BoolVar(&(fc.NaSe.Cached), "nase-cached", false, "Flag to indicate, whether to use a cache for NaSe. (Env: NASE_CACHED)")

	// trigger node tls configuration
	flag.StringVar(&(fc.Trigger.Cert), "trigger-cert", "", "Certificate for trigger node connection. (Env: TRIGGER_CERT)")
	flag.StringVar(&(fc.Trigger.Key), "trigger-key", "", "Key file for trigger node connection. (Env: TRIGGER_KEY)")
	flag.Parse()

	// override with ENV variables
	if err := env.Parse(&fc); err != nil {
		panic(err)
	}

	// validate some stuff
	if fc.Location.Lat > 90 || fc.Location.Lat < -90 {
		flag.Usage()
		log.Fatal().Msgf("Given latitude %f is not within latitude range from -90 to 90.", fc.Location.Lat)
	}

	if fc.Location.Lng > 180 || fc.Location.Lat < -180 {
		flag.Usage()
		log.Fatal().Msgf("Given longitutde %f is not within latitude range from -180 to 180.", fc.Location.Lng)
	}

	if fc.Storage.Adaptor != "remote" && fc.Storage.Adaptor != "badgerdb" && fc.Storage.Adaptor != "memory" && fc.Storage.Adaptor != "dynamo" {
		flag.Usage()
		log.Fatal().Msgf("Given storage adaptor %s is not one of: \"remote\", \"badgerdb\", \"memory\", \"dynamo\".", fc.Storage.Adaptor)
	}

	if fc.Log.Handler != "dev" && fc.Log.Handler != "prod" {
		flag.Usage()
		log.Fatal().Msgf("Given log handler %s is not one of: \"dev\", \"prod\".", fc.Log.Handler)
	}

	if fc.Log.Level != "debug" && fc.Log.Level != "info" && fc.Log.Level != "warn" && fc.Log.Level != "errors" && fc.Log.Level != "fatal" && fc.Log.Level != "panic" {
		flag.Usage()
		log.Fatal().Msgf("Given log level %s is not one of: \"debug\", \"info\" ,\"warn\", \"error\", \"fatal\", \"panic\".", fc.Log.Level)
	}

	if fc.General.nodeID == "" {
		fc.General.nodeID = geohash.Encode(fc.Location.Lat, fc.Location.Lng)
	}

	return
}

func main() {
	var err error

	fc := parseArgs()

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

	var n fred.NameService
	n, err = etcdnase.NewNameService(fc.General.nodeID, []string{fc.NaSe.Host}, fc.NaSe.Cert, fc.NaSe.Key, fc.NaSe.CA)

	if err != nil {
		log.Err(err).Msg(err.(*errors.Error).ErrorStack())
		panic(err)
	}

	if fc.NaSe.Cached {
		n, err = nasecache.NewNameServiceCache(n)

		if err != nil {
			log.Err(err).Msg(err.(*errors.Error).ErrorStack())
			panic(err)
		}
	}

	f := fred.New(&fred.Config{
		Store:             store,
		Client:            c,
		NaSe:              n,
		PeeringHost:       fc.Peering.Host,
		PeeringHostProxy:  fc.Peering.HostProxy,
		ExternalHost:      fc.Server.Host,
		ExternalHostProxy: fc.Server.Proxy,
		TriggerCert:       fc.Trigger.Cert,
		TriggerKey:        fc.Trigger.Key,
	})

	log.Debug().Msg("Starting Interconnection Server...")
	is := peering.NewServer(fc.Peering.Host, f.I)

	log.Debug().Msg("Starting GRPC Server for Client (==Externalconnection)...")
	isProxied := fc.Server.Proxy != "" && fc.Server.Host != fc.Server.Proxy
	es := api.NewServer(fc.Server.Host, f.E, fc.Server.Cert, fc.Server.Key, fc.Server.CA, isProxied, fc.Server.Proxy)

	// TODO this code should live somewhere where it is called every n seconds, but for testing purposes the easiest way
	// TODO to simulate an internet shutdown is via killing a node, so testing once at startup should be enough
	missedItems := n.RequestNodeStatus(n.GetNodeID())
	if missedItems != nil {
		log.Warn().Msg("NodeStatus: This node was offline has missed some updates, getting them from other nodes")
		for _, item := range missedItems {
			nodeID, addr := n.GetNodeWithBiggerExpiry(item.Keygroup)
			if addr == "" {
				log.Error().Msgf("NodeStatus: Was not able to find node that can provide item %s, skipping it...", item.Keygroup)
				continue
			}
			log.Info().Msgf("Getting item of KG %s ID %s from Node %s @ %s", string(item.Keygroup), item.ID, string(nodeID), addr)
			item, err := c.SendGetItem(addr, item.Keygroup, item.ID)
			if err != nil {
				log.Err(err).Msg("Was not able to get Items from node")
			}
			expiry, _ := n.GetExpiry(item.Keygroup)
			err = store.Update(string(item.Keygroup), item.ID, item.Val, expiry)
			if err != nil {
				log.Error().Msgf("Could not update missed item %s", item.ID)
			}
		}
	} else {
		log.Debug().Msg("NodeStatus: No updates were missed by this node.")
	}

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
