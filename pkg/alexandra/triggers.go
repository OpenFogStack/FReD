package alexandra

import (
	"context"

	alexandraProto "git.tu-berlin.de/mcc-fred/fred/proto/middleware"
)

func (s *Server) GetKeygroupTriggers(ctx context.Context, request *alexandraProto.GetKeygroupTriggerRequest) (*alexandraProto.GetKeygroupTriggerResponse, error) {
	return s.clientsMgr.GetClientTo(s.lighthouse).client.GetKeygroupTriggers(ctx, request)
}

func (s *Server) AddTrigger(ctx context.Context, request *alexandraProto.AddTriggerRequest) (*alexandraProto.StatusResponse, error) {
	return s.clientsMgr.GetClientTo(s.lighthouse).client.AddTrigger(ctx, request)
}

func (s *Server) RemoveTrigger(ctx context.Context, request *alexandraProto.RemoveTriggerRequest) (*alexandraProto.StatusResponse, error) {
	return s.clientsMgr.GetClientTo(s.lighthouse).client.RemoveTrigger(ctx, request)
}

