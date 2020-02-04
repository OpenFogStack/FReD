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
	i data.Service
	k keygroup.Service
	r replication.Service
}

// New creates a new handler for internal request (i.e. from peer nodes or the naming service).
func New(i data.Service, k keygroup.Service, r replication.Service) Handler {
	return &handler{
		i: i,
		k: k,
		r: r,
	}
}

// HandleCreateKeygroup handles requests to the CreateKeygroup endpoint of the internal interface.
func (h *handler) HandleCreateKeygroup(k keygroup.Keygroup, nodes []replication.Node) error {
	if err := h.k.Create(keygroup.Keygroup{
		Name: k.Name,
	}); err != nil {
		log.Err(err).Msg("Inthandler cannot create keygroup with keygroup service")
		return err
	}

	if err := h.i.CreateKeygroup(data.Item{
		Keygroup: k.Name,
	}); err != nil {
		log.Err(err).Msg("Inthandler cannot create keygroup with data service")
		return err
	}

	if err := h.r.CreateKeygroup(keygroup.Keygroup{
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
		if err := h.r.AddReplica(kg, node, nil, false); err != nil {
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
	if err := h.k.Delete(keygroup.Keygroup{
		Name: k.Name,
	}); err != nil {
		log.Err(err).Msg("Inthandler cannot delete keygroup with keygroup service")
		return err
	}

	if err := h.i.DeleteKeygroup(data.Item{
		Keygroup: k.Name,
	}); err != nil {
		log.Err(err).Msg("Inthandler cannot delete keygroup with data service")
		return err
	}

	if err := h.r.DeleteKeygroup(keygroup.Keygroup{
		Name: k.Name,
	}); err != nil {
		log.Err(err).Msg("Inthandler cannot delete keygroup with replication service")
		return err
	}

	return nil
}

// HandleUpdate handles requests to the Update endpoint of the internal interface.
func (h *handler) HandleUpdate(i data.Item) error {
	if !h.k.Exists(keygroup.Keygroup{
		Name: i.Keygroup,
	}) {
		return errors.New(errors.StatusNotFound, "inthandler: keygroup does not exist")
	}

	return h.i.Update(i)
}

// HandleDelete handles requests to the Delete endpoint of the internal interface.
func (h *handler) HandleDelete(i data.Item) error {
	if !h.k.Exists(keygroup.Keygroup{
		Name: i.Keygroup,
	}) {
		return errors.New(errors.StatusNotFound, "inthandler: keygroup does not exist")
	}

	return h.i.Delete(i)
}

func (h *handler) HandleAddReplica(k keygroup.Keygroup, n replication.Node) error {
	return h.r.AddReplica(k, n, nil,false)
}

func (h *handler) HandleRemoveReplica(k keygroup.Keygroup, n replication.Node) error {
	return h.r.RemoveReplica(k, n, false)
}

func (h *handler) HandleAddNode(n replication.Node) error {
	return h.r.AddNode(n, false)
}

func (h *handler) HandleRemoveNode(n replication.Node) error {
	return h.r.RemoveNode(n, false)
}

func (h *handler) HandleIntroduction(introducer replication.Node, self replication.Node, node []replication.Node) error {
	log.Debug().Msgf("HandleIntroduction from inthandler called: introducer=%#v, self=%#v, nodes=%#v", introducer, self, node)
	err := h.r.Seed(self)

	if err != nil {
		return err
	}

	// +1 to account for the introducer
	e := make([]string, len(node)+1)
	ec := 0

	// add introducer (the node that sent the introduction) to our list of known nodes
	if err := h.r.AddNode(introducer, false); err != nil {
		log.Err(err).Msgf("inthandler: cannot add a new node %#v)", introducer)
		e[ec] = fmt.Sprintf("%v", err)
		ec++
	}

	// add the list of nodes that we were told about by the introducer
	for _, node := range node {
		if err := h.r.AddNode(node, false); err != nil {
			log.Err(err).Msgf("inthandler: cannot add a new node %#v)", node)
			e[ec] = fmt.Sprintf("%v", err)
			ec++
		}
	}

	if ec > 0 {
		return errors.New(errors.StatusInternalError, fmt.Sprintf("exthandler: %v", e))
	}

	return nil
}

func (h *handler) HandleDetroduction() error {
	err := h.r.Unseed()

	if err != nil {
		return err
	}

	nodes, err := h.r.GetNodes()

	if err != nil {
		return err
	}

	e := make([]string, len(nodes))
	ec := 0

	for _, node := range nodes {
		if err := h.r.RemoveNode(node, false); err != nil {
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
