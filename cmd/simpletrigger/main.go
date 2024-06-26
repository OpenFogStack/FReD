package main

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"

	"git.tu-berlin.de/mcc-fred/fred/pkg/grpcutil"
	"git.tu-berlin.de/mcc-fred/fred/proto/trigger"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
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
func (s *Server) PutItemTrigger(_ context.Context, request *trigger.PutItemTriggerRequest) (*trigger.Empty, error) {
	log.Debug().Msgf("Trigger Node has rcvd PutItem. In: %+v", request)

	s.log = append(s.log, LogEntry{
		Op:  "put",
		Kg:  request.Keygroup,
		ID:  request.Id,
		Val: request.Val,
	})

	return &trigger.Empty{}, nil
}

// DeleteItemTrigger calls this Method on the Inthandler
func (s *Server) DeleteItemTrigger(_ context.Context, request *trigger.DeleteItemTriggerRequest) (*trigger.Empty, error) {
	log.Debug().Msgf("Trigger Node has rcvd DeleteItem. In: %+v", request)

	s.log = append(s.log, LogEntry{
		Op: "del",
		Kg: request.Keygroup,
		ID: request.Id,
	})

	return &trigger.Empty{}, nil
}

func startServer(cert string, key string, ca string, skipVerify bool, host string, wsHost string) {
	creds, _, err := grpcutil.GetCredsFromConfig(cert, key, []string{ca}, false, skipVerify, &tls.Config{ClientAuth: tls.RequireAndVerifyClientCert})

	if err != nil {
		log.Fatal().Err(err).Msg("Failed to get credentials")
	}

	s := &Server{grpc.NewServer(grpc.Creds(creds)), []LogEntry{}}

	lis, err := net.Listen("tcp", host)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to listen")
	}

	trigger.RegisterTriggerNodeServer(s.Server, s)

	log.Debug().Msgf("Interconnection Server is listening on %s", host)

	go func() {
		err := s.Server.Serve(lis)
		defer s.GracefulStop()
		if err != nil {
			log.Fatal().Msgf("Interconnection Server exited: %s", err.Error())
		}
	}()

	// start a http server that lets us see what the trigger node has received
	http.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
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
	loghandler := flag.String("log-handler", "dev", "dev=>pretty, prod=>json")
	cert := flag.String("cert", "", "certificate file for grpc server")
	key := flag.String("key", "", "key file for grpc server")
	ca := flag.String("ca-file", "", "CA root for grpc server")
	skipVerify := flag.Bool("skip-verify", false, "Skip verification of client certificates")

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

	startServer(*cert, *key, *ca, *skipVerify, *host, *wsHost)
}
