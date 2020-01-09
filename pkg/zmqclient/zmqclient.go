package zmqclient

import (
	"encoding/json"
	"fmt"
	"net"

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
func (c *Client) SendCreateKeygroup(ip net.IP, port int, kgname commons.KeygroupName) (err error) {
	req, err := json.Marshal(&zmqcommon.Request{
		Keygroup: string(kgname),
	})

	if err != nil {
		return
	}

	err = c.sendMessage(zmqcommon.CreateKeygroup, ip, port, req)
	return
}

// SendDeleteKeygroup sends the message to the specified node.
func (c *Client) SendDeleteKeygroup(ip net.IP, port int, kgname commons.KeygroupName) (err error) {
	req, err := json.Marshal(&zmqcommon.Request{
		Keygroup: string(kgname),
	})

	if err != nil {
		return
	}

	err = c.sendMessage(zmqcommon.DeleteKeygroup, ip, port, req)
	return
}

// SendUpdate sends a PUT message to the specified node.
func (c *Client) SendUpdate(ip net.IP, port int, kgname commons.KeygroupName, id, value string) (err error) {
	req, err := json.Marshal(&zmqcommon.Request{
		Keygroup: string(kgname),
		ID:       id,
		Value:    value,
	})

	if err != nil {
		return
	}

	err = c.sendMessage(zmqcommon.PutItem, ip, port, req)
	return
}

// SendDelete sends the message to the specified node.
func (c *Client) SendDelete(ip net.IP, port int, kgname commons.KeygroupName, id string) (err error) {
	req, err := json.Marshal(&zmqcommon.Request{
		Keygroup: string(kgname),
		ID:       id,
	})

	if err != nil {
		return
	}

	err = c.sendMessage(zmqcommon.DeleteItem, ip, port, req)
	return
}

// sendMessage to the specified IP.
func (c *Client) sendMessage(msType byte, ip net.IP, port int, msg []byte) (err error) {
	endpoint := fmt.Sprintf("%s:%d", ip, port)
	cSender, exists := c.senders[endpoint]
	if !exists {
		log.Debug().Msgf("Created a new Socket to send to node %s:%d \n", ip, port)
		cSender = *NewSender(ip, port)
		c.senders[endpoint] = cSender
		// If the controller also needs to listen to answers
		// the sender needs to be passed to the controller
		//err = c.poller.Add(cSender.GetSocket())
	}

	if err != nil {
		return err
	}

	log.Debug().Bytes("msg", msg).Msgf("ZMQClient is sending a new message: ip=%d, msType=%v", ip, msType)
	err = cSender.SendMessageWithType(msType, msg)
	return
}

// SendAddNode sends the message to the specified node.
func (c *Client) SendAddNode(ip net.IP, port int, nodeID replication.ID, nodeIP net.IP, nodePort int) (err error) {
	req, err := json.Marshal(&zmqcommon.ReplicationRequest{
		Node: struct {
			ID   string
			IP   string
			Port string
		}{
			ID: string(nodeID),
			IP: string(nodeIP),
			Port: string(nodePort),
		},
	})

	if err != nil {
		return
	}

	err = c.sendMessage(zmqcommon.AddNode, ip, port, req)
	return
}

// SendRemoveNode sends the message to the specified node.
func (c *Client) SendRemoveNode(ip net.IP, port int, nodeID replication.ID) (err error) {
	req, err := json.Marshal(&zmqcommon.ReplicationRequest{
		Node: struct {
			ID   string
			IP   string
			Port string
		}{
			ID: string(nodeID),
		},
	})

	if err != nil {
		return
	}

	err = c.sendMessage(zmqcommon.RemoveNode, ip, port, req)
	return
}

// SendAddReplica sends the message to the specified node.
func (c *Client) SendAddReplica(ip net.IP, port int, kgname commons.KeygroupName, nodeID replication.ID) (err error) {
	req, err := json.Marshal(&zmqcommon.ReplicationRequest{
		Keygroup: string(kgname),
		Node: struct {
			ID   string
			IP   string
			Port string
		}{
			ID: string(nodeID),
		},
	})

	if err != nil {
		return
	}

	err = c.sendMessage(zmqcommon.AddReplica, ip, port, req)
	return
}

// SendRemoveReplica sends the message to the specified node.
func (c *Client) SendRemoveReplica(ip net.IP, port int, kgname commons.KeygroupName, nodeID replication.ID) (err error) {
	req, err := json.Marshal(&zmqcommon.ReplicationRequest{
		Keygroup: string(kgname),
		Node: struct {
			ID   string
			IP   string
			Port string
		}{
			ID: string(nodeID),
		},
	})

	if err != nil {
		return
	}

	err = c.sendMessage(zmqcommon.RemoveReplica, ip, port, req)
	return
}
