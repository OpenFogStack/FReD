package api

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"net"
	"strconv"
	"time"

	"git.tu-berlin.de/mcc-fred/fred/pkg/fred"
	"git.tu-berlin.de/mcc-fred/fred/pkg/grpcutil"
	"git.tu-berlin.de/mcc-fred/fred/proto/client"
	"git.tu-berlin.de/mcc-fred/vclock"
	"github.com/go-errors/errors"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

// Server handles GRPC Requests and calls the according functions of the exthandler
type Server struct {
	e          *fred.ExtHandler
	roots      *x509.CertPool
	isProxied  bool
	skipVerify bool
	proxyHost  string
	proxyPort  string
	*grpc.Server
}

// Roles map the internal grpc representation of rbac Roles to the representation within fred
var (
	Roles = map[client.UserRole]fred.Role{
		client.UserRole_ReadKeygroup:       fred.ReadKeygroup,
		client.UserRole_WriteKeygroup:      fred.WriteKeygroup,
		client.UserRole_ConfigureReplica:   fred.ConfigureReplica,
		client.UserRole_ConfigureTrigger:   fred.ConfigureTrigger,
		client.UserRole_ConfigureKeygroups: fred.ConfigureKeygroups,
	}
)

// CheckCert checks the certificate from the given gRPC context for validity and returns the Common Name
func (s *Server) CheckCert(ctx context.Context) (string, error) {
	// get peer information
	p, ok := peer.FromContext(ctx)

	if !ok {
		return "", errors.Errorf("no peer found")
	}

	// get TLS credential information for this connection
	tlsAuth, ok := p.AuthInfo.(credentials.TLSInfo)

	if !ok {
		return "", errors.Errorf("unexpected peer transport credentials")
	}

	// check that the certificate exists
	if len(tlsAuth.State.VerifiedChains) == 0 || len(tlsAuth.State.VerifiedChains[0]) == 0 {
		return "", errors.Errorf("could not verify peer certificate: %v", tlsAuth.State)
	}

	host, _, err := net.SplitHostPort(p.Addr.String())

	if err != nil {
		return "", errors.New(err)
	}

	// set host to empty, so we don't check it in the certificate
	if s.skipVerify {
		host = ""
	}

	// verify the certificate:
	// IF we are not proxied and communicate with the client directly:
	// 1) it should be issued by a CA in our root CA pool
	// 2) any intermediates are valid for us
	// 3) the certificate should be valid for client authentication
	// 4) the certificate should have the clients address as a SAN
	if !s.isProxied {
		_, err = tlsAuth.State.VerifiedChains[0][0].Verify(x509.VerifyOptions{
			Roots:         s.roots,
			CurrentTime:   time.Now(),
			Intermediates: x509.NewCertPool(),
			KeyUsages:     []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
			DNSName:       host,
		})

		if err != nil {
			return "", errors.New(err)
		}

		// Check subject common name exists and return it as the user name for that client
		name := tlsAuth.State.VerifiedChains[0][0].Subject.CommonName

		if name == "" {
			return "", errors.Errorf("invalid subject common name")
		}

		log.Trace().Msgf("CheckCert: GRPC Context Certificate Name: %s", name)

		return name, nil
	}
	// ELSE we sit behind a proxy and the proxy should be the one tunneling the gRPC connection to us
	// hence if we can ensure that it is indeed the proxy that is talking to us and not someone who has found
	// their way into the network, we can be sure that the proxy/LB has checked the certificate
	// in this case, the proxy will give the user name to us as a header (thanks, proxy!)

	if !s.skipVerify && host != s.proxyHost {
		return "", errors.Errorf("node is proxied but got request not from proxy (%s instead of %s)", host, s.proxyHost)
	}

	_, err = tlsAuth.State.VerifiedChains[0][0].Verify(x509.VerifyOptions{
		Roots:         s.roots,
		CurrentTime:   time.Now(),
		Intermediates: x509.NewCertPool(),
		KeyUsages:     []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
	})

	if err != nil {
		return "", errors.New(err)
	}

	md, ok := metadata.FromIncomingContext(ctx)

	if !ok {
		return "", errors.Errorf("no metadata could be found for proxied request")
	}

	u := md.Get("user")

	if len(u) != 1 {
		return "", errors.Errorf("invalid user header for proxied request")
	}

	return u[0], nil
}

// NewServer creates a new Server for requests from Fred Clients
func NewServer(host string, handler *fred.ExtHandler, certFile string, keyFile string, caFile string, skipVerify bool, isProxied bool, proxy string) *Server {
	creds, rootCAs, err := grpcutil.GetCredsFromConfig(certFile, keyFile, []string{caFile}, false, skipVerify, &tls.Config{ClientAuth: tls.RequireAndVerifyClientCert})

	if err != nil {
		log.Fatal().Err(err).Msg("API server: Failed to generate credentials")
		return nil
	}

	proxyHost, proxyPort, err := net.SplitHostPort(proxy)

	if isProxied && err != nil {
		log.Fatal().Err(err).Msg("API server: Failed to parse proxy host and port")
		return nil
	}

	s := &Server{
		e:          handler,
		roots:      rootCAs,
		isProxied:  isProxied,
		proxyHost:  proxyHost,
		proxyPort:  proxyPort,
		skipVerify: skipVerify,
		Server: grpc.NewServer(
			grpc.Creds(creds),
		),
	}

	lis, err := net.Listen("tcp", host)

	if err != nil {
		log.Fatal().Err(err).Msg("API server: Failed to listen")
		return nil
	}

	client.RegisterClientServer(s.Server, s)

	log.Info().Msgf("API Server is listening on %s", host)

	go func() {
		err := s.Server.Serve(lis)

		// if Serve returns without an error, we probably intentionally closed it
		if err != nil {
			log.Fatal().Msgf("API Server exited: %s", err.Error())
		}
	}()

	return s
}

// Close closes the grpc server for internal communication.
func (s *Server) Close() error {
	s.Server.GracefulStop()
	return nil
}

// CreateKeygroup calls this method on the exthandler
func (s *Server) CreateKeygroup(ctx context.Context, request *client.CreateKeygroupRequest) (*client.Empty, error) {

	log.Trace().Msgf("API Server has rcvd CreateKeygroup. In: %+v", request)
	log.Debug().Msgf("Create Keygroup keygroup: %s, mutable: %t, expiry: %d", request.Keygroup, request.Mutable, request.Expiry)

	user, err := s.CheckCert(ctx)

	if err != nil {
		return nil, err
	}

	log.Trace().Msgf("API Server has rcvd CreateKeygroup. User: %s", user)

	err = s.e.HandleCreateKeygroup(user, fred.Keygroup{Name: fred.KeygroupName(request.Keygroup), Mutable: request.Mutable, Expiry: int(request.Expiry)})

	if err != nil {
		return nil, err
	}

	return &client.Empty{}, nil
}

// DeleteKeygroup calls this method on the exthandler
func (s *Server) DeleteKeygroup(ctx context.Context, request *client.DeleteKeygroupRequest) (*client.Empty, error) {

	log.Trace().Msgf("API Server has rcvd DeleteKeygroup. In: %+v", request)
	log.Debug().Msgf("Delete Keygroup keygroup: %s", request.Keygroup)

	user, err := s.CheckCert(ctx)

	if err != nil {
		return nil, err
	}

	log.Trace().Msgf("API Server has rcvd DeleteKeygroup. User: %s", user)

	err = s.e.HandleDeleteKeygroup(user, fred.Keygroup{Name: fred.KeygroupName(request.Keygroup)})

	if err != nil {
		return nil, err
	}

	return &client.Empty{}, nil
}

// Read calls this method on the exthandler
func (s *Server) Read(ctx context.Context, request *client.ReadRequest) (*client.ReadResponse, error) {
	log.Trace().Msgf("API Server has rcvd Read. In: %+v", request)
	log.Debug().Msgf("Read: keygroup: %s, id: %s, versions: %+v", request.Keygroup, request.Id, request.Versions)

	user, err := s.CheckCert(ctx)

	if err != nil {
		return nil, err
	}

	log.Trace().Msgf("API Server has rcvd Read. User: %s", user)

	versions := make([]vclock.VClock, len(request.Versions))
	for i, v := range request.Versions {
		versions[i] = v.Version
	}

	res, err := s.e.HandleRead(user, fred.Item{Keygroup: fred.KeygroupName(request.Keygroup), ID: request.Id}, versions)

	if err != nil {
		log.Debug().Msgf("API Server is returning error: %+v", err)
		return nil, err

	}

	data := make([]*client.Item, len(res))

	for i := 0; i < len(res); i++ {
		data[i] = &client.Item{
			Id:  res[i].ID,
			Val: res[i].Val,
			Version: &client.Version{
				Version: res[i].Version.GetMap(),
			},
		}
	}

	return &client.ReadResponse{
		Data: data,
	}, nil

}

// Scan calls this method on the exthandler
func (s *Server) Scan(ctx context.Context, request *client.ScanRequest) (*client.ScanResponse, error) {
	log.Trace().Msgf("API Server has rcvd Scan. In: %+v", request)
	log.Debug().Msgf("Scan: keygroup: %s, id: %s, count: %d", request.Keygroup, request.Id, request.Count)

	user, err := s.CheckCert(ctx)

	if err != nil {
		return nil, err
	}

	log.Trace().Msgf("API Server has rcvd Scan. User: %s", user)

	res, err := s.e.HandleScan(user, fred.Item{Keygroup: fred.KeygroupName(request.Keygroup), ID: request.Id}, request.Count)

	if err != nil {
		log.Debug().Msgf("API Server is returning error: %+v", err)
		return nil, err

	}

	data := make([]*client.Item, len(res))

	for i := 0; i < len(res); i++ {
		data[i] = &client.Item{
			Id:  res[i].ID,
			Val: res[i].Val,
			Version: &client.Version{
				Version: res[i].Version.GetMap(),
			},
		}
	}

	return &client.ScanResponse{
		Data: data,
	}, nil

}

// Keys calls this method on the exthandler
func (s *Server) Keys(ctx context.Context, request *client.KeysRequest) (*client.KeysResponse, error) {
	log.Trace().Msgf("API Server has rcvd Keys. In: %+v", request)
	log.Debug().Msgf("Keys: keygroup: %s, id: %s, count: %d", request.Keygroup, request.Id, request.Count)

	user, err := s.CheckCert(ctx)

	if err != nil {
		return nil, err
	}

	log.Trace().Msgf("API Server has rcdv Keys. User: %s", user)

	res, err := s.e.HandleKeys(user, fred.Item{Keygroup: fred.KeygroupName(request.Keygroup), ID: request.Id}, request.Count)

	if err != nil {
		log.Debug().Msgf("API Server is returning error: %+v", err)
		return nil, err

	}

	keys := make([]*client.Key, len(res))

	for i := 0; i < len(res); i++ {
		keys[i] = &client.Key{
			Id: res[i].ID,
			Version: &client.Version{
				Version: res[i].Version.GetMap(),
			},
		}
	}

	return &client.KeysResponse{
		Keys: keys,
	}, nil

}

// Append calls this method on the exthandler
func (s *Server) Append(ctx context.Context, request *client.AppendRequest) (*client.AppendResponse, error) {
	log.Trace().Msgf("API Server has rcvd Append. In: %+v", request)
	log.Debug().Msgf("Append: keygroup: %s, id: %d", request.Keygroup, request.Id)

	user, err := s.CheckCert(ctx)

	if err != nil {
		return nil, err
	}

	log.Trace().Msgf("API Server has rcvd Append. User: %s", user)

	res, err := s.e.HandleAppend(user, fred.Item{Keygroup: fred.KeygroupName(request.Keygroup), ID: strconv.FormatUint(request.Id, 10), Val: request.Data})

	if err != nil {
		return nil, err
	}

	return &client.AppendResponse{
		Id: res.ID,
	}, nil
}

// Update calls this method on the exthandler
func (s *Server) Update(ctx context.Context, request *client.UpdateRequest) (*client.UpdateResponse, error) {

	log.Trace().Msgf("API Server has rcvd Update. In: %+v", request)
	log.Debug().Msgf("Update: keygroup: %s, id: %s, versions: %+v", request.Keygroup, request.Id, request.Versions)

	user, err := s.CheckCert(ctx)

	if err != nil {
		return nil, err
	}

	log.Trace().Msgf("API Server has rcvd Update. User: %s", user)

	var versions []vclock.VClock
	for _, v := range request.Versions {
		versions = append(versions, v.Version)
	}

	i, err := s.e.HandleUpdate(user, fred.Item{Keygroup: fred.KeygroupName(request.Keygroup), ID: request.Id, Val: request.Data}, versions)

	if err != nil {
		return nil, err
	}

	return &client.UpdateResponse{
		Version: &client.Version{Version: i.Version.GetMap()},
	}, nil
}

// Delete calls this method on the exthandler
func (s *Server) Delete(ctx context.Context, request *client.DeleteRequest) (*client.DeleteResponse, error) {
	log.Trace().Msgf("API Server has rcvd Delete. In: %+v", request)
	log.Debug().Msgf("Delete: keygroup: %s, id: %s, versions: %+v", request.Keygroup, request.Id, request.Versions)

	user, err := s.CheckCert(ctx)

	if err != nil {
		return nil, err
	}

	log.Trace().Msgf("API Server has rcvd Delete. User: %s", user)

	var versions []vclock.VClock
	for _, v := range request.Versions {
		versions = append(versions, v.Version)
	}

	i, err := s.e.HandleDelete(user, fred.Item{Keygroup: fred.KeygroupName(request.Keygroup), ID: request.Id}, versions)

	if err != nil {
		return nil, err
	}

	return &client.DeleteResponse{
		Version: &client.Version{Version: i.Version.GetMap()},
	}, nil
}

// AddReplica calls this method on the exthandler
func (s *Server) AddReplica(ctx context.Context, request *client.AddReplicaRequest) (*client.Empty, error) {
	log.Trace().Msgf("API Server has rcvd AddReplica. In: %+v", request)
	log.Debug().Msgf("AddReplica: keygroup: %s, node: %s, expiry: %d", request.Keygroup, request.NodeId, request.Expiry)

	user, err := s.CheckCert(ctx)

	if err != nil {
		return nil, err
	}

	log.Trace().Msgf("API Server has rcvd AddReplica. User: %s", user)

	err = s.e.HandleAddReplica(user, fred.Keygroup{Name: fred.KeygroupName(request.Keygroup), Expiry: int(request.Expiry)}, fred.Node{ID: fred.NodeID(request.NodeId)})

	if err != nil {
		return nil, err
	}

	return &client.Empty{}, nil
}

// GetKeygroupInfo calls this method on the exthandler
func (s *Server) GetKeygroupInfo(ctx context.Context, request *client.GetKeygroupInfoRequest) (*client.GetKeygroupInfoResponse, error) {
	log.Trace().Msgf("API Server has rcvd GetKeygroupReplica. In: %+v", request)
	log.Debug().Msgf("GetKeygroupInfo: keygroup: %s", request.Keygroup)

	user, err := s.CheckCert(ctx)

	if err != nil {
		log.Debug().Msgf("API Server is returning error: %+v", err)
		return nil, err
	}

	log.Trace().Msgf("API Server has rcvd GetKeygroupInfo. User: %s", user)

	m, n, e, err := s.e.HandleGetKeygroupInfo(user, fred.Keygroup{Name: fred.KeygroupName(request.Keygroup)})

	log.Trace().Msgf("... received replicas: %+v", n)
	if err != nil {
		return nil, err
	}

	// Copy only the interesting values into a new array
	replicas := make([]*client.KeygroupReplica, len(n))

	for i := 0; i < len(n); i++ {
		replicas[i] = &client.KeygroupReplica{
			NodeId: string(n[i].ID),
			Expiry: int64(e[n[i].ID]),
			Host:   n[i].Host,
		}
	}

	return &client.GetKeygroupInfoResponse{
		Mutable: m,
		Replica: replicas,
	}, nil

}

// RemoveReplica calls this method on the exthandler
func (s *Server) RemoveReplica(ctx context.Context, request *client.RemoveReplicaRequest) (*client.Empty, error) {
	log.Trace().Msgf("API Server has rcvd RemoveReplica. In: %+v", request)
	log.Debug().Msgf("RemoveReplica: keygroup: %s, node: %s", request.Keygroup, request.NodeId)

	user, err := s.CheckCert(ctx)

	if err != nil {
		return nil, err
	}

	log.Trace().Msgf("API Server has rcvd RemoveReplica. User: %s", user)

	err = s.e.HandleRemoveReplica(user, fred.Keygroup{Name: fred.KeygroupName(request.Keygroup)}, fred.Node{ID: fred.NodeID(request.NodeId)})

	if err != nil {
		return nil, err
	}

	return &client.Empty{}, nil
}

// GetReplica calls this method on the exthandler
func (s *Server) GetReplica(ctx context.Context, request *client.GetReplicaRequest) (*client.GetReplicaResponse, error) {
	log.Trace().Msgf("API Server has rcvd GetReplica. In: %+v", request)
	log.Debug().Msgf("GetReplica: node: %s", request.NodeId)

	user, err := s.CheckCert(ctx)

	if err != nil {
		return nil, err
	}

	log.Trace().Msgf("API Server has rcvd GetReplica. User: %s", user)

	res, err := s.e.HandleGetReplica(user, fred.Node{ID: fred.NodeID(request.NodeId)})

	if err != nil {
		return nil, err
	}

	return &client.GetReplicaResponse{
		NodeId: string(res.ID), Host: res.Host,
	}, nil
}

// GetAllReplica calls this method on the exthandler
func (s *Server) GetAllReplica(ctx context.Context, request *client.Empty) (*client.GetAllReplicaResponse, error) {
	log.Trace().Msgf("API Server has rcvd GetAllReplica. In: %+v", request)
	log.Debug().Msgf("GetAllReplica")

	user, err := s.CheckCert(ctx)

	if err != nil {
		return nil, err
	}

	log.Trace().Msgf("API Server has rcvd GetAllReplica. User: %s", user)

	res, err := s.e.HandleGetAllReplica(user)

	if err != nil {
		log.Debug().Msgf("API Server is returning error: %+v", err)
		return nil, err
	}

	replicas := make([]*client.Replica, len(res))

	for i := 0; i < len(res); i++ {
		replicas[i] = &client.Replica{
			NodeId: string(res[i].ID),
			Host:   res[i].Host,
		}
	}

	return &client.GetAllReplicaResponse{Replicas: replicas}, nil
}

// GetKeygroupTriggers calls this method on the exthandler
func (s *Server) GetKeygroupTriggers(ctx context.Context, request *client.GetKeygroupTriggerRequest) (*client.GetKeygroupTriggerResponse, error) {
	log.Trace().Msgf("API Server has rcvd GetKeygroupTriggers. In: %+v", request)
	log.Debug().Msgf("GetKeygroupTriggers: keygroup: %s", request.Keygroup)

	user, err := s.CheckCert(ctx)

	if err != nil {
		return nil, err
	}

	log.Trace().Msgf("API Server has rcvd GetKeygroupTriggers. User: %s", user)

	res, err := s.e.HandleGetKeygroupTriggers(user, fred.Keygroup{Name: fred.KeygroupName(request.Keygroup)})

	if err != nil {
		log.Debug().Msgf("API Server is returning error: %+v", err)
		return nil, err
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
func (s *Server) AddTrigger(ctx context.Context, request *client.AddTriggerRequest) (*client.Empty, error) {
	log.Trace().Msgf("API Server has rcvd AddTrigger. In: %+v", request)
	log.Debug().Msgf("AddTrigger: keygroup: %s, triggerId: %s, triggerHost: %s", request.Keygroup, request.TriggerId, request.TriggerHost)

	user, err := s.CheckCert(ctx)

	if err != nil {
		return nil, err
	}

	log.Trace().Msgf("API Server has rcvd AddTrigger. User: %s", user)

	err = s.e.HandleAddTrigger(user, fred.Keygroup{Name: fred.KeygroupName(request.Keygroup)}, fred.Trigger{ID: request.TriggerId, Host: request.TriggerHost})

	if err != nil {
		return nil, err
	}

	return &client.Empty{}, nil
}

// RemoveTrigger calls this method on the exthandler
func (s *Server) RemoveTrigger(ctx context.Context, request *client.RemoveTriggerRequest) (*client.Empty, error) {
	log.Trace().Msgf("API Server has rcvd RemoveTrigger. In: %+v", request)
	log.Debug().Msgf("RemoveTrigger: keygroup: %s, triggerId: %s", request.Keygroup, request.TriggerId)

	user, err := s.CheckCert(ctx)

	if err != nil {
		return nil, err
	}

	log.Trace().Msgf("API Server has rcvd RemoveTrigger. User: %s", user)

	err = s.e.HandleRemoveTrigger(user, fred.Keygroup{Name: fred.KeygroupName(request.Keygroup)}, fred.Trigger{ID: request.TriggerId})

	if err != nil {
		return nil, err
	}

	return &client.Empty{}, nil
}

// AddUser calls this method on the exthandler
func (s *Server) AddUser(ctx context.Context, request *client.AddUserRequest) (*client.Empty, error) {
	log.Trace().Msgf("API Server has rcvd AddUser. In: %+v", request)
	log.Debug().Msgf("AddUser: user: %s, keygroup: %s, role: %s", request.User, request.Keygroup, request.Role)

	user, err := s.CheckCert(ctx)

	if err != nil {
		return nil, err
	}

	log.Trace().Msgf("API Server has rcvd AddUser. User: %s", user)

	err = s.e.HandleAddUser(user, request.User, fred.Keygroup{Name: fred.KeygroupName(request.Keygroup)}, Roles[request.Role])

	if err != nil {
		return nil, err
	}

	return &client.Empty{}, nil
}

// RemoveUser calls this method on the exthandler
func (s *Server) RemoveUser(ctx context.Context, request *client.RemoveUserRequest) (*client.Empty, error) {
	log.Trace().Msgf("API Server has rcvd RemoveUser. In: %+v", request)
	log.Debug().Msgf("RemoveUser: user: %s, keygroup: %s, role: %s", request.User, request.Keygroup, request.Role)

	user, err := s.CheckCert(ctx)

	if err != nil {
		return nil, err
	}

	log.Trace().Msgf("API Server has rcvd RemoveUser. User: %s", user)

	err = s.e.HandleRemoveUser(user, request.User, fred.Keygroup{Name: fred.KeygroupName(request.Keygroup)}, Roles[request.Role])

	if err != nil {
		return nil, err
	}

	return &client.Empty{}, nil
}
