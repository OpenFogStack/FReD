package zmqclient

import (
	"encoding/json"
	"fmt"

	"github.com/rs/zerolog/log"

	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/commons"
	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/replication"
	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/zmqcommon"
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
func (c *Client) SendCreateKeygroup(addr replication.Address, port int, kgname commons.KeygroupName) (err error) {
	req, err := json.Marshal(&zmqcommon.Request{
		Keygroup: string(kgname),
	})

	if err != nil {
		return
	}

	err = c.sendMessage(zmqcommon.CreateKeygroup, addr, port, req)
	return
}

// SendDeleteKeygroup sends the message to the specified node.
func (c *Client) SendDeleteKeygroup(addr replication.Address, port int, kgname commons.KeygroupName) (err error) {
	req, err := json.Marshal(&zmqcommon.Request{
		Keygroup: string(kgname),
	})

	if err != nil {
		return
	}

	err = c.sendMessage(zmqcommon.DeleteKeygroup, addr, port, req)
	return
}

// SendUpdate sends a PUT message to the specified node.
func (c *Client) SendUpdate(addr replication.Address, port int, kgname commons.KeygroupName, id, value string) (err error) {
	req, err := json.Marshal(&zmqcommon.Request{
		Keygroup: string(kgname),
		ID:       id,
		Value:    value,
	})

	if err != nil {
		return
	}

	err = c.sendMessage(zmqcommon.PutItem, addr, port, req)
	return
}

// SendDelete sends the message to the specified node.
func (c *Client) SendDelete(addr replication.Address, port int, kgname commons.KeygroupName, id string) (err error) {
	req, err := json.Marshal(&zmqcommon.Request{
		Keygroup: string(kgname),
		ID:       id,
	})

	if err != nil {
		return
	}

	err = c.sendMessage(zmqcommon.DeleteItem, addr, port, req)
	return
}

// sendMessage to the specified Addr.
func (c *Client) sendMessage(msType byte, addr replication.Address, port int, msg []byte) (err error) {
	endpoint := fmt.Sprintf("%s:%d", addr, port)
	cSender, exists := c.senders[endpoint]
	if !exists {
		log.Debug().Msgf("Created a new Socket to send to node %s:%d \n", addr, port)
		cSender = *NewSender(addr, port)
		c.senders[endpoint] = cSender
		// If the controller also needs to listen to answers
		// the sender needs to be passed to the controller
		//err = c.poller.Add(cSender.GetSocket())
	}

	if err != nil {
		return err
	}

	log.Debug().Bytes("msg", msg).Msgf("ZMQClient is sending a new message: addr=%d, msType=%v", addr, msType)
	err = cSender.SendMessageWithType(msType, msg)
	return
}

// SendAddNode sends the message to the specified node.
func (c *Client) SendAddNode(addr replication.Address, port int, node replication.Node) (err error) {
	req, err := json.Marshal(&zmqcommon.ReplicationRequest{
		Node: node,
	})

	if err != nil {
		return
	}

	err = c.sendMessage(zmqcommon.AddNode, addr, port, req)
	return
}

// SendRemoveNode sends the message to the specified node.
func (c *Client) SendRemoveNode(addr replication.Address, port int, node replication.Node) (err error) {
	req, err := json.Marshal(&zmqcommon.ReplicationRequest{
		Node: node,
	})

	if err != nil {
		return
	}

	err = c.sendMessage(zmqcommon.RemoveNode, addr, port, req)
	return
}

// SendAddReplica sends the message to the specified node.
func (c *Client) SendAddReplica(addr replication.Address, port int, kgname commons.KeygroupName, node replication.Node) (err error) {
	req, err := json.Marshal(&zmqcommon.ReplicationRequest{
		Keygroup: kgname,
		Node: node,
	})

	if err != nil {
		return
	}

	err = c.sendMessage(zmqcommon.AddReplica, addr, port, req)
	return
}

// SendRemoveReplica sends the message to the specified node.
func (c *Client) SendRemoveReplica(addr replication.Address, port int, kgname commons.KeygroupName, node replication.Node) (err error) {
	req, err := json.Marshal(&zmqcommon.ReplicationRequest{
		Keygroup: kgname,
		Node: node,
	})

	if err != nil {
		return
	}

	err = c.sendMessage(zmqcommon.RemoveReplica, addr, port, req)
	return
}
