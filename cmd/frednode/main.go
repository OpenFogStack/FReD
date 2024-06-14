package main

import (
	"flag"
	"os"
	"os/signal"
	"runtime"
	"runtime/pprof"
	"strings"
	"syscall"

	"git.tu-berlin.de/mcc-fred/fred/pkg/dynamo"
	"git.tu-berlin.de/mcc-fred/fred/pkg/storageclient"
	"github.com/caarlos0/env/v9"
	"github.com/go-errors/errors"
	"github.com/mmcloughlin/geohash"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"git.tu-berlin.de/mcc-fred/fred/pkg/api"
	"git.tu-berlin.de/mcc-fred/fred/pkg/badgerdb"
	"git.tu-berlin.de/mcc-fred/fred/pkg/etcdnase"
	"git.tu-berlin.de/mcc-fred/fred/pkg/fred"
	"git.tu-berlin.de/mcc-fred/fred/pkg/peering"
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
		Host          string `env:"HOST"`
		AdvertiseHost string `env:"ADVERTISE_HOST"`
		Proxy         string `env:"PROXY"`
		Cert          string `env:"CERT"`
		Key           string `env:"KEY"`
		CA            string `env:"CA_FILE"`
		SkipVerify    bool   `env:"SKIP_VERIFY"`
	}
	Storage struct {
		Adaptor string `env:"STORAGE_ADAPTOR"`
	}
	Peering struct {
		Host             string `env:"PEERING_HOST"`
		AdvertiseHost    string `env:"PEERING_ADVERTISE_HOST"`
		Proxy            string `env:"PEERING_PROXY"`
		Cert             string `env:"PEERING_CERT"`
		Key              string `env:"PEERING_KEY"`
		CA               string `env:"PEERING_CA"`
		AsyncReplication bool   `env:"PEERING_ASYNC_REPLICATION"`
		SkipVerify       bool   `env:"PEERING_SKIP_VERIFY"`
	}
	Log struct {
		Level   string `env:"LOG_LEVEL"`
		Handler string `env:"LOG_HANDLER"`
	}
	NaSe struct {
		Host       string `env:"NASE_HOST"`
		Cert       string `env:"NASE_CERT"`
		Key        string `env:"NASE_KEY"`
		CA         string `env:"NASE_CA"`
		SkipVerify bool   `env:"NASE_SKIP_VERIFY"`
		Cached     bool   `env:"NASE_CACHED"`
	}
	RemoteStore struct {
		Host       string `env:"REMOTE_STORAGE_HOST"`
		Cert       string `env:"REMOTE_STORAGE_CERT"`
		Key        string `env:"REMOTE_STORAGE_KEY"`
		CA         string `env:"REMOTE_STORAGE_CA"`
		SkipVerify bool   `env:"REMOTE_STORAGE_SKIP_VERIFY"`
	}
	DynamoDB struct {
		Table       string `env:"DYNAMODB_TABLE"`
		Region      string `env:"DYNAMODB_REGION"`
		PublicKey   string `env:"DYNAMODB_PUBLIC_KEY"`
		PrivateKey  string `env:"DYNAMODB_PRIVATE_KEY"`
		Endpoint    string `env:"DYNAMODB_ENDPOINT"`
		CreateTable bool   `env:"DYNAMODB_CREATETABLE"`
	}
	Bdb struct {
		Path string `env:"BADGERDB_PATH"`
	}
	Trigger struct {
		Cert       string `env:"TRIGGER_CERT"`
		Key        string `env:"TRIGGER_KEY"`
		CA         string `env:"TRIGGER_CA"`
		SkipVerify bool   `env:"TRIGGER_SKIP_VERIFY"`
		Async      bool   `env:"TRIGGER_ASYNC"`
	}
	Profiling struct {
		CPUProfPath string `env:"PROFILING_CPU_PATH"`
		MemProfPath string `env:"PROFILING_MEM_PATH"`
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
	flag.StringVar(&(fc.Server.AdvertiseHost), "advertise-host", "", "Publicly reachable host address of server for external connections. (Env: ADVERTISE_HOST)")
	flag.StringVar(&(fc.Server.Proxy), "host-proxy", "", "Publicly reachable host address of server for external connections (if behind a proxy). (Env: PROXY)")
	flag.StringVar(&(fc.Server.Cert), "cert", "", "Certificate for external connections. (Env: CERT)")
	flag.StringVar(&(fc.Server.Key), "key", "", "Key file for external connections. (Env: KEY)")
	flag.StringVar(&(fc.Server.CA), "ca-file", "", "Certificate authority root certificate file for external connections. (Env: CA_FILE)")
	flag.BoolVar(&(fc.Server.SkipVerify), "skip-verify", false, "Skip verification of client certificates. (Env: SKIP_VERIFY)")

	// peering configuration
	// this is the address that grpc will bind to (locally)
	flag.StringVar(&(fc.Peering.Host), "peer-host", "", "local address of this peering server. (Env: PEERING_HOST)")
	flag.StringVar(&(fc.Peering.AdvertiseHost), "peer-advertise-host", "", "Publicly reachable address of this peering server. (Env: PEERING_ADVERTISE_HOST)")
	// this is the address that the node will advertise to nase
	flag.StringVar(&(fc.Peering.Proxy), "peer-host-proxy", "", "Publicly reachable address of this peering server (if behind a proxy). (Env: PEERING_PROXY)")
	flag.StringVar(&(fc.Peering.Cert), "peer-cert", "", "Certificate for peering connection. (Env: PEERING_CERT)")
	flag.StringVar(&(fc.Peering.Key), "peer-key", "", "Key file for peering connection. (Env: PEERING_KEY)")
	flag.StringVar(&(fc.Peering.CA), "peer-ca", "", "Certificate authority root certificate file for peering connections. (Env: PEERING_CA)")
	flag.BoolVar(&(fc.Peering.AsyncReplication), "peer-async-replication", false, "Enable asynchronous replication. Experimental. (Env: PEERING_ASYNC_REPLICATION)")
	flag.BoolVar(&(fc.Peering.SkipVerify), "peer-skip-verify", false, "Skip verification of client certificates. (Env: PEERING_SKIP_VERIFY)")

	// storage configuration
	flag.StringVar(&(fc.Storage.Adaptor), "adaptor", "", "Storage adaptor, can be \"remote\", \"badgerdb\", \"memory\", \"dynamo\". (Env: STORAGE_ADAPTOR)")

	flag.StringVar(&(fc.RemoteStore.Host), "remote-storage-host", "", "Host address of GRPC Server for storage connection. (Env: REMOTE_STORAGE_HOST)")
	flag.StringVar(&(fc.RemoteStore.Cert), "remote-storage-cert", "", "Certificate for storage connection. (Env: REMOTE_STORAGE_CERT)")
	flag.StringVar(&(fc.RemoteStore.Key), "remote-storage-key", "", "Key file for storage connection. (Env: REMOTE_STORAGE_KEY)")
	flag.StringVar(&(fc.RemoteStore.CA), "remote-storage-ca", "", "Comma-separated list of CA certificate files for storage connection. (Env: REMOTE_STORAGE_KEY)")

	flag.StringVar(&(fc.DynamoDB.Table), "dynamo-table", "", "AWS table for DynamoDB storage backend. (Env: DYNAMODB_TABLE)")
	flag.StringVar(&(fc.DynamoDB.Region), "dynamo-region", "", "AWS region for DynamoDB storage backend. (Env: DYNAMODB_REGION)")
	flag.StringVar(&(fc.DynamoDB.PublicKey), "dynamo-public-key", "", "AWS public key for DynamoDB storage backend. (Env: DYNAMODB_PUBLIC_KEY)")
	flag.StringVar(&(fc.DynamoDB.PrivateKey), "dynamo-private-key", "", "AWS private key for DynamoDB storage backend. (Env: DYNAMODB_PRIVATE_KEY)")
	flag.StringVar(&(fc.DynamoDB.Endpoint), "dynamo-endpoint", "", "Endpoint of local DynamoDB instance. Leave empty to use DynamoDB on AWS. (Env: DYNAMODB_ENDPOINT)")
	flag.BoolVar(&(fc.DynamoDB.CreateTable), "dynamo-create-table", false, "Create table in DynamoDB on first run rather than using existing table. (Env: DYNAMODB_CREATETABLE)")

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
	flag.BoolVar(&(fc.NaSe.SkipVerify), "nase-skip-verify", false, "Skip verification of etcd certificates. (Env: NASE_SKIP_VERIFY)")
	flag.BoolVar(&(fc.NaSe.Cached), "nase-cached", false, "Flag to indicate, whether to use a cache for NaSe. (Env: NASE_CACHED)")

	// trigger node tls configuration
	flag.StringVar(&(fc.Trigger.Cert), "trigger-cert", "", "Certificate for trigger node connection. (Env: TRIGGER_CERT)")
	flag.StringVar(&(fc.Trigger.Key), "trigger-key", "", "Key file for trigger node connection. (Env: TRIGGER_KEY)")
	flag.StringVar(&(fc.Trigger.CA), "trigger-ca", "", "Comma-separated list of CA certificate files for trigger node connection. (Env: TRIGGER_CA)")
	flag.BoolVar(&(fc.Trigger.SkipVerify), "trigger-skip-verify", false, "Skip verification of client certificates. (Env: TRIGGER_SKIP_VERIFY)")
	flag.BoolVar(&(fc.Trigger.Async), "trigger-async", false, "Enable asynchronous trigger replication. Experimental. (Env: TRIGGER_ASYNC)")

	flag.StringVar(&(fc.Profiling.CPUProfPath), "cpuprofile", "", "Enable CPU profiling and specify path for pprof output")
	flag.StringVar(&(fc.Profiling.MemProfPath), "memprofile", "", "Enable memory profiling and specify path for pprof output")

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

	if fc.Log.Level != "debug" && fc.Log.Level != "info" && fc.Log.Level != "warn" && fc.Log.Level != "error" && fc.Log.Level != "fatal" && fc.Log.Level != "panic" {
		flag.Usage()
		log.Fatal().Msgf("Given log level %s is not one of: \"debug\", \"info\" ,\"warn\", \"error\", \"fatal\", \"panic\".", fc.Log.Level)
	}

	if fc.Server.AdvertiseHost != "" && fc.Server.Proxy != "" {
		flag.Usage()
		log.Fatal().Msgf("You can only specify one of: \"host\", \"host-proxy\".")
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
	switch fc.Log.Handler {
	case "dev":
		log.Logger = log.Output(
			zerolog.ConsoleWriter{
				Out:     os.Stderr,
				NoColor: false,
			},
		)

		zerolog.DisableSampling(true)

		log.Info().Msgf("Current configuration:\n%+v", fc)
	case "prod":
		// add millisecond fraction to log
		// is also slightly faster
		zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMs
	default:
		log.Fatal().Msg("Log Handler has to be either dev or prod")
	}

	// https://github.com/influxdata/influxdb/blob/master/cmd/influxd/internal/profile/profile.go
	var prof struct {
		cpu *os.File
		mem *os.File
	}

	if fc.Profiling.CPUProfPath != "" {
		prof.cpu, err = os.Create(fc.Profiling.CPUProfPath)

		if err != nil {
			log.Fatal().Msgf("cpuprofile: %v", err)
		}

		err = pprof.StartCPUProfile(prof.cpu)

		if err != nil {
			log.Fatal().Msgf("cpuprofile: %v", err)
		}
	}

	if fc.Profiling.CPUProfPath != "" {
		prof.mem, err = os.Create(fc.Profiling.MemProfPath)

		if err != nil {
			log.Fatal().Msgf("memprofile: %v", err)
		}

		runtime.MemProfileRate = 4096
	}

	// Uncomment to print json config
	// log.Debug().Msgf("Configuration: %is", (func() string {
	// 	is, _ := json.MarshalIndent(fc, "", "    ")
	// 	return string(is)
	// })())
	log.Info().Msgf("Current configuration:\n%+v", fc)

	switch fc.Log.Level {
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

	var store fred.Store

	switch fc.Storage.Adaptor {
	case "badgerdb":
		log.Debug().Msgf("badgerdb struct is: %+v", fc.Bdb)
		store = badgerdb.New(fc.Bdb.Path)
	case "memory":
		store = badgerdb.NewMemory()
	case "remote":
		store = storageclient.NewClient(fc.RemoteStore.Host, fc.RemoteStore.Cert, fc.RemoteStore.Key, strings.Split(fc.RemoteStore.CA, ","), fc.RemoteStore.SkipVerify)
	case "dynamo":
		store, err = dynamo.New(fc.DynamoDB.Table, fc.DynamoDB.Region, fc.DynamoDB.Endpoint, fc.DynamoDB.CreateTable)
		if err != nil {
			log.Fatal().Msgf("could not open a dynamo connection: %s", err.(*errors.Error).ErrorStack())
		}
	default:
		log.Fatal().Msg("unknown storage backend")
	}

	if fc.Server.AdvertiseHost == "" {
		fc.Server.AdvertiseHost = fc.Server.Host
	}

	if fc.Peering.AdvertiseHost == "" {
		fc.Peering.AdvertiseHost = fc.Peering.Host
	}

	log.Debug().Msg("Starting Interconnection Client...")
	c := peering.NewClient(fc.Peering.Cert, fc.Peering.Key, fc.Peering.CA, fc.Peering.SkipVerify)

	log.Debug().Msg("Starting NaSe Client...")

	var n fred.NameService
	n, err = etcdnase.NewNameService(fc.General.nodeID, []string{fc.NaSe.Host}, fc.NaSe.Cert, fc.NaSe.Key, fc.NaSe.CA, fc.NaSe.SkipVerify, fc.NaSe.Cached)

	if err != nil {
		log.Err(err).Msg(err.(*errors.Error).ErrorStack())
		panic(err)
	}

	f := fred.New(&fred.Config{
		Store:                   store,
		Client:                  c,
		NaSe:                    n,
		PeeringHost:             fc.Peering.AdvertiseHost,
		PeeringHostProxy:        fc.Peering.Proxy,
		PeeringAsyncReplication: fc.Peering.AsyncReplication,
		ExternalHost:            fc.Server.AdvertiseHost,
		ExternalHostProxy:       fc.Server.Proxy,
		TriggerCert:             fc.Trigger.Cert,
		TriggerKey:              fc.Trigger.Key,
		TriggerCA:               strings.Split(fc.Trigger.CA, ","),
		TriggerAsync:            fc.Trigger.Async,
	})

	log.Debug().Msg("Starting Interconnection Server...")
	is := peering.NewServer(fc.Peering.Host, f.I, fc.Peering.Cert, fc.Peering.Key, fc.Peering.CA, fc.Peering.SkipVerify)

	log.Debug().Msg("Starting GRPC Server for Client (==Externalconnection)...")
	isProxied := fc.Server.Proxy != "" && fc.Server.Host != fc.Server.Proxy
	es := api.NewServer(fc.Server.Host, f.E, fc.Server.Cert, fc.Server.Key, fc.Server.CA, fc.Server.SkipVerify, isProxied, fc.Server.Proxy)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit,
		os.Interrupt,
		syscall.SIGTERM)

	<-quit
	log.Info().Msg("FReD Node Closing Now!")
	c.Destroy()
	log.Err(is.Close()).Msg("closing peering server")
	log.Err(es.Close()).Msg("closing api server")
	log.Err(store.Close()).Msg("closing database")

	if prof.cpu != nil {
		pprof.StopCPUProfile()
		err = prof.cpu.Close()
		log.Err(err).Msg("stopping cpu profile")
		prof.cpu = nil
	}

	if prof.mem != nil {
		err = pprof.Lookup("heap").WriteTo(prof.mem, 0)
		log.Err(err).Msg("stopping mem profile")
		err = prof.mem.Close()
		log.Err(err).Msg("stopping mem profile")
		prof.mem = nil
	}
}
