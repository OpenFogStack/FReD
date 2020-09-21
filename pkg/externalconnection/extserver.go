package externalconnection

import (
	"context"
	"github.com/rs/zerolog/log"
	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/fred"
	"google.golang.org/grpc"
	"net"
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

	RegisterClientServer(s.Server, s)

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

func statusResponseFromError(err error) (*StatusResponse, error) {
	if err == nil {
		return &StatusResponse{Status: EnumStatus_OK}, nil
	}
	log.Debug().Msgf("ExtServer is returning error: %#v", err)
	return &StatusResponse{Status: EnumStatus_ERROR, ErrorMessage: err.Error()}, err

}

// CreateKeygroup calls this method on the exthandler
func (s *Server) CreateKeygroup(ctx context.Context, request *CreateKeygroupRequest) (*StatusResponse, error) {
	log.Debug().Msgf("ExtServer has rcvd CreateKeygroup. In: %#v", request)
	err := s.e.HandleCreateKeygroup(fred.Keygroup{Name: fred.KeygroupName(request.Keygroup)})
	return statusResponseFromError(err)
}

// DeleteKeygroup calls this method on the exthandler
func (s *Server) DeleteKeygroup(ctx context.Context, request *DeleteKeygroupRequest) (*StatusResponse, error) {
	log.Debug().Msgf("ExtServer has rcvd DeleteKeygroup. In: %#v", request)
	err := s.e.HandleDeleteKeygroup(fred.Keygroup{Name: fred.KeygroupName(request.Keygroup)})
	return statusResponseFromError(err)
}

// Read calls this method on the exthandler
func (s *Server) Read(ctx context.Context, request *ReadRequest) (*ReadResponse, error) {
	log.Debug().Msgf("ExtServer has rcvd Read. In: %#v", request)
	res, err := s.e.HandleRead(fred.Item{Keygroup: fred.KeygroupName(request.Keygroup), ID: request.Id})
	if err != nil {
		log.Debug().Msgf("ExtServer is returning error: %#v", err)
		return &ReadResponse{}, err
	}
	return &ReadResponse{Data: res.Val}, nil

}

// Update calls this method on the exthandler
func (s *Server) Update(ctx context.Context, request *UpdateRequest) (*StatusResponse, error) {
	log.Debug().Msgf("ExtServer has rcvd DeleteKeygroup. In: %#v", request)
	err := s.e.HandleUpdate(fred.Item{Keygroup: fred.KeygroupName(request.Keygroup), ID: request.Id, Val: request.Data})
	return statusResponseFromError(err)
}

// Delete calls this method on the exthandler
func (s *Server) Delete(ctx context.Context, request *DeleteRequest) (*StatusResponse, error) {
	log.Debug().Msgf("ExtServer has rcvd Delete. In: %#v", request)
	err := s.e.HandleDelete(fred.Item{Keygroup: fred.KeygroupName(request.Keygroup), ID: request.Id})
	return statusResponseFromError(err)
}

// AddReplica calls this method on the exthandler
func (s *Server) AddReplica(ctx context.Context, request *AddReplicaRequest) (*StatusResponse, error) {
	log.Debug().Msgf("ExtServer has rcvd AddReplica. In: %#v", request)
	err := s.e.HandleAddReplica(fred.Keygroup{Name: fred.KeygroupName(request.Keygroup)}, fred.Node{ID: fred.NodeID(request.NodeId)})
	return statusResponseFromError(err)
}

// GetKeygroupReplica calls this method on the exthandler
func (s *Server) GetKeygroupReplica(ctx context.Context, request *GetKeygroupReplicaRequest) (*GetKeygroupReplicaResponse, error) {
	log.Debug().Msgf("ExtServer has rcvd GetKeygroupReplica. In: %#v", request)
	res, err := s.e.HandleGetKeygroupReplica(fred.Keygroup{Name: fred.KeygroupName(request.Keygroup)})
	// Copy only the interesting values into a new array
	nodes := make([]string, len(res))
	for i := 0; i < len(res); i++ {
		nodes[i] = string(res[i].ID)
	}
	if err != nil {
		log.Debug().Msgf("ExtServer is returning error: %#v", err)
		return &GetKeygroupReplicaResponse{}, err
	}
	return &GetKeygroupReplicaResponse{NodeId: nodes}, nil

}

// RemoveReplica calls this method on the exthandler
func (s *Server) RemoveReplica(ctx context.Context, request *RemoveReplicaRequest) (*StatusResponse, error) {
	log.Debug().Msgf("ExtServer has rcvd RemoveReplica. In: %#v", request)
	err := s.e.HandleRemoveReplica(fred.Keygroup{Name: fred.KeygroupName(request.Keygroup)}, fred.Node{ID: fred.NodeID(request.NodeId)})
	return statusResponseFromError(err)
}

func replicaResponseFromNode(n fred.Node) *GetReplicaResponse {
	return &GetReplicaResponse{NodeId: string(n.ID), Addr: n.Addr.Addr, Port: int32(n.Port)}
}

// GetReplica calls this method on the exthandler
func (s *Server) GetReplica(ctx context.Context, request *GetReplicaRequest) (*GetReplicaResponse, error) {
	log.Debug().Msgf("ExtServer has rcvd GetReplica. In: %#v", request)
	res, err := s.e.HandleGetReplica(fred.Node{ID: fred.NodeID(request.NodeId)})
	return replicaResponseFromNode(res), err
}

// GetAllReplica calls this method on the exthandler
func (s *Server) GetAllReplica(ctx context.Context, request *GetAllReplicaRequest) (*GetAllReplicaResponse, error) {
	log.Debug().Msgf("ExtServer has rcvd GetAllReplica. In: %#v", request)
	res, err := s.e.HandleGetAllReplica()
	if err != nil {
		log.Debug().Msgf("ExtServer is returning error: %#v", err)
		return &GetAllReplicaResponse{}, err
	}
	replicas := make([]*GetReplicaResponse, len(res))
	for i := 0; i < len(res); i++ {
		replicas[i] = replicaResponseFromNode(res[i])
	}
	return &GetAllReplicaResponse{Replicas: replicas}, nil
}

// GetKeygroupTriggers calls this method on the exthandler
func (s *Server) GetKeygroupTriggers(ctx context.Context, request *GetKeygroupTriggerRequest) (*GetKeygroupTriggerResponse, error) {
	log.Debug().Msgf("ExtServer has rcvd GetKeygroupTriggers. In: %#v", request)
	res, err := s.e.HandleGetKeygroupTriggers(fred.Keygroup{Name: fred.KeygroupName(request.Keygroup)})
	if err != nil {
		log.Debug().Msgf("ExtServer is returning error: %#v", err)
		return &GetKeygroupTriggerResponse{}, err
	}
	triggers := make([]*Trigger, len(res))
	for i := 0; i < len(res); i++ {
		triggers[i] = &Trigger{
			Id:   res[i].ID,
			Host: res[i].Host,
		}
	}
	return &GetKeygroupTriggerResponse{Triggers: triggers}, nil
}

// AddTrigger calls this method on the exthandler
func (s *Server) AddTrigger(ctx context.Context, request *AddTriggerRequest) (*StatusResponse, error) {
	log.Debug().Msgf("ExtServer has rcvd AddTrigger. In: %#v", request)
	err := s.e.HandleAddTrigger(fred.Keygroup{Name: fred.KeygroupName(request.Keygroup)}, fred.Trigger{ID: request.TriggerId, Host: request.TriggerHost})
	return statusResponseFromError(err)
}

// RemoveTrigger calls this method on the exthandler
func (s *Server) RemoveTrigger(ctx context.Context, request *RemoveTriggerRequest) (*StatusResponse, error) {
	log.Debug().Msgf("ExtServer has rcvd RemoveTrigger. In: %#v", request)
	err := s.e.HandleRemoveTrigger(fred.Keygroup{Name: fred.KeygroupName(request.Keygroup)}, fred.Trigger{ID: request.TriggerId})
	return statusResponseFromError(err)
}
