package peering

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"net"
	"os"

	"git.tu-berlin.de/mcc-fred/fred/proto/peering"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/peer"

	"git.tu-berlin.de/mcc-fred/fred/pkg/fred"
)

// Server is a grpc server that let's peers access the internal handler.
type Server struct {
	i    *fred.IntHandler
	host string
	*grpc.Server
}

// NewServer creates a new Server for communication to the inthandler from other nodes
func NewServer(host string, handler *fred.IntHandler, certFile string, keyFile string, caFile string) *Server {
	if certFile == "" {
		log.Fatal().Msg("peering server: no certificate file given")
	}

	if keyFile == "" {
		log.Fatal().Msg("peering server: no key file given")
	}

	if caFile == "" {
		log.Fatal().Msg("peering server: no root certificate file given")
	}

	// Load server's certificate and private key
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)

	if err != nil {
		log.Fatal().Msgf("peering server: could not load key pair: %v", err)
	}

	// Create a new cert pool and add our own CA certificate
	rootCAs := x509.NewCertPool()

	loaded, err := os.ReadFile(caFile)

	if err != nil {
		log.Fatal().Msgf("peering server: unexpected missing certfile: %v", err)
	}

	rootCAs.AppendCertsFromPEM(loaded)

	// Create the credentials and return it
	config := &tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    rootCAs,
		MinVersion:   tls.VersionTLS12,
	}

	s := &Server{
		i:      handler,
		host:   host,
		Server: grpc.NewServer(grpc.Creds(credentials.NewTLS(config))),
	}

	lis, err := net.Listen("tcp", host)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to listen")
		return nil
	}

	peering.RegisterNodeServer(s.Server, s)

	log.Debug().Msgf("Peering Server is listening on %s", host)

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
	log.Info().Msgf("Peering server has rcvd CreateKeygroup. In: %+v", request)
	err := s.i.HandleCreateKeygroup(fred.Keygroup{Name: fred.KeygroupName(request.Keygroup)})
	if err != nil {
		return nil, err
	}
	return &peering.Empty{}, nil
}

// DeleteKeygroup calls this Method on the Inthandler
func (s *Server) DeleteKeygroup(_ context.Context, request *peering.DeleteKeygroupRequest) (*peering.Empty, error) {
	log.Info().Msgf("Peering server has rcvd DeleteKeygroup. In: %+v", request)
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

	log.Info().Msgf("Peering server has rcvd PutItem. In: %+v from %s to %s", request, other, s.host)

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

// AppendItem calls HandleAppend on the Inthandler
func (s *Server) AppendItem(_ context.Context, request *peering.AppendItemRequest) (*peering.Empty, error) {
	log.Info().Msgf("Peering server has rcvd AppendItem. In: %+v", request)

	err := s.i.HandleAppend(fred.Item{
		Keygroup: fred.KeygroupName(request.Keygroup),
		ID:       request.Id,
		Val:      request.Data,
	})

	if err != nil {
		return nil, err
	}

	return &peering.Empty{}, nil
}

// GetItem has no implementation
func (s *Server) GetItem(_ context.Context, request *peering.GetItemRequest) (*peering.GetItemResponse, error) {
	log.Info().Msgf("Peering server has rcvd GetItem. In: %+v", request)
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

	return &peering.GetItemResponse{
		Data: data,
	}, nil
}

// GetAllItems has no implementation
func (s *Server) GetAllItems(_ context.Context, request *peering.GetAllItemsRequest) (*peering.GetAllItemsResponse, error) {
	log.Info().Msgf("Peering server has rcvd GetItem. In: %+v", request)
	data, err := s.i.HandleGetAllItems(fred.Keygroup{
		Name: fred.KeygroupName(request.Keygroup),
	})
	if err != nil {
		return nil, err
	}

	d := make([]*peering.Data, len(data))

	for i, item := range data {
		d[i] = &peering.Data{
			Id:      item.ID,
			Val:     item.Val,
			Version: item.Version,
		}
	}

	return &peering.GetAllItemsResponse{
		Data: d,
	}, nil
}
