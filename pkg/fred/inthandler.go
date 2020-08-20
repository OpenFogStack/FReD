package fred

import (
	"github.com/go-errors/errors"
	"github.com/rs/zerolog/log"
)

type inthandler struct {
	s *storeService
	r *replicationService
	t *triggerService
}

// newInthandler creates a new handler for internal request (i.e. from peer nodes or the naming service).
func newInthandler(s *storeService, r *replicationService, t *triggerService) *inthandler {
	return &inthandler{
		s: s,
		r: r,
		t: t,
	}
}

// HandleCreateKeygroup handles requests to the CreateKeygroup endpoint of the internal interface.
func (h *inthandler) HandleCreateKeygroup(k Keygroup) error {
	if err := h.s.createKeygroup(k.Name); err != nil {
		log.Err(err).Msg(err.(*errors.Error).ErrorStack())
		return errors.Errorf("error creating keygroup")
	}

	if err := h.t.createKeygroup(k); err != nil {
		log.Err(err).Msg(err.(*errors.Error).ErrorStack())
		return errors.Errorf("error creating keygroup")
	}

	return nil
}

// HandleDeleteKeygroup handles requests to the DeleteKeygroup endpoint of the internal interface.
func (h *inthandler) HandleDeleteKeygroup(k Keygroup) error {
	if err := h.s.deleteKeygroup(k.Name); err != nil {
		log.Err(err).Msg(err.(*errors.Error).ErrorStack())
		return errors.Errorf("error deleting keygroup")
	}

	if err := h.r.deleteKeygroup(Keygroup{
		Name: k.Name,
	}); err != nil {
		log.Err(err).Msg(err.(*errors.Error).ErrorStack())
		return errors.Errorf("error deleting keygroup")
	}

	if err := h.t.deleteKeygroup(k.Name); err != nil {
		log.Err(err).Msg(err.(*errors.Error).ErrorStack())
		return errors.Errorf("error deleting keygroup")
	}

	return nil
}

// HandleUpdate handles requests to the Update endpoint of the internal interface.
func (h *inthandler) HandleUpdate(i Item) error {
	if err := h.s.update(i); err != nil {
		log.Err(err).Msg(err.(*errors.Error).ErrorStack())
		return errors.Errorf("error updating item")
	}

	if err := h.t.triggerUpdate(i); err != nil {
		log.Err(err).Msg(err.(*errors.Error).ErrorStack())
		return errors.Errorf("error updating item")
	}

	return nil
}

// HandleDelete handles requests to the Delete endpoint of the internal interface.
func (h *inthandler) HandleDelete(i Item) error {
	if err := h.s.delete(i.Keygroup, i.ID); err != nil {
		log.Err(err).Msg(err.(*errors.Error).ErrorStack())
		return errors.Errorf("error deleting item")
	}

	if err := h.t.triggerDelete(i); err != nil {
		log.Err(err).Msg(err.(*errors.Error).ErrorStack())
		return errors.Errorf("error deleting item")
	}

	return nil
}

func (h *inthandler) HandleAddReplica(k Keygroup, n Node) error {
	return h.r.addReplica(k, n, false)
}

func (h *inthandler) HandleRemoveReplica(k Keygroup, n Node) error {
	return h.r.removeReplica(k, n, false)
}
