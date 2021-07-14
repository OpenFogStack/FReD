package alexandra

import (
	"context"

	alexandraProto "git.tu-berlin.de/mcc-fred/fred/proto/middleware"
)

func (s *Server) Read(ctx context.Context, request *alexandraProto.ReadRequest) (*alexandraProto.ReadResponse, error) {
	return s.clientsMgr.GetClientTo(s.lighthouse).client.Read(ctx, request)
}

func (s *Server) Scan(ctx context.Context, request *alexandraProto.ScanRequest) (*alexandraProto.ScanResponse, error) {
	return s.clientsMgr.GetClientTo(s.lighthouse).client.Scan(ctx, request)
}

func (s *Server) Update(ctx context.Context, request *alexandraProto.UpdateRequest) (*alexandraProto.StatusResponse, error) {
	return s.clientsMgr.GetClientTo(s.lighthouse).client.Update(ctx, request)
}

func (s *Server) Delete(ctx context.Context, request *alexandraProto.DeleteRequest) (*alexandraProto.StatusResponse, error) {
	return s.clientsMgr.GetClientTo(s.lighthouse).client.Delete(ctx, request)
}

func (s *Server) Append(ctx context.Context, request *alexandraProto.AppendRequest) (*alexandraProto.AppendResponse, error) {
	return s.clientsMgr.GetClientTo(s.lighthouse).client.Append(ctx, request)
}
