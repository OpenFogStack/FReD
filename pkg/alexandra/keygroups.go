package alexandra

import (
	"context"

	alexandraProto "git.tu-berlin.de/mcc-fred/fred/proto/middleware"
)

func (s *Server) CreateKeygroup(ctx context.Context, request *alexandraProto.CreateKeygroupRequest) (*alexandraProto.StatusResponse, error) {
	return s.clientsMgr.GetClientTo(s.lighthouse).client.CreateKeygroup(ctx, request)
}

func (s *Server) DeleteKeygroup(ctx context.Context, request *alexandraProto.DeleteKeygroupRequest) (*alexandraProto.StatusResponse, error) {
	return s.clientsMgr.GetClientTo(s.lighthouse).client.DeleteKeygroup(ctx, request)
}
