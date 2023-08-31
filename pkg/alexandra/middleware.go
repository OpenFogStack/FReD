package alexandra

import (
	"context"
	"fmt"

	api "git.tu-berlin.de/mcc-fred/fred/proto/client"
	"git.tu-berlin.de/mcc-fred/fred/proto/middleware"
	"git.tu-berlin.de/mcc-fred/vclock"
	"github.com/rs/zerolog/log"
)

// Scan issues a scan request from the client to the middleware. The request is forwarded to FReD and incoming items are
// checked for their versions by comparing locally cached versions (if any). The local cache is also updated
// (if applicable).
func (m *Middleware) Scan(ctx context.Context, req *middleware.ScanRequest) (*middleware.ScanResponse, error) {

	m.vcache.cRLock(req.Keygroup, req.Id)
	defer m.vcache.cRUnlock(req.Keygroup, req.Id)

	res, err := m.clientsMgr.getLightHouse().Client.Scan(ctx, &api.ScanRequest{
		Keygroup: req.Keygroup,
		Id:       req.Id,
		Count:    req.Count,
	})

	if err != nil {
		return nil, err
	}

	data := make([]*middleware.Data, len(res.Data))

	for i, datum := range res.Data {
		data[i] = &middleware.Data{
			Id:   datum.Id,
			Data: datum.Val,
		}

		err = m.vcache.add(req.Keygroup, datum.Id, datum.Version.Version)

		if err != nil {
			return nil, err
		}
	}

	return &middleware.ScanResponse{Data: data}, nil
}

// Read reads a datum from FReD. Read data are placed in cache (if not in there already). If multiple versions of a
// datum exist, all versions will be returned to the client so that it can choose one. If the read data is outdated
// compared to seen versions, an error is returned.
func (m *Middleware) Read(_ context.Context, req *middleware.ReadRequest) (*middleware.ReadResponse, error) {
	log.Debug().Msgf("Alexandra has rcvd Read")

	m.vcache.cRLock(req.Keygroup, req.Id)
	defer m.vcache.cRUnlock(req.Keygroup, req.Id)

	vals, versions, err := m.clientsMgr.readFromAnywhere(req)

	// TODO: serialize access

	if err != nil {
		log.Error().Err(err)
		return nil, err
	}

	known, err := m.vcache.get(req.Keygroup, req.Id)

	if err != nil {
		return nil, err
	}

	for _, seen := range known {
		// TODO: is there a more elegant solution to this?
		// essentially we need to check the read data to cover all versions we have seen so far
		// this means that for every version we have seen so far, there must be at least one version in the read data
		// (preferably exactly one version) that is equal or newer to that seen version.
		covered := false
		for _, read := range versions {
			if seen.Compare(read, vclock.Descendant|vclock.Equal) {
				covered = true
				break
			}
		}

		if !covered {
			log.Error().Msgf("Alexandra Read has seen version %v is not covered by given versions %+v", seen, versions)
			return nil, fmt.Errorf("seen version %v is not covered by given versions %+v", seen, versions)
		}
	}

	log.Debug().Msgf("Alexandra Read key %s in kg %s: got vals %v versions %v", req.Id, req.Keygroup, vals, versions)

	for i := range versions {
		log.Debug().Msgf("Alexandra Read: putting version %v in cache for %s", versions[i], req.Id)
		err = m.vcache.add(req.Keygroup, req.Id, versions[i])
		if err != nil {
			log.Error().Err(err)
			return nil, err
		}
	}

	items := make([]*middleware.Item, len(vals))

	for i := range vals {
		items[i] = &middleware.Item{
			Val:     vals[i],
			Version: versions[i].GetMap(),
		}
	}

	log.Info().Msgf("Read: old %+v new %+v", known, versions)

	return &middleware.ReadResponse{
		Items: items,
	}, nil
}

// Update updates a datum in FReD. This could either be a value that has previously been read (if the datum is in cache)
// or a spontaneous write.
//
// If write-follows-read (i.e., datum can be found in cache), all versions of that datum THAT ARE KNOWN AT THE POINT OF
// THE UPDATE are superseded by the write.
// The assumption is that the client has merged conflicting values.
//
// If spontaneous write (i.e., datum cannot be found in cache), we assume an empty vector clock in the cache and send
// that to FReD. If there is a newer (any) data item in FReD already, this will fail.
func (m *Middleware) Update(ctx context.Context, req *middleware.UpdateRequest) (*middleware.Empty, error) {
	log.Debug().Msgf("Alexandra has rcvd Update")

	m.vcache.cLock(req.Keygroup, req.Id)
	defer m.vcache.cUnlock(req.Keygroup, req.Id)

	c, err := m.clientsMgr.getClient(req.Keygroup)

	if err != nil {
		log.Error().Err(err)
		return nil, err
	}

	known, err := m.vcache.get(req.Keygroup, req.Id)

	if err != nil {
		log.Error().Err(err)
		return nil, err
	}

	log.Debug().Msgf("Alexandra Update: know versions %v for %s", known, req.Id)

	if len(known) == 0 {
		// spontaneous write!
		// that means that we don't have any version cached yet
		// what we will do is send an empty vector clock (NOT AN EMPTY VECTOR CLOCK LIST!) to FReD
		// FReD will understand that there are currently no entries, which will still be greater than the existing
		// versions (i.e., no versions)
		known = []vclock.VClock{{}}
	}

	v, err := c.updateVersions(ctx, req.Keygroup, req.Id, req.Data, known)

	if err != nil {
		log.Error().Err(err)
		return nil, err
	}

	log.Debug().Msgf("Alexandra Update: new version %v for %s", v, req.Id)

	err = m.vcache.supersede(req.Keygroup, req.Id, known, v)

	if err != nil {
		log.Error().Err(err)
		return nil, err
	}

	log.Info().Msgf("Update: old %+v new %+v", known, v)

	return &middleware.Empty{}, nil
}

// Delete deletes a datum in FReD (it is actually only tombstoned, but this is irrelevant to the middleware or client).
// This could either be a value that has previously been read (if the datum is in cache) or a spontaneous delete.
//
// If write-follows-read (i.e., datum can be found in cache), all versions of that datum are superseded by the write.
// The assumption is that the client has merged conflicting values.
//
// If spontaneous delete (i.e., datum cannot be found in cache), we assume an empty vector clock in the cache and send
// that to FReD. If there is a newer (any) data item in FReD already, this will fail.
func (m *Middleware) Delete(ctx context.Context, req *middleware.DeleteRequest) (*middleware.Empty, error) {
	log.Debug().Msgf("Alexandra has rcvd Delete")

	m.vcache.cLock(req.Keygroup, req.Id)
	defer m.vcache.cUnlock(req.Keygroup, req.Id)

	c, err := m.clientsMgr.getClient(req.Keygroup)

	if err != nil {
		log.Error().Err(err)
		return nil, err
	}

	known, err := m.vcache.get(req.Keygroup, req.Id)

	if err != nil {
		log.Error().Err(err)
		return nil, err
	}

	if len(known) == 0 {
		// spontaneous delete!
		known = []vclock.VClock{{}}
	}

	v, err := c.deleteVersions(ctx, req.Keygroup, req.Id, known)

	if err != nil {
		log.Error().Err(err)
		return nil, err
	}

	err = m.vcache.supersede(req.Keygroup, req.Id, known, v)

	if err != nil {
		log.Error().Err(err)
		return nil, err
	}

	return &middleware.Empty{}, nil
}

// Append appends a new datum to an immutable keygroup in FReD. As data cannot be changed, no versions are necessary.
// Thus, the request is only passed through to FReD without caching it.
// FReD's append endpoint requires a unique ID for a datum. ALExANDRA automatically uses a Unix nanosecond timestamp for
// this.
func (m *Middleware) Append(ctx context.Context, req *middleware.AppendRequest) (*middleware.AppendResponse, error) {

	c, err := m.clientsMgr.getClient(req.Keygroup)

	if err != nil {
		return nil, err
	}
	res, err := c.append(ctx, req.Keygroup, req.Data)

	if err != nil {
		return nil, err
	}

	return &middleware.AppendResponse{Id: res.Id}, err
}

// Notify notifies the middleware about a version of a datum that the client has seen by bypassing the middleware. This
// is required to capture external causality.
func (m *Middleware) Notify(_ context.Context, req *middleware.NotifyRequest) (*middleware.Empty, error) {


	m.vcache.cLock(req.Keygroup, req.Id)
	defer m.vcache.cUnlock(req.Keygroup, req.Id)

	err := m.vcache.add(req.Keygroup, req.Id, req.Version)

	if err != nil {
		return nil, err
	}

	return &middleware.Empty{}, nil
}

// ChooseReplica allows a client to choose a particular note to send requests to for a keygroup. This will override the
// fastest node if exists
func (m *Middleware) ChooseReplica(_ context.Context, req *middleware.ChooseReplicaRequest) (*middleware.Empty, error) {
	log.Debug().Msgf("AlexandraServer has rcdv ChooseReplica: %+v", req)
	err := m.clientsMgr.setPreferred(req.Keygroup, req.NodeId)

	if err != nil {
		return nil, err
	}

	return &middleware.Empty{}, nil
}

// CreateKeygroup creates the keygroup and also adds the first node (This is two operations in the eye of FReD:
// CreateKeygroup and AddReplica)
func (m *Middleware) CreateKeygroup(ctx context.Context, req *middleware.CreateKeygroupRequest) (*middleware.Empty, error) {
	log.Debug().Msgf("AlexandraServer has rcdv CreateKeygroup: %+v", req)
	getReplica, err := m.clientsMgr.getFastestClient().getReplica(ctx, req.FirstNodeId)

	if err != nil {
		return nil, err
	}

	log.Debug().Msgf("CreateKeygroup: using node %s (addr=%s)", getReplica.NodeId, getReplica.Host)

	_, err = m.clientsMgr.getClientTo(getReplica.Host, getReplica.NodeId).createKeygroup(ctx, req.Keygroup, req.Mutable, req.Expiry)

	if err != nil {
		return nil, err
	}

	return &middleware.Empty{}, err
}

// DeleteKeygroup deletes a keygroup from FReD.
func (m *Middleware) DeleteKeygroup(ctx context.Context, req *middleware.DeleteKeygroupRequest) (*middleware.Empty, error) {
	client, err := m.clientsMgr.getClient(req.Keygroup)
	if err != nil {
		return nil, err
	}
	log.Debug().Msgf("DeleteKeygroup: using node %+v", client.nodeID)

	_, err = client.deleteKeygroup(ctx, req.Keygroup)

	if err != nil {
		return nil, err
	}

	return &middleware.Empty{}, err
}

// AddReplica lets the client explicitly add a new replica for a keygroup. In the future, this should happen
// automatically.
func (m *Middleware) AddReplica(ctx context.Context, req *middleware.AddReplicaRequest) (*middleware.Empty, error) {
	c, err := m.clientsMgr.getClient(req.Keygroup)

	if err != nil {
		return nil, err
	}

	defer m.clientsMgr.updateKeygroupClients(req.Keygroup)

	_, err = c.Client.AddReplica(ctx, &api.AddReplicaRequest{
		Keygroup: req.Keygroup,
		NodeId:   req.NodeId,
		Expiry:   req.Expiry,
	})

	if err != nil {
		return nil, err
	}

	return &middleware.Empty{}, err
}

// RemoveReplica lets the client explicitly remove a new replica for a keygroup. In the future, this should happen
// automatically.
func (m *Middleware) RemoveReplica(ctx context.Context, req *middleware.RemoveReplicaRequest) (*middleware.Empty, error) {
	c, err := m.clientsMgr.getClient(req.Keygroup)

	if err != nil {
		return nil, err
	}

	defer m.clientsMgr.updateKeygroupClients(req.Keygroup)

	_, err = c.Client.RemoveReplica(ctx, &api.RemoveReplicaRequest{
		Keygroup: req.Keygroup,
		NodeId:   req.NodeId,
	})
	if err != nil {
		return nil, err
	}

	return &middleware.Empty{}, err
}

// GetReplica returns information about a specific FReD node. In the future, this API will be removed as ALExANDRA
// handles data replication.
func (m *Middleware) GetReplica(ctx context.Context, req *middleware.GetReplicaRequest) (*middleware.GetReplicaResponse, error) {
	res, err := m.clientsMgr.getFastestClient().Client.GetReplica(ctx, &api.GetReplicaRequest{NodeId: req.NodeId})

	if err != nil {
		return nil, err
	}

	return &middleware.GetReplicaResponse{NodeId: res.NodeId, Host: res.Host}, err
}

// GetAllReplica returns a list of all FReD nodes. In the future, this API will be removed as ALExANDRA handles data
// replication.
func (m *Middleware) GetAllReplica(ctx context.Context, _ *middleware.GetAllReplicaRequest) (*middleware.GetAllReplicaResponse, error) {
	res, err := m.clientsMgr.getFastestClient().Client.GetAllReplica(ctx, &api.Empty{})

	if err != nil {
		return nil, err
	}

	replicas := make([]*middleware.GetReplicaResponse, len(res.Replicas))
	for i, replica := range res.Replicas {
		replicas[i] = &middleware.GetReplicaResponse{
			NodeId: replica.NodeId,
			Host:   replica.Host,
		}
	}

	return &middleware.GetAllReplicaResponse{Replicas: replicas}, err
}

// GetKeygroupInfo returns a list of all FReD nodes that replicate a given keygroup. In the future, this API will be
// removed as ALExANDRA handles data replication.
func (m *Middleware) GetKeygroupInfo(ctx context.Context, req *middleware.GetKeygroupInfoRequest) (*middleware.GetKeygroupInfoResponse, error) {
	c, err := m.clientsMgr.getClient(req.Keygroup)

	if err != nil {
		return nil, err
	}

	res, err := c.Client.GetKeygroupInfo(ctx, &api.GetKeygroupInfoRequest{Keygroup: req.Keygroup})

	if err != nil {
		return nil, err
	}

	replicas := make([]*middleware.KeygroupReplica, len(res.Replica))
	for i, replica := range res.Replica {
		replicas[i] = &middleware.KeygroupReplica{
			NodeId: replica.NodeId,
			Host:   replica.Host,
		}
	}

	return &middleware.GetKeygroupInfoResponse{
		Mutable: res.Mutable,
		Replica: replicas,
	}, err
}

// GetKeygroupTriggers returns a list of trigger nodes for a keygroup.
func (m *Middleware) GetKeygroupTriggers(ctx context.Context, req *middleware.GetKeygroupTriggerRequest) (*middleware.GetKeygroupTriggerResponse, error) {
	res, err := m.clientsMgr.getLightHouse().Client.GetKeygroupTriggers(ctx, &api.GetKeygroupTriggerRequest{
		Keygroup: req.Keygroup,
	})

	if err != nil {
		return nil, err
	}

	triggers := make([]*middleware.Trigger, len(res.Triggers))
	for i, trigger := range res.Triggers {
		triggers[i] = &middleware.Trigger{
			Id:   trigger.Id,
			Host: trigger.Host,
		}
	}
	return &middleware.GetKeygroupTriggerResponse{Triggers: triggers}, nil
}

// AddTrigger adds a new trigger to a keygroup.
func (m *Middleware) AddTrigger(ctx context.Context, req *middleware.AddTriggerRequest) (*middleware.Empty, error) {
	_, err := m.clientsMgr.getLightHouse().Client.AddTrigger(ctx, &api.AddTriggerRequest{
		Keygroup:    req.Keygroup,
		TriggerId:   req.TriggerId,
		TriggerHost: req.TriggerHost,
	})

	if err != nil {
		return nil, err
	}

	return &middleware.Empty{}, err
}

// RemoveTrigger removes a trigger node for a keygroup.
func (m *Middleware) RemoveTrigger(ctx context.Context, req *middleware.RemoveTriggerRequest) (*middleware.Empty, error) {
	_, err := m.clientsMgr.getLightHouse().Client.RemoveTrigger(ctx, &api.RemoveTriggerRequest{
		Keygroup:  req.Keygroup,
		TriggerId: req.TriggerId,
	})

	if err != nil {
		return nil, err
	}

	return &middleware.Empty{}, err
}

// AddUser adds permissions to access a keygroup for a particular user to FReD.
func (m *Middleware) AddUser(ctx context.Context, req *middleware.UserRequest) (*middleware.Empty, error) {
	_, err := m.clientsMgr.getLightHouse().Client.AddUser(ctx, &api.AddUserRequest{
		User:     req.User,
		Keygroup: req.Keygroup,
		Role:     api.UserRole(req.Role),
	})

	if err != nil {
		return nil, err
	}

	return &middleware.Empty{}, err
}

// RemoveUser removes permissions to access a keygroup for a particular user from FReD.
func (m *Middleware) RemoveUser(ctx context.Context, req *middleware.UserRequest) (*middleware.Empty, error) {
	_, err := m.clientsMgr.getLightHouse().Client.RemoveUser(ctx, &api.RemoveUserRequest{
		User:     req.User,
		Keygroup: req.Keygroup,
		Role:     api.UserRole(req.Role),
	})

	if err != nil {
		return nil, err
	}

	return &middleware.Empty{}, err
}
