package zmqtointhandler

import (
	"github.com/rs/zerolog/log"

	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/fred"
	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/zmqcommon"
	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/zmqserver"
)

type zmqToInthandler struct {
	i fred.IntHandler
}

// New creates a new zmqToInthandler that uses the given handler.
// It translates calls from the ZMQ Layer to the Inthandler
func New(h fred.IntHandler) (l zmqserver.MessageHandler) {
	l = &zmqToInthandler{
		i: h,
	}

	return l
}

// HandleCreateKeygroup handles requests to the CreateKeygroup endpoint of the internal zmqclient interface.
func (l *zmqToInthandler) HandleCreateKeygroup(req *zmqcommon.KeygroupRequest, from string) {
	err := l.i.HandleCreateKeygroup(fred.Keygroup{Name: req.Keygroup}, req.Nodes)

	if err != nil {
		log.Err(err).Msg("error in HandleCreateKeygroup")
	}
	// TODO Error handling: send a reply message if necessary, the identity of the sender is in req.From
}

// HandlePutValueIntoKeygroup handles requests to the Update endpoint of the internal zmqclient interface.
func (l *zmqToInthandler) HandlePutValueIntoKeygroup(req *zmqcommon.DataRequest, from string) {
	err := l.i.HandleUpdate(fred.Item{
		Keygroup: fred.KeygroupName(req.Keygroup),
		ID:       req.ID,
		Data:     req.Value,
	})
	if err != nil {
		log.Err(err).Msg("error in HandlePutValueIntoKeygroup")
	}
}

// HandleDeleteFromKeygroup handles requests to the Delete endpoint of the internal zmqclient interface.
func (l *zmqToInthandler) HandleDeleteFromKeygroup(req *zmqcommon.DataRequest, from string) {
	err := l.i.HandleDelete(fred.Item{
		Keygroup: fred.KeygroupName(req.Keygroup),
		ID:       req.ID,
	})

	if err != nil {
		log.Err(err).Msg("error in HandleDeleteFromKeygroup")
	}
}

// HandleDeleteKeygroup handles requests to the DeleteKeygroup endpoint of the internal zmqclient interface.
func (l *zmqToInthandler) HandleDeleteKeygroup(req *zmqcommon.KeygroupRequest, from string) {
	err := l.i.HandleDeleteKeygroup(fred.Keygroup{Name: req.Keygroup})

	if err != nil {
		log.Err(err).Msg("error in HandleDeleteKeygroup")
	}
}

func (l *zmqToInthandler) HandleAddReplica(req *zmqcommon.ReplicationRequest, from string) {
	err := l.i.HandleAddReplica(fred.Keygroup{Name: req.Keygroup}, req.Node)

	if err != nil {
		log.Err(err).Msg("error in HandleAddReplica")
	}
}

func (l *zmqToInthandler) HandleRemoveReplica(req *zmqcommon.ReplicationRequest, from string) {
	err := l.i.HandleRemoveReplica(fred.Keygroup{Name: req.Keygroup}, req.Node)

	if err != nil {
		log.Err(err).Msg("error in HandleRemoveReplica")
	}
}
