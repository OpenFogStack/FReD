package storageclient

import (
	"context"
	"crypto/tls"
	"io"

	"github.com/go-errors/errors"
	"github.com/rs/zerolog/log"
	"gitlab.tu-berlin.de/mcc-fred/fred/proto/storage"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// Client to a grpc server
type Client struct {
	dbClient storage.DatabaseClient
	con      *grpc.ClientConn
}

// NewClient Client creates a new Client to communicate with a GRpc server
func NewClient(host, certFile, keyFile string) *Client {
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)

	if err != nil {
		log.Fatal().Err(err).Msg("Cannot load certificates")
		return nil
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
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
		if err == io.EOF {
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
func (c *Client) Update(kg string, id string, val string, expiry int) error {
	response, err := c.dbClient.Update(context.Background(), &storage.UpdateItem{Keygroup: kg, Val: val, Id: id, Expiry: int64(expiry)})
	log.Debug().Err(err).Msgf("StorageClient: Update in: %#v,%#v,%#v out: %#v", kg, id, val, response)

	if err != nil {
		return errors.New(err)
	}

	return nil
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
		if err == io.EOF {
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

// Exists calls the same method on the remote server
func (c *Client) Exists(kg string, id string) bool {
	response, err := c.dbClient.Exists(context.Background(), &storage.Key{Keygroup: kg, Id: id})
	log.Debug().Err(err).Msgf("StorageClient: Exists in: %#v %#v out: %#v", kg, id, response)

	if err != nil {
		return false
	}

	return response.Success
}

// Destroy destroys the connection
func (c *Client) Destroy() {
	_ = c.con.Close()
}
