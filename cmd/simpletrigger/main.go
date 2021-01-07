package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"

	"git.tu-berlin.de/mcc-fred/fred/proto/trigger"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// LogEntry is one entry of the trigger node log of operations that it has received.
type LogEntry struct {
	Op  string `json:"op"`
	Kg  string `json:"kg"`
	ID  string `json:"id"`
	Val string `json:"val"`
}

// Server is a grpc server that let's peers access the internal handler.
type Server struct {
	*grpc.Server
	log []LogEntry
}

// PutItemTrigger calls HandleUpdate on the Inthandler
func (s *Server) PutItemTrigger(_ context.Context, request *trigger.PutItemTriggerRequest) (*trigger.TriggerResponse, error) {
	log.Debug().Msgf("Trigger Node has rcvd PutItem. In: %#v", request)

	s.log = append(s.log, LogEntry{
		Op:  "put",
		Kg:  request.Keygroup,
		ID:  request.Id,
		Val: request.Val,
	})

	return &trigger.TriggerResponse{Status: trigger.EnumTriggerStatus_TRIGGER_OK}, nil
}

// DeleteItemTrigger calls this Method on the Inthandler
func (s *Server) DeleteItemTrigger(_ context.Context, request *trigger.DeleteItemTriggerRequest) (*trigger.TriggerResponse, error) {
	log.Debug().Msgf("Trigger Node has rcvd DeleteItem. In: %#v", request)

	s.log = append(s.log, LogEntry{
		Op: "del",
		Kg: request.Keygroup,
		ID: request.Id,
	})

	return &trigger.TriggerResponse{Status: trigger.EnumTriggerStatus_TRIGGER_OK}, nil
}

func startServer(cert string, key string, ca string, host string, wsHost string) {
	// Load server's certificate and private key
	serverCert, err := tls.LoadX509KeyPair(cert, key)

	if err != nil {
		log.Fatal().Msgf("could not load key pair: %v", err)
	}

	// Create a new cert pool and add our own CA certificate
	rootCAs := x509.NewCertPool()

	loaded, err := ioutil.ReadFile(ca)

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

	s := &Server{grpc.NewServer(grpc.Creds(credentials.NewTLS(config))), []LogEntry{}}

	lis, err := net.Listen("tcp", host)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to listen")
	}

	trigger.RegisterTriggerNodeServer(s.Server, s)

	log.Debug().Msgf("Interconnection Server is listening on %s", host)

	go func() {
		defer s.GracefulStop()
		log.Fatal().Err(s.Server.Serve(lis)).Msg("Interconnection Server")
	}()

	// start a http server that let's us see what the trigger node has received
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		l, err := json.Marshal(s.log)

		if err != nil {
			log.Err(err).Msg("error getting logs")
		}

		_, err = fmt.Fprint(w, string(l))

		if err != nil {
			log.Err(err).Msg("error getting logs")
		}

	})

	log.Fatal().Err(http.ListenAndServe(wsHost, nil))
}

func main() {

	host := flag.String("host", ":3333", "Host for the server to listen to")
	wsHost := flag.String("wshost", ":80", "Host for the server to listen to")
	loglevel := flag.String("loglevel", "dev", "dev=>pretty, prod=>json")
	cert := flag.String("cert", "", "certificate file for grpc server")
	key := flag.String("key", "", "key file for grpc server")
	ca := flag.String("ca-file", "", "CA root for grpc server")

	flag.Parse()

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

	startServer(*cert, *key, *ca, *host, *wsHost)
}
