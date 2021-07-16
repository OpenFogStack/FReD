package fred

import (
	"github.com/go-errors/errors"
	"github.com/rs/zerolog/log"
)

type inthandler struct {
	s *storeService
	r *replicationService
	t *triggerService
	n NameService
}

// newInthandler creates a new handler for internal request (i.e. from peer nodes or the naming service).
func newInthandler(s *storeService, r *replicationService, t *triggerService, n NameService) *inthandler {
	return &inthandler{
		s: s,
		r: r,
		t: t,
		n: n,
	}
}

// HandleGet handles requests to the Get endpoint of the internal interface.
func (h *inthandler) HandleGet(i Item) (Item, error) {
	data, err := h.s.read(i.Keygroup, i.ID)

	if err != nil {
		log.Err(err).Msg(err.(*errors.Error).ErrorStack())
		return Item{}, errors.Errorf("error reading item")
	}

	return Item{
		Keygroup: i.Keygroup,
		ID:       i.ID,
		Val:      data,
	}, nil
}

// HandleGet handles requests to the Get endpoint of the internal interface.
func (h *inthandler) HandleGetAllItems(k Keygroup) ([]Item, error) {
	data, err := h.s.readAll(k.Name)

	if err != nil {
		log.Err(err).Msg(err.(*errors.Error).ErrorStack())
		return nil, errors.Errorf("error reading all keygroup items")
	}

	return data, nil
}

// HandleCreateKeygroup handles requests to the CreateKeygroup endpoint of the internal interface.
func (h *inthandler) HandleCreateKeygroup(k Keygroup) error {
	if err := h.s.createKeygroup(k.Name); err != nil {
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

	return nil
}

// HandleUpdate handles requests to the Update endpoint of the internal interface.
func (h *inthandler) HandleUpdate(i Item) error {

	expiry, err := h.n.GetExpiry(i.Keygroup)

	if err != nil {
		return err
	}

	if err := h.s.update(i, false, expiry); err != nil {
		log.Err(err).Msg(err.(*errors.Error).ErrorStack())
		return errors.Errorf("error updating item")
	}

	if err := h.t.triggerUpdate(i); err != nil {
		log.Err(err).Msg(err.(*errors.Error).ErrorStack())
		return errors.Errorf("error updating item")
	}

	return nil
}

// HandleAppend handles requests to the Append endpoint of the internal interface.
func (h *inthandler) HandleAppend(i Item) error {

	expiry, err := h.n.GetExpiry(i.Keygroup)

	if err != nil {
		return err
	}

	if err := h.s.update(i, true, expiry); err != nil {
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
