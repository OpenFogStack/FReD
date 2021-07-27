package alexandra

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"net"

	alexandraProto "git.tu-berlin.de/mcc-fred/fred/proto/middleware"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// Server listens to GRPC requests from clients (and sends them to the relevant Fred Node etc.)
// The implementation is split up into different files in this folder.
type Server struct {
	roots      *x509.CertPool
	isProxied  bool
	proxyHost  string
	clientsMgr *ClientsMgr
	lighthouse string
	lis        net.Listener
	*grpc.Server
}

// NewServer creates a new Server for requests from Alexandra Clients
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
		NewClientsManager(nodesCert, nodesKey, lighthouse),
		lighthouse,
		lis,
		grpc.NewServer(
			grpc.Creds(credentials.NewTLS(config)),
		),
	}

	alexandraProto.RegisterMiddlewareServer(s.Server, s)

	log.Debug().Msgf("Alexandra Server is listening on %s", host)

	return s
}

func (s *Server) ServeBlocking() {
	log.Fatal().Err(s.Server.Serve(s.lis)).Msg("Alexandra Server")

}
