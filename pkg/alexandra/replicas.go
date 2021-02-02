package alexandra

import (
	"context"

	alexandraProto "git.tu-berlin.de/mcc-fred/fred/proto/middleware"
)

func (s *Server) AddReplica(ctx context.Context, request *alexandraProto.AddReplicaRequest) (*alexandraProto.StatusResponse, error) {
	return s.clientsMgr.GetClientTo(s.lighthouse).client.AddReplica(ctx, request)
}

func (s *Server) RemoveReplica(ctx context.Context, request *alexandraProto.RemoveReplicaRequest) (*alexandraProto.StatusResponse, error) {
	return s.clientsMgr.GetClientTo(s.lighthouse).client.RemoveReplica(ctx, request)
}

func (s *Server) GetReplica(ctx context.Context, request *alexandraProto.GetReplicaRequest) (*alexandraProto.GetReplicaResponse, error) {
	return s.clientsMgr.GetClientTo(s.lighthouse).client.GetReplica(ctx, request)
}

func (s *Server) GetAllReplica(ctx context.Context, request *alexandraProto.GetAllReplicaRequest) (*alexandraProto.GetAllReplicaResponse, error) {
	return s.clientsMgr.GetClientTo(s.lighthouse).client.GetAllReplica(ctx, request)
}

func (s *Server) GetKeygroupReplica(ctx context.Context, request *alexandraProto.GetKeygroupReplicaRequest) (*alexandraProto.GetKeygroupReplicaResponse, error) {
	return s.clientsMgr.GetClientTo(s.lighthouse).client.GetKeygroupReplica(ctx, request)
}


