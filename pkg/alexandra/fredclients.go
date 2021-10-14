package alexandra

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"time"

	api "git.tu-berlin.de/mcc-fred/fred/proto/client"
	"github.com/DistributedClocks/GoVector/govec/vclock"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// Alpha is the value used for the exponential moving average. range=[0;1] with higher value => discount older observations faster
const alphaItemSpeed = float32(0.8)

type Client struct {
	host      string
	nodeID    string
	Client    api.ClientClient
	conn      *grpc.ClientConn
	ReadSpeed float32
}

func newClient(nodeID string, host string, certFile string, keyFile string, caCert string) *Client {

	if certFile == "" {
		log.Fatal().Msg("fredclient: no certificate file given")
	}

	if keyFile == "" {
		log.Fatal().Msg("fredclient: no key file given")
	}

	cert, err := tls.LoadX509KeyPair(certFile, keyFile)

	if err != nil {
		log.Fatal().Err(err).Str("certFile", certFile).Str("keyFile", keyFile).Msg("Cannot load certificates for new FredClient")
		return nil
	}

	// Create a new cert pool and add our own CA certificate
	rootCAs := x509.NewCertPool()

	loaded, err := ioutil.ReadFile(caCert)

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
		return &Client{Client: api.NewClientClient(conn)}
	}
	log.Info().Msgf("Creating a connection to fred node: %s", host)
	return &Client{
		Client:    api.NewClientClient(conn),
		conn:      conn,
		ReadSpeed: -1,
		host:      host,
		nodeID:    nodeID,
	}
}

// updateItemSpeed saves a moving average of how long it takes for a fred node to respond.
// the expectation is that read, writes, deletes and appends of items in keygroups should give an indication on whether
// a node is fast to reach for operations on items or not
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

func (c *Client) createKeygroup(ctx context.Context, keygroup string, mutable bool, expiry int64) (*api.Empty, error) {
	res, err := c.Client.CreateKeygroup(ctx, &api.CreateKeygroupRequest{
		Keygroup: keygroup,
		Mutable:  mutable,
		Expiry:   expiry,
	})
	return res, err
}

func (c *Client) deleteKeygroup(ctx context.Context, keygroup string) (*api.Empty, error) {
	res, err := c.Client.DeleteKeygroup(ctx, &api.DeleteKeygroupRequest{Keygroup: keygroup})
	return res, err
}

// Read also updates the moving average item speed
func (c *Client) read(ctx context.Context, keygroup string, id string) (*api.ReadResponse, error) {
	start := time.Now()
	res, err := c.Client.Read(ctx, &api.ReadRequest{
		Keygroup: keygroup,
		Id:       id,
	})
	if err == nil {
		elapsed := time.Since(start)
		c.updateItemSpeed(elapsed)
	}
	return res, err
}

/*
// Update also updates the moving average item speed
func (c *Client) update(ctx context.Context, keygroup string, id string, data string) (vclock.VClock, error) {
	start := time.Now()
	res, err := c.Client.Update(ctx, &api.UpdateRequest{
		Keygroup: keygroup,
		Id:       id,
		Data:     data,
	})

	if err != nil {
		return nil, err
	}

	elapsed := time.Since(start)
	c.updateItemSpeed(elapsed)

	return res.Version.Version, nil
}*/

// UpdateVersions also updates the moving average item speed
func (c *Client) updateVersions(ctx context.Context, keygroup string, id string, data string, versions []vclock.VClock) (vclock.VClock, error) {
	v := make([]*api.Version, len(versions))

	for i, vvector := range versions {
		v[i] = &api.Version{
			Version: vvector.GetMap(),
		}
	}

	start := time.Now()
	res, err := c.Client.Update(ctx, &api.UpdateRequest{
		Keygroup: keygroup,
		Id:       id,
		Data:     data,
		Versions: v,
	})

	if err != nil {
		return nil, err
	}

	elapsed := time.Since(start)
	c.updateItemSpeed(elapsed)

	return res.Version.Version, nil
}

/*
// Delete also updates the moving average item speed
func (c *Client) delete(ctx context.Context, keygroup string, id string) (vclock.VClock, error) {
	start := time.Now()
	res, err := c.Client.Delete(ctx, &api.DeleteRequest{
		Keygroup: keygroup,
		Id:       id,
	})

	if err != nil {
		return nil, err
	}

	elapsed := time.Since(start)
	c.updateItemSpeed(elapsed)

	return res.Version.Version, nil
}*/

// DeleteVersions also updates the moving average item speed
func (c *Client) deleteVersions(ctx context.Context, keygroup string, id string, versions []vclock.VClock) (vclock.VClock, error) {
	v := make([]*api.Version, len(versions))

	for i, vvector := range versions {
		v[i] = &api.Version{
			Version: vvector.GetMap(),
		}
	}

	start := time.Now()
	res, err := c.Client.Delete(ctx, &api.DeleteRequest{
		Keygroup: keygroup,
		Id:       id,
		Versions: v,
	})

	if err != nil {
		return nil, err
	}

	elapsed := time.Since(start)
	c.updateItemSpeed(elapsed)

	return res.Version.Version, nil
}

// Append also updates the moving average item speed
func (c *Client) append(ctx context.Context, keygroup string, data string) (*api.AppendResponse, error) {
	start := time.Now()
	id := uint64(time.Now().UnixNano())
	res, err := c.Client.Append(ctx, &api.AppendRequest{
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

func (c *Client) getKeygroupReplica(ctx context.Context, keygroup string) (*api.GetKeygroupInfoResponse, error) {
	res, err := c.Client.GetKeygroupInfo(ctx, &api.GetKeygroupInfoRequest{Keygroup: keygroup})
	return res, err
}

func (c *Client) getReplica(ctx context.Context, nodeID string) (*api.GetReplicaResponse, error) {
	return c.Client.GetReplica(ctx, &api.GetReplicaRequest{NodeId: nodeID})
}
