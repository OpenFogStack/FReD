package client

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"os"

	"git.tu-berlin.de/mcc-fred/fred/proto/middleware"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type AlexandraClient struct {
	client middleware.MiddlewareClient
}

func NewAlexandraClient(address string) AlexandraClient {
	cert, err := tls.LoadX509KeyPair("/cert/client.crt", "/cert/client.key")

	if err != nil {
		log.Fatal().Err(err).Msg("Cannot load certificates")
		return AlexandraClient{}
	}

	// Create a new cert pool and add our own CA certificate
	rootCAs := x509.NewCertPool()

	loaded, err := os.ReadFile("/cert/ca.crt")

	if err != nil {
		log.Fatal().Msgf("unexpected missing certfile: %v", err)
	}

	rootCAs.AppendCertsFromPEM(loaded)

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		MinVersion:   tls.VersionTLS12,
		RootCAs:      rootCAs,
	}

	tc := credentials.NewTLS(tlsConfig)

	// IP of middleware in our test setup
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(tc))

	if err != nil {
		log.Fatal().Err(err).Msg("Cannot create Grpc connection")
		return AlexandraClient{}
	}
	c := middleware.NewMiddlewareClient(conn)

	return AlexandraClient{
		client: c,
	}
}

func (c *AlexandraClient) dealWithResponse(operation string, err error, expectError bool) {
	// Got error but expected none
	if err != nil && !expectError {
		log.Fatal().Err(err).Msgf("%s got Error but expected no error", operation)
	} else if err == nil && expectError {
		// Got no error but expected error
		log.Fatal().Msgf("%s got no error but expected an error", operation)
	}
}

func (c *AlexandraClient) CreateKeygroup(firstNodeID string, kgname string, mutable bool, expiry int64, expectError bool) {
	log.Debug().Msgf("CreateKeygroup: %s, %s, %t, %d", firstNodeID, kgname, mutable, expiry)
	_, err := c.client.CreateKeygroup(context.Background(), &middleware.CreateKeygroupRequest{
		Keygroup:    kgname,
		Mutable:     mutable,
		Expiry:      expiry,
		FirstNodeId: firstNodeID,
	})
	// res.status
	c.dealWithResponse("CreateKeygroup", err, expectError)
}

func (c *AlexandraClient) Update(kgname, id, data string, expectError bool) {
	log.Debug().Msgf("Update: %s, %s, %s", kgname, id, data)
	_, err := c.client.Update(context.Background(), &middleware.UpdateRequest{
		Keygroup: kgname,
		Id:       id,
		Data:     data,
	})
	c.dealWithResponse("Update", err, expectError)
}

func (c *AlexandraClient) Read(keygroup, id string, minExpiry int64, expectError bool) []string {
	log.Debug().Msgf("Read: %s, %s, %d", keygroup, id, minExpiry)

	res, err := c.client.Read(context.Background(), &middleware.ReadRequest{
		Keygroup:  keygroup,
		Id:        id,
		MinExpiry: minExpiry,
	})

	c.dealWithResponse("Read", err, expectError)

	if err != nil {
		return nil
	}

	vals := make([]string, len(res.Items))

	for i := range res.Items {
		vals[i] = res.Items[i].Val
	}

	return vals
}

func (c *AlexandraClient) AddKeygroupReplica(keygroup, node string, expiry int64, expectError bool) {
	log.Debug().Msgf("AddKeygroupReplica: %s, %s, %d", keygroup, node, expiry)
	_, err := c.client.AddReplica(context.Background(), &middleware.AddReplicaRequest{
		Keygroup: keygroup,
		NodeId:   node,
		Expiry:   expiry,
	})
	c.dealWithResponse("AddKeygroupReplica", err, expectError)
}
