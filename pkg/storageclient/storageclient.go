package storageclient

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"os"

	"git.tu-berlin.de/mcc-fred/fred/proto/storage"
	"github.com/DistributedClocks/GoVector/govec/vclock"
	"github.com/go-errors/errors"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// Client to a grpc server
type Client struct {
	dbClient storage.DatabaseClient
	con      *grpc.ClientConn
}

// NewClient Client creates a new Client to communicate with a GRpc server
func NewClient(host, certFile string, keyFile string, caFiles []string) *Client {
	if certFile == "" {
		log.Fatal().Msg("Remote storage client: no certificate file given")
	}

	if keyFile == "" {
		log.Fatal().Msg("Remote storage client: no key file given")
	}

	if len(caFiles) == 0 {
		log.Fatal().Msg("Remote storage client: no root certificate files given")
	}

	cert, err := tls.LoadX509KeyPair(certFile, keyFile)

	if err != nil {
		log.Fatal().Err(err).Msg("Remote storage client: cannot load certificates")

		return nil
	}

	// Create a new cert pool and add our own CA certificate
	rootCAs, err := x509.SystemCertPool()

	if err != nil {
		log.Fatal().Err(err).Msg("Remote storage client: cannot load root certificates")
		return nil
	}

	for _, f := range caFiles {
		loaded, err := os.ReadFile(f)

		if err != nil {
			log.Fatal().Msgf("Remote storage client: unexpected missing certfile: %v", err)
		}

		rootCAs.AppendCertsFromPEM(loaded)
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		MinVersion:   tls.VersionTLS12,
		RootCAs:      rootCAs,
	}

	tc := credentials.NewTLS(tlsConfig)

	conn, err := grpc.Dial(host, grpc.WithTransportCredentials(tc))

	if err != nil {
		log.Fatal().Err(err).Msg("Remote storage client: cannot create Grpc connection")

		return &Client{dbClient: storage.NewDatabaseClient(conn)}
	}
	log.Info().Msgf("Remote storage client: creating a connection to remote storage: %s", host)

	return &Client{dbClient: storage.NewDatabaseClient(conn), con: conn}
}

// Close closes the connection to the storage service.
func (c *Client) Close() error {
	return c.con.Close()
}

// Read calls the same method on the remote server
func (c *Client) Read(kg, id string) ([]string, []vclock.VClock, bool, error) {
	res, err := c.dbClient.Read(context.Background(), &storage.ReadRequest{Keygroup: kg, Id: id})
	log.Debug().Err(err).Msgf("StorageClient: Read in: %+v %+v out: %+v", kg, id, res)

	if err != nil {
		return nil, nil, false, errors.New(err)
	}

	if len(res.Items) == 0 {
		return []string{}, []vclock.VClock{}, false, nil
	}

	vals := make([]string, len(res.Items))
	vvectors := make([]vclock.VClock, len(res.Items))

	for i, item := range res.Items {
		vals[i] = item.Val
		vvectors[i] = item.Version
	}

	return vals, vvectors, true, nil
}

// ReadSome calls the same method on the remote server
func (c *Client) ReadSome(kg, id string, count uint64) ([]string, []string, []vclock.VClock, error) {
	res, err := c.dbClient.Scan(context.Background(), &storage.ScanRequest{
		Keygroup: kg,
		Start:    id,
		Count:    count,
	})

	if err != nil {
		log.Err(err).Msgf("StorageClient: Error in Scan in: %+v count %d", kg, count)
		return nil, nil, nil, errors.New(err)
	}

	keys := make([]string, len(res.Items))
	vals := make([]string, len(res.Items))
	vvectors := make([]vclock.VClock, len(res.Items))

	for i, item := range res.Items {
		keys[i] = item.Id
		vals[i] = item.Val
		vvectors[i] = item.Version
	}

	return keys, vals, vvectors, nil
}

// ReadAll calls the same method on the remote server
func (c *Client) ReadAll(kg string) ([]string, []string, []vclock.VClock, error) {
	res, err := c.dbClient.ReadAll(context.Background(), &storage.ReadAllRequest{
		Keygroup: kg,
	})

	if err != nil {
		log.Err(err).Msgf("StorageClient: Error in ReadAll in: %+v", kg)
		return nil, nil, nil, errors.New(err)
	}

	keys := make([]string, len(res.Items))
	vals := make([]string, len(res.Items))
	vvectors := make([]vclock.VClock, len(res.Items))

	for i, item := range res.Items {
		keys[i] = item.Id
		vals[i] = item.Val
		vvectors[i] = item.Version
	}

	return keys, vals, vvectors, nil
}

// Update calls the same method on the remote server

func (c *Client) Update(kg string, id string, val string, expiry int, vvector vclock.VClock) error {
	response, err := c.dbClient.Update(context.Background(), &storage.UpdateRequest{
		Keygroup: kg,
		Val:      val,
		Id:       id,
		Expiry:   int64(expiry),
		Version:  vvector.GetMap(),
	})

	log.Debug().Err(err).Msgf("StorageClient: Update in: %+v,%+v,%+v,%+v out: %+v", kg, id, val, vvector, response)

	if err != nil {
		return errors.New(err)
	}

	return nil
}

// Append calls the same method on the remote server
func (c *Client) Append(kg string, id string, val string, expiry int) error {
	response, err := c.dbClient.Append(context.Background(), &storage.AppendRequest{
		Keygroup: kg,
		Id:       id,
		Val:      val,
		Expiry:   int64(expiry)},
	)
	log.Debug().Err(err).Msgf("StorageClient: Append in: %+v,%+v out: %+v", kg, val, response)

	if err != nil {
		return errors.New(err)
	}

	return nil
}

// Delete calls the same method on the remote server
func (c *Client) Delete(kg string, id string, vvector vclock.VClock) error {
	response, err := c.dbClient.Delete(context.Background(), &storage.DeleteRequest{
		Keygroup: kg,
		Id:       id,
		Version:  vvector.GetMap(),
	})

	log.Debug().Err(err).Msgf("StorageClient: Delete in: %+v %+v %+v out: %+v", kg, id, vvector, response)

	if err != nil {
		return errors.New(err)
	}

	return nil
}

// CreateKeygroup calls the same method on the remote server
func (c *Client) CreateKeygroup(kg string) error {
	keygroup := &storage.CreateKeygroupRequest{Keygroup: kg}
	response, err := c.dbClient.CreateKeygroup(context.Background(), keygroup)
	log.Debug().Err(err).Msgf("StorageClient: CreateKeygroup in: %+v out: %+v", kg, response)

	if err != nil {
		return errors.New(err)
	}

	return nil
}

// DeleteKeygroup calls the same method on the remote server
func (c *Client) DeleteKeygroup(kg string) error {
	response, err := c.dbClient.DeleteKeygroup(context.Background(), &storage.DeleteKeygroupRequest{Keygroup: kg})
	log.Debug().Err(err).Msgf("StorageClient: DeleteKeygroup in: %+v out: %+v", kg, response)

	if err != nil {
		return errors.New(err)
	}

	return nil
}

// IDs calls the same method on the remote server
func (c *Client) IDs(kg string) ([]string, error) {
	res, err := c.dbClient.IDs(context.Background(), &storage.IDsRequest{Keygroup: kg})
	if err != nil {
		log.Err(err).Msgf("StorageClient: Error in IDs in: %+v", kg)
		return nil, errors.New(err)
	}

	return res.Ids, nil
}

// Exists calls the same method on the remote server.
func (c *Client) Exists(kg string, id string) bool {
	response, err := c.dbClient.Exists(context.Background(), &storage.ExistsRequest{
		Keygroup: kg,
		Id:       id,
	})
	log.Debug().Err(err).Msgf("StorageClient: Exists in: %+v %+v out: %+v", kg, id, response)

	if err != nil {
		return false
	}

	return response.Exists
}

// ExistsKeygroup calls the same method on the remote server.
func (c *Client) ExistsKeygroup(kg string) bool {
	response, err := c.dbClient.ExistsKeygroup(context.Background(), &storage.ExistsKeygroupRequest{Keygroup: kg})
	log.Debug().Err(err).Msgf("StorageClient: ExistsKeygroup in: %+v out: %+v", kg, response)

	if err != nil {
		return false
	}

	return response.Exists
}

// Destroy destroys the connection
func (c *Client) Destroy() {
	_ = c.con.Close()
}

// AddKeygroupTrigger calls the same method on the remote server.
func (c *Client) AddKeygroupTrigger(kg string, id string, host string) error {
	keygroupTrigger := &storage.AddKeygroupTriggerRequest{
		Keygroup: kg,

		Id:   id,
		Host: host,
	}
	response, err := c.dbClient.AddKeygroupTrigger(context.Background(), keygroupTrigger)
	log.Debug().Err(err).Msgf("StorageClient: AddKeygroupTrigger in: %+v out: %+v", kg, response)

	if err != nil {
		return errors.New(err)
	}

	return nil
}

// DeleteKeygroupTrigger calls the same method on the remote server.
func (c *Client) DeleteKeygroupTrigger(kg string, id string) error {
	keygroupTrigger := &storage.DeleteKeygroupTriggerRequest{
		Keygroup: kg,
		Id:       id,
	}
	response, err := c.dbClient.DeleteKeygroupTrigger(context.Background(), keygroupTrigger)
	log.Debug().Err(err).Msgf("StorageClient: DeleteKeygroupTrigger in: %+v out: %+v", kg, response)

	if err != nil {
		return errors.New(err)
	}

	return nil
}

// GetKeygroupTrigger calls the same method on the remote server.
func (c *Client) GetKeygroupTrigger(kg string) (map[string]string, error) {
	keygroup := &storage.GetKeygroupTriggerRequest{
		Keygroup: kg,
	}
	res, err := c.dbClient.GetKeygroupTrigger(context.Background(), keygroup)
	log.Debug().Err(err).Msgf("StorageClient: GetKeygroupTrigger in: %+v", kg)

	if err != nil {
		return nil, errors.New(err)
	}

	triggers := make(map[string]string)

	for _, t := range res.Triggers {
		triggers[t.Id] = t.Host
	}

	return triggers, nil
}
