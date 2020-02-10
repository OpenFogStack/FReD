package storage

import (
	"context"
	"fmt"
	"io"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"

	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/commons"
	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/data"
)

// Client to a grpc server
type Client struct {
	dbClient DatabaseClient
	con      grpc.ClientConn
}

// NewClient Client creates a new Client to communicate with a GRpc server
func NewClient(host string, port int) *Client {
	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", host, port), grpc.WithInsecure())

	if err != nil {
		log.Fatal().Err(err).Msg("Cannot create Grpc connection")
		return &Client{dbClient: NewDatabaseClient(conn)}
	}
	log.Info().Msgf("Creating a connection to remote storage: %s:%d", host, port)
	return &Client{dbClient: NewDatabaseClient(conn), con: *conn}
}

// Read calls the same method on the remote server
func (c Client) Read(kg commons.KeygroupName, id string) (string, error) {
	response, err := c.dbClient.Read(context.Background(), &Key{Keygroup: string(kg), Id: id})
	log.Debug().Err(err).Msgf("StorageClient: Read in: %#v %#v out: %#v", kg, id, response)
	return response.Data, err
}

// ReadAll calls the same method on the remote server
func (c Client) ReadAll(kg commons.KeygroupName) ([]data.Item, error) {
	stream, err := c.dbClient.ReadAll(context.Background(), &Keygroup{Keygroup: string(kg)})
	if err != nil {
		log.Err(err).Msgf("StorageClient: Error in ReadAll in: %#v", kg)
		return nil, err
	}
	var responses []data.Item
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
		responses = append(responses, *RPCItemToDataItem(in))
	}
	return responses, nil
}

// Update calls the same method on the remote server
func (c Client) Update(i data.Item) error {
	response, err := c.dbClient.Update(context.Background(), &Item{Keygroup: KeygroupObjectToString(i.Keygroup), Data: i.Data, Id: i.ID})
	log.Debug().Err(err).Msgf("StorageClient: Update in: %#v out: %#v", i, response)
	return err
}

// Delete calls the same method on the remote server
func (c Client) Delete(kg commons.KeygroupName, id string) error {
	response, err := c.dbClient.Delete(context.Background(), &Key{Keygroup: string(kg), Id: id})
	log.Debug().Err(err).Msgf("StorageClient: Delete in: %#v %#v out: %#v", kg, id, response)
	return err
}

// CreateKeygroup calls the same method on the remote server
func (c Client) CreateKeygroup(kg commons.KeygroupName) error {
	keygroup := &Keygroup{Keygroup: string(kg)}
	response, err := c.dbClient.CreateKeygroup(context.Background(), keygroup)
	log.Debug().Err(err).Msgf("StorageClient: CreateKeygroup in: %#v out: %#v", kg, response)
	return err
}

// DeleteKeygroup calls the same method on the remote server
func (c Client) DeleteKeygroup(kg commons.KeygroupName) error {
	response, err := c.dbClient.DeleteKeygroup(context.Background(), &Keygroup{Keygroup: string(kg)})
	log.Debug().Err(err).Msgf("StorageClient: DeleteKeygroup in: %#v out: %#v", kg, response)
	return err
}

// IDs calls the same method on the remote server
func (c Client) IDs(kg commons.KeygroupName) ([]data.Item, error) {
	stream, err := c.dbClient.IDs(context.Background(), &Keygroup{Keygroup: string(kg)})
	if err != nil {
		log.Err(err).Msgf("StorageClient: Error in IDs in: %#v", kg)
		return nil, err
	}
	var responses []data.Item
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
		responses = append(responses, *RPCKeyToItem(in))
	}
	return responses, nil
}

// Exists calls the same method on the remote server
func (c Client) Exists(kg commons.KeygroupName, id string) bool {
	response, err := c.dbClient.Delete(context.Background(), &Key{Keygroup: string(kg), Id: id})
	log.Debug().Err(err).Msgf("StorageClient: Exists in: %#v %#v out: %#v", kg, id, response)
	return response.Success
}

// Destroy destroys the connection
func (c Client) Destroy() {
	c.con.Close()
}
