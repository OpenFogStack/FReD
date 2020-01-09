package inthandler

import (
	"errors"

	"github.com/rs/zerolog/log"

	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/data"
	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/keygroup"
)

type handler struct {
	i data.Service
	k keygroup.Service
}

// New creates a new handler for internal request (i.e. from peer nodes or the naming service).
func New(i data.Service, k keygroup.Service) Handler {
	return &handler{
		i: i,
		k: k,
	}
}

// HandleCreateKeygroup handles requests to the CreateKeygroup endpoint of the internal interface.
func (h *handler) HandleCreateKeygroup(k keygroup.Keygroup) error {
	if err := h.k.Create(keygroup.Keygroup{
		Name: k.Name,
	}); err != nil {
		log.Err(err).Msg("Inthandler can not create keygroup with keygroup service")
		return err
	}

	if err := h.i.CreateKeygroup(data.Item{
		Keygroup: k.Name,
	}); err != nil {
		log.Err(err).Msg("Inthandler can not create keygroup with data service")
		return err
	}

	return nil
}

// HandleDeleteKeygroup handles requests to the DeleteKeygroup endpoint of the internal interface.
func (h *handler) HandleDeleteKeygroup(k keygroup.Keygroup) error {
	if err := h.k.Delete(keygroup.Keygroup{
		Name: k.Name,
	}); err != nil {
		log.Err(err).Msg("Inthandler can not delete keygroup with keygroup service")
		return err
	}

	if err := h.i.DeleteKeygroup(data.Item{
		Keygroup: k.Name,
	}); err != nil {
		log.Err(err).Msg("Inthandler can not delete keygroup with data service")
		return err
	}

	return nil
}

// HandleUpdate handles requests to the Update endpoint of the internal interface.
func (h *handler) HandleUpdate(i data.Item) error {
	if !h.k.Exists(keygroup.Keygroup{
		Name: i.Keygroup,
	}) {
		return errors.New("inthandler: keygroup does not exist")
	}

	return h.i.Update(i)
}

// HandleDelete handles requests to the Delete endpoint of the internal interface.
func (h *handler) HandleDelete(i data.Item) error {
	if !h.k.Exists(keygroup.Keygroup{
		Name: i.Keygroup,
	}) {
		return errors.New("inthandler: keygroup does not exist")
	}

	return h.i.Delete(i)
}
