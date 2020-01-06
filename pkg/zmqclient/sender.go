package zmqclient

import (
	"fmt"
	"github.com/zeromq/goczmq"
	"github.com/rs/zerolog/log"
)

// Sender can send zmqclient messages to a zmqclient socket (both synchronously and asynchronously).
type Sender struct {
	socket *goczmq.Sock
}

// NewSender creates a zmqclient Sender to send messages to the specified ip and port.
func NewSender(ip string, port int) (sen *Sender) {
	// Create a dealer socket and connect it to the router.
	dealer, err := goczmq.NewDealer(fmt.Sprintf("tcp://%s:%d", ip, port))
	if err != nil {
		log.Error().Err(err).Msg("Can not create ZMQ Dealer")
	}
	log.Printf("Sender has created a dealer to tcp://%s:%d\n", ip, port)
	sen = &Sender{dealer}
	return
}

// Destroy the sender
func (s *Sender) Destroy() {
	s.socket.Destroy()
}

// GetSocket ...
func (s *Sender) GetSocket() (socket *goczmq.Sock) {
	return s.socket
}

// SendBytes sends a message composed of bytes.
func (s *Sender) SendBytes(data []byte) (err error) {
	// TODO if the message is too big (maybe 500mb and more?) it should be sent as a multipart message
	//log.Printf("Sending:: %#v", data)
	err = s.socket.SendFrame(data, goczmq.FlagNone)
	return
}

// SendMessageWithType ...
func (s *Sender) SendMessageWithType(msType byte, data []byte) (err error) {
	// TODO some sanity checks, whatevs
	err = s.socket.SendFrame([]byte{msType}, goczmq.FlagMore)
	if err != nil {
		return
	}
	err = s.SendBytes(data)
	return
}
