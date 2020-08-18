package interconnection

import (
	"context"

	"github.com/rs/zerolog/log"

	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/fred"
)

// Server is the server that calls the inthandler
type Server struct {
	intHandler fred.IntHandler
}

// NewServer creates a new Server for communication to the inthandler from other nodes
func NewServer(handler *fred.IntHandler) *Server {
	return &Server{intHandler: *handler}
}

func createResponse(err error) (*StatusResponse, error) {
	if err == nil {
		return &StatusResponse{Status: EnumStatus_OK}, nil
	}
	// TODO Behave differently depending on error message
	// if error == leveldb.ErrNotFound {
	//
	// }
	return &StatusResponse{Status: EnumStatus_ERROR, ErrorMessage: err.Error()}, err
}

// CreateKeygroup calls this Method on the Inthandler
func (s Server) CreateKeygroup(ctx context.Context, request *CreateKeygroupRequest) (*StatusResponse, error) {
	log.Debug().Msgf("InterServer has rcvd CreateKeygroup. In: %#v", request)
	err := s.intHandler.HandleCreateKeygroup(fred.Keygroup{Name: fred.KeygroupName(request.Keygroup)}, nil)
	return createResponse(err)
}

// DeleteKeygroup calls this Method on the Inthandler
func (s Server) DeleteKeygroup(ctx context.Context, request *DeleteKeygroupRequest) (*StatusResponse, error) {
	log.Debug().Msgf("InterServer has rcvd DeleteKeygroup. In: %#v", request)
	err := s.intHandler.HandleDeleteKeygroup(fred.Keygroup{Name: fred.KeygroupName(request.Keygroup)})
	return createResponse(err)
}

// PutItem calls HandleUpdate on the Inthandler
func (s Server) PutItem(ctx context.Context, request *PutItemRequest) (*StatusResponse, error) {
	log.Debug().Msgf("InterServer has rcvd PutItem. In: %#v", request)
	err := s.intHandler.HandleUpdate(fred.Item{
		Keygroup: fred.KeygroupName(request.Keygroup),
		ID:       request.Id,
		Val:      request.Data,
	})
	return createResponse(err)
}

// GetItem has no implementation
func (s Server) GetItem(ctx context.Context, request *GetItemRequest) (*GetItemResponse, error) {
	panic("implement me")
}

// GetAllItems has no implementation
func (s Server) GetAllItems(ctx context.Context, request *GetAllItemsRequest) (*GetAllItemsResponse, error) {
	panic("implement me")
}

// DeleteItem calls this Method on the Inthandler
func (s Server) DeleteItem(ctx context.Context, request *DeleteItemRequest) (*StatusResponse, error) {
	log.Debug().Msgf("InterServer has rcvd DeleteItem. In: %#v", request)
	err := s.intHandler.HandleDelete(fred.Item{
		Keygroup: fred.KeygroupName(request.Keygroup),
		ID:       request.Id,
	})
	return createResponse(err)
}

// AddReplica calls this Method on the Inthandler
func (s Server) AddReplica(ctx context.Context, request *AddReplicaRequest) (*StatusResponse, error) {
	log.Debug().Msgf("InterServer has rcvd AddReplica. In: %#v", request)
	err := s.intHandler.HandleAddReplica(fred.Keygroup{Name: fred.KeygroupName(request.Keygroup)}, fred.Node{ID: fred.NodeID(request.NodeId)})
	return createResponse(err)
}

// RemoveReplica calls this Method on the Inthandler
func (s Server) RemoveReplica(ctx context.Context, request *RemoveReplicaRequest) (*StatusResponse, error) {
	log.Debug().Msgf("InterServer has rcvd RemoveReplica. In: %#v", request)
	err := s.intHandler.HandleRemoveReplica(fred.Keygroup{Name: fred.KeygroupName(request.Keygroup)}, fred.Node{ID: fred.NodeID(request.NodeId)})
	return createResponse(err)
}
