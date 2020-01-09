package exthandler

import (
	"errors"
	"fmt"

	"github.com/rs/zerolog/log"

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
		log.Err(err).Msg("Exthandler cannot create keygroup with keygroup service")
		return err
	}

	if err := h.i.CreateKeygroup(data.Item{
		Keygroup: k.Name,
	}); err != nil {
		log.Err(err).Msg("Exthandler cannot create keygroup with data service")
		return err
	}

	if err := h.r.RelayCreateKeygroup(k); err != nil {
		log.Err(err).Msg("Exthandler cannot create keygroup with replication service")
		return err
	}

	return nil
}

// HandleDeleteKeygroup handles requests to the DeleteKeygroup endpoint of the client interface.
func (h *handler) HandleDeleteKeygroup(k keygroup.Keygroup) error {
	if err := h.k.Delete(keygroup.Keygroup{
		Name: k.Name,
	}); err != nil {
		log.Err(err).Msg("Exthandler cannot delete keygroup with keygroup service")
		return err
	}

	if err := h.i.DeleteKeygroup(data.Item{
		Keygroup: k.Name,
	}); err != nil {
		log.Err(err).Msg("Exthandler cannot delete keygroup with data service")
		return err
	}

	if err := h.r.RelayDeleteKeygroup(keygroup.Keygroup{
		Name: k.Name,
	}); err != nil {
		log.Err(err).Msg("Exthandler cannot delete keygroup with replication service")
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
		log.Err(err).Msg("Exthandler cannot relay update with data service")
		return err
	}

	if err := h.r.RelayUpdate(i); err != nil {
		log.Err(err).Msg("Exthandler cannot delete keygroup with replication service")
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
		log.Err(err).Msg("Exthandler cannot delete data item with data service")
		return err
	}

	if err := h.r.RelayDelete(i); err != nil {
		log.Err(err).Msg("Exthandler cannot delete data item with data service")
		return err
	}

	return nil
}

// HandleAddKeygroupReplica handles requests to the AddKeygroupReplica endpoint of the client interface.
func (h *handler) HandleAddKeygroupReplica(k keygroup.Keygroup, n replication.Node) error {
	if !h.k.Exists(keygroup.Keygroup{
		Name: k.Name,
	}) {
		return errors.New("exthandler: keygroup does not exist")
	}

	if err := h.r.AddReplica(k, n); err != nil {
		log.Err(err).Msg("Exthandler cannot add a new keygroup replica")
		return err
	}

	return nil
}

// HandleGetKeygroupReplica handles requests to the GetKeygroupReplica endpoint of the client interface.
func (h *handler) HandleGetKeygroupReplica(k keygroup.Keygroup) ([]replication.Node, error) {
	if !h.k.Exists(keygroup.Keygroup{
		Name: k.Name,
	}) {
		return nil, errors.New("exthandler: keygroup does not exist")
	}

	return h.r.GetReplica(k)
}

// HandleRemoveKeygroupReplica handles requests to the RemoveKeygroupReplica( endpoint of the client interface.
func (h *handler) HandleRemoveKeygroupReplica(k keygroup.Keygroup, n replication.Node) error {
	if !h.k.Exists(keygroup.Keygroup{
		Name: k.Name,
	}) {
		return errors.New("exthandler: keygroup does not exist")
	}

	if err := h.r.RemoveReplica(k, n); err != nil {
		return err
	}

	return nil
}

// HandleAddReplica handles requests to the AddReplica endpoint of the client interface.
func (h *handler) HandleAddReplica(n []replication.Node) error {
	e := make([]string, len(n))
	ec := 0

	for _, node := range n {
		if err := h.r.AddNode(node); err != nil {
			log.Err(err).Msgf("Exthandler can no add a new replica node. (node=%v)", node)
			e[ec] = fmt.Sprintf("%v", err)
			ec++
		}
	}

	if ec > 0 {
		return fmt.Errorf("exthandler: %v", e)
	}

	return nil
}

// HandleGetReplica handles requests to the GetReplica endpoint of the client interface.
func (h *handler) HandleGetReplica() ([]replication.Node, error) {
	return h.r.GetNodes()
}

// HandleRemoveReplica handles requests to the RemoveReplica endpoint of the client interface.
func (h *handler) HandleRemoveReplica(n replication.Node) error {
	return h.r.RemoveNode(n)
}