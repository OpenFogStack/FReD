package fred

import (
	"fmt"

	"github.com/rs/zerolog/log"
)

type inthandler struct {
	i *storeService
	r *replicationService
}

// newInthandler creates a new handler for internal request (i.e. from peer nodes or the naming service).
func newInthandler(i *storeService, r *replicationService) *inthandler {
	return &inthandler{
		i: i,
		r: r,
	}
}

// HandleCreateKeygroup handles requests to the CreateKeygroup endpoint of the internal interface.
func (h *inthandler) HandleCreateKeygroup(k Keygroup, nodes []Node) error {
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

	kg := Keygroup{
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
		return newError(StatusInternalError, fmt.Sprintf("exthandler: %v", e))
	}

	return nil
}

// HandleDeleteKeygroup handles requests to the DeleteKeygroup endpoint of the internal interface.
func (h *inthandler) HandleDeleteKeygroup(k Keygroup) error {
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
func (h *inthandler) HandleUpdate(i Item) error {
	return h.i.Update(i)
}

// HandleDelete handles requests to the Delete endpoint of the internal interface.
func (h *inthandler) HandleDelete(i Item) error {
	return h.i.Delete(i.Keygroup, i.ID)
}

func (h *inthandler) HandleAddReplica(k Keygroup, n Node) error {
	return h.replService.AddReplica(k, n, nil, false)
}

func (h *inthandler) HandleRemoveReplica(k Keygroup, n Node) error {
	return h.r.RemoveReplica(k, n, false)
}
