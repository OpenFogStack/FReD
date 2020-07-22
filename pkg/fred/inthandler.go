package fred

import (
	"github.com/go-errors/errors"
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
		log.Err(err).Msg(err.(*errors.Error).ErrorStack())
		return errors.Errorf("error creating keygroup")
	}

	if err := h.replService.CreateKeygroup(keygroup.Keygroup{
		Name: k.Name,
	}); err != nil {
		log.Err(err).Msg(err.(*errors.Error).ErrorStack())

		return errors.Errorf("error creating keygroup")
	}

	kg := Keygroup{
		Name: k.Name,
	}

	var ec int

	for _, node := range nodes {
		if err := h.r.AddReplica(kg, node, nil, false); err != nil {
			log.Err(err).Msg(err.(*errors.Error).ErrorStack())
			ec++
		}
	}

	if ec > 0 {
		return errors.Errorf("error updating %d nodes", ec)
	}

	return nil
}

// HandleDeleteKeygroup handles requests to the DeleteKeygroup endpoint of the internal interface.
func (h *inthandler) HandleDeleteKeygroup(k Keygroup) error {
	if err := h.dataService.DeleteKeygroup(k.Name); err != nil {
		log.Err(err).Msg(err.(*errors.Error).ErrorStack())
		return errors.Errorf("error deleting keygroup")
	}

	if err := h.replService.DeleteKeygroup(keygroup.Keygroup{
		Name: k.Name,
	}); err != nil {
		log.Err(err).Msg(err.(*errors.Error).ErrorStack())
		return errors.Errorf("error deleting keygroup")
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
