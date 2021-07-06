package alexandra

import (
	"crypto/tls"

	alexandraProto "git.tu-berlin.de/mcc-fred/fred/proto/middleware"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type ClientsMgr struct {
	clients                 map[string]*Client
	clientsCert, clientsKey string
}

type Client struct {
	client alexandraProto.MiddlewareClient
	conn   *grpc.ClientConn
}

func newClientsManager(clientsCert, clientsKey string) *ClientsMgr {
	return &ClientsMgr{
		clients:     make(map[string]*Client),
		clientsCert: clientsCert,
		clientsKey:  clientsKey,
	}
}

func (m *ClientsMgr) GetClientTo(host string) (client *Client) {
	client = m.clients[host]
	if client != nil {
		return
	}
	client = newClient(host, m.clientsCert, m.clientsKey)
	m.clients[host] = client
	return
}

func newClient(host, certFile, keyFile string) *Client {
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)

	if err != nil {
		log.Fatal().Err(err).Str("certFile", certFile).Str("keyFile", keyFile).Msg("Cannot load certificates for new FredClient")
		return nil
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
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
