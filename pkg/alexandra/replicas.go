package alexandra

import (
	"context"

	fredClients "git.tu-berlin.de/mcc-fred/fred/proto/client"
	alexandraProto "git.tu-berlin.de/mcc-fred/fred/proto/middleware"
)

func (s *Server) AddReplica(ctx context.Context, request *alexandraProto.AddReplicaRequest) (*alexandraProto.Empty, error) {
	_, err := s.clientsMgr.GetFastestClient().Client.AddReplica(ctx, &fredClients.AddReplicaRequest{
		Keygroup: request.Keygroup,
		NodeId:   request.NodeId,
		Expiry:   request.Expiry,
	})

	if err != nil {
		return nil, err
	}

	return &alexandraProto.Empty{}, err
}

func (s *Server) RemoveReplica(ctx context.Context, request *alexandraProto.RemoveReplicaRequest) (*alexandraProto.Empty, error) {
	_, err := s.clientsMgr.GetFastestClient().Client.RemoveReplica(ctx, &fredClients.RemoveReplicaRequest{
		Keygroup: request.Keygroup,
		NodeId:   request.NodeId,
	})
	if err != nil {
		return nil, err
	}

	return &alexandraProto.Empty{}, err
}

func (s *Server) GetReplica(ctx context.Context, request *alexandraProto.GetReplicaRequest) (*alexandraProto.GetReplicaResponse, error) {
	res, err := s.clientsMgr.GetFastestClient().Client.GetReplica(ctx, &fredClients.GetReplicaRequest{NodeId: request.NodeId})

	if err != nil {
		return nil, err
	}

	return &alexandraProto.GetReplicaResponse{NodeId: res.NodeId, Host: res.Host}, err
}

func (s *Server) GetAllReplica(ctx context.Context, request *alexandraProto.GetAllReplicaRequest) (*alexandraProto.GetAllReplicaResponse, error) {
	res, err := s.clientsMgr.GetFastestClient().Client.GetAllReplica(ctx, &fredClients.Empty{})

	if err != nil {
		return nil, err
	}

	replicas := make([]*alexandraProto.GetReplicaResponse, len(res.Replicas))
	for i, replica := range res.Replicas {
		replicas[i] = &alexandraProto.GetReplicaResponse{
			NodeId: replica.NodeId,
			Host:   replica.Host,
		}
	}

	return &alexandraProto.GetAllReplicaResponse{Replicas: replicas}, err
}

func (s *Server) GetKeygroupReplica(ctx context.Context, request *alexandraProto.GetKeygroupReplicaRequest) (*alexandraProto.GetKeygroupReplicaResponse, error) {
	res, err := s.clientsMgr.GetFastestClient().Client.GetKeygroupReplica(ctx, &fredClients.GetKeygroupReplicaRequest{Keygroup: request.Keygroup})

	if err != nil {
		return nil, err
	}

	replicas := make([]*alexandraProto.KeygroupReplica, len(res.Replica))
	for i, replica := range res.Replica {
		replicas[i] = &alexandraProto.KeygroupReplica{
			NodeId: replica.NodeId,
			Host:   replica.Host,
		}
	}

	return &alexandraProto.GetKeygroupReplicaResponse{Replica: replicas}, err
}
