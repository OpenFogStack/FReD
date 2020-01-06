package exthandler

import (
	"errors"

	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/data"
	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/keygroup"
	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/replication"
)

type handler struct {
	i data.Service
	k keygroup.Service
	r replication.Service
}

// New creates a new handler for external request (i.e. from clients).
func New(i data.Service, k keygroup.Service, r replication.Service) Handler {
	return &handler{
		i: i,
		k: k,
		r: r,
	}
}

// HandleCreateKeygroup handles requests to the CreateKeygroup endpoint of the client interface.
func (h *handler) HandleCreateKeygroup(k keygroup.Keygroup) error {
	if err := h.k.Create(keygroup.Keygroup{
		Name: k.Name,
	}); err != nil {
		return err
	}

	if err := h.i.CreateKeygroup(data.Item{
		Keygroup: k.Name,
	}); err != nil {
		return err
	}

	return nil
}

// HandleDeleteKeygroup handles requests to the DeleteKeygroup endpoint of the client interface.
func (h *handler) HandleDeleteKeygroup(k keygroup.Keygroup) error {
	if err := h.k.Delete(keygroup.Keygroup{
		Name: k.Name,
	}); err != nil {
		return err
	}

	if err := h.i.DeleteKeygroup(data.Item{
		Keygroup: k.Name,
	}); err != nil {
		return err
	}

	if err := h.r.RelayDeleteKeygroup(keygroup.Keygroup{
		Name: k.Name,
	}); err != nil {
		return err
	}

	return nil
}

// HandleRead handles requests to the Read endpoint of the client interface.
func (h *handler) HandleRead(i data.Item) (data.Item, error) {
	if !h.k.Exists(keygroup.Keygroup{
		Name: i.Keygroup,
	}) {
		return i, errors.New("exthandler: keygroup does not exist")
	}

	return h.i.Read(i)
}

// HandleUpdate handles requests to the Update endpoint of the client interface.
func (h *handler) HandleUpdate(i data.Item) error {
	if !h.k.Exists(keygroup.Keygroup{
		Name: i.Keygroup,
	}) {
		return errors.New("exthandler: keygroup does not exist")
	}

	if err := h.i.Update(i); err != nil {
		return err
	}

	if err := h.r.RelayUpdate(i); err != nil {
		return err
	}

	return nil
}

// HandleDelete handles requests to the Delete endpoint of the client interface.
func (h *handler) HandleDelete(i data.Item) error {
	if !h.k.Exists(keygroup.Keygroup{
		Name: i.Keygroup,
	}) {
		return errors.New("exthandler: keygroup does not exist")
	}

	if err := h.i.Delete(i); err != nil {
		return err
	}

	if err := h.r.RelayDelete(i); err != nil {
		return err
	}

	return nil
}
