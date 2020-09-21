package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gitlab.tu-berlin.de/mcc-fred/fred/proto/trigger"
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

func main() {

	host := flag.String("host", ":3333", "Host for the server to listen to")
	wsHost := flag.String("wshost", ":80", "Host for the server to listen to")
	loglevel := flag.String("loglevel", "dev", "dev=>pretty, prod=>json")
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

	s := &Server{grpc.NewServer(), []LogEntry{}}

	lis, err := net.Listen("tcp", *host)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to listen")
	}

	trigger.RegisterTriggerNodeServer(s.Server, s)

	log.Debug().Msgf("Interconnection Server is listening on %s", *host)

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

	log.Fatal().Err(http.ListenAndServe(*wsHost, nil))
}
