package fred

import (
	"github.com/go-errors/errors"
	"github.com/rs/zerolog/log"
)

type exthandler struct {
	s *storeService
	r *replicationService
	t *triggerService
}

// newExthandler creates a new handler for external request (i.e. from clients).
func newExthandler(s *storeService, r *replicationService, t *triggerService) *exthandler {
	return &exthandler{
		s: s,
		r: r,
		t: t,
	}
}

// HandleCreateKeygroup handles requests to the CreateKeygroup endpoint of the client interface.
func (h *exthandler) HandleCreateKeygroup(k Keygroup) error {
	if err := h.s.createKeygroup(k.Name); err != nil {
		log.Err(err).Msg(err.(*errors.Error).ErrorStack())
		return errors.Errorf("error creating keygroup")
	}

	if err := h.t.createKeygroup(k); err != nil {
		log.Err(err).Msg(err.(*errors.Error).ErrorStack())
		return errors.Errorf("error creating keygroup")
	}

	if err := h.r.createKeygroup(k); err != nil {
		log.Err(err).Msg(err.(*errors.Error).ErrorStack())
		return errors.Errorf("error creating keygroup")
	}

	return nil
}

// HandleDeleteKeygroup handles requests to the DeleteKeygroup endpoint of the client interface.
func (h *exthandler) HandleDeleteKeygroup(k Keygroup) error {
	if err := h.s.deleteKeygroup(k.Name); err != nil {
		log.Err(err).Msg(err.(*errors.Error).ErrorStack())
		return errors.Errorf("error deleting keygroup")
	}

	if err := h.t.deleteKeygroup(k.Name); err != nil {
		log.Err(err).Msg(err.(*errors.Error).ErrorStack())
		return errors.Errorf("error deleting keygroup")
	}

	if err := h.r.relayDeleteKeygroup(Keygroup{
		Name: k.Name,
	}); err != nil {
		log.Err(err).Msg(err.(*errors.Error).ErrorStack())
		return errors.Errorf("error deleting keygroup")
	}

	return nil
}

// HandleRead handles requests to the Read endpoint of the client interface.
func (h *exthandler) HandleRead(i Item) (Item, error) {
	result, err := h.s.read(i.Keygroup, i.ID)

	if err != nil {
		log.Error().Msgf("Error in AddReplica is: %#v", err)
		log.Err(err).Msg(err.(*errors.Error).ErrorStack())
		return i, errors.Errorf("error reading item %s from keygroup %s", i.ID, i.Keygroup)
	}

	// KeygroupStore result in passed object to return an Item and not only the result string
	i.Val = result
	return i, nil
}

// HandleUpdate handles requests to the Update endpoint of the client interface.
func (h *exthandler) HandleUpdate(i Item) error {
	if err := h.s.update(i); err != nil {
		log.Printf("%#v", err)
		log.Err(err).Msg(err.(*errors.Error).ErrorStack())
		return errors.Errorf("error updating item")
	}

	if err := h.r.relayUpdate(i); err != nil {
		log.Err(err).Msg(err.(*errors.Error).ErrorStack())
		return errors.Errorf("error updating item")
	}

	if err := h.t.triggerUpdate(i); err != nil {
		log.Err(err).Msg(err.(*errors.Error).ErrorStack())
		return errors.Errorf("error updating item")
	}

	return nil
}

// HandleDelete handles requests to the Delete endpoint of the client interface.
func (h *exthandler) HandleDelete(i Item) error {
	if err := h.s.delete(i.Keygroup, i.ID); err != nil {
		log.Err(err).Msg(err.(*errors.Error).ErrorStack())
		return errors.Errorf("error deleting item")
	}

	if err := h.r.relayDelete(i); err != nil {
		log.Err(err).Msg(err.(*errors.Error).ErrorStack())
		return errors.Errorf("error deleting item")
	}

	if err := h.t.triggerDelete(i); err != nil {
		log.Err(err).Msg(err.(*errors.Error).ErrorStack())
		return errors.Errorf("error deleting item")
	}

	return nil
}

// HandleAddReplica handles requests to the AddKeygroupReplica endpoint of the client interface.
func (h *exthandler) HandleAddReplica(k Keygroup, n Node) error {
	if err := h.r.addReplica(k, n, true); err != nil {
		log.Error().Msgf("Error in AddReplica is: %#v", err)
		log.Err(err).Msg(err.(*errors.Error).ErrorStack())
		return errors.Errorf("error adding replica")
	}

	return nil
}

// HandleGetKeygroupReplica handles requests to the GetKeygroupReplica endpoint of the client interface.
func (h *exthandler) HandleGetKeygroupReplica(k Keygroup) ([]Node, error) {
	return h.r.getReplica(k)
}

// HandleRemoveReplica handles requests to the RemoveKeygroupReplica endpoint of the client interface.
func (h *exthandler) HandleRemoveReplica(k Keygroup, n Node) error {
	if err := h.r.removeReplica(k, n, true); err != nil {
		log.Err(err).Msg(err.(*errors.Error).ErrorStack())
		return errors.Errorf("error removing replica")
	}

	return nil
}

// HandleAddReplica handles requests to the AddKeygroupTrigger endpoint of the client interface.
func (h *exthandler) HandleAddTriggers(k Keygroup, t Trigger) error {
	if err := h.t.addTrigger(k, t); err != nil {
		log.Error().Msgf("Error in AddTrigger is: %#v", err)
		log.Err(err).Msg(err.(*errors.Error).ErrorStack())
		return errors.Errorf("error adding replica")
	}
	return nil
}

// HandleGetKeygroupTriggers handles requests to the GetKeygroupTrigger endpoint of the client interface.
func (h *exthandler) HandleGetKeygroupTriggers(k Keygroup) ([]Trigger, error) {
	return h.t.getTrigger(k)
}

// HandleRemoveTrigger handles requests to the RemoveKeygroupTrigger endpoint of the client interface.
func (h *exthandler) HandleRemoveTrigger(k Keygroup, t Trigger) error {
	if err := h.t.removeTrigger(k, t); err != nil {
		log.Err(err).Msg(err.(*errors.Error).ErrorStack())
		return errors.Errorf("error removing trigger")
	}

	return nil
}

// HandleAddNode handles requests to the AddReplica endpoint of the client interface.
func (h *exthandler) HandleAddNode(n []Node) error {
	// TODO remove? IDK
	panic("Is it still necessary to add a node? The node can do this on startup by itself")
	// e := make([]string, len(n))
	// ec := 0
	//
	// for _, node := range n {
	// 	if err := h.replService.AddNode(node, true); err != nil {
	// 		log.Err(err).Msgf("Exthandler can not add a new replica node. (node=%#v)", node)
	// 		e[ec] = fmt.Sprintf("%v", err)
	// 		ec++
	// 	}
	// }
	//
	// if ec > 0 {
	// 	return errors.New(errors.StatusInternalError, fmt.Sprintf("exthandler: %v", e))
	// }

	return nil
}

// HandleGetReplica handles requests to the GetAllReplica endpoint of the client interface.
func (h *exthandler) HandleGetReplica(n Node) (Node, error) {
	return h.r.getNode(n)
}

// HandleGetAllReplica handles requests to the GetAllReplica endpoint of the client interface.
func (h *exthandler) HandleGetAllReplica() ([]Node, error) {
	return h.r.getNodes()
}

// HandleRemoveNode handles requests to the RemoveReplica endpoint of the client interface.
func (h *exthandler) HandleRemoveNode(n Node) error {
	// TODO why would this be called and what to do now.
	//return h.replService.RemoveNode(n, true)
	return nil
}