package fred

import (
	"github.com/go-errors/errors"
	"github.com/rs/zerolog/log"
)

type IntHandler struct {
	s *storeService
	r *replicationService
	t *triggerService
	n NameService
}

// newInthandler creates a new handler for internal request (i.e. from peer nodes or the naming service).
func newInthandler(s *storeService, r *replicationService, t *triggerService, n NameService) *IntHandler {
	return &IntHandler{
		s: s,
		r: r,
		t: t,
		n: n,
	}
}

// HandleGet handles requests to the Get endpoint of the internal interface.
func (h *IntHandler) HandleGet(i Item) ([]Item, error) {
	data, err := h.s.read(i.Keygroup, i.ID)

	if err != nil {
		log.Err(err).Msg(err.(*errors.Error).ErrorStack())
		return nil, errors.Errorf("error reading item")
	}

	return data, nil
}

// HandleGetAllItems handles requests to the Get endpoint of the internal interface.
func (h *IntHandler) HandleGetAllItems(k Keygroup) ([]Item, error) {
	data, err := h.s.readAll(k.Name)

	if err != nil {
		log.Err(err).Msg(err.(*errors.Error).ErrorStack())
		return nil, errors.Errorf("error reading all keygroup items")
	}

	return data, nil
}

// HandleCreateKeygroup handles requests to the CreateKeygroup endpoint of the internal interface.
func (h *IntHandler) HandleCreateKeygroup(k Keygroup) error {
	if err := h.s.createKeygroup(k.Name); err != nil {
		log.Err(err).Msg(err.(*errors.Error).ErrorStack())
		return errors.Errorf("error creating keygroup")
	}

	return nil
}

// HandleDeleteKeygroup handles requests to the DeleteKeygroup endpoint of the internal interface.
func (h *IntHandler) HandleDeleteKeygroup(k Keygroup) error {
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
func (h *IntHandler) HandleUpdate(i Item) error {

	expiry, err := h.n.GetExpiry(i.Keygroup)

	if err != nil {
		return err
	}

	if err := h.s.addVersion(i, i.Version, expiry); err != nil {
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
func (h *IntHandler) HandleAppend(i Item) error {

	expiry, err := h.n.GetExpiry(i.Keygroup)

	if err != nil {
		return err
	}

	if err := h.s.append(i, expiry); err != nil {
		log.Err(err).Msg(err.(*errors.Error).ErrorStack())
		return errors.Errorf("error updating item")
	}

	if err := h.t.triggerUpdate(i); err != nil {
		log.Err(err).Msg(err.(*errors.Error).ErrorStack())
		return errors.Errorf("error updating item")
	}

	return nil
}
