package fred

import (
	"github.com/go-errors/errors"
	"github.com/rs/zerolog/log"
)

type inthandler struct {
	s *storeService
	r *replicationService
}

// newInthandler creates a new handler for internal request (i.e. from peer nodes or the naming service).
func newInthandler(s *storeService, r *replicationService) *inthandler {
	return &inthandler{
		s: s,
		r: r,
	}
}

// HandleCreateKeygroup handles requests to the CreateKeygroup endpoint of the internal interface.
func (h *inthandler) HandleCreateKeygroup(k Keygroup, nodes []Node) error {
	if err := h.s.createKeygroup(k.Name); err != nil {
		log.Err(err).Msg(err.(*errors.Error).ErrorStack())
		return errors.Errorf("error creating keygroup")
	}

	if err := h.r.CreateKeygroup(Keygroup{
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
	if err := h.s.deleteKeygroup(k.Name); err != nil {
		log.Err(err).Msg(err.(*errors.Error).ErrorStack())
		return errors.Errorf("error deleting keygroup")
	}

	if err := h.r.DeleteKeygroup(Keygroup{
		Name: k.Name,
	}); err != nil {
		log.Err(err).Msg(err.(*errors.Error).ErrorStack())
		return errors.Errorf("error deleting keygroup")
	}

	return nil
}

// HandleUpdate handles requests to the Update endpoint of the internal interface.
func (h *inthandler) HandleUpdate(i Item) error {
	return h.s.update(i)
}

// HandleDelete handles requests to the Delete endpoint of the internal interface.
func (h *inthandler) HandleDelete(i Item) error {
	return h.s.delete(i.Keygroup, i.ID)
}

func (h *inthandler) HandleAddReplica(k Keygroup, n Node) error {
	return h.r.AddReplica(k, n, nil, false)
}

func (h *inthandler) HandleRemoveReplica(k Keygroup, n Node) error {
	return h.r.RemoveReplica(k, n, false)
}
