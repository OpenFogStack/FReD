package api

import (
	"context"
	"net"

	"github.com/rs/zerolog/log"
	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/fred"
	"gitlab.tu-berlin.de/mcc-fred/fred/proto/client"
	"google.golang.org/grpc"
)

// Server handles GRPC Requests and calls the according functions of the exthandler
type Server struct {
	e fred.ExtHandler
	*grpc.Server
}

// NewServer creates a new Server for requests from Fred Clients
func NewServer(host string, handler fred.ExtHandler) *Server {

	s := &Server{handler, grpc.NewServer()}

	lis, err := net.Listen("tcp", host)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to listen")
		return nil
	}

	client.RegisterClientServer(s.Server, s)

	log.Debug().Msgf("Externalconnection Server is listening on %s", host)

	go func() {
		log.Fatal().Err(s.Server.Serve(lis)).Msg("Externalconnection Server")
	}()

	return s
}

// Close closes the grpc server for internal communication.
func (s *Server) Close() error {
	s.Server.GracefulStop()
	return nil
}

func statusResponseFromError(err error) (*client.StatusResponse, error) {
	if err == nil {
		return &client.StatusResponse{Status: client.EnumStatus_OK}, nil
	}
	log.Debug().Msgf("ExtServer is returning error: %#v", err)
	return &client.StatusResponse{Status: client.EnumStatus_ERROR, ErrorMessage: err.Error()}, err

}

// CreateKeygroup calls this method on the exthandler
func (s *Server) CreateKeygroup(_ context.Context, request *client.CreateKeygroupRequest) (*client.StatusResponse, error) {
	log.Debug().Msgf("ExtServer has rcvd CreateKeygroup. In: %#v", request)
	err := s.e.HandleCreateKeygroup(fred.Keygroup{Name: fred.KeygroupName(request.Keygroup), Mutable: request.Mutable, Expiry: int(request.Expiry)})
	return statusResponseFromError(err)
}

// DeleteKeygroup calls this method on the exthandler
func (s *Server) DeleteKeygroup(_ context.Context, request *client.DeleteKeygroupRequest) (*client.StatusResponse, error) {
	log.Debug().Msgf("ExtServer has rcvd DeleteKeygroup. In: %#v", request)
	err := s.e.HandleDeleteKeygroup(fred.Keygroup{Name: fred.KeygroupName(request.Keygroup)})
	return statusResponseFromError(err)
}

// Read calls this method on the exthandler
func (s *Server) Read(_ context.Context, request *client.ReadRequest) (*client.ReadResponse, error) {
	log.Debug().Msgf("ExtServer has rcvd Read. In: %#v", request)
	res, err := s.e.HandleRead(fred.Item{Keygroup: fred.KeygroupName(request.Keygroup), ID: request.Id})
	if err != nil {
		log.Debug().Msgf("ExtServer is returning error: %#v", err)
		return &client.ReadResponse{}, err
	}
	return &client.ReadResponse{Data: res.Val}, nil

}

// Update calls this method on the exthandler
func (s *Server) Update(_ context.Context, request *client.UpdateRequest) (*client.StatusResponse, error) {
	log.Debug().Msgf("ExtServer has rcvd DeleteKeygroup. In: %#v", request)
	err := s.e.HandleUpdate(fred.Item{Keygroup: fred.KeygroupName(request.Keygroup), ID: request.Id, Val: request.Data})
	return statusResponseFromError(err)
}

// Delete calls this method on the exthandler
func (s *Server) Delete(_ context.Context, request *client.DeleteRequest) (*client.StatusResponse, error) {
	log.Debug().Msgf("ExtServer has rcvd Delete. In: %#v", request)
	err := s.e.HandleDelete(fred.Item{Keygroup: fred.KeygroupName(request.Keygroup), ID: request.Id})
	return statusResponseFromError(err)
}

// AddReplica calls this method on the exthandler
func (s *Server) AddReplica(_ context.Context, request *client.AddReplicaRequest) (*client.StatusResponse, error) {
	log.Debug().Msgf("ExtServer has rcvd AddReplica. In: %#v", request)
	err := s.e.HandleAddReplica(fred.Keygroup{Name: fred.KeygroupName(request.Keygroup), Expiry: int(request.Expiry)}, fred.Node{ID: fred.NodeID(request.NodeId)})
	return statusResponseFromError(err)
}

// GetKeygroupReplica calls this method on the exthandler
func (s *Server) GetKeygroupReplica(_ context.Context, request *client.GetKeygroupReplicaRequest) (*client.GetKeygroupReplicaResponse, error) {
	log.Debug().Msgf("ExtServer has rcvd GetKeygroupReplica. In: %#v", request)
	n, e, err := s.e.HandleGetKeygroupReplica(fred.Keygroup{Name: fred.KeygroupName(request.Keygroup)})
	// Copy only the interesting values into a new array
	nodes := make([]string, len(n))
	expiries := make([]int64, len(n))
	for i := 0; i < len(n); i++ {
		nodes[i] = string(n[i].ID)
		expiries[i] = int64(e[n[i].ID])
	}
	if err != nil {
		log.Debug().Msgf("ExtServer is returning error: %#v", err)
		return &client.GetKeygroupReplicaResponse{}, err
	}
	return &client.GetKeygroupReplicaResponse{NodeId: nodes, Expiry: expiries}, nil

}

// RemoveReplica calls this method on the exthandler
func (s *Server) RemoveReplica(_ context.Context, request *client.RemoveReplicaRequest) (*client.StatusResponse, error) {
	log.Debug().Msgf("ExtServer has rcvd RemoveReplica. In: %#v", request)
	err := s.e.HandleRemoveReplica(fred.Keygroup{Name: fred.KeygroupName(request.Keygroup)}, fred.Node{ID: fred.NodeID(request.NodeId)})
	return statusResponseFromError(err)
}

func replicaResponseFromNode(n fred.Node) *client.GetReplicaResponse {
	return &client.GetReplicaResponse{NodeId: string(n.ID), Addr: n.Addr.Addr, Port: int32(n.Port)}
}

// GetReplica calls this method on the exthandler
func (s *Server) GetReplica(_ context.Context, request *client.GetReplicaRequest) (*client.GetReplicaResponse, error) {
	log.Debug().Msgf("ExtServer has rcvd GetReplica. In: %#v", request)
	res, err := s.e.HandleGetReplica(fred.Node{ID: fred.NodeID(request.NodeId)})
	return replicaResponseFromNode(res), err
}

// GetAllReplica calls this method on the exthandler
func (s *Server) GetAllReplica(_ context.Context, request *client.GetAllReplicaRequest) (*client.GetAllReplicaResponse, error) {
	log.Debug().Msgf("ExtServer has rcvd GetAllReplica. In: %#v", request)
	res, err := s.e.HandleGetAllReplica()
	if err != nil {
		log.Debug().Msgf("ExtServer is returning error: %#v", err)
		return &client.GetAllReplicaResponse{}, err
	}
	replicas := make([]*client.GetReplicaResponse, len(res))
	for i := 0; i < len(res); i++ {
		replicas[i] = replicaResponseFromNode(res[i])
	}
	return &client.GetAllReplicaResponse{Replicas: replicas}, nil
}

// GetKeygroupTriggers calls this method on the exthandler
func (s *Server) GetKeygroupTriggers(_ context.Context, request *client.GetKeygroupTriggerRequest) (*client.GetKeygroupTriggerResponse, error) {
	log.Debug().Msgf("ExtServer has rcvd GetKeygroupTriggers. In: %#v", request)
	res, err := s.e.HandleGetKeygroupTriggers(fred.Keygroup{Name: fred.KeygroupName(request.Keygroup)})
	if err != nil {
		log.Debug().Msgf("ExtServer is returning error: %#v", err)
		return &client.GetKeygroupTriggerResponse{}, err
	}
	triggers := make([]*client.Trigger, len(res))
	for i := 0; i < len(res); i++ {
		triggers[i] = &client.Trigger{
			Id:   res[i].ID,
			Host: res[i].Host,
		}
	}
	return &client.GetKeygroupTriggerResponse{Triggers: triggers}, nil
}

// AddTrigger calls this method on the exthandler
func (s *Server) AddTrigger(_ context.Context, request *client.AddTriggerRequest) (*client.StatusResponse, error) {
	log.Debug().Msgf("ExtServer has rcvd AddTrigger. In: %#v", request)
	err := s.e.HandleAddTrigger(fred.Keygroup{Name: fred.KeygroupName(request.Keygroup)}, fred.Trigger{ID: request.TriggerId, Host: request.TriggerHost})
	return statusResponseFromError(err)
}

// RemoveTrigger calls this method on the exthandler
func (s *Server) RemoveTrigger(_ context.Context, request *client.RemoveTriggerRequest) (*client.StatusResponse, error) {
	log.Debug().Msgf("ExtServer has rcvd RemoveTrigger. In: %#v", request)
	err := s.e.HandleRemoveTrigger(fred.Keygroup{Name: fred.KeygroupName(request.Keygroup)}, fred.Trigger{ID: request.TriggerId})
	return statusResponseFromError(err)
}
