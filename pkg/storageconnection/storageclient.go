package storageconnection

import (
	"context"
	"fmt"
	"io"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

// Client to a grpc server
type Client struct {
	dbClient DatabaseClient
	con      *grpc.ClientConn
}

// NewClient Client creates a new Client to communicate with a GRpc server
func NewClient(host string, port int) *Client {
	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", host, port), grpc.WithInsecure())

	if err != nil {
		log.Fatal().Err(err).Msg("Cannot create Grpc connection")
		return &Client{dbClient: NewDatabaseClient(conn)}
	}
	log.Info().Msgf("Creating a connection to remote storage: %s:%d", host, port)
	return &Client{dbClient: NewDatabaseClient(conn), con: conn}
}

// Read calls the same method on the remote server
func (c *Client) Read(kg string, id string) (string, error) {
	response, err := c.dbClient.Read(context.Background(), &Key{Keygroup: kg, Id: id})
	log.Debug().Err(err).Msgf("StorageClient: Read in: %#v %#v out: %#v", kg, id, response)
	return response.Val, err
}

// ReadAll calls the same method on the remote server
func (c *Client) ReadAll(kg string) (map[string]string, error) {
	stream, err := c.dbClient.ReadAll(context.Background(), &Keygroup{Keygroup: kg})
	if err != nil {
		log.Err(err).Msgf("StorageClient: Error in ReadAll in: %#v", kg)
		return nil, err
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
			return nil, err
		}

		responses[in.Id] = in.Val

	}
	return responses, nil
}

// Update calls the same method on the remote server
func (c *Client) Update(kg string, id string, val string) error {
	response, err := c.dbClient.Update(context.Background(), &Item{Keygroup: kg, Val: val, Id: id})
	log.Debug().Err(err).Msgf("StorageClient: Update in: %#v,%#v,%#v out: %#v", kg, id, val, response)
	return err
}

// Delete calls the same method on the remote server
func (c *Client) Delete(kg string, id string) error {
	response, err := c.dbClient.Delete(context.Background(), &Key{Keygroup: kg, Id: id})
	log.Debug().Err(err).Msgf("StorageClient: Delete in: %#v %#v out: %#v", kg, id, response)
	return err
}

// CreateKeygroup calls the same method on the remote server
func (c *Client) CreateKeygroup(kg string) error {
	keygroup := &Keygroup{Keygroup: kg}
	response, err := c.dbClient.CreateKeygroup(context.Background(), keygroup)
	log.Debug().Err(err).Msgf("StorageClient: CreateKeygroup in: %#v out: %#v", kg, response)
	return err
}

// DeleteKeygroup calls the same method on the remote server
func (c *Client) DeleteKeygroup(kg string) error {
	response, err := c.dbClient.DeleteKeygroup(context.Background(), &Keygroup{Keygroup: kg})
	log.Debug().Err(err).Msgf("StorageClient: DeleteKeygroup in: %#v out: %#v", kg, response)
	return err
}

// IDs calls the same method on the remote server
func (c *Client) IDs(kg string) ([]string, error) {
	stream, err := c.dbClient.IDs(context.Background(), &Keygroup{Keygroup: kg})
	if err != nil {
		log.Err(err).Msgf("StorageClient: Error in IDs in: %#v", kg)
		return nil, err
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
			return nil, err
		}
		responses = append(responses, in.Id)
	}
	return responses, nil
}

// Exists calls the same method on the remote server
func (c *Client) Exists(kg string, id string) bool {
	response, err := c.dbClient.Delete(context.Background(), &Key{Keygroup: kg, Id: id})
	log.Debug().Err(err).Msgf("StorageClient: Exists in: %#v %#v out: %#v", kg, id, response)
	return response.Success
}

// Destroy destroys the connection
func (c *Client) Destroy() {
	_ = c.con.Close()
}
