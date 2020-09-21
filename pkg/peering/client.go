package peering

import (
	"context"
	"fmt"

	"github.com/go-errors/errors"
	"github.com/rs/zerolog/log"
	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/fred"
	"gitlab.tu-berlin.de/mcc-fred/fred/proto/peering"
	"google.golang.org/grpc"
)

// Client is an peering client to communicate with peers.
type Client struct{}

// NewClient creates a new empty client to communicate with peers.
func NewClient() *Client {
	return &Client{}
}

// createConnAndClient creates a new connection to a server.
// Maybe it could be useful to reuse these?
// IDK whether this would be faster to store them in a map
func (c *Client) getConnAndClient(host fred.Address, port int) (client peering.NodeClient, conn *grpc.ClientConn) {
	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", host.Addr, port), grpc.WithInsecure())
	if err != nil {
		log.Fatal().Err(err).Msg("Cannot create Grpc connection")
		return nil, nil
	}
	log.Info().Msgf("Interclient: Created Connection to %s:%d", host.Addr, port)
	client = peering.NewNodeClient(conn)
	return
}

// logs the response and returns the correct error message
func (c *Client) dealWithStatusResponse(res *peering.StatusResponse, err error, from string) error {
	if res != nil {
		log.Debug().Msgf("Interclient got Response from %s, Status %s with Message %s and Error %s", from, res.Status, res.ErrorMessage, err)
	} else {
		log.Debug().Msgf("Interclient got empty Response from %s", from)
	}

	if err != nil || res == nil {
		return errors.New(err)
	}

	if res.Status == peering.EnumStatus_ERROR {
		return errors.New(res.ErrorMessage)
	}

	return nil

}

// Destroy currently does nothing, but might delete open connections if they are implemented
func (c *Client) Destroy() {
}

// SendCreateKeygroup sends this command to the server at this address
func (c *Client) SendCreateKeygroup(addr fred.Address, port int, kgname fred.KeygroupName, expiry int) error {
	client, conn := c.getConnAndClient(addr, port)
	res, err := client.CreateKeygroup(context.Background(), &peering.CreateKeygroupRequest{Keygroup: string(kgname), Expiry: int64(expiry)})

	err = c.dealWithStatusResponse(res, err, "CreateKeygroup")

	if err != nil {
		return err
	}

	return conn.Close()
}

// SendDeleteKeygroup sends this command to the server at this address
func (c *Client) SendDeleteKeygroup(addr fred.Address, port int, kgname fred.KeygroupName) error {
	client, conn := c.getConnAndClient(addr, port)
	res, err := client.DeleteKeygroup(context.Background(), &peering.DeleteKeygroupRequest{Keygroup: string(kgname)})

	err = c.dealWithStatusResponse(res, err, "DeleteKeygroup")

	if err != nil {
		return err
	}

	return conn.Close()
}

// SendUpdate sends this command to the server at this address
func (c *Client) SendUpdate(addr fred.Address, port int, kgname fred.KeygroupName, id string, value string) error {
	client, conn := c.getConnAndClient(addr, port)
	res, err := client.PutItem(context.Background(), &peering.PutItemRequest{
		Keygroup: string(kgname),
		Id:       id,
		Data:     value,
	})

	err = c.dealWithStatusResponse(res, err, "SendUpdate")

	if err != nil {
		return err
	}

	return conn.Close()
}

// SendDelete sends this command to the server at this address
func (c *Client) SendDelete(addr fred.Address, port int, kgname fred.KeygroupName, id string) error {
	client, _ := c.getConnAndClient(addr, port)
	res, err := client.DeleteItem(context.Background(), &peering.DeleteItemRequest{
		Keygroup: string(kgname),
		Id:       id,
	})
	return c.dealWithStatusResponse(res, err, "DeleteItem")
}

// SendAddReplica sends this command to the server at this address
func (c *Client) SendAddReplica(addr fred.Address, port int, kgname fred.KeygroupName, node fred.Node, expiry int) error {
	client, _ := c.getConnAndClient(addr, port)
	res, err := client.AddReplica(context.Background(), &peering.AddReplicaRequest{
		NodeId:   string(node.ID),
		Keygroup: string(kgname),
		Expiry:   int64(expiry),
	})
	return c.dealWithStatusResponse(res, err, "AddReplica")
}

// SendRemoveReplica sends this command to the server at this address
func (c *Client) SendRemoveReplica(addr fred.Address, port int, kgname fred.KeygroupName, node fred.Node) error {
	client, _ := c.getConnAndClient(addr, port)
	res, err := client.RemoveReplica(context.Background(), &peering.RemoveReplicaRequest{
		NodeId:   string(node.ID),
		Keygroup: string(kgname),
	})
	return c.dealWithStatusResponse(res, err, "RemoveReplica")
}
