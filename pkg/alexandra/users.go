package alexandra

import (
	"context"

	fredClients "git.tu-berlin.de/mcc-fred/fred/proto/client"
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

func (s *Server) AddUser(ctx context.Context, request *alexandraProto.UserRequest) (*alexandraProto.Empty, error) {
	_, err := s.clientsMgr.GetClientTo(s.lighthouse).Client.AddUser(ctx, &fredClients.UserRequest{
		User:     request.User,
		Keygroup: request.Keygroup,
		Role:     fredClients.UserRole(request.Role),
	})

	if err != nil {
		return nil, err
	}

	return &alexandraProto.Empty{}, err
}

func (s *Server) RemoveUser(ctx context.Context, request *alexandraProto.UserRequest) (*alexandraProto.Empty, error) {
	_, err := s.clientsMgr.GetClientTo(s.lighthouse).Client.RemoveUser(ctx, &fredClients.UserRequest{
		User:     request.User,
		Keygroup: request.Keygroup,
		Role:     fredClients.UserRole(request.Role),
	})

	if err != nil {
		return nil, err
	}

	return &alexandraProto.Empty{}, err
}
