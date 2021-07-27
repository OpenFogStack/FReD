package alexandra

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"time"

	fredClients "git.tu-berlin.de/mcc-fred/fred/proto/client"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// Alpha is the value used for the exponential moving average. range=[0;1] with higher value => discount older observations faster
const alphaItemSpeed = float32(0.8)

type Client struct {
	Client    fredClients.ClientClient
	conn      *grpc.ClientConn
	ReadSpeed float32
}

func NewClient(host, certFile, keyFile string) *Client {
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)

	if err != nil {
		log.Fatal().Err(err).Str("certFile", certFile).Str("keyFile", keyFile).Msg("Cannot load certificates for new FredClient")
		return nil
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
		RootCAs:      rootCAs,
	}

	tc := credentials.NewTLS(tlsConfig)

	conn, err := grpc.Dial(host, grpc.WithTransportCredentials(tc))

	if err != nil {
		log.Fatal().Err(err).Msgf("Cannot create Grpc connection to client %s", host)
		return &Client{Client: fredClients.NewClientClient(conn)}
	}
	log.Info().Msgf("Creating a connection to fred node: %s", host)
	return &Client{Client: fredClients.NewClientClient(conn), conn: conn, ReadSpeed: -1}
}

// updateItemSpeed saves a moving average of how long it takes for a fred node to respond.
// the expectation is that read, writes, deletes and appends of items in keygroups should give an indication on whether
// a node ist fast to reach for operations on items or not
// see https://en.wikipedia.org/wiki/Moving_average#Exponential_moving_average
func (c *Client) updateItemSpeed(elapsed time.Duration) {
	elapsedMs := float32(elapsed.Milliseconds())
	if c.ReadSpeed == -1 {
		// Read speed was not initialized
		c.ReadSpeed = elapsedMs
	} else {
		c.ReadSpeed = alphaItemSpeed*elapsedMs + (1-alphaItemSpeed)*c.ReadSpeed
	}
}

func (c *Client) CreateKeygroup(ctx context.Context, keygroup string, mutable bool, expiry int64) (*fredClients.StatusResponse, error) {
	res, err := c.Client.CreateKeygroup(ctx, &fredClients.CreateKeygroupRequest{
		Keygroup: keygroup,
		Mutable:  mutable,
		Expiry:   expiry,
	})
	return res, err
}

func (c *Client) DeleteKeygroup(ctx context.Context, keygroup string) (*fredClients.StatusResponse, error) {
	res, err := c.Client.DeleteKeygroup(ctx, &fredClients.DeleteKeygroupRequest{Keygroup: keygroup})
	return res, err
}

// Read also updates the moving average item speed
func (c *Client) Read(ctx context.Context, keygroup string, id string) (*fredClients.ReadResponse, error) {
	start := time.Now()
	res, err := c.Client.Read(ctx, &fredClients.ReadRequest{
		Keygroup: keygroup,
		Id:       id,
	})
	if err == nil {
		elapsed := time.Since(start)
		c.updateItemSpeed(elapsed)
	}
	return res, err
}

// Update also updates the moving average item speed
func (c *Client) Update(ctx context.Context, keygroup string, id string, data string) (*fredClients.StatusResponse, error) {
	start := time.Now()
	res, err := c.Client.Update(ctx, &fredClients.UpdateRequest{
		Keygroup: keygroup,
		Id:       id,
		Data:     data,
	})
	if err == nil {
		elapsed := time.Since(start)
		c.updateItemSpeed(elapsed)
	}
	return res, err
}

// Delete also updates the moving average item speed
func (c *Client) Delete(ctx context.Context, keygroup string, id string) (*fredClients.StatusResponse, error) {
	start := time.Now()
	res, err := c.Client.Delete(ctx, &fredClients.DeleteRequest{
		Keygroup: keygroup,
		Id:       id,
	})
	if err == nil {
		elapsed := time.Since(start)
		c.updateItemSpeed(elapsed)
	}
	return res, err
}

// Append also updates the moving average item speed
func (c *Client) Append(ctx context.Context, keygroup string, data string) (*fredClients.AppendResponse, error) {
	start := time.Now()
	res, err := c.Client.Append(ctx, &fredClients.AppendRequest{
		Keygroup: keygroup,
		Data:     data,
	})
	if err == nil {
		elapsed := time.Since(start)
		c.updateItemSpeed(elapsed)
	}
	return res, err
}

func (c *Client) AddReplica(ctx context.Context, keygroup string, nodeID string, expiry int64) (*fredClients.StatusResponse, error) {
	res, err := c.Client.AddReplica(ctx, &fredClients.AddReplicaRequest{
		Keygroup: keygroup,
		NodeId:   nodeID,
		Expiry:   expiry,
	})
	return res, err
}

func (c *Client) GetKeygroupReplica(ctx context.Context, keygroup string) (*fredClients.GetKeygroupReplicaResponse, error) {
	res, err := c.Client.GetKeygroupReplica(ctx, &fredClients.GetKeygroupReplicaRequest{Keygroup: keygroup})
	return res, err
}

func (c *Client) RemoveReplica(ctx context.Context, keygroup string, nodeID string) (*fredClients.StatusResponse, error) {
	return c.Client.RemoveReplica(ctx, &fredClients.RemoveReplicaRequest{
		Keygroup: keygroup,
		NodeId:   nodeID,
	})
}

func (c *Client) GetReplica(ctx context.Context, nodeID string) (*fredClients.GetReplicaResponse, error) {
	return c.Client.GetReplica(ctx, &fredClients.GetReplicaRequest{NodeId: nodeID})
}

func (c *Client) GetAllReplica(ctx context.Context) (*fredClients.GetAllReplicaResponse, error) {
	return c.Client.GetAllReplica(ctx, &fredClients.GetAllReplicaRequest{})
}

func (c *Client) GetKeygroupTriggers(ctx context.Context, keygroup string) (*fredClients.GetKeygroupTriggerResponse, error) {
	return c.Client.GetKeygroupTriggers(ctx, &fredClients.GetKeygroupTriggerRequest{Keygroup: keygroup})
}

func (c *Client) AddTrigger(ctx context.Context, keygroup string, triggerID string, triggerHost string) (*fredClients.StatusResponse, error) {
	return c.Client.AddTrigger(ctx, &fredClients.AddTriggerRequest{
		Keygroup:    keygroup,
		TriggerId:   triggerID,
		TriggerHost: triggerHost,
	})
}

func (c *Client) RemoveTrigger(ctx context.Context, keygroup, triggerID string) (*fredClients.StatusResponse, error) {
	return c.Client.RemoveTrigger(ctx, &fredClients.RemoveTriggerRequest{
		Keygroup:  keygroup,
		TriggerId: triggerID,
	})
}

func (c *Client) AddUser(ctx context.Context, user, keygroup string, role fredClients.UserRole) (*fredClients.StatusResponse, error) {
	return c.Client.AddUser(ctx, &fredClients.UserRequest{
		User:     user,
		Keygroup: keygroup,
		Role:     role,
	})
}

func (c *Client) RemoveUser(ctx context.Context, user, keygroup string, role fredClients.UserRole) (*fredClients.StatusResponse, error) {
	return c.Client.RemoveUser(ctx, &fredClients.UserRequest{
		User:     user,
		Keygroup: keygroup,
		Role:     role,
	})
}
