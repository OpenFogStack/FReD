package alexandra

import (
	"context"

	alexandraProto "git.tu-berlin.de/mcc-fred/fred/proto/middleware"
	"github.com/rs/zerolog/log"
)

// CreateKeygroup creates the keygroup and also adds the first node (This is two operations in the eye of FReD: CreateKeygroup and AddReplica)
func (s *Server) CreateKeygroup(ctx context.Context, request *alexandraProto.CreateKeygroupRequest) (*alexandraProto.Empty, error) {
	log.Debug().Msgf("AlexandraServer has rcdv CreateKeygroup: %#v", request)
	getReplica, err := s.clientsMgr.GetFastestClient().GetReplica(ctx, request.FirstNodeId)
	if err != nil {
		return nil, err
	}
	log.Debug().Msgf("CreateKeygroup: using node %s (addr=%s)", getReplica.NodeId, getReplica.Host)

	_, err = s.clientsMgr.GetClientTo(getReplica.Host).CreateKeygroup(ctx, request.Keygroup, request.Mutable, request.Expiry)

	if err != nil {
		return nil, err
	}

	return &alexandraProto.Empty{}, err
}

func (s *Server) DeleteKeygroup(ctx context.Context, request *alexandraProto.DeleteKeygroupRequest) (*alexandraProto.Empty, error) {
	client, err := s.clientsMgr.GetFastestClientWithKeygroup(request.Keygroup, 1)
	if err != nil {
		return nil, err
	}
	log.Debug().Msgf("DeleteKeygroup: using node %#v", client)

	_, err = client.DeleteKeygroup(ctx, request.Keygroup)

	if err != nil {
		return nil, err
	}

	return &alexandraProto.Empty{}, err
}
