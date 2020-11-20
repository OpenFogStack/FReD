package fred

import (
	"github.com/go-errors/errors"
	"github.com/rs/zerolog/log"
)

type exthandler struct {
	s *storeService
	r *replicationService
	t *triggerService
	a *authService
	n NameService
}

// newExthandler creates a new handler for client request (i.e. from clients).
func newExthandler(s *storeService, r *replicationService, t *triggerService, a *authService, n NameService) *exthandler {
	return &exthandler{
		s: s,
		r: r,
		t: t,
		a: a,
		n: n,
	}
}

// HandleCreateKeygroup handles requests to the CreateKeygroup endpoint of the client interface.
func (h *exthandler) HandleCreateKeygroup(user string, k Keygroup) error {
	if err := h.r.createKeygroup(k); err != nil {
		log.Err(err).Msg(err.(*errors.Error).ErrorStack())
		return errors.Errorf("error creating keygroup")
	}

	if err := h.s.createKeygroup(k.Name); err != nil {
		log.Err(err).Msg(err.(*errors.Error).ErrorStack())
		return errors.Errorf("error creating keygroup")
	}

	// when a user creates a keygroup, they should have all rights for that keygroup
	err := h.a.addRoles(user, []Role{ReadKeygroup, WriteKeygroup, ConfigureReplica, ConfigureTrigger, ConfigureKeygroups}, k.Name)

	if err != nil {
		return err
	}

	return nil
}

// HandleDeleteKeygroup handles requests to the DeleteKeygroup endpoint of the client interface.
func (h *exthandler) HandleDeleteKeygroup(user string, k Keygroup) error {
	allowed, err := h.a.isAllowed(user, DeleteKeygroup, k.Name)

	if err != nil || !allowed {
		return errors.Errorf("user %s cannot delete keygroup %s", user, k.Name)
	}

	if err := h.s.deleteKeygroup(k.Name); err != nil {
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
func (h *exthandler) HandleRead(user string, i Item) (Item, error) {
	allowed, err := h.a.isAllowed(user, Read, i.Keygroup)

	if err != nil || !allowed {
		return Item{}, errors.Errorf("user %s cannot read from keygroup %s", user, i.Keygroup)
	}

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

// HandleAppend handles requests to the Append endpoint of the client interface.
func (h *exthandler) HandleAppend(user string, i Item) (Item, error) {
	allowed, err := h.a.isAllowed(user, Update, i.Keygroup)

	if err != nil || !allowed {
		return i, errors.Errorf("user %s cannot update in keygroup %s", user, i.Keygroup)
	}

	log.Debug().Msgf("checking if keygroup %s is mutable...", i.Keygroup)
	m, err := h.n.IsMutable(i.Keygroup)

	if err != nil {
		return i, err
	}

	if m {
		return i, errors.Errorf("cannot append %s to mutable keygroup", i.ID)
	}

	log.Debug().Msgf("keygroup %s is immutable, proceeding...", i.Keygroup)

	expiry, err := h.n.GetExpiry(i.Keygroup)

	if err != nil {
		return i, err
	}

	result, err := h.s.append(i, expiry)

	if err != nil {
		log.Printf("%#v", err)
		log.Err(err).Msg(err.(*errors.Error).ErrorStack())
		return i, errors.Errorf("error updating item")
	}

	if err := h.r.relayUpdate(result); err != nil {
		log.Err(err).Msg(err.(*errors.Error).ErrorStack())
		return result, errors.Errorf("error updating item")
	}

	if err := h.t.triggerUpdate(result); err != nil {
		log.Err(err).Msg(err.(*errors.Error).ErrorStack())
		return result, errors.Errorf("error updating item")
	}

	return result, nil
}

// HandleUpdate handles requests to the Update endpoint of the client interface.
func (h *exthandler) HandleUpdate(user string, i Item) error {
	allowed, err := h.a.isAllowed(user, Update, i.Keygroup)

	if err != nil || !allowed {
		return errors.Errorf("user %s cannot update in keygroup %s", user, i.Keygroup)
	}

	log.Debug().Msgf("checking if keygroup %s is mutable...", i.Keygroup)
	m, err := h.n.IsMutable(i.Keygroup)

	if err != nil {
		return err
	}

	if !m {
		return errors.Errorf("cannot update item %s because keygroup is immutable", i.ID)
	}

	log.Debug().Msgf("keygroup %s is mutable, proceeding...", i.Keygroup)

	expiry, err := h.n.GetExpiry(i.Keygroup)

	if err != nil {
		return err
	}

	if err := h.s.update(i, expiry); err != nil {
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
func (h *exthandler) HandleDelete(user string, i Item) error {
	allowed, err := h.a.isAllowed(user, Delete, i.Keygroup)

	if err != nil || !allowed {
		return errors.Errorf("user %s cannot delete keygroup %s", user, i.Keygroup)
	}

	m, err := h.n.IsMutable(i.Keygroup)

	if err != nil {
		return err
	}

	if !m {
		return errors.Errorf("cannot update item %s because keygroup is immutable", i.ID)
	}

	if !h.s.exists(i) {
		return errors.Errorf("item does not exist so it cannot be deleted.")
	}

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
func (h *exthandler) HandleAddReplica(user string, k Keygroup, n Node) error {
	allowed, err := h.a.isAllowed(user, AddReplica, k.Name)

	if err != nil || !allowed {
		return errors.Errorf("user %s cannot add replica to keygroup %s", user, k.Name)
	}

	if err := h.r.addReplica(k, n, true); err != nil {
		log.Error().Msgf("Error in AddReplica is: %#v", err)
		log.Err(err).Msg(err.(*errors.Error).ErrorStack())
		return errors.Errorf("error adding replica")
	}

	return nil
}

// HandleGetKeygroupReplica handles requests to the GetKeygroupReplica endpoint of the client interface.
func (h *exthandler) HandleGetKeygroupReplica(user string, k Keygroup) ([]Node, map[NodeID]int, error) {
	allowed, err := h.a.isAllowed(user, GetReplica, k.Name)

	if err != nil || !allowed {
		return nil, nil, errors.Errorf("user %s cannot get replica for keygroup %s", user, k.Name)
	}

	return h.r.getReplica(k)
}

// HandleRemoveReplica handles requests to the RemoveKeygroupReplica endpoint of the client interface.
func (h *exthandler) HandleRemoveReplica(user string, k Keygroup, n Node) error {
	allowed, err := h.a.isAllowed(user, RemoveReplica, k.Name)

	if err != nil || !allowed {
		return errors.Errorf("user %s cannot remove replica of keygroup %s", user, k.Name)
	}

	if err := h.r.removeReplica(k, n, true); err != nil {
		log.Err(err).Msg(err.(*errors.Error).ErrorStack())
		return errors.Errorf("error removing replica")
	}

	return nil
}

// HandleAddReplica handles requests to the AddKeygroupTrigger endpoint of the client interface.
func (h *exthandler) HandleAddTrigger(user string, k Keygroup, t Trigger) error {
	allowed, err := h.a.isAllowed(user, AddTrigger, k.Name)

	if err != nil || !allowed {
		return errors.Errorf("user %s cannot add trigger to keygroup %s", user, k.Name)
	}

	if err := h.t.addTrigger(k, t); err != nil {
		log.Error().Msgf("Error in AddTrigger is: %#v", err)
		log.Err(err).Msg(err.(*errors.Error).ErrorStack())
		return errors.Errorf("error adding replica")
	}
	return nil
}

// HandleGetKeygroupTriggers handles requests to the GetKeygroupTrigger endpoint of the client interface.
func (h *exthandler) HandleGetKeygroupTriggers(user string, k Keygroup) ([]Trigger, error) {
	allowed, err := h.a.isAllowed(user, GetTrigger, k.Name)

	if err != nil || !allowed {
		return nil, errors.Errorf("user %s cannot get triggers for keygroup %s", user, k.Name)
	}

	return h.t.getTrigger(k)
}

// HandleRemoveTrigger handles requests to the RemoveKeygroupTrigger endpoint of the client interface.
func (h *exthandler) HandleRemoveTrigger(user string, k Keygroup, t Trigger) error {
	allowed, err := h.a.isAllowed(user, RemoveTrigger, k.Name)

	if err != nil || !allowed {
		return errors.Errorf("user %s cannot remove triggers from keygroup %s", user, k.Name)
	}

	if err := h.t.removeTrigger(k, t); err != nil {
		log.Err(err).Msg(err.(*errors.Error).ErrorStack())
		return errors.Errorf("error removing trigger")
	}

	return nil
}

// HandleGetReplica handles requests to the GetReplica endpoint of the client interface.
func (h *exthandler) HandleGetReplica(user string, n Node) (Node, error) {

	return h.r.getNode(n)
}

// HandleGetAllReplica handles requests to the GetAllReplica endpoint of the client interface.
func (h *exthandler) HandleGetAllReplica(user string) ([]Node, error) {

	return h.r.getNodes()
}

// AddUser adds permissions to a keygroup to a new user.
func (h *exthandler) HandleAddUser(user string, newuser string, k Keygroup, r Role) error {
	allowed, err := h.a.isAllowed(user, AddUser, k.Name)

	if err != nil || !allowed {
		return errors.Errorf("user %s cannot add user permissions to keygroup %s", user, k.Name)
	}

	return h.a.addRoles(newuser, []Role{r}, k.Name)
}

// RemoveUser removes permissions to a keygroup to a user.
func (h *exthandler) HandleRemoveUser(user string, newuser string, k Keygroup, r Role) error {
	allowed, err := h.a.isAllowed(user, RemoveUser, k.Name)

	if err != nil || !allowed {
		return errors.Errorf("user %s cannot remove user permissions to keygroup %s", user, k.Name)
	}

	return h.a.revokeRoles(newuser, []Role{r}, k.Name)
}
