package memoryzmq

import (
	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/commons"
	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/data"
	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/inthandler"
	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/keygroup"
	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/zmqcommon"
	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/zmqserver"
)

type localMemoryMessageHandler struct {
	i inthandler.Handler
}

// New creates a new localMemoryMessageHandler that uses the given handler.
func New(h inthandler.Handler) (l zmqserver.MessageHandler) {
	l = &localMemoryMessageHandler{
		i: h,
	}

	return l
}

// HandleCreateKeygroup handles requests to the CreateKeygroup endpoint of the internal zmqclient interface.
func (l *localMemoryMessageHandler) HandleCreateKeygroup(req *zmqcommon.Request, from string) {
	_ = l.i.HandleCreateKeygroup(keygroup.Keygroup{Name: commons.KeygroupName(req.Keygroup)})
	// TODO Error handling: send a reply message if necessary, the identity of the sender is in req.From
}

// HandlePutValueIntoKeygroup handles requests to the Update endpoint of the internal zmqclient interface.
func (l *localMemoryMessageHandler) HandlePutValueIntoKeygroup(req *zmqcommon.Request, from string) {
	_ = l.i.HandleUpdate(data.Item{
		Keygroup: commons.KeygroupName(req.Keygroup),
		ID:       req.ID,
	})
}

// HandleDeleteFromKeygroup handles requests to the Delete endpoint of the internal zmqclient interface.
func (l *localMemoryMessageHandler) HandleDeleteFromKeygroup(req *zmqcommon.Request, from string) {
	_ = l.i.HandleDelete(data.Item{
		Keygroup: commons.KeygroupName(req.Keygroup),
		ID:       req.ID,
	})
}

// HandleDeleteKeygroup handles requests to the DeleteKeygroup endpoint of the internal zmqclient interface.
func (l *localMemoryMessageHandler) HandleDeleteKeygroup(req *zmqcommon.Request, from string) {
	_ = l.i.HandleDeleteKeygroup(keygroup.Keygroup{Name: commons.KeygroupName(req.Keygroup)})
}

func (l *localMemoryMessageHandler) HandleAddNode(req *zmqcommon.ReplicationRequest, src string) {
	panic("implement me")
}

func (l *localMemoryMessageHandler) HandleRemoveNode(req *zmqcommon.ReplicationRequest, src string) {
	panic("implement me")
}

func (l *localMemoryMessageHandler) HandleAddReplica(req *zmqcommon.ReplicationRequest, src string) {
	panic("implement me")
}

func (l *localMemoryMessageHandler) HandleRemoveReplica(req *zmqcommon.ReplicationRequest, src string) {
	panic("implement me")
}