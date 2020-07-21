package zmq

import (
	"encoding/json"
	"fmt"

	"github.com/rs/zerolog/log"

	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/fred"
)

// Client : Linter wants a comment here. Linter is dumb.
type Client struct {
	senders map[string]Sender
}

// NewClient creates a new Client.
func NewClient() (client *Client) {
	client = &Client{senders: make(map[string]Sender)}
	return
}

// Destroy the server.
func (c *Client) Destroy() {
	for _, sender := range c.senders {
		sender.Destroy()
	}
}

// SendCreateKeygroup sends the message to the specified node.
func (c *Client) SendCreateKeygroup(addr fred.Address, port int, kgname fred.KeygroupName) (err error) {
	req, err := json.Marshal(&DataRequest{
		Keygroup: string(kgname),
	})

	if err != nil {
		return
	}

	err = c.sendMessage(CreateKeygroup, addr, port, req)
	return
}

// SendDeleteKeygroup sends the message to the specified node.
func (c *Client) SendDeleteKeygroup(addr fred.Address, port int, kgname fred.KeygroupName) (err error) {
	req, err := json.Marshal(&DataRequest{
		Keygroup: string(kgname),
	})

	if err != nil {
		return
	}

	err = c.sendMessage(DeleteKeygroup, addr, port, req)
	return
}

// SendUpdate sends a PUT message to the specified node.
func (c *Client) SendUpdate(addr fred.Address, port int, kgname fred.KeygroupName, id, value string) (err error) {
	req, err := json.Marshal(&DataRequest{
		Keygroup: string(kgname),
		ID:       id,
		Value:    value,
	})

	if err != nil {
		return
	}

	err = c.sendMessage(PutItem, addr, port, req)
	return
}

// SendDelete sends the message to the specified node.
func (c *Client) SendDelete(addr fred.Address, port int, kgname fred.KeygroupName, id string) (err error) {
	req, err := json.Marshal(&DataRequest{
		Keygroup: string(kgname),
		ID:       id,
	})

	if err != nil {
		return
	}

	err = c.sendMessage(DeleteItem, addr, port, req)
	return
}

// sendMessage to the specified Addr.
func (c *Client) sendMessage(msType byte, addr fred.Address, port int, msg []byte) (err error) {
	endpoint := fmt.Sprintf("%s:%d", addr.Addr, port)
	cSender, exists := c.senders[endpoint]
	if !exists {
		log.Debug().Msgf("Created a new Socket to send to node %s:%d \n", addr.Addr, port)
		cSender = *NewSender(addr, port)
		c.senders[endpoint] = cSender
		// If the controller also needs to listen to answers
		// the sender needs to be passed to the controller
		//err = c.poller.Add(cSender.GetSocket())
	}

	if err != nil {
		return err
	}

	log.Debug().Bytes("msg", msg).Msgf("ZMQClient is sending a new message: addr=%s, msType=%#X", addr.Addr, msType)
	err = cSender.SendMessageWithType(msType, msg)
	return
}

// SendAddNode sends the message to the specified node. Receiver should add the node to the list of known nodes.
func (c *Client) SendAddNode(addr fred.Address, port int, node fred.Node) (err error) {
	req, err := json.Marshal(&ReplicationRequest{
		Node: node,
	})

	if err != nil {
		return
	}

	err = c.sendMessage(AddNode, addr, port, req)
	return
}

// SendRemoveNode sends the message to the specified node.
func (c *Client) SendRemoveNode(addr fred.Address, port int, node fred.Node) (err error) {
	req, err := json.Marshal(&ReplicationRequest{
		Node: node,
	})

	if err != nil {
		return
	}

	err = c.sendMessage(RemoveNode, addr, port, req)
	return
}

// SendAddReplica sends the message to the specified node.
func (c *Client) SendAddReplica(addr fred.Address, port int, kgname fred.KeygroupName, node fred.Node) (err error) {
	req, err := json.Marshal(&ReplicationRequest{
		Keygroup: kgname,
		Node:     node,
	})

	if err != nil {
		return
	}

	err = c.sendMessage(AddReplica, addr, port, req)
	return
}

// SendRemoveReplica sends the message to the specified node.
func (c *Client) SendRemoveReplica(addr fred.Address, port int, kgname fred.KeygroupName, node fred.Node) (err error) {
	req, err := json.Marshal(&ReplicationRequest{
		Keygroup: kgname,
		Node:     node,
	})

	if err != nil {
		return
	}

	err = c.sendMessage(RemoveReplica, addr, port, req)
	return
}
