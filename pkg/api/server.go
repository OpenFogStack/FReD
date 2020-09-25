package api

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"net"
	"strings"
	"time"

	"github.com/go-errors/errors"
	"github.com/rs/zerolog/log"
	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/fred"
	"gitlab.tu-berlin.de/mcc-fred/fred/proto/client"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/peer"
)

// Server handles GRPC Requests and calls the according functions of the exthandler
type Server struct {
	e     fred.ExtHandler
	roots *x509.CertPool
	*grpc.Server
}

var (
	roles = map[client.UserRole]fred.Role{
		client.UserRole_ReadKeygroup:       fred.ReadKeygroup,
		client.UserRole_WriteKeygroup:      fred.WriteKeygroup,
		client.UserRole_ConfigureReplica:   fred.ConfigureReplica,
		client.UserRole_ConfigureTrigger:   fred.ConfigureTrigger,
		client.UserRole_ConfigureKeygroups: fred.ConfigureKeygroups,
	}
)

// checkCert checks the certificate from the given gRPC context for validity and returns the Common Name
func (s *Server) checkCert(ctx context.Context) (name string, err error) {
	// get peer information
	p, ok := peer.FromContext(ctx)

	if !ok {
		return name, errors.Errorf("no peer found")
	}

	// get TLS credential information for this connection
	tlsAuth, ok := p.AuthInfo.(credentials.TLSInfo)

	if !ok {
		return name, errors.Errorf("unexpected peer transport credentials")
	}

	// check that the certificate exists
	if len(tlsAuth.State.VerifiedChains) == 0 || len(tlsAuth.State.VerifiedChains[0]) == 0 {
		return name, errors.Errorf("could not verify peer certificate: %v", tlsAuth.State)
	}

	// verify the certificate:
	// 1) it should be issued by a CA in our root CA pool
	// 2) any intermediates are valid for us
	// 3) the certificate should be valid for client authentication
	// 4) the certificate should have the clients address as a SAN
	_, err = tlsAuth.State.VerifiedChains[0][0].Verify(x509.VerifyOptions{
		Roots:         s.roots,
		CurrentTime:   time.Now(),
		Intermediates: x509.NewCertPool(),
		KeyUsages:     []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
		// p.Addr is of form "IP:PORT" so split it by ":" char and use the first entry
		DNSName: strings.Split(p.Addr.String(), ":")[0],
	})

	if err != nil {
		return name, errors.New(err)
	}

	// Check subject common name exists and return it as the user name for that client
	if name = tlsAuth.State.VerifiedChains[0][0].Subject.CommonName; name == "" {
		return name, errors.Errorf("invalid subject common name")
	}

	log.Debug().Msgf("CN: %s", name)

	return name, nil
}

// NewServer creates a new Server for requests from Fred Clients
func NewServer(host string, handler fred.ExtHandler, cert string, key string, caCert string) *Server {
	// Load server's certificate and private key
	serverCert, err := tls.LoadX509KeyPair(cert, key)

	if err != nil {
		log.Fatal().Msgf("could not load key pair: %v", err)
		return nil
	}

	// Create a new cert pool and add our own CA certificate
	rootCAs := x509.NewCertPool()

	loaded, err := ioutil.ReadFile(caCert)

	if err != nil {
		log.Fatal().Msgf("unexpected missing certfile: %v", err)
	}

	rootCAs.AppendCertsFromPEM(loaded)
	// Create the credentials and return it
	config := &tls.Config{
		Certificates: []tls.Certificate{serverCert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    rootCAs,
		MinVersion:   tls.VersionTLS12,
	}

	s := &Server{
		handler,
		rootCAs,
		grpc.NewServer(
			grpc.Creds(credentials.NewTLS(config)),
		),
	}

	lis, err := net.Listen("tcp", host)

	if err != nil {
		log.Fatal().Err(err).Msg("Failed to listen")
		return nil
	}

	client.RegisterClientServer(s.Server, s)

	log.Debug().Msgf("Externalconnection Server is listening on %s", host)

	go func() {
		log.Fatal().Err(s.Server.Serve(lis)).Msg("Externalconnection Server")
	}()

	return s
}

// Close closes the grpc server for internal communication.
func (s *Server) Close() error {
	s.Server.GracefulStop()
	return nil
}

func statusResponseFromError(err error) (*client.StatusResponse, error) {
	if err == nil {
		return &client.StatusResponse{Status: client.EnumStatus_OK}, nil
	}

	log.Debug().Msgf("ExtServer is returning error: %#v", err)

	return &client.StatusResponse{Status: client.EnumStatus_ERROR, ErrorMessage: err.Error()}, err

}

// CreateKeygroup calls this method on the exthandler
func (s *Server) CreateKeygroup(ctx context.Context, request *client.CreateKeygroupRequest) (*client.StatusResponse, error) {

	log.Debug().Msgf("ExtServer has rcvd CreateKeygroup. In: %#v", request)

	user, err := s.checkCert(ctx)

	if err != nil {
		return statusResponseFromError(err)
	}

	err = s.e.HandleCreateKeygroup(user, fred.Keygroup{Name: fred.KeygroupName(request.Keygroup), Mutable: request.Mutable, Expiry: int(request.Expiry)})

	return statusResponseFromError(err)
}

// DeleteKeygroup calls this method on the exthandler
func (s *Server) DeleteKeygroup(ctx context.Context, request *client.DeleteKeygroupRequest) (*client.StatusResponse, error) {

	log.Debug().Msgf("ExtServer has rcvd DeleteKeygroup. In: %#v", request)

	user, err := s.checkCert(ctx)

	if err != nil {
		return statusResponseFromError(err)
	}

	err = s.e.HandleDeleteKeygroup(user, fred.Keygroup{Name: fred.KeygroupName(request.Keygroup)})

	return statusResponseFromError(err)
}

// Read calls this method on the exthandler
func (s *Server) Read(ctx context.Context, request *client.ReadRequest) (*client.ReadResponse, error) {
	log.Debug().Msgf("ExtServer has rcvd Read. In: %#v", request)

	user, err := s.checkCert(ctx)

	if err != nil {
		_, err = statusResponseFromError(err)
		return nil, err
	}

	res, err := s.e.HandleRead(user, fred.Item{Keygroup: fred.KeygroupName(request.Keygroup), ID: request.Id})

	if err != nil {
		log.Debug().Msgf("ExtServer is returning error: %#v", err)
		return &client.ReadResponse{}, err

	}
	return &client.ReadResponse{Data: res.Val}, nil

}

// Update calls this method on the exthandler
func (s *Server) Update(ctx context.Context, request *client.UpdateRequest) (*client.StatusResponse, error) {

	log.Debug().Msgf("ExtServer has rcvd DeleteKeygroup. In: %#v", request)

	user, err := s.checkCert(ctx)

	if err != nil {
		_, err = statusResponseFromError(err)
		return nil, err
	}

	err = s.e.HandleUpdate(user, fred.Item{Keygroup: fred.KeygroupName(request.Keygroup), ID: request.Id, Val: request.Data})

	return statusResponseFromError(err)
}

// Delete calls this method on the exthandler
func (s *Server) Delete(ctx context.Context, request *client.DeleteRequest) (*client.StatusResponse, error) {
	log.Debug().Msgf("ExtServer has rcvd Delete. In: %#v", request)

	user, err := s.checkCert(ctx)

	if err != nil {
		_, err = statusResponseFromError(err)
		return nil, err
	}

	err = s.e.HandleDelete(user, fred.Item{Keygroup: fred.KeygroupName(request.Keygroup), ID: request.Id})

	return statusResponseFromError(err)
}

// AddReplica calls this method on the exthandler
func (s *Server) AddReplica(ctx context.Context, request *client.AddReplicaRequest) (*client.StatusResponse, error) {
	log.Debug().Msgf("ExtServer has rcvd AddReplica. In: %#v", request)

	user, err := s.checkCert(ctx)

	if err != nil {
		_, err = statusResponseFromError(err)
		return nil, err
	}

	err = s.e.HandleAddReplica(user, fred.Keygroup{Name: fred.KeygroupName(request.Keygroup), Expiry: int(request.Expiry)}, fred.Node{ID: fred.NodeID(request.NodeId)})

	return statusResponseFromError(err)
}

// GetKeygroupReplica calls this method on the exthandler
func (s *Server) GetKeygroupReplica(ctx context.Context, request *client.GetKeygroupReplicaRequest) (*client.GetKeygroupReplicaResponse, error) {
	log.Debug().Msgf("ExtServer has rcvd GetKeygroupReplica. In: %#v", request)

	user, err := s.checkCert(ctx)

	if err != nil {
		_, err = statusResponseFromError(err)
		return nil, err
	}

	n, e, err := s.e.HandleGetKeygroupReplica(user, fred.Keygroup{Name: fred.KeygroupName(request.Keygroup)})

	// Copy only the interesting values into a new array
	replicas := make([]*client.KeygroupReplica, len(n))

	for i := 0; i < len(n); i++ {
		replicas[i] = &client.KeygroupReplica{
			NodeId: string(n[i].ID),
			Expiry: int64(e[n[i].ID]),
		}
	}

	if err != nil {
		log.Debug().Msgf("ExtServer is returning error: %#v", err)
		return &client.GetKeygroupReplicaResponse{}, err
	}

	return &client.GetKeygroupReplicaResponse{
		Replica: replicas,
	}, nil

}

// RemoveReplica calls this method on the exthandler
func (s *Server) RemoveReplica(ctx context.Context, request *client.RemoveReplicaRequest) (*client.StatusResponse, error) {
	log.Debug().Msgf("ExtServer has rcvd RemoveReplica. In: %#v", request)

	user, err := s.checkCert(ctx)

	if err != nil {
		_, err = statusResponseFromError(err)
		return nil, err
	}

	err = s.e.HandleRemoveReplica(user, fred.Keygroup{Name: fred.KeygroupName(request.Keygroup)}, fred.Node{ID: fred.NodeID(request.NodeId)})

	return statusResponseFromError(err)
}

func replicaResponseFromNode(n fred.Node) *client.GetReplicaResponse {
	return &client.GetReplicaResponse{NodeId: string(n.ID), Addr: n.Addr.Addr, Port: int32(n.Port)}
}

// GetReplica calls this method on the exthandler
func (s *Server) GetReplica(ctx context.Context, request *client.GetReplicaRequest) (*client.GetReplicaResponse, error) {
	log.Debug().Msgf("ExtServer has rcvd GetReplica. In: %#v", request)

	user, err := s.checkCert(ctx)

	if err != nil {
		_, err = statusResponseFromError(err)
		return nil, err
	}

	res, err := s.e.HandleGetReplica(user, fred.Node{ID: fred.NodeID(request.NodeId)})

	return replicaResponseFromNode(res), err
}

// GetAllReplica calls this method on the exthandler
func (s *Server) GetAllReplica(ctx context.Context, request *client.GetAllReplicaRequest) (*client.GetAllReplicaResponse, error) {
	log.Debug().Msgf("ExtServer has rcvd GetAllReplica. In: %#v", request)

	user, err := s.checkCert(ctx)

	if err != nil {
		_, err = statusResponseFromError(err)
		return nil, err
	}

	res, err := s.e.HandleGetAllReplica(user)

	if err != nil {
		log.Debug().Msgf("ExtServer is returning error: %#v", err)
		return &client.GetAllReplicaResponse{}, err
	}

	replicas := make([]*client.GetReplicaResponse, len(res))

	for i := 0; i < len(res); i++ {
		replicas[i] = replicaResponseFromNode(res[i])
	}

	return &client.GetAllReplicaResponse{Replicas: replicas}, nil
}

// GetKeygroupTriggers calls this method on the exthandler
func (s *Server) GetKeygroupTriggers(ctx context.Context, request *client.GetKeygroupTriggerRequest) (*client.GetKeygroupTriggerResponse, error) {
	log.Debug().Msgf("ExtServer has rcvd GetKeygroupTriggers. In: %#v", request)

	user, err := s.checkCert(ctx)

	if err != nil {
		_, err = statusResponseFromError(err)
		return nil, err
	}

	res, err := s.e.HandleGetKeygroupTriggers(user, fred.Keygroup{Name: fred.KeygroupName(request.Keygroup)})
	if err != nil {
		log.Debug().Msgf("ExtServer is returning error: %#v", err)
		return &client.GetKeygroupTriggerResponse{}, err
	}
	triggers := make([]*client.Trigger, len(res))
	for i := 0; i < len(res); i++ {
		triggers[i] = &client.Trigger{
			Id:   res[i].ID,
			Host: res[i].Host,
		}
	}
	return &client.GetKeygroupTriggerResponse{Triggers: triggers}, nil
}

// AddTrigger calls this method on the exthandler
func (s *Server) AddTrigger(ctx context.Context, request *client.AddTriggerRequest) (*client.StatusResponse, error) {
	log.Debug().Msgf("ExtServer has rcvd AddTrigger. In: %#v", request)

	user, err := s.checkCert(ctx)

	if err != nil {
		_, err = statusResponseFromError(err)
		return nil, err
	}

	err = s.e.HandleAddTrigger(user, fred.Keygroup{Name: fred.KeygroupName(request.Keygroup)}, fred.Trigger{ID: request.TriggerId, Host: request.TriggerHost})

	return statusResponseFromError(err)
}

// RemoveTrigger calls this method on the exthandler
func (s *Server) RemoveTrigger(ctx context.Context, request *client.RemoveTriggerRequest) (*client.StatusResponse, error) {
	log.Debug().Msgf("ExtServer has rcvd RemoveTrigger. In: %#v", request)

	user, err := s.checkCert(ctx)

	if err != nil {
		_, err = statusResponseFromError(err)
		return nil, err
	}

	err = s.e.HandleRemoveTrigger(user, fred.Keygroup{Name: fred.KeygroupName(request.Keygroup)}, fred.Trigger{ID: request.TriggerId})

	return statusResponseFromError(err)
}

// AddUser calls this method on the exthandler
func (s *Server) AddUser(ctx context.Context, request *client.UserRequest) (*client.StatusResponse, error) {
	log.Debug().Msgf("ExtServer has rcvd AddUser. In: %#v", request)

	user, err := s.checkCert(ctx)

	if err != nil {
		_, err = statusResponseFromError(err)
		return nil, err
	}

	err = s.e.HandleAddUser(user, request.User, fred.Keygroup{Name: fred.KeygroupName(request.Keygroup)}, roles[request.Role])

	return statusResponseFromError(err)
}

// RemoveUser calls this method on the exthandler
func (s *Server) RemoveUser(ctx context.Context, request *client.UserRequest) (*client.StatusResponse, error) {
	log.Debug().Msgf("ExtServer has rcvd RemoveUser. In: %#v", request)

	user, err := s.checkCert(ctx)

	if err != nil {
		_, err = statusResponseFromError(err)
		return nil, err
	}

	err = s.e.HandleRemoveUser(user, request.User, fred.Keygroup{Name: fred.KeygroupName(request.Keygroup)}, roles[request.Role])

	return statusResponseFromError(err)
}
