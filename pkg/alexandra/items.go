package alexandra

import (
	"context"
	"math/rand"

	"git.tu-berlin.de/mcc-fred/fred/proto/client"
	alexandraProto "git.tu-berlin.de/mcc-fred/fred/proto/middleware"
	"github.com/rs/zerolog/log"
)

func (s *Server) Scan(ctx context.Context, request *alexandraProto.ScanRequest) (*alexandraProto.ScanResponse, error) {
	res, err := s.clientsMgr.GetClientTo(s.lighthouse).Client.Scan(ctx, &client.ScanRequest{
		Keygroup: request.Keygroup,
		Id:       request.Id,
		Count:    request.Count,
	})

	if err != nil {
		return nil, err
	}

	data := make([]*alexandraProto.Data, len(res.Data))

	for i, datum := range res.Data {
		data[i] = &alexandraProto.Data{
			Id:   datum.Id,
			Data: datum.Val,
		}
	}
	return &alexandraProto.ScanResponse{Data: data}, err
}

func (s *Server) Read(ctx context.Context, request *alexandraProto.ReadRequest) (*alexandraProto.ReadResponse, error) {
	log.Debug().Msgf("Alexandra has rcvd Read")

	return s.clientsMgr.ReadFromAnywhere(ctx, request)
}

func (s *Server) Update(ctx context.Context, request *alexandraProto.UpdateRequest) (*alexandraProto.Empty, error) {
	var c *Client
	var err error
	if rand.Float32() > UseSlowerNodeProb {
		c, err = s.clientsMgr.GetFastestClientWithKeygroup(request.Keygroup, 1)
	} else {
		c, err = s.clientsMgr.GetRandomClientWithKeygroup(request.Keygroup, 1)
	}

	if err != nil {
		return nil, err
	}

	_, err = c.Update(ctx, request.Keygroup, request.Id, request.Data)
	// Convert from FredClient response to AlexandraResponse

	if err != nil {
		return nil, err
	}

	return &alexandraProto.Empty{}, err
}

func (s *Server) Delete(ctx context.Context, request *alexandraProto.DeleteRequest) (*alexandraProto.Empty, error) {
	var c *Client
	var err error
	if rand.Float32() > UseSlowerNodeProb {
		c, err = s.clientsMgr.GetFastestClientWithKeygroup(request.Keygroup, 1)
	} else {
		c, err = s.clientsMgr.GetRandomClientWithKeygroup(request.Keygroup, 1)
	}

	if err != nil {
		return nil, err
	}

	_, err = c.Delete(ctx, request.Keygroup, request.Id)

	if err != nil {
		return nil, err
	}

	return &alexandraProto.Empty{}, err
}

func (s *Server) Append(ctx context.Context, request *alexandraProto.AppendRequest) (*alexandraProto.AppendResponse, error) {
	var c *Client
	var err error
	if rand.Float32() > UseSlowerNodeProb {
		c, err = s.clientsMgr.GetFastestClientWithKeygroup(request.Keygroup, 1)
	} else {
		c, err = s.clientsMgr.GetRandomClientWithKeygroup(request.Keygroup, 1)
	}

	if err != nil {
		return nil, err
	}
	res, err := c.Append(ctx, request.Keygroup, request.Data)

	if err != nil {
		return nil, err
	}

	return &alexandraProto.AppendResponse{Id: res.Id}, err
}
