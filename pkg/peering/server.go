package peering

import (
	"context"
	"net"

	"github.com/rs/zerolog/log"
	"gitlab.tu-berlin.de/mcc-fred/fred/proto/peering"
	"google.golang.org/grpc"

	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/fred"
)

// Server is a grpc server that let's peers access the internal handler.
type Server struct {
	i fred.IntHandler
	*grpc.Server
}

// NewServer creates a new Server for communication to the inthandler from other nodes
func NewServer(host string, handler fred.IntHandler) *Server {

	s := &Server{handler, grpc.NewServer()}

	lis, err := net.Listen("tcp", host)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to listen")
		return nil
	}

	peering.RegisterNodeServer(s.Server, s)

	log.Debug().Msgf("Interconnection Server is listening on %s", host)

	go func() {
		log.Fatal().Err(s.Server.Serve(lis)).Msg("Interconnection Server")
	}()

	return s
}

// Close closes the grpc server for internal communication.
func (s *Server) Close() error {
	s.Server.GracefulStop()
	return nil
}

func createResponse(err error) (*peering.StatusResponse, error) {
	if err == nil {
		return &peering.StatusResponse{Status: peering.EnumStatus_OK}, nil
	}
	// TODO Behave differently depending on error message
	// if error == leveldb.ErrNotFound {
	//
	// }
	return &peering.StatusResponse{Status: peering.EnumStatus_ERROR, ErrorMessage: err.Error()}, err
}

// CreateKeygroup calls this Method on the Inthandler
func (s *Server) CreateKeygroup(_ context.Context, request *peering.CreateKeygroupRequest) (*peering.StatusResponse, error) {
	log.Debug().Msgf("InterServer has rcvd CreateKeygroup. In: %#v", request)
	err := s.i.HandleCreateKeygroup(fred.Keygroup{Name: fred.KeygroupName(request.Keygroup)})
	return createResponse(err)
}

// DeleteKeygroup calls this Method on the Inthandler
func (s *Server) DeleteKeygroup(_ context.Context, request *peering.DeleteKeygroupRequest) (*peering.StatusResponse, error) {
	log.Debug().Msgf("InterServer has rcvd DeleteKeygroup. In: %#v", request)
	err := s.i.HandleDeleteKeygroup(fred.Keygroup{Name: fred.KeygroupName(request.Keygroup)})
	return createResponse(err)
}

// PutItem calls HandleUpdate on the Inthandler
func (s *Server) PutItem(_ context.Context, request *peering.PutItemRequest) (*peering.StatusResponse, error) {
	log.Debug().Msgf("InterServer has rcvd PutItem. In: %#v", request)
	err := s.i.HandleUpdate(fred.Item{
		Keygroup: fred.KeygroupName(request.Keygroup),
		ID:       request.Id,
		Val:      request.Data,
	})
	return createResponse(err)
}

// GetItem has no implementation
func (s *Server) GetItem(_ context.Context, _ *peering.GetItemRequest) (*peering.GetItemResponse, error) {
	panic("implement me")
}

// GetAllItems has no implementation
func (s *Server) GetAllItems(_ context.Context, _ *peering.GetAllItemsRequest) (*peering.GetAllItemsResponse, error) {
	panic("implement me")
}

// DeleteItem calls this Method on the Inthandler
func (s *Server) DeleteItem(_ context.Context, request *peering.DeleteItemRequest) (*peering.StatusResponse, error) {
	log.Debug().Msgf("InterServer has rcvd DeleteItem. In: %#v", request)
	err := s.i.HandleDelete(fred.Item{
		Keygroup: fred.KeygroupName(request.Keygroup),
		ID:       request.Id,
	})
	return createResponse(err)
}

// AddReplica calls this Method on the Inthandler
func (s *Server) AddReplica(_ context.Context, request *peering.AddReplicaRequest) (*peering.StatusResponse, error) {
	log.Debug().Msgf("InterServer has rcvd AddReplica. In: %#v", request)
	err := s.i.HandleAddReplica(fred.Keygroup{Name: fred.KeygroupName(request.Keygroup), Expiry: int(request.Expiry)}, fred.Node{ID: fred.NodeID(request.NodeId)})
	return createResponse(err)
}

// RemoveReplica calls this Method on the Inthandler
func (s *Server) RemoveReplica(_ context.Context, request *peering.RemoveReplicaRequest) (*peering.StatusResponse, error) {
	log.Debug().Msgf("InterServer has rcvd RemoveReplica. In: %#v", request)
	err := s.i.HandleRemoveReplica(fred.Keygroup{Name: fred.KeygroupName(request.Keygroup)}, fred.Node{ID: fred.NodeID(request.NodeId)})
	return createResponse(err)
}
