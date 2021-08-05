package alexandratest

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"

	alexandraProto "git.tu-berlin.de/mcc-fred/fred/proto/middleware"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type alexandraClient struct {
	client alexandraProto.MiddlewareClient
}

func newAlexandraClient(address string, clientCert string, clientKey string, caCert string) *alexandraClient {
	if clientCert == "" {
		log.Fatal().Msg("alexandra client: no certificate file given")
	}

	if clientKey == "" {
		log.Fatal().Msg("alexandra client: no key file given")
	}

	if caCert == "" {
		log.Fatal().Msg("alexandra client: no root certificate file given")
	}

	cert, err := tls.LoadX509KeyPair(clientCert, clientKey)

	if err != nil {
		log.Fatal().Err(err).Msg("alexandraClient: Cannot load certificates")
		return nil
	}

	// Create a new cert pool and add our own CA certificate
	rootCAs := x509.NewCertPool()

	loaded, err := ioutil.ReadFile(caCert)

	if err != nil {
		log.Fatal().Msgf("alexandraClient: unexpected missing certfile: %v", err)
	}

	rootCAs.AppendCertsFromPEM(loaded)

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		MinVersion:   tls.VersionTLS12,
		RootCAs:      rootCAs,
	}

	tc := credentials.NewTLS(tlsConfig)

	// IP of alexandraProto in our test setup
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(tc))

	if err != nil {
		log.Fatal().Err(err).Msg("Cannot create Grpc connection")
		return nil
	}
	c := alexandraProto.NewMiddlewareClient(conn)

	return &alexandraClient{
		client: c,
	}
}

func (c *alexandraClient) createKeygroup(firstNodeID string, kgname string, mutable bool, expiry int64) error {
	log.Debug().Msgf("createKeygroup: %s, %s, %t, %d", firstNodeID, kgname, mutable, expiry)
	_, err := c.client.CreateKeygroup(context.Background(), &alexandraProto.CreateKeygroupRequest{
		Keygroup:    kgname,
		Mutable:     mutable,
		Expiry:      expiry,
		FirstNodeId: firstNodeID,
	})

	if err != nil {
		return err
	}

	return nil
}

func (c *alexandraClient) update(kgname, id, data string) error {
	log.Debug().Msgf("update: %s, %s, %s", kgname, id, data)
	_, err := c.client.Update(context.Background(), &alexandraProto.UpdateRequest{
		Keygroup: kgname,
		Id:       id,
		Data:     data,
	})

	if err != nil {
		return err
	}

	return nil
}

func (c *alexandraClient) read(keygroup, id string, minExpiry int64) (string, error) {
	log.Debug().Msgf("read: %s, %s, %d", keygroup, id, minExpiry)
	res, err := c.client.Read(context.Background(), &alexandraProto.ReadRequest{
		Keygroup:  keygroup,
		Id:        id,
		MinExpiry: minExpiry,
	})

	if err != nil {
		return "", err
	}

	return res.Data, nil
}

func (c *alexandraClient) addKeygroupReplica(keygroup, node string, expiry int64) error {
	log.Debug().Msgf("addKeygroupReplica: %s, %s, %d", keygroup, node, expiry)
	_, err := c.client.AddReplica(context.Background(), &alexandraProto.AddReplicaRequest{
		Keygroup: keygroup,
		NodeId:   node,
		Expiry:   expiry,
	})

	if err != nil {
		return err
	}

	return nil
}
