package storageclient

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"io"
	"io/ioutil"

	"git.tu-berlin.de/mcc-fred/fred/proto/storage"
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
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)

	if err != nil {
		log.Fatal().Err(err).Msg("Cannot load certificates")

		return nil
	}

	// Create a new cert pool and add our own CA certificate
	rootCAs, err := x509.SystemCertPool()

	if err != nil {
		log.Fatal().Err(err).Msg("Cannot load root certificates")
		return nil
	}

	for _, f := range caFiles {
		loaded, err := ioutil.ReadFile(f)

		if err != nil {
			log.Fatal().Msgf("unexpected missing certfile: %v", err)
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
		log.Fatal().Err(err).Msg("Cannot create Grpc connection")

		return &Client{dbClient: storage.NewDatabaseClient(conn)}
	}
	log.Info().Msgf("Creating a connection to remote storage: %s", host)

	return &Client{dbClient: storage.NewDatabaseClient(conn), con: conn}
}

// Close closes the connection to the storage service.
func (c *Client) Close() error {
	return c.con.Close()
}

// Read calls the same method on the remote server
func (c *Client) Read(kg string, id string) (string, error) {
	response, err := c.dbClient.Read(context.Background(), &storage.Key{Keygroup: kg, Id: id})
	log.Debug().Err(err).Msgf("StorageClient: Read in: %#v %#v out: %#v", kg, id, response)

	if err != nil {
		return "", errors.New(err)
	}

	return response.Val, nil
}

// Scan calls the same method on the remote server
func (c *Client) ReadSome(kg string, id string, count uint64) (map[string]string, error) {
	stream, err := c.dbClient.Scan(context.Background(), &storage.ScanRequest{
		Key:   &storage.Key{Keygroup: kg, Id: id},
		Count: count,
	})
	if err != nil {
		log.Err(err).Msgf("StorageClient: Error in Scan in: %#v count %d", kg, count)
		return nil, errors.New(err)
	}
	responses := make(map[string]string)

	for {
		in, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			// read done.
			break
		}
		if err != nil {
			log.Err(err).Msg("StorageClient: Error in Scan while receiving a Item")
			return nil, errors.New(err)
		}

		responses[in.Id] = in.Val

	}
	return responses, nil
}

// ReadAll calls the same method on the remote server
func (c *Client) ReadAll(kg string) (map[string]string, error) {
	stream, err := c.dbClient.ReadAll(context.Background(), &storage.Keygroup{Keygroup: kg})
	if err != nil {
		log.Err(err).Msgf("StorageClient: Error in ReadAll in: %#v", kg)
		return nil, errors.New(err)
	}
	responses := make(map[string]string)

	for {
		in, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			// read done.
			break
		}
		if err != nil {
			log.Err(err).Msg("StorageClient: Error in ReadAll while receiving a Item")
			return nil, errors.New(err)
		}

		responses[in.Id] = in.Val

	}
	return responses, nil
}

// Update calls the same method on the remote server
func (c *Client) Update(kg string, id string, val string, append bool, expiry int) error {
	response, err := c.dbClient.Update(context.Background(), &storage.UpdateItem{
		Keygroup: kg,
		Val:      val,
		Id:       id,
		Append:   append,
		Expiry:   int64(expiry)})
	log.Debug().Err(err).Msgf("StorageClient: Update in: %#v,%#v,%#v out: %#v", kg, id, val, response)

	if err != nil {
		return errors.New(err)
	}

	return nil
}

// Append calls the same method on the remote server
func (c *Client) Append(kg string, val string, expiry int) (string, error) {
	response, err := c.dbClient.Append(context.Background(), &storage.AppendItem{Keygroup: kg, Val: val, Expiry: int64(expiry)})
	log.Debug().Err(err).Msgf("StorageClient: Append in: %#v,%#v out: %#v", kg, val, response)

	if err != nil {
		return "", errors.New(err)
	}

	return response.Id, nil
}

// Delete calls the same method on the remote server
func (c *Client) Delete(kg string, id string) error {
	response, err := c.dbClient.Delete(context.Background(), &storage.Key{Keygroup: kg, Id: id})
	log.Debug().Err(err).Msgf("StorageClient: Delete in: %#v %#v out: %#v", kg, id, response)

	if err != nil {
		return errors.New(err)
	}

	return nil
}

// CreateKeygroup calls the same method on the remote server
func (c *Client) CreateKeygroup(kg string) error {
	keygroup := &storage.Keygroup{Keygroup: kg}
	response, err := c.dbClient.CreateKeygroup(context.Background(), keygroup)
	log.Debug().Err(err).Msgf("StorageClient: CreateKeygroup in: %#v out: %#v", kg, response)

	if err != nil {
		return errors.New(err)
	}

	return nil
}

// DeleteKeygroup calls the same method on the remote server
func (c *Client) DeleteKeygroup(kg string) error {
	response, err := c.dbClient.DeleteKeygroup(context.Background(), &storage.Keygroup{Keygroup: kg})
	log.Debug().Err(err).Msgf("StorageClient: DeleteKeygroup in: %#v out: %#v", kg, response)

	if err != nil {
		return errors.New(err)
	}

	return nil
}

// IDs calls the same method on the remote server
func (c *Client) IDs(kg string) ([]string, error) {
	stream, err := c.dbClient.IDs(context.Background(), &storage.Keygroup{Keygroup: kg})
	if err != nil {
		log.Err(err).Msgf("StorageClient: Error in IDs in: %#v", kg)
		return nil, errors.New(err)
	}
	var responses []string
	for {
		in, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			// read done.
			break
		}
		if err != nil {
			log.Err(err).Msg("StorageClient: Error in IDs while receiving a Item")
			return nil, errors.New(err)
		}
		responses = append(responses, in.Id)
	}
	return responses, nil
}

// Exists calls the same method on the remote server.
func (c *Client) Exists(kg string, id string) bool {
	response, err := c.dbClient.Exists(context.Background(), &storage.Key{Keygroup: kg, Id: id})
	log.Debug().Err(err).Msgf("StorageClient: Exists in: %#v %#v out: %#v", kg, id, response)

	if err != nil {
		return false
	}

	return response.Success
}

// ExistsKeygroup calls the same method on the remote server.
func (c *Client) ExistsKeygroup(kg string) bool {
	response, err := c.dbClient.ExistsKeygroup(context.Background(), &storage.Keygroup{Keygroup: kg})
	log.Debug().Err(err).Msgf("StorageClient: ExistsKeygroup in: %#v out: %#v", kg, response)

	if err != nil {
		return false
	}

	return response.Success
}

// Destroy destroys the connection
func (c *Client) Destroy() {
	_ = c.con.Close()
}

// AddKeygroupTrigger calls the same method on the remote server.
func (c *Client) AddKeygroupTrigger(kg string, id string, host string) error {
	keygroupTrigger := &storage.KeygroupTrigger{
		Keygroup: kg,
		Trigger: &storage.Trigger{
			Id:   id,
			Host: host,
		},
	}
	response, err := c.dbClient.AddKeygroupTrigger(context.Background(), keygroupTrigger)
	log.Debug().Err(err).Msgf("StorageClient: AddKeygroupTrigger in: %#v out: %#v", kg, response)

	if err != nil {
		return errors.New(err)
	}

	return nil
}

// DeleteKeygroupTrigger calls the same method on the remote server.
func (c *Client) DeleteKeygroupTrigger(kg string, id string) error {
	keygroupTrigger := &storage.KeygroupTrigger{
		Keygroup: kg,
		Trigger: &storage.Trigger{
			Id: id,
		},
	}
	response, err := c.dbClient.DeleteKeygroupTrigger(context.Background(), keygroupTrigger)
	log.Debug().Err(err).Msgf("StorageClient: DeleteKeygroupTrigger in: %#v out: %#v", kg, response)

	if err != nil {
		return errors.New(err)
	}

	return nil
}

// GetKeygroupTrigger calls the same method on the remote server.
func (c *Client) GetKeygroupTrigger(kg string) (map[string]string, error) {
	keygroup := &storage.Keygroup{
		Keygroup: kg,
	}
	stream, err := c.dbClient.GetKeygroupTrigger(context.Background(), keygroup)
	log.Debug().Err(err).Msgf("StorageClient: DeleteKeygroupTrigger in: %#v", kg)

	if err != nil {
		return nil, errors.New(err)
	}

	responses := make(map[string]string)

	for {
		in, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			// read done.
			break
		}
		if err != nil {
			log.Err(err).Msg("StorageClient: Error in GetKeygroupTrigger while receiving a Trigger")
			return nil, errors.New(err)
		}

		responses[in.Id] = in.Host
	}

	return responses, nil
}
