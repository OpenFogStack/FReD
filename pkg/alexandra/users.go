package alexandra

import (
	"context"

	alexandraProto "git.tu-berlin.de/mcc-fred/fred/proto/middleware"
)

// the roles map the internal grpc representation of rbac roles to the representation within fred
//var (
//	roles = map[alexandraProto.UserRole]fred.Role{
//		alexandraProto.UserRole_ReadKeygroup:       fred.ReadKeygroup,
//		alexandraProto.UserRole_WriteKeygroup:      fred.WriteKeygroup,
//		alexandraProto.UserRole_ConfigureReplica:   fred.ConfigureReplica,
//		alexandraProto.UserRole_ConfigureTrigger:   fred.ConfigureTrigger,
//		alexandraProto.UserRole_ConfigureKeygroups: fred.ConfigureKeygroups,
//	}
//)

func (s *Server) AddUser(ctx context.Context, request *alexandraProto.UserRequest) (*alexandraProto.StatusResponse, error) {
	return s.clientsMgr.GetClientTo(s.lighthouse).client.AddUser(ctx, request)
}

func (s *Server) RemoveUser(ctx context.Context, request *alexandraProto.UserRequest) (*alexandraProto.StatusResponse, error) {
	return s.clientsMgr.GetClientTo(s.lighthouse).client.RemoveUser(ctx, request)
}
