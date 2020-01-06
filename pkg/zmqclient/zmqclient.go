package zmqclient

import (
	"encoding/json"
	"fmt"
	"net"

	"github.com/rs/zerolog/log"

	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/zmqcommon"
)

// Client : Linter wants a comment here. Linter is dumb.
type Client struct {
	senders map[string]Sender
}

// NewClient creates a new Client.
func NewClient() (client *Client){
	client = &Client{senders: make(map[string]Sender)}
	return
}

// SendCreateKeygroup sends the message to the specified node.
func (c *Client) SendCreateKeygroup(ip net.IP, port int, kgname string) (err error) {
	req, err := json.Marshal(&zmqcommon.Request{
		Keygroup: kgname,
	})

	if err != nil {
		return
	}

	err = c.sendMessage(zmqcommon.CreateKeygroup, ip, port, req)
	return
}

// SendDeleteKeygroup sends the message to the specified node.
func (c *Client) SendDeleteKeygroup(ip net.IP, port int, kgname string) (err error) {
	req, err := json.Marshal(&zmqcommon.Request{
		Keygroup: kgname,
	})

	if err != nil {
		return
	}

	err = c.sendMessage(zmqcommon.DeleteKeygroup, ip, port, req)
	return
}

// SendUpdate sends a PUT message to the specified node.
func (c *Client) SendUpdate(ip net.IP, port int, kgname, kgid, value string) (err error) {
	req, err := json.Marshal(&zmqcommon.Request{
		Keygroup: kgname,
		ID:       kgid,
		Value:    value,
	})

	if err != nil {
		return
	}

	err = c.sendMessage(zmqcommon.PutItem, ip, port, req)
	return
}

// SendDelete sends the message to the specified node.
func (c *Client) SendDelete(ip net.IP, port int, kgname, kgid string) (err error) {
	req, err := json.Marshal(&zmqcommon.Request{
		Keygroup: kgname,
		ID:       kgid,
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

	err = cSender.SendMessageWithType(msType, msg)
	return
}

// Destroy the server.
func (c *Client) Destroy() {
	for _,sender := range c.senders {
		sender.Destroy()
	}
}