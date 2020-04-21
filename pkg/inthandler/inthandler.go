package inthandler

import (
	"fmt"

	"github.com/rs/zerolog/log"

	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/data"
	errors "gitlab.tu-berlin.de/mcc-fred/fred/pkg/errors"
	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/keygroup"
	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/replication"
)

type handler struct {
	dataService data.Service
	replService replication.Service
}

// New creates a new handler for internal request (dataService.e. from peer nodes or the naming service).
func New(i data.Service, r replication.Service) Handler {
	return &handler{
		dataService: i,
		replService: r,
	}
}

// HandleCreateKeygroup handles requests to the CreateKeygroup endpoint of the internal interface.
func (h *handler) HandleCreateKeygroup(k keygroup.Keygroup, nodes []replication.Node) error {
	if err := h.dataService.CreateKeygroup(k.Name); err != nil {
		log.Err(err).Msg("Inthandler cannot create keygroup with data service")
		return err
	}

	if err := h.replService.CreateKeygroup(keygroup.Keygroup{
		Name: k.Name,
	}); err != nil {
		log.Err(err).Msg("Inthandler cannot create keygroup with replication service")
		return err
	}

	kg := keygroup.Keygroup{
		Name: k.Name,
	}

	e := make([]string, len(nodes))
	ec := 0

	for _, node := range nodes {
		if err := h.replService.AddReplica(kg, node, nil, false); err != nil {
			log.Err(err).Msgf("inthandler: cannot remove node %#v)", node)
			e[ec] = fmt.Sprintf("%v", err)
			ec++
		}
	}

	if ec > 0 {
		return errors.New(errors.StatusInternalError, fmt.Sprintf("exthandler: %v", e))
	}

	return nil
}

// HandleDeleteKeygroup handles requests to the DeleteKeygroup endpoint of the internal interface.
func (h *handler) HandleDeleteKeygroup(k keygroup.Keygroup) error {
	if err := h.dataService.DeleteKeygroup(k.Name); err != nil {
		log.Err(err).Msg("Inthandler cannot delete keygroup with data service")
		return err
	}

	if err := h.replService.DeleteKeygroup(keygroup.Keygroup{
		Name: k.Name,
	}); err != nil {
		log.Err(err).Msg("Inthandler cannot delete keygroup with replication service")
		return err
	}

	return nil
}

// HandleUpdate handles requests to the Update endpoint of the internal interface.
func (h *handler) HandleUpdate(i data.Item) error {
	return h.dataService.Update(i)
}

// HandleDelete handles requests to the Delete endpoint of the internal interface.
func (h *handler) HandleDelete(i data.Item) error {
	return h.dataService.Delete(i.Keygroup, i.ID)
}

func (h *handler) HandleAddReplica(k keygroup.Keygroup, n replication.Node) error {
	return h.replService.AddReplica(k, n, nil, false)
}

func (h *handler) HandleRemoveReplica(k keygroup.Keygroup, n replication.Node) error {
	return h.replService.RemoveReplica(k, n, false)
}
