package zmq

import (
	"fmt"

	"github.com/zeromq/goczmq"
)

// Receiver can receive zmqclient messages on a zmqclient socket and respond to them (if necessary).
type Receiver struct {
	socket *goczmq.Sock
}

// NewReceiver creates a zmqclient Receiver that listens on the given port.
func NewReceiver(id string, port int) (rec *Receiver, err error) {
	// Create a router socket and bind it.
	r, err := goczmq.NewRouter(fmt.Sprintf("tcp://0.0.0.0:%d", port))

	if err != nil {
		return nil, err
	}

	r.SetIdentity(fmt.Sprintf("receiver-%s", id))

	rec = &Receiver{socket: r}

	return
}

// GetSocket of receiver.
func (r *Receiver) GetSocket() (socket *goczmq.Sock) {
	return r.socket
}

// ReplyTo a sender that has sent a request that needs an answer.
func (r *Receiver) ReplyTo(id string, msType byte, data []byte) (err error) {
	err = r.socket.SendFrame([]byte(id), goczmq.FlagMore)

	if err != nil {
		return
	}

	err = r.socket.SendFrame([]byte{msType}, goczmq.FlagMore)

	if err != nil {
		return
	}

	// TODO if data is too big maybe split this up
	err = r.socket.SendFrame(data, goczmq.FlagNone)
	return
}

// Destroy the receiver.
func (r *Receiver) Destroy() {
	r.socket.Destroy()
}
