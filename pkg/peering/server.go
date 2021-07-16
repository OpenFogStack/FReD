package peering

import (
	"context"
	"net"

	"git.tu-berlin.de/mcc-fred/fred/proto/peering"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"

	"git.tu-berlin.de/mcc-fred/fred/pkg/fred"
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
		err := s.Server.Serve(lis)

		// if Serve returns without an error, we probably intentionally closed it
		if err != nil {
			log.Fatal().Msgf("Interconnection Server exited: %s", err.Error())
		}
	}()

	return s
}

// Close closes the grpc server for internal communication.
func (s *Server) Close() error {
	s.Server.GracefulStop()
	return nil
}

// CreateKeygroup calls this Method on the Inthandler
func (s *Server) CreateKeygroup(_ context.Context, request *peering.CreateKeygroupRequest) (*peering.Empty, error) {
	log.Info().Msgf("InterServer has rcvd CreateKeygroup. In: %#v", request)
	err := s.i.HandleCreateKeygroup(fred.Keygroup{Name: fred.KeygroupName(request.Keygroup)})
	if err != nil {
		return nil, err
	}
	return &peering.Empty{}, nil
}

// DeleteKeygroup calls this Method on the Inthandler
func (s *Server) DeleteKeygroup(_ context.Context, request *peering.DeleteKeygroupRequest) (*peering.Empty, error) {
	log.Info().Msgf("InterServer has rcvd DeleteKeygroup. In: %#v", request)
	err := s.i.HandleDeleteKeygroup(fred.Keygroup{Name: fred.KeygroupName(request.Keygroup)})
	if err != nil {
		return nil, err
	}
	return &peering.Empty{}, nil
}

// PutItem calls HandleUpdate on the Inthandler
func (s *Server) PutItem(_ context.Context, request *peering.PutItemRequest) (*peering.Empty, error) {
	log.Info().Msgf("InterServer has rcvd PutItem. In: %#v", request)
	err := s.i.HandleUpdate(fred.Item{
		Keygroup: fred.KeygroupName(request.Keygroup),
		ID:       request.Id,
		Val:      request.Data,
	})
	if err != nil {
		return nil, err
	}
	return &peering.Empty{}, nil
}

// AppendItem calls HandleAppend on the Inthandler
func (s *Server) AppendItem(_ context.Context, request *peering.AppendItemRequest) (*peering.Empty, error) {
	log.Info().Msgf("InterServer has rcvd AppendItem. In: %#v", request)

	err := s.i.HandleAppend(fred.Item{
		Keygroup: fred.KeygroupName(request.Keygroup),
		ID:       request.Id,
		Val:      request.Data,
	})

	if err != nil {
		return nil, err
	}

	return &peering.Empty{}, nil
}

// GetItem has no implementation
func (s *Server) GetItem(_ context.Context, request *peering.GetItemRequest) (*peering.GetItemResponse, error) {
	log.Info().Msgf("InterServer has rcvd GetItem. In: %#v", request)
	data, err := s.i.HandleGet(fred.Item{
		Keygroup: fred.KeygroupName(request.Keygroup),
		ID:       request.Id,
	})

	if err != nil {
		return nil, err
	}
	return &peering.GetItemResponse{
		Data: data.Val,
	}, nil
}

// GetAllItems has no implementation
func (s *Server) GetAllItems(_ context.Context, request *peering.GetAllItemsRequest) (*peering.GetAllItemsResponse, error) {
	log.Info().Msgf("InterServer has rcvd GetItem. In: %#v", request)
	data, err := s.i.HandleGetAllItems(fred.Keygroup{
		Name: fred.KeygroupName(request.Keygroup),
	})
	if err != nil {
		return nil, err
	}

	d := make([]*peering.Data, len(data))

	for i, item := range data {
		d[i] = &peering.Data{
			Id:   item.ID,
			Data: item.Val,
		}
	}

	return &peering.GetAllItemsResponse{
		Data: d,
	}, nil
}

// DeleteItem calls this Method on the Inthandler
func (s *Server) DeleteItem(_ context.Context, request *peering.DeleteItemRequest) (*peering.Empty, error) {
	log.Info().Msgf("InterServer has rcvd DeleteItem. In: %#v", request)
	err := s.i.HandleDelete(fred.Item{
		Keygroup: fred.KeygroupName(request.Keygroup),
		ID:       request.Id,
	})
	if err != nil {
		return nil, err
	}
	return &peering.Empty{}, nil
}

// AddReplica calls this Method on the Inthandler
func (s *Server) AddReplica(_ context.Context, request *peering.AddReplicaRequest) (*peering.Empty, error) {
	log.Info().Msgf("InterServer has rcvd AddReplica. In: %#v", request)
	err := s.i.HandleAddReplica(fred.Keygroup{Name: fred.KeygroupName(request.Keygroup), Expiry: int(request.Expiry)}, fred.Node{ID: fred.NodeID(request.NodeId)})
	if err != nil {
		return nil, err
	}
	return &peering.Empty{}, nil
}

// RemoveReplica calls this Method on the Inthandler
func (s *Server) RemoveReplica(_ context.Context, request *peering.RemoveReplicaRequest) (*peering.Empty, error) {
	log.Info().Msgf("InterServer has rcvd RemoveReplica. In: %#v", request)
	err := s.i.HandleRemoveReplica(fred.Keygroup{Name: fred.KeygroupName(request.Keygroup)}, fred.Node{ID: fred.NodeID(request.NodeId)})
	if err != nil {
		return nil, err
	}
	return &peering.Empty{}, nil
}
