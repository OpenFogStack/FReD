package inthandler

import (
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
	panic("implement me")
}

// HandleDeleteKeygroup handles requests to the DeleteKeygroup endpoint of the internal interface.
func (h *handler) HandleDeleteKeygroup(k keygroup.Keygroup) error {
	panic("implement me")
}

// HandleRead handles requests to the Read endpoint of the internal interface.
func (h *handler) HandleRead(i data.Item) (data.Item, error) {
	panic("implement me")
}

// HandleUpdate handles requests to the Update endpoint of the internal interface.
func (h *handler) HandleUpdate(i data.Item) error {
	panic("implement me")
}

// HandleDelete handles requests to the Delete endpoint of the internal interface.
func (h *handler) HandleDelete(i data.Item) error {
	panic("implement me")
}
