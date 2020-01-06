package zmqserver

import (
	"encoding/json"

	"github.com/rs/zerolog/log"
	"github.com/zeromq/goczmq"

	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/zmqcommon"
)

// Server is a ZMQ server that accepts incoming requests within the fred system.
type Server struct {
	poller          *goczmq.Poller
	receiver        *Receiver
	handler         MessageHandler
	continueRunning bool
}

// Setup creates a controller and runs it.
func Setup(port int, id string, handler MessageHandler) (server *Server, err error) {
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

	server = &Server{
		poller:          p,
		receiver:        r,
		handler:         handler,
		continueRunning: true,
	}

	go pollForever(server)

	return
}

func pollForever(c *Server) error {
	for c.continueRunning {
		newMessageSocket := c.poller.Wait(10_000)
		if newMessageSocket == nil {
			//return errors.New("there was no new message for 10 seconds, shutting down")
			continue
		}

		// Receiver has got a new message
		request, err := newMessageSocket.RecvMessage()

		// TODO error handling
		if err != nil {
			log.Err(err).Msg("pollForever has received an error during its receive")
		}

		// src = identity of socket from dealer
		src := string(request[0])
		// the first byte of answer tells whether an answer is expected
		answerType := request[1][0]
		// the rest is the message we got
		msg := request[2]

		// identity of sender can be either:
		// - our own receiver socket. This means another node wants to initiate a conservation with us
		// - another socket. This must be a sender socket (since we have no other sockets)
		//   This means we have at some time sent a request to another receiver socket and the other node has replied to our request
		//   Currently this never happens because we dont send answers to requests
		//   But it should be expanded to handle error replies ect.
		// We dont want to send the answer in the current thread because that would block polling
		if newMessageSocket.Identity() == c.receiver.GetSocket().Identity() {
			switch answerType {
			case zmqcommon.CreateKeygroup: // Create keygroup
				var req = &zmqcommon.Request{}
				err = json.Unmarshal(msg, &req)
				go c.handler.HandleCreateKeygroup(req, src)
			case zmqcommon.DeleteKeygroup: // Delete keygroup
				var req = &zmqcommon.Request{}
				err = json.Unmarshal(msg, &req)
				go c.handler.HandleDeleteKeygroup(req, src)
			//case 0x12: // Get from Keygroup
			//	var req = &Request{}
			//	err = json.Unmarshal(msg, &req)
			//	go c.handler.HandleGetValueFromKeygroup(req, src)
			case zmqcommon.PutItem: // Put into keygroup
				var req = &zmqcommon.Request{}
				err = json.Unmarshal(msg, &req)
				go c.handler.HandlePutValueIntoKeygroup(req, src)
			case zmqcommon.DeleteItem: // Delete in Keygroup
				var req = &zmqcommon.Request{}
				err = json.Unmarshal(msg, &req)
				go c.handler.HandleDeleteFromKeygroup(req, src)
			}
		} else {
			// Not necessary because we only need eventual consistency, so we dont receive answers to our questions
			//switch answerType {
			//case 0x12: // Answer to a get request received
			//	var res = &Response{}
			//	err = json.Unmarshal(msg, &res)
			//	log.Println("Yeah! My get request was answered! ")
			//	log.Println(res)
			//}
		}
	}
	c.destroy()
	return nil
}

// Destroy the server
func (c *Server) destroy() {
	c.poller.Destroy()
	c.receiver.Destroy()
}

// Shutdown the server
func (c *Server) Shutdown() {
	c.continueRunning = false
}