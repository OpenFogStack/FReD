package zmqserver

import (
	"encoding/json"
	"github.com/zeromq/goczmq"
	"log"
)

// Controller is a ZMQ server that accepts incoming requests within the fred system.
type Controller struct {
	poller   *goczmq.Poller
	receiver *Receiver
	senders  map[string]Sender
	handler  MessageHandler
}

// Setup creates a controller and runs it.
func Setup(port string, id string, handler MessageHandler) (err error) {
	p, err := goczmq.NewPoller()

	if err != nil {
		return
	}

	r, err := NewReceiver(id, port)

	if err != nil {
		return
	}

	err = p.Add(r.GetSocket())

	if err != nil {
		return
	}

	controller := &Controller{
		poller:   p,
		receiver: r,
		senders:  make(map[string]Sender),
		handler:  handler,
	}

	//if anything happens, we should be destroying this controller...
	defer controller.Destroy()

	return pollForever(controller)
}

func pollForever(c *Controller) error {
	for {
		newMessageSocket := c.poller.Wait(10_000)
		if newMessageSocket == nil {
			//return errors.New("there was no new message for 10 seconds, shutting down")
			continue
		}

		// Receiver has got a new message
		request, err := newMessageSocket.RecvMessage()

		// TODO error handling
		if err != nil {
			log.Println(err)
		}

		// src = identity of socket from dealer
		src := string(request[0])
		recvd := request[1]
		msg := recvd[1:]

		// the first byte of answer tells whether an answer is expected
		answerType := recvd[0]

		// identity of sender can be either:
		// - our own receiver socket. This means another node wants to initiate a conservation with us
		// - another socket. This must be a sender socket (since we have no other sockets)
		//   This means we have at some time sent a request to another receiver socket and the other node has replied to our request
		//   Currently this only handles cases where a reply is expected (eg get request)
		//   But it should be expanded to also expect error replies and handle them
		// We dont want to send the answer in the current thread because that would block polling
		if newMessageSocket.Identity() == c.receiver.GetSocket().Identity() {
			switch answerType {
			case CreateKeygroup: // Create keygroup
				var req = &Request{}
				err = json.Unmarshal(msg, &req)
				go c.handler.HandleCreateKeygroup(req, src)
			case DeleteKeygroup: // Delete keygroup
				var req = &Request{}
				err = json.Unmarshal(msg, &req)
				go c.handler.HandleDeleteKeygroup(req, src)
			//case 0x12: // Get from Keygroup
			//	var req = &Request{}
			//	err = json.Unmarshal(msg, &req)
			//	go c.handler.HandleGetValueFromKeygroup(req, src)
			case PutItem: // Put into keygroup
				var req = &Request{}
				err = json.Unmarshal(msg, &req)
				go c.handler.HandlePutValueIntoKeygroup(req, src)
			case DeleteItem: // Delete in Keygroup
				var req = &Request{}
				err = json.Unmarshal(msg, &req)
				go c.handler.HandleDeleteFromKeygroup(req, src)
			}
		} else {
			// Not necessary because we only need eventual consistency
			//switch answerType {
			//case 0x12: // Answer to a get request received
			//	var res = &Response{}
			//	err = json.Unmarshal(msg, &res)
			//	log.Println("Yeah! My get request was answered! ")
			//	log.Println(res)
			//}
		}
	}
}

// Destroy the controller
func (c *Controller) Destroy() {
	c.poller.Destroy()
	c.receiver.Destroy()
}

// sendMessage to the specified IP
func (c *Controller) sendMessage(msType byte, ip string, msg []byte) (err error) {
	cSender, exists := c.senders[ip]
	if !exists {
		c.senders[ip] = *NewSender(ip, 5555)
		cSender = c.senders[ip]
		err = c.poller.Add(cSender.GetSocket())
	}

	if err != nil {
		return err
	}

	err = cSender.SendMessageWithType(msType, msg)
	return
}

// SendPut that answers a getRequest
func (c *Controller) SendPut(ip, kgname, kgid, value string) (err error) {
	req, err := json.Marshal(&Request{
		Keygroup: kgname,
		ID:       kgid,
		Value:    value,
	})

	if err != nil {
		return
	}

	err = c.sendMessage(PutItem, ip, req)
	return
}

// SendGetReply that answers a getRequest.
//func (c *Controller) SendGetReply(to string, kgname, id, value string) (err error) {
//	rep, err := json.Marshal(&Response{
//		Keygroup: kgname,
//		ID:       id,
//		Value:    value,
//	})
//
//	if err != nil {
//		return
//	}
//
//	err = c.receiver.ReplyTo(to, 0x12, rep)
//	return
//}
