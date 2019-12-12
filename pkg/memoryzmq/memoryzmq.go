package memoryzmq

import (
	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/data"
	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/keygroup"
	"log"

	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/inthandler"
	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/zmqserver"
)

type localMemoryMessageHandler struct {
	i inthandler.Handler
}

// New creates a new localMemoryMessageHandler that uses the given handler.
func New(h inthandler.Handler) (l zmqserver.MessageHandler) {
	l = &localMemoryMessageHandler{
		i:h,
	}

	return l
}

// HandleCreateKeygroup handles requests to the CreateKeygroup endpoint of the internal zmq interface.
func (l *localMemoryMessageHandler) HandleCreateKeygroup(req *zmqserver.Request, from string) {
	_ = l.i.HandleCreateKeygroup(keygroup.Keygroup{Name: req.Keygroup})
	// TODO Error handling: send a reply message if necessary, the identity of the sender is in req.From
}

// HandleGetValueFromKeygroup handles requests to the Read endpoint of the internal zmq interface.
func (l *localMemoryMessageHandler) HandleGetValueFromKeygroup(req *zmqserver.Request, from string) {
	_, err := l.i.HandleRead(data.Item{
		Keygroup: req.Keygroup,
		ID:       req.ID,
	})

	if err == nil {
		log.Fatal("Get from App failed! ", err)
		return
	}

	//C.SendGetAnswer(from, req.Keygroup, val)
}

// HandlePutValueIntoKeygroup handles requests to the Update endpoint of the internal zmq interface.
func (l *localMemoryMessageHandler) HandlePutValueIntoKeygroup(req *zmqserver.Request, from string) {
	_ = l.i.HandleUpdate(data.Item{
		Keygroup: req.Keygroup,
		ID:       req.ID,
	})
}

// HandleDeleteFromKeygroup handles requests to the Delete endpoint of the internal zmq interface.
func (l *localMemoryMessageHandler) HandleDeleteFromKeygroup(req *zmqserver.Request, from string) {
	_ = l.i.HandleDelete(data.Item{
		Keygroup: req.Keygroup,
		ID:       req.ID,
	})
}

// HandleDeleteKeygroup handles requests to the DeleteKeygroup endpoint of the internal zmq interface.
func (l *localMemoryMessageHandler) HandleDeleteKeygroup(req *zmqserver.Request, from string) {
	_ = l.i.HandleDeleteKeygroup(keygroup.Keygroup{Name: req.Keygroup})
}
