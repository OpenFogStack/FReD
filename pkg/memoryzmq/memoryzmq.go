package memoryzmq

import (
	"github.com/rs/zerolog/log"

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
func (l *localMemoryMessageHandler) HandleCreateKeygroup(req *zmqcommon.KeygroupRequest, from string) {
	err := l.i.HandleCreateKeygroup(keygroup.Keygroup{Name: req.Keygroup}, req.Nodes)

	log.Err(err).Msg("")
	// TODO Error handling: send a reply message if necessary, the identity of the sender is in req.From
}

// HandlePutValueIntoKeygroup handles requests to the Update endpoint of the internal zmqclient interface.
func (l *localMemoryMessageHandler) HandlePutValueIntoKeygroup(req *zmqcommon.DataRequest, from string) {
	err := l.i.HandleUpdate(data.Item{
		Keygroup: commons.KeygroupName(req.Keygroup),
		ID:       req.ID,
		Data:     req.Value,
	})

	log.Err(err).Msg("")
}

// HandleDeleteFromKeygroup handles requests to the Delete endpoint of the internal zmqclient interface.
func (l *localMemoryMessageHandler) HandleDeleteFromKeygroup(req *zmqcommon.DataRequest, from string) {
	err := l.i.HandleDelete(data.Item{
		Keygroup: commons.KeygroupName(req.Keygroup),
		ID:       req.ID,
	})

	log.Err(err).Msg("")
}

// HandleDeleteKeygroup handles requests to the DeleteKeygroup endpoint of the internal zmqclient interface.
func (l *localMemoryMessageHandler) HandleDeleteKeygroup(req *zmqcommon.KeygroupRequest, from string) {
	err := l.i.HandleDeleteKeygroup(keygroup.Keygroup{Name: req.Keygroup})

	log.Err(err).Msg("")
}

func (l *localMemoryMessageHandler) HandleAddNode(req *zmqcommon.ReplicationRequest, from string) {
	err := l.i.HandleAddNode(req.Node)

	log.Err(err).Msg("")
}

func (l *localMemoryMessageHandler) HandleRemoveNode(req *zmqcommon.ReplicationRequest, from string) {
	err := l.i.HandleRemoveNode(req.Node)

	log.Err(err).Msg("")
}

func (l *localMemoryMessageHandler) HandleAddReplica(req *zmqcommon.ReplicationRequest, from string) {
	err := l.i.HandleAddReplica(keygroup.Keygroup{Name: req.Keygroup}, req.Node)

	log.Err(err).Msg("")
}

func (l *localMemoryMessageHandler) HandleRemoveReplica(req *zmqcommon.ReplicationRequest, from string) {
	err := l.i.HandleRemoveReplica(keygroup.Keygroup{Name: req.Keygroup}, req.Node)

	log.Err(err).Msg("")
}

func (l *localMemoryMessageHandler) HandleIntroduction(req *zmqcommon.IntroductionRequest, src string) {
	err := l.i.HandleIntroduction(req.Self, req.Other, req.Node)

	log.Err(err).Msg("")
}

func (l *localMemoryMessageHandler) HandleDetroduction(req *zmqcommon.IntroductionRequest, src string) {
	err := l.i.HandleDetroduction()

	log.Err(err).Msg("")
}