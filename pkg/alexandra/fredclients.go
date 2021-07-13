package alexandra

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"

	alexandraProto "git.tu-berlin.de/mcc-fred/fred/proto/middleware"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type ClientsMgr struct {
	clients                            map[string]*Client
	clientsCert, clientsKey, clientsCA string
}

type Client struct {
	client alexandraProto.MiddlewareClient
	conn   *grpc.ClientConn
}

func newClientsManager(clientsCert, clientsKey, clientsCA string) *ClientsMgr {
	return &ClientsMgr{
		clients:     make(map[string]*Client),
		clientsCert: clientsCert,
		clientsKey:  clientsKey,
		clientsCA:   clientsCA,
	}
}

func (m *ClientsMgr) GetClientTo(host string) (client *Client) {
	client = m.clients[host]
	if client != nil {
		return
	}
	client = newClient(host, m.clientsCert, m.clientsKey, m.clientsCA)
	m.clients[host] = client
	return
}

func newClient(host, certFile, keyFile, caFile string) *Client {
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)

	if err != nil {
		log.Fatal().Err(err).Str("certFile", certFile).Str("keyFile", keyFile).Msg("Cannot load certificates for new FredClient")
		return nil
	}

	// Create a new cert pool and add our own CA certificate
	rootCAs := x509.NewCertPool()

	loaded, err := ioutil.ReadFile(caFile)

	if err != nil {
		log.Fatal().Msgf("unexpected missing certfile: %v", err)
	}

	rootCAs.AppendCertsFromPEM(loaded)

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      rootCAs,
	}

	tc := credentials.NewTLS(tlsConfig)

	conn, err := grpc.Dial(host, grpc.WithTransportCredentials(tc))

	if err != nil {
		log.Fatal().Err(err).Msgf("Cannot create Grpc connection to client %s", host)
		return &Client{client: alexandraProto.NewMiddlewareClient(conn)}
	}
	log.Info().Msgf("Creating a connection to fred node: %s", host)
	return &Client{client: alexandraProto.NewMiddlewareClient(conn), conn: conn}
}
