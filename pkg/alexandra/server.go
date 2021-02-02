package alexandra

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"net"
	"time"

	alexandraProto "git.tu-berlin.de/mcc-fred/fred/proto/middleware"
	"github.com/go-errors/errors"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/peer"
)

// Server listens to GRPC requests from clients (and sends them to the relevant Fred Node etc.)
// The implementation is split up into different files in this folder.
type Server struct {
	roots     *x509.CertPool
	isProxied bool
	proxyHost string
	clientsMgr *ClientsMgr
	lighthouse string
	lis        net.Listener
	*grpc.Server
}

// NewServer creates a new Server for requests from Fred Clients
func NewServer(host string, caCert string, serverCert string, serverKey string, nodesCert string, nodesKey string, lighthouse string, isProxied bool, proxyHost string) *Server {
	// Load server's certificate and private key
	loadedServerCert, err := tls.LoadX509KeyPair(serverCert, serverKey)

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
		Certificates: []tls.Certificate{loadedServerCert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    rootCAs,
		MinVersion:   tls.VersionTLS12,
	}

	lis, err := net.Listen("tcp", host)

	if err != nil {
		log.Fatal().Err(err).Msg("Failed to listen")
		return nil
	}

	s := &Server{
		rootCAs,
		isProxied,
		proxyHost,
		newClientsManager(nodesCert, nodesKey),
		lighthouse,
		lis,
		grpc.NewServer(
			grpc.Creds(credentials.NewTLS(config)),
		),
	}

	alexandraProto.RegisterClientServer(s.Server, s)

	log.Debug().Msgf("Alexandra Server is listening on %s", host)

	return s
}

func (s *Server) ServeBlocking() {
	log.Fatal().Err(s.Server.Serve(s.lis)).Msg("Alexandra Server")

}

// CheckCert checks the certificate from the given gRPC context for validity and returns the Common Name
func (s *Server) CheckCert(ctx context.Context) (name string, err error) {
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

	host, _, err := net.SplitHostPort(p.Addr.String())

	if err != nil {
		return name, errors.New(err)
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
			return name, errors.New(err)
		}
	} else {
		// ELSE we sit behind a proxy and the proxy should be the one tunneling the gRPC connection to us
		// hence if we can ensure that it is indeed the proxy that is talking to us and not someone who has found
		// their way into the network, we can be sure that the proxy/LB has checked the certificate
		_, err = tlsAuth.State.VerifiedChains[0][0].Verify(x509.VerifyOptions{
			Roots:         s.roots,
			CurrentTime:   time.Now(),
			Intermediates: x509.NewCertPool(),
			KeyUsages:     []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
		})

		if err != nil {
			return name, errors.New(err)
		}
	}

	// Check subject common name exists and return it as the user name for that client
	if name = tlsAuth.State.VerifiedChains[0][0].Subject.CommonName; name == "" {
		return name, errors.Errorf("invalid subject common name")
	}

	log.Debug().Msgf("CheckCert: GRPC Context Certificate Name: %s", name)

	return name, nil
}
