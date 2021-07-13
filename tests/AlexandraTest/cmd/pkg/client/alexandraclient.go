package client

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"

	alexandra "git.tu-berlin.de/mcc-fred/fred/proto/middleware"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type AlexandraClient struct {
	client alexandra.MiddlewareClient
}

func NewAlexandraClient(address string) AlexandraClient {
	cert, err := tls.LoadX509KeyPair("/cert/client.crt", "/cert/client.key")

	if err != nil {
		log.Fatal().Err(err).Msg("Cannot load certificates")
		return AlexandraClient{}
	}

	// Create a new cert pool and add our own CA certificate
	rootCAs := x509.NewCertPool()

	loaded, err := ioutil.ReadFile("/cert/ca.crt")

	if err != nil {
		log.Fatal().Msgf("unexpected missing certfile: %v", err)
	}

	rootCAs.AppendCertsFromPEM(loaded)

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		MinVersion:   tls.VersionTLS12,
	}

	tc := credentials.NewTLS(tlsConfig)

	// IP of alexandra in our test setup
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(tc))

	if err != nil {
		log.Fatal().Err(err).Msg("Cannot create Grpc connection")
		return AlexandraClient{}
	}
	c := alexandra.NewMiddlewareClient(conn)

	return AlexandraClient{
		client: c,
	}
}

func (c *AlexandraClient) dealWithResponse(operation string, status alexandra.EnumStatus, err error, expectError bool) {
	// Got error but expected none
	if err != nil && !expectError {
		log.Fatal().Err(err).Msgf("%s got Error but expected no error (status %s)", operation, status)
	} else if err == nil && expectError {
		// Got no error but expected error
		log.Fatal().Msgf("%s got no error but expected an error (status %s)", operation, status)
	}
}

func (c *AlexandraClient) CreateKeygroup(kgname string, mutable bool, expiry int64, expectError bool) {
	res, err := c.client.CreateKeygroup(context.Background(), &alexandra.CreateKeygroupRequest{
		Keygroup: kgname,
		Mutable:  mutable,
		Expiry:   expiry,
	})
	c.dealWithResponse("CreateKeygroup", res.Status, err, expectError)
}

func (c *AlexandraClient) Update(kgname, id, data string, expectError bool) {
	res, err := c.client.Update(context.Background(), &alexandra.UpdateRequest{
		Keygroup: kgname,
		Id:       id,
		Data:     data,
	})
	c.dealWithResponse("Update", res.Status, err, expectError)
}
