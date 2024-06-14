package peering

import (
	"context"
	"crypto/tls"
	"io"
	"net"

	"git.tu-berlin.de/mcc-fred/fred/pkg/grpcutil"
	"git.tu-berlin.de/mcc-fred/fred/proto/peering"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"

	"git.tu-berlin.de/mcc-fred/fred/pkg/fred"
)

// Server is a grpc server that lets peers access the internal handler.
type Server struct {
	i    *fred.IntHandler
	host string
	*grpc.Server
}

// NewServer creates a new Server for communication to the inthandler from other nodes
func NewServer(host string, handler *fred.IntHandler, certFile string, keyFile string, caFile string, skipVerify bool) *Server {

	creds, _, err := grpcutil.GetCredsFromConfig(certFile, keyFile, []string{caFile}, false, skipVerify, &tls.Config{ClientAuth: tls.RequireAndVerifyClientCert})
	if err != nil {
		log.Fatal().Err(err).Msg("peering server: Cannot create TLS credentials")
		return nil
	}

	s := &Server{
		i:      handler,
		host:   host,
		Server: grpc.NewServer(grpc.Creds(creds)),
	}

	lis, err := net.Listen("tcp", host)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to listen")
		return nil
	}

	peering.RegisterNodeServer(s.Server, s)

	log.Trace().Msgf("Peering Server is listening on %s", host)

	go func() {
		err := s.Server.Serve(lis)

		// if Serve returns without an error, we probably intentionally closed it
		if err != nil {
			log.Fatal().Msgf("Peering Server exited: %s", err.Error())
		}
	}()

	return s
}

// Close closes the grpc server for internal communication.
func (s *Server) Close() error {
	s.Server.GracefulStop()
	return nil
}

// CreateKeygroup calls this Method on the Inthandler
func (s *Server) CreateKeygroup(_ context.Context, request *peering.CreateKeygroupRequest) (*peering.Empty, error) {
	log.Trace().Msgf("Peering server has rcvd CreateKeygroup. In: %+v", request)
	log.Debug().Msgf("Peering server has rcvd CreateKeygroup for keygroup %s", request.Keygroup)

	err := s.i.HandleCreateKeygroup(fred.Keygroup{Name: fred.KeygroupName(request.Keygroup)})
	if err != nil {
		return nil, err
	}
	return &peering.Empty{}, nil
}

// DeleteKeygroup calls this Method on the Inthandler
func (s *Server) DeleteKeygroup(_ context.Context, request *peering.DeleteKeygroupRequest) (*peering.Empty, error) {
	log.Trace().Msgf("Peering server has rcvd DeleteKeygroup. In: %+v", request)
	log.Debug().Msgf("Peering server has rcvd DeleteKeygroup for keygroup %s", request.Keygroup)

	err := s.i.HandleDeleteKeygroup(fred.Keygroup{Name: fred.KeygroupName(request.Keygroup)})
	if err != nil {
		return nil, err
	}
	return &peering.Empty{}, nil
}

// PutItem calls HandleUpdate on the Inthandler
func (s *Server) PutItem(ctx context.Context, request *peering.PutItemRequest) (*peering.Empty, error) {
	other := ""
	p, ok := peer.FromContext(ctx)
	if ok {
		other = p.Addr.String()
	}

	log.Trace().Msgf("Peering server has rcvd PutItem. In: %+v from %s to %s", request, other, s.host)
	log.Debug().Msgf("Peering server has rcvd PutItem for keygroup %s and id %s (append = %v)", request.Keygroup, request.Id, request.Append)

	if request.Append {
		err := s.i.HandleAppend(fred.Item{
			Keygroup: fred.KeygroupName(request.Keygroup),
			ID:       request.Id,
			Val:      request.Val,
		})
		if err != nil {
			return nil, err
		}
		return &peering.Empty{}, nil
	}

	err := s.i.HandleUpdate(fred.Item{
		Keygroup:   fred.KeygroupName(request.Keygroup),
		ID:         request.Id,
		Val:        request.Val,
		Version:    request.Version,
		Tombstoned: request.Tombstoned,
	})

	if err != nil {
		return nil, err
	}
	return &peering.Empty{}, nil
}

// DeleteItem calls this Method on the Inthandler
func (s *Server) StreamPut(p peering.Node_StreamPutServer) error {
	log.Trace().Msgf("Peering server has rcvd StreamPut. In: %+v", s)
	log.Debug().Msgf("Peering server has rcvd StreamPut")

	for item, err := p.Recv(); err != io.EOF; item, err = p.Recv() {
		if err != nil {
			return err
		}

		log.Trace().Msgf("Peering server has rcvd StreamPut for keygroup %s and id %s", item.Keygroup, item.Id)

		if item.Append {
			err = s.i.HandleAppend(fred.Item{
				Keygroup: fred.KeygroupName(item.Keygroup),
				ID:       item.Id,
				Val:      item.Val,
			})
			if err != nil {
				return err
			}
			continue
		}

		err = s.i.HandleUpdate(fred.Item{
			Keygroup:   fred.KeygroupName(item.Keygroup),
			ID:         item.Id,
			Val:        item.Val,
			Version:    item.Version,
			Tombstoned: item.Tombstoned,
		})

		if err != nil {
			return err
		}
	}

	return nil
}

// GetItem has no implementation
func (s *Server) GetItem(_ context.Context, request *peering.GetItemRequest) (*peering.ItemResponse, error) {
	log.Trace().Msgf("Peering server has rcvd GetItem. In: %+v", request)
	log.Debug().Msgf("Peering server has rcvd GetItem for keygroup %s and id %s", request.Keygroup, request.Id)

	items, err := s.i.HandleGet(fred.Item{
		Keygroup: fred.KeygroupName(request.Keygroup),
		ID:       request.Id,
	})

	if err != nil {
		return nil, err
	}

	data := make([]*peering.Data, len(items))

	for i, item := range items {
		data[i] = &peering.Data{
			Id:      item.ID,
			Val:     item.Val,
			Version: item.Version,
		}
	}

	return &peering.ItemResponse{
		Data: data,
	}, nil
}

// GetAllItems has no implementation
func (s *Server) GetAllItems(request *peering.GetAllItemsRequest, server peering.Node_GetAllItemsServer) error {
	log.Trace().Msgf("Peering server has rcvd GetItem. In: %+v", request)
	log.Debug().Msgf("Peering server has rcvd GetAllItems for keygroup %s", request.Keygroup)

	data, err := s.i.HandleGetAllItems(fred.Keygroup{
		Name: fred.KeygroupName(request.Keygroup),
	})

	if err != nil {
		return err
	}

	for _, item := range data {
		if err := server.Send(&peering.ItemResponse{
			Data: []*peering.Data{
				{
					Id:      item.ID,
					Val:     item.Val,
					Version: item.Version,
				},
			},
		}); err != nil {
			return err
		}
	}

	return nil
}
