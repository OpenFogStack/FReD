package zmqclient

import (
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/zeromq/goczmq"

	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/replication"
)

// Sender can send zmqclient messages to a zmqclient socket (both synchronously and asynchronously).
type Sender struct {
	socket *goczmq.Sock
}

// NewSender creates a zmqclient Sender to send messages to the specified addr and port.
func NewSender(addr replication.Address, port int) (sen *Sender) {
	// Create a dealer socket and connect it to the router.
	dealer, err := goczmq.NewDealer(fmt.Sprintf("tcp://%s:%d", addr.Addr, port))
	if err != nil {
		log.Error().Err(err).Msg("cannot create ZMQ Dealer")
	}
	log.Printf("Sender has created a dealer to tcp://%s:%d\n", addr.Addr, port)
	sen = &Sender{dealer}
	return
}

// Destroy the sender
func (s *Sender) Destroy() {
	s.socket.Destroy()
}

// GetSocket TODO
func (s *Sender) GetSocket() (socket *goczmq.Sock) {
	return s.socket
}

// SendBytes sends a message composed of bytes.
func (s *Sender) SendBytes(data []byte) (err error) {
	// TODO if the message is too big (maybe 500mb and more?) it should be sent as a multipart message
	// log.Printf("Sending:: %#v", data)
	err = s.socket.SendFrame(data, goczmq.FlagNone)
	return
}

// SendMessageWithType TODO
func (s *Sender) SendMessageWithType(msType byte, data []byte) (err error) {
	// TODO some sanity checks, whatevs
	log.Debug().Msgf("Sending message type %#v to %s", msType, s.socket.Identity())
	err = s.socket.SendFrame([]byte{msType}, goczmq.FlagMore)
	if err != nil {
		return
	}
	err = s.SendBytes(data)
	return
}
