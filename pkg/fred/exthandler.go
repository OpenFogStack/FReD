package fred

import (
	"github.com/DistributedClocks/GoVector/govec/vclock"
	"github.com/go-errors/errors"
	"github.com/rs/zerolog/log"
)

type ExtHandler struct {
	s *storeService
	r *replicationService
	t *triggerService
	a *authService
	n NameService
}

// newExthandler creates a new handler for client request (i.e. from clients).
func newExthandler(s *storeService, r *replicationService, t *triggerService, a *authService, n NameService) *ExtHandler {
	return &ExtHandler{
		s: s,
		r: r,
		t: t,
		a: a,
		n: n,
	}
}

// HandleCreateKeygroup handles requests to the CreateKeygroup endpoint of the client interface.
func (h *ExtHandler) HandleCreateKeygroup(user string, k Keygroup) error {

	if err := h.r.createKeygroup(k); err != nil {
		log.Debug().Msg(err.(*errors.Error).ErrorStack())

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
func (h *ExtHandler) HandleDeleteKeygroup(user string, k Keygroup) error {
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
func (h *ExtHandler) HandleRead(user string, i Item, versions []vclock.VClock) ([]Item, error) {
	allowed, err := h.a.isAllowed(user, Read, i.Keygroup)

	if err != nil || !allowed {
		return nil, errors.Errorf("user %s cannot read from keygroup %s", user, i.Keygroup)
	}

	// if read request has a version, return all items newer than this version
	// if read request does not have a version, return all items with all versions
	// TODO: decide what to do with tombstoned items? right now we just don't show them

	var r []Item
	if len(versions) == 0 {
		r, err = h.s.read(i.Keygroup, i.ID)
	} else {
		r, err = h.s.readVersion(i.Keygroup, i.ID, versions)
	}

	result := make([]Item, 0, len(r))
	for _, it := range r {
		if !it.Tombstoned {
			result = append(result, it)
		}
	}

	if err != nil {
		log.Err(err).Msg(err.(*errors.Error).ErrorStack())
		return nil, errors.Errorf("error reading item %s from keygroup %s", i.ID, i.Keygroup)
	}

	if len(result) == 0 {
		return nil, errors.Errorf("item %s not found in keygroup %s", i.ID, i.Keygroup)
	}

	return result, nil
}

// HandleScan handles requests to the Scan endpoint of the client interface.
func (h *ExtHandler) HandleScan(user string, i Item, count uint64) ([]Item, error) {
	allowed, err := h.a.isAllowed(user, Read, i.Keygroup)

	if err != nil || !allowed {
		return nil, errors.Errorf("user %s cannot read from keygroup %s", user, i.Keygroup)
	}

	if count <= 0 {
		return nil, errors.Errorf("count must be at least 1, got %d", count)
	}

	// always return all versions
	result, err := h.s.scan(i.Keygroup, i.ID, count)

	if err != nil {
		log.Error().Msgf("Error in Scan is: %#v", err)
		// This prints the error stack whenever a item is not found, nobody cares about this...
		// log.Err(err).Msg(err.(*errors.Error).ErrorStack())
		return nil, errors.Errorf("error scanning %d items starting at %s from keygroup %s", count, i.ID, i.Keygroup)
	}

	return result, nil
}

// HandleAppend handles requests to the Append endpoint of the client interface.
func (h *ExtHandler) HandleAppend(user string, i Item) (Item, error) {
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

	log.Debug().Msgf("...keygroup %s is immutable", i.Keygroup)

	expiry, err := h.n.GetExpiry(i.Keygroup)

	if err != nil {
		return i, err
	}

	err = h.s.append(i, expiry)

	if err != nil {
		log.Printf("%#v", err)
		log.Err(err).Msg(err.(*errors.Error).ErrorStack())
		return i, errors.Errorf("error updating item")
	}

	if err := h.r.relayAppend(i); err != nil {
		log.Err(err).Msg(err.(*errors.Error).ErrorStack())
		return i, errors.Errorf("error updating item")
	}

	if err := h.t.triggerUpdate(i); err != nil {
		log.Err(err).Msg(err.(*errors.Error).ErrorStack())
		return i, errors.Errorf("error updating item")
	}

	return i, nil
}

// HandleUpdate handles requests to the Update endpoint of the client interface.
func (h *ExtHandler) HandleUpdate(user string, i Item, versions []vclock.VClock) (Item, error) {
	allowed, err := h.a.isAllowed(user, Update, i.Keygroup)

	if err != nil || !allowed {
		return i, errors.Errorf("user %s cannot update in keygroup %s", user, i.Keygroup)
	}

	log.Debug().Msgf("checking if keygroup %s is mutable...", i.Keygroup)
	m, err := h.n.IsMutable(i.Keygroup)

	if err != nil {
		return i, err
	}

	if !m {
		return i, errors.Errorf("cannot update item %s because keygroup is immutable", i.ID)
	}

	log.Debug().Msgf("...keygroup %s is mutable", i.Keygroup)

	expiry, err := h.n.GetExpiry(i.Keygroup)

	if err != nil {
		return i, err
	}

	// if update request has a list of versions, all versions that are equal or less than those versions will be overwritten
	// else only the local counter will be incremented
	if versions == nil {
		i.Version, err = h.s.update(i, expiry)

		if err != nil {
			log.Printf("%#v", err)
			log.Err(err).Msg(err.(*errors.Error).ErrorStack())
			return i, errors.Errorf("error updating item")
		}
	} else {
		i.Version, err = h.s.updateVersions(i, versions, expiry)

		if err != nil {
			log.Printf("%#v", err)
			log.Err(err).Msg(err.(*errors.Error).ErrorStack())
			return i, errors.Errorf("error updating item")
		}
	}

	if err := h.r.relayUpdate(i); err != nil {
		log.Err(err).Msg(err.(*errors.Error).ErrorStack())
		return i, errors.Errorf("error updating item")
	}

	if err := h.t.triggerUpdate(i); err != nil {
		log.Err(err).Msg(err.(*errors.Error).ErrorStack())
		return i, errors.Errorf("error updating item")
	}

	return i, nil
}

// HandleDelete handles requests to the Delete endpoint of the client interface.
func (h *ExtHandler) HandleDelete(user string, i Item, versions []vclock.VClock) (Item, error) {
	allowed, err := h.a.isAllowed(user, Delete, i.Keygroup)

	if err != nil || !allowed {
		return i, errors.Errorf("user %s cannot delete keygroup %s", user, i.Keygroup)
	}

	m, err := h.n.IsMutable(i.Keygroup)

	if err != nil {
		return i, err
	}

	if !m {
		return i, errors.Errorf("cannot update item %s because keygroup is immutable", i.ID)
	}

	if !h.s.exists(i) {
		return i, errors.Errorf("item does not exist so it cannot be deleted.")
	}

	i.Tombstoned = true

	// if delete request has a list of versions, all versions that are equal or less than those versions will be overwritten with tombstone
	// else only the local counter will be incremented
	if versions == nil {
		i.Version, err = h.s.tombstone(i)

		if err != nil {
			log.Printf("%#v", err)
			log.Err(err).Msg(err.(*errors.Error).ErrorStack())
			return i, errors.Errorf("error updating item")
		}
	} else {
		i.Version, err = h.s.tombstoneVersions(i, versions)

		if err != nil {
			log.Printf("%#v", err)
			log.Err(err).Msg(err.(*errors.Error).ErrorStack())
			return i, errors.Errorf("error updating item")
		}
	}

	if err := h.r.relayUpdate(i); err != nil {
		log.Err(err).Msg(err.(*errors.Error).ErrorStack())
		return i, errors.Errorf("error deleting item")
	}

	if err := h.t.triggerDelete(i); err != nil {
		log.Err(err).Msg(err.(*errors.Error).ErrorStack())
		return i, errors.Errorf("error deleting item")
	}

	return i, nil
}

// HandleAddReplica handles requests to the AddKeygroupReplica endpoint of the client interface.
func (h *ExtHandler) HandleAddReplica(user string, k Keygroup, n Node) error {
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
func (h *ExtHandler) HandleGetKeygroupReplica(user string, k Keygroup) ([]Node, map[NodeID]int, error) {
	allowed, err := h.a.isAllowed(user, GetReplica, k.Name)

	if err != nil || !allowed {
		log.Err(err).Msg("Exthandler uer is not allowed to GetKeygroupReplica")
		return nil, nil, errors.Errorf("user %s cannot get replica for keygroup %s", user, k.Name)
	}

	return h.r.getReplicaExternal(k)
}

// HandleRemoveReplica handles requests to the RemoveKeygroupReplica endpoint of the client interface.
func (h *ExtHandler) HandleRemoveReplica(user string, k Keygroup, n Node) error {
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

// HandleAddTrigger handles requests to the AddKeygroupTrigger endpoint of the client interface.
func (h *ExtHandler) HandleAddTrigger(user string, k Keygroup, t Trigger) error {
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
func (h *ExtHandler) HandleGetKeygroupTriggers(user string, k Keygroup) ([]Trigger, error) {
	allowed, err := h.a.isAllowed(user, GetTrigger, k.Name)

	if err != nil || !allowed {
		return nil, errors.Errorf("user %s cannot get triggers for keygroup %s", user, k.Name)
	}

	return h.t.getTrigger(k)
}

// HandleRemoveTrigger handles requests to the RemoveKeygroupTrigger endpoint of the client interface.
func (h *ExtHandler) HandleRemoveTrigger(user string, k Keygroup, t Trigger) error {
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
func (h *ExtHandler) HandleGetReplica(user string, n Node) (Node, error) {

	return h.r.getNodeExternal(n)
}

// HandleGetAllReplica handles requests to the GetAllReplica endpoint of the client interface.
func (h *ExtHandler) HandleGetAllReplica(user string) ([]Node, error) {

	return h.r.getNodesExternalAdress()
}

// HandleAddUser adds permissions to a keygroup to a new user.
func (h *ExtHandler) HandleAddUser(user string, newuser string, k Keygroup, r Role) error {
	allowed, err := h.a.isAllowed(user, AddUser, k.Name)

	if err != nil || !allowed {
		return errors.Errorf("user %s cannot add user permissions to keygroup %s", user, k.Name)
	}

	return h.a.addRoles(newuser, []Role{r}, k.Name)
}

// HandleRemoveUser removes permissions to a keygroup to a user.
func (h *ExtHandler) HandleRemoveUser(user string, newuser string, k Keygroup, r Role) error {
	allowed, err := h.a.isAllowed(user, RemoveUser, k.Name)

	if err != nil || !allowed {
		return errors.Errorf("user %s cannot remove user permissions to keygroup %s", user, k.Name)
	}

	return h.a.revokeRoles(newuser, []Role{r}, k.Name)
}
