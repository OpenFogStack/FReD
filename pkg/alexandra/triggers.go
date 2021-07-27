package alexandra

import (
	"context"

	fredClients "git.tu-berlin.de/mcc-fred/fred/proto/client"
	alexandraProto "git.tu-berlin.de/mcc-fred/fred/proto/middleware"
)

func (s *Server) GetKeygroupTriggers(ctx context.Context, request *alexandraProto.GetKeygroupTriggerRequest) (*alexandraProto.GetKeygroupTriggerResponse, error) {
	res, err := s.clientsMgr.GetClientTo(s.lighthouse).Client.GetKeygroupTriggers(ctx, &fredClients.GetKeygroupTriggerRequest{
		Keygroup: request.Keygroup,
	})

	if err != nil {
		return nil, err
	}

	triggers := make([]*alexandraProto.Trigger, len(res.Triggers))
	for i, trigger := range res.Triggers {
		triggers[i] = &alexandraProto.Trigger{
			Id:   trigger.Id,
			Host: trigger.Host,
		}
	}
	return &alexandraProto.GetKeygroupTriggerResponse{Triggers: triggers}, nil
}

func (s *Server) AddTrigger(ctx context.Context, request *alexandraProto.AddTriggerRequest) (*alexandraProto.Empty, error) {
	_, err := s.clientsMgr.GetClientTo(s.lighthouse).Client.AddTrigger(ctx, &fredClients.AddTriggerRequest{
		Keygroup:    request.Keygroup,
		TriggerId:   request.TriggerId,
		TriggerHost: request.TriggerHost,
	})

	if err != nil {
		return nil, err
	}

	return &alexandraProto.Empty{}, err
}

func (s *Server) RemoveTrigger(ctx context.Context, request *alexandraProto.RemoveTriggerRequest) (*alexandraProto.Empty, error) {
	_, err := s.clientsMgr.GetClientTo(s.lighthouse).Client.RemoveTrigger(ctx, &fredClients.RemoveTriggerRequest{
		Keygroup:  request.Keygroup,
		TriggerId: request.TriggerId,
	})

	if err != nil {
		return nil, err
	}

	return &alexandraProto.Empty{}, err
}
