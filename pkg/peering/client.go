package peering

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"os"
	"sync"

	"git.tu-berlin.de/mcc-fred/fred/pkg/fred"
	"git.tu-berlin.de/mcc-fred/fred/proto/peering"
	"git.tu-berlin.de/mcc-fred/vclock"
	"github.com/go-errors/errors"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// Client is an peering client to communicate with peers.
type Client struct {
	conn        map[string]peering.NodeClient
	connLock    sync.RWMutex
	credentials credentials.TransportCredentials
}

// NewClient creates a new empty client to communicate with peers.
func NewClient(certFile string, keyFile string, caFile string) *Client {
	if certFile == "" {
		log.Fatal().Msg("peering client: no certificate file given")
	}

	if keyFile == "" {
		log.Fatal().Msg("peering client: no key file given")
	}

	if caFile == "" {
		log.Fatal().Msg("peering client: no root certificate file given")
	}

	cert, err := tls.LoadX509KeyPair(certFile, keyFile)

	if err != nil {
		log.Fatal().Err(err).Msg("peering client: Cannot load certificates")

		return nil
	}

	// Create a new cert pool and add our own CA certificate
	rootCAs, err := x509.SystemCertPool()

	if err != nil {
		log.Fatal().Err(err).Msg("peering client: Cannot load root certificates")
		return nil
	}

	loaded, err := os.ReadFile(caFile)

	if err != nil {
		log.Fatal().Msgf("peering client: unexpected missing certfile: %v", err)
	}

	rootCAs.AppendCertsFromPEM(loaded)

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		MinVersion:   tls.VersionTLS12,
		RootCAs:      rootCAs,
	}

	return &Client{
		conn:        make(map[string]peering.NodeClient),
		connLock:    sync.RWMutex{},
		credentials: credentials.NewTLS(tlsConfig),
	}
}

// getClient creates a new connection to a server or uses an existing one.
func (c *Client) getClient(host string) (peering.NodeClient, error) {

	c.connLock.RLock()
	client, ok := c.conn[host]
	c.connLock.RUnlock()

	if !ok {
		c.connLock.Lock()
		conn, err := grpc.Dial(host, grpc.WithTransportCredentials(c.credentials))

		if err != nil {
			c.connLock.Unlock()
			log.Error().Err(err).Msg("peering client: Cannot create Grpc connection")
			return nil, errors.New(err)
		}

		log.Debug().Msgf("peering client: Created Connection to %s", host)

		client = peering.NewNodeClient(conn)

		c.conn[host] = client
		c.connLock.Unlock()
	}

	return client, nil
}

// Destroy currently does nothing, but might delete open connections if they are implemented
func (c *Client) Destroy() {
}

// SendCreateKeygroup sends this command to the server at this address
func (c *Client) SendCreateKeygroup(host string, kgname fred.KeygroupName, expiry int) error {
	client, err := c.getClient(host)

	if err != nil {
		return errors.New(err)
	}

	_, err = client.CreateKeygroup(context.Background(), &peering.CreateKeygroupRequest{Keygroup: string(kgname), Expiry: int64(expiry)})

	if err != nil {
		return errors.New(err)
	}
	return nil
}

// SendDeleteKeygroup sends this command to the server at this address
func (c *Client) SendDeleteKeygroup(host string, kgname fred.KeygroupName) error {
	client, err := c.getClient(host)

	if err != nil {
		return errors.New(err)
	}

	_, err = client.DeleteKeygroup(context.Background(), &peering.DeleteKeygroupRequest{Keygroup: string(kgname)})

	if err != nil {
		return errors.New(err)
	}
	return nil
}

// SendUpdate sends this command to the server at this address
func (c *Client) SendUpdate(host string, kgname fred.KeygroupName, id string, value string, tombstoned bool, vvector vclock.VClock) error {
	client, err := c.getClient(host)

	if err != nil {
		return errors.New(err)
	}

	_, err = client.PutItem(context.Background(), &peering.PutItemRequest{
		Keygroup:   string(kgname),
		Id:         id,
		Val:        value,
		Tombstoned: tombstoned,
		Version:    vvector,
	})

	if err != nil {
		return errors.New(err)
	}
	return nil
}

// SendAppend sends this command to the server at this address
func (c *Client) SendAppend(host string, kgname fred.KeygroupName, id string, value string) error {
	client, err := c.getClient(host)

	if err != nil {
		return errors.New(err)
	}

	_, err = client.AppendItem(context.Background(), &peering.AppendItemRequest{
		Keygroup: string(kgname),
		Id:       id,
		Data:     value,
	})

	if err != nil {
		return errors.New(err)
	}
	return nil
}

// SendGetItem sends this command to the server at this address
func (c *Client) SendGetItem(host string, kgname fred.KeygroupName, id string) ([]fred.Item, error) {
	client, err := c.getClient(host)

	if err != nil {
		return nil, errors.New(err)
	}

	res, err := client.GetItem(context.Background(), &peering.GetItemRequest{
		Keygroup: string(kgname),
		Id:       id,
	})

	if err != nil {
		return nil, errors.New(err)
	}

	items := make([]fred.Item, len(res.Data))

	for i, item := range res.Data {
		items[i] = fred.Item{
			Keygroup: kgname,
			ID:       id,
			Val:      item.Val,
			Version:  item.Version,
		}
	}

	return items, nil
}

// SendGetAllItems sends this command to the server at this address
func (c *Client) SendGetAllItems(host string, kgname fred.KeygroupName) ([]fred.Item, error) {
	client, err := c.getClient(host)

	if err != nil {
		return nil, err
	}

	res, err := client.GetAllItems(context.Background(), &peering.GetAllItemsRequest{
		Keygroup: string(kgname),
	})

	if err != nil {
		return nil, errors.New(err)
	}

	d := make([]fred.Item, len(res.Data))

	for i, item := range res.Data {
		d[i] = fred.Item{
			Keygroup: kgname,
			ID:       item.Id,
			Val:      item.Val,
			Version:  item.Version,
		}
	}

	return d, nil
}
