package alexandra

import (
	"context"
	"fmt"

	api "git.tu-berlin.de/mcc-fred/fred/proto/client"
	alexandraProto "git.tu-berlin.de/mcc-fred/fred/proto/middleware"
	"github.com/DistributedClocks/GoVector/govec/vclock"
	"github.com/rs/zerolog/log"
)

// Scan issues a scan request from the client to the middleware. The request is forwarded to FReD and incoming items are
// checked for their versions by comparing locally cached versions (if any). The local cache is also updated
// (if applicable).
func (s *Server) Scan(ctx context.Context, req *alexandraProto.ScanRequest) (*alexandraProto.ScanResponse, error) {
	res, err := s.clientsMgr.getClientTo(s.lighthouse).Client.Scan(ctx, &api.ScanRequest{
		Keygroup: req.Keygroup,
		Id:       req.Id,
		Count:    req.Count,
	})

	if err != nil {
		return nil, err
	}

	data := make([]*alexandraProto.Data, len(res.Data))

	for i, datum := range res.Data {
		data[i] = &alexandraProto.Data{
			Id:   datum.Id,
			Data: datum.Val,
		}

		err = s.cache.add(req.Keygroup, req.Id, datum.Version.Version)
	}

	return &alexandraProto.ScanResponse{Data: data}, err
}

// Read reads a datum from FReD. Read data are placed in cache (if not in there already). If multiple versions of a
// datum exist, all versions will be returned to the client so that it can choose one. If the read data is outdated
// compared to seen versions, an error is returned.
func (s *Server) Read(ctx context.Context, req *alexandraProto.ReadRequest) (*alexandraProto.ReadResponse, error) {
	log.Debug().Msgf("Alexandra has rcvd Read")

	vals, versions, err := s.clientsMgr.readFromAnywhere(ctx, req)

	if err != nil {
		return nil, err
	}

	known, err := s.cache.get(req.Keygroup, req.Id)

	if err != nil {
		return nil, err
	}

	for _, read := range versions {
		for _, seen := range known {
			if seen.Compare(read, vclock.Ancestor) {
				return nil, fmt.Errorf("read version %v is older than seen version %v", read, seen)
			}
		}
	}

	for i := range versions {
		err = s.cache.add(req.Keygroup, req.Id, versions[i])
		if err != nil {
			return nil, err
		}
	}

	items := make([]*alexandraProto.Item, len(vals))

	for i := range vals {
		items[i] = &alexandraProto.Item{
			Val:     vals[i],
			Version: versions[i].GetMap(),
		}
	}

	return &alexandraProto.ReadResponse{
		Items: items,
	}, nil
}

// Update updates a datum in FReD. This could either be a value that has previously been read (if the datum is in cache)
// or a spontaneous write.
//
// If write-follows-read (i.e., datum can be found in cache), all versions of that datum are superseded by the write.
// The assumption is that the client has merged conflicting values.
//
// If spontaneous write (i.e., datum cannot be found in cache) TODO
func (s *Server) Update(ctx context.Context, req *alexandraProto.UpdateRequest) (*alexandraProto.Empty, error) {
	c, err := s.clientsMgr.getClient(req.Keygroup, UseSlowerNodeProb)

	if err != nil {
		return nil, err
	}

	known, err := s.cache.get(req.Keygroup, req.Id)

	if err != nil {
		return nil, err
	}

	if len(known) == 0 {
		// spontaneous write!
		// that means that we don't have any version cached yet
		// what we will do is send an empty vector clock (NOT AN EMPTY VECTOR CLOCK LIST!) to FReD
		// FReD will understand that there are currently no entries, which will still be greater than the existing
		// versions (i.e., no versions)

		v, err := c.updateVersions(ctx, req.Keygroup, req.Id, req.Data, []vclock.VClock{{}})

		if err != nil {
			return nil, err
		}

		err = s.cache.supersede(req.Keygroup, req.Id, v)

		if err != nil {
			return nil, err
		}

		return &alexandraProto.Empty{}, nil
	}

	v, err := c.updateVersions(ctx, req.Keygroup, req.Id, req.Data, known)

	if err != nil {
		return nil, err
	}

	err = s.cache.supersede(req.Keygroup, req.Id, v)

	if err != nil {
		return nil, err
	}

	return &alexandraProto.Empty{}, nil
}

// Delete deletes a datum in FReD (it is actually only tombstoned, but this is irrelevant to the middleware or client).
// This could either be a value that has previously been read (if the datum is in cache) or a spontaneous delete.
//
// If write-follows-read (i.e., datum can be found in cache), all versions of that datum are superseded by the write.
// The assumption is that the client has merged conflicting values.
//
// If spontaneous delete (i.e., datum cannot be found in cache) TODO
func (s *Server) Delete(ctx context.Context, req *alexandraProto.DeleteRequest) (*alexandraProto.Empty, error) {
	c, err := s.clientsMgr.getClient(req.Keygroup, UseSlowerNodeProb)

	if err != nil {
		return nil, err
	}

	known, err := s.cache.get(req.Keygroup, req.Id)

	if err != nil {
		return nil, err
	}

	if len(known) == 0 {
		// spontaneous delete!

		v, err := c.deleteVersions(ctx, req.Keygroup, req.Id, []vclock.VClock{{}})

		if err != nil {
			return nil, err
		}

		err = s.cache.supersede(req.Keygroup, req.Id, v)

		if err != nil {
			return nil, err
		}

		return &alexandraProto.Empty{}, nil
	}

	v, err := c.deleteVersions(ctx, req.Keygroup, req.Id, known)

	if err != nil {
		return nil, err
	}

	err = s.cache.supersede(req.Keygroup, req.Id, v)

	if err != nil {
		return nil, err
	}

	return &alexandraProto.Empty{}, nil
}

// Append appends a new datum to an immutable keygroup in FReD. As data cannot be changed, no versions are necessary.
// Thus, the request is only passed through to FReD without caching it.
// FReD's append endpoint requires a unique ID for a datum. ALExANDRA automatically uses a Unix nanosecond timestamp for
// this.
func (s *Server) Append(ctx context.Context, req *alexandraProto.AppendRequest) (*alexandraProto.AppendResponse, error) {
	c, err := s.clientsMgr.getClient(req.Keygroup, UseSlowerNodeProb)

	if err != nil {
		return nil, err
	}
	res, err := c.append(ctx, req.Keygroup, req.Data)

	if err != nil {
		return nil, err
	}

	return &alexandraProto.AppendResponse{Id: res.Id}, err
}

// Notify notifies the middleware about a version of a datum that the client has seen by bypassing the middleware. This
// is required to capture external causality.
func (s *Server) Notify(_ context.Context, req *alexandraProto.NotifyRequest) (*alexandraProto.NotifyResponse, error) {
	err := s.cache.add(req.Keygroup, req.Id, req.Version)

	if err != nil {
		return nil, err
	}

	return &alexandraProto.NotifyResponse{}, nil
}

// CreateKeygroup creates the keygroup and also adds the first node (This is two operations in the eye of FReD:
// CreateKeygroup and AddReplica)
func (s *Server) CreateKeygroup(ctx context.Context, req *alexandraProto.CreateKeygroupRequest) (*alexandraProto.Empty, error) {
	log.Debug().Msgf("AlexandraServer has rcdv CreateKeygroup: %#v", req)
	getReplica, err := s.clientsMgr.getFastestClient().getReplica(ctx, req.FirstNodeId)
	if err != nil {
		return nil, err
	}
	log.Debug().Msgf("CreateKeygroup: using node %s (addr=%s)", getReplica.NodeId, getReplica.Host)

	_, err = s.clientsMgr.getClientTo(getReplica.Host).createKeygroup(ctx, req.Keygroup, req.Mutable, req.Expiry)

	if err != nil {
		return nil, err
	}

	return &alexandraProto.Empty{}, err
}

// DeleteKeygroup deletes a keygroup from FReD.
func (s *Server) DeleteKeygroup(ctx context.Context, req *alexandraProto.DeleteKeygroupRequest) (*alexandraProto.Empty, error) {
	client, err := s.clientsMgr.getFastestClientWithKeygroup(req.Keygroup, 1)
	if err != nil {
		return nil, err
	}
	log.Debug().Msgf("DeleteKeygroup: using node %#v", client)

	_, err = client.deleteKeygroup(ctx, req.Keygroup)

	if err != nil {
		return nil, err
	}

	return &alexandraProto.Empty{}, err
}

// AddReplica lets the client explicitly add a new replica for a keygroup. In the future, this should happen
// automatically.
func (s *Server) AddReplica(ctx context.Context, req *alexandraProto.AddReplicaRequest) (*alexandraProto.Empty, error) {
	_, err := s.clientsMgr.getFastestClient().Client.AddReplica(ctx, &api.AddReplicaRequest{
		Keygroup: req.Keygroup,
		NodeId:   req.NodeId,
		Expiry:   req.Expiry,
	})

	if err != nil {
		return nil, err
	}

	return &alexandraProto.Empty{}, err
}

// RemoveReplica lets the client explicitly remove a new replica for a keygroup. In the future, this should happen
// automatically.
func (s *Server) RemoveReplica(ctx context.Context, req *alexandraProto.RemoveReplicaRequest) (*alexandraProto.Empty, error) {
	_, err := s.clientsMgr.getFastestClient().Client.RemoveReplica(ctx, &api.RemoveReplicaRequest{
		Keygroup: req.Keygroup,
		NodeId:   req.NodeId,
	})
	if err != nil {
		return nil, err
	}

	return &alexandraProto.Empty{}, err
}

// GetReplica returns information about a specific FReD node. In the future, this API will be removed as ALExANDRA
// handles data replication.
func (s *Server) GetReplica(ctx context.Context, req *alexandraProto.GetReplicaRequest) (*alexandraProto.GetReplicaResponse, error) {
	res, err := s.clientsMgr.getFastestClient().Client.GetReplica(ctx, &api.GetReplicaRequest{NodeId: req.NodeId})

	if err != nil {
		return nil, err
	}

	return &alexandraProto.GetReplicaResponse{NodeId: res.NodeId, Host: res.Host}, err
}

// GetAllReplica returns a list of all FReD nodes. In the future, this API will be removed as ALExANDRA handles data
//// replication.
func (s *Server) GetAllReplica(ctx context.Context, _ *alexandraProto.GetAllReplicaRequest) (*alexandraProto.GetAllReplicaResponse, error) {
	res, err := s.clientsMgr.getFastestClient().Client.GetAllReplica(ctx, &api.Empty{})

	if err != nil {
		return nil, err
	}

	replicas := make([]*alexandraProto.GetReplicaResponse, len(res.Replicas))
	for i, replica := range res.Replicas {
		replicas[i] = &alexandraProto.GetReplicaResponse{
			NodeId: replica.NodeId,
			Host:   replica.Host,
		}
	}

	return &alexandraProto.GetAllReplicaResponse{Replicas: replicas}, err
}

// GetKeygroupReplica returns a list of all FReD nodes that replicate a given keygroup. In the future, this API will be
// removed as ALExANDRA handles data replication.
func (s *Server) GetKeygroupReplica(ctx context.Context, req *alexandraProto.GetKeygroupReplicaRequest) (*alexandraProto.GetKeygroupReplicaResponse, error) {
	res, err := s.clientsMgr.getFastestClient().Client.GetKeygroupReplica(ctx, &api.GetKeygroupReplicaRequest{Keygroup: req.Keygroup})

	if err != nil {
		return nil, err
	}

	replicas := make([]*alexandraProto.KeygroupReplica, len(res.Replica))
	for i, replica := range res.Replica {
		replicas[i] = &alexandraProto.KeygroupReplica{
			NodeId: replica.NodeId,
			Host:   replica.Host,
		}
	}

	return &alexandraProto.GetKeygroupReplicaResponse{Replica: replicas}, err
}

// GetKeygroupTriggers returns a list of trigger nodes for a keygroup.
func (s *Server) GetKeygroupTriggers(ctx context.Context, req *alexandraProto.GetKeygroupTriggerRequest) (*alexandraProto.GetKeygroupTriggerResponse, error) {
	res, err := s.clientsMgr.getClientTo(s.lighthouse).Client.GetKeygroupTriggers(ctx, &api.GetKeygroupTriggerRequest{
		Keygroup: req.Keygroup,
	})

	if err != nil {
		return nil, err
	}

	triggers := make([]*alexandraProto.Trigger, len(res.Triggers))
	for i, trigger := range res.Triggers {
		triggers[i] = &alexandraProto.Trigger{
			Id:   trigger.Id,
			Host: trigger.Host,
		}
	}
	return &alexandraProto.GetKeygroupTriggerResponse{Triggers: triggers}, nil
}

// AddTrigger adds a new trigger to a keygroup.
func (s *Server) AddTrigger(ctx context.Context, req *alexandraProto.AddTriggerRequest) (*alexandraProto.Empty, error) {
	_, err := s.clientsMgr.getClientTo(s.lighthouse).Client.AddTrigger(ctx, &api.AddTriggerRequest{
		Keygroup:    req.Keygroup,
		TriggerId:   req.TriggerId,
		TriggerHost: req.TriggerHost,
	})

	if err != nil {
		return nil, err
	}

	return &alexandraProto.Empty{}, err
}

// RemoveTrigger removes a trigger node for a keygroup.
func (s *Server) RemoveTrigger(ctx context.Context, req *alexandraProto.RemoveTriggerRequest) (*alexandraProto.Empty, error) {
	_, err := s.clientsMgr.getClientTo(s.lighthouse).Client.RemoveTrigger(ctx, &api.RemoveTriggerRequest{
		Keygroup:  req.Keygroup,
		TriggerId: req.TriggerId,
	})

	if err != nil {
		return nil, err
	}

	return &alexandraProto.Empty{}, err
}

// AddUser adds permissions to access a keygroup for a particular user to FReD.
func (s *Server) AddUser(ctx context.Context, req *alexandraProto.UserRequest) (*alexandraProto.Empty, error) {
	_, err := s.clientsMgr.getClientTo(s.lighthouse).Client.AddUser(ctx, &api.AddUserRequest{
		User:     req.User,
		Keygroup: req.Keygroup,
		Role:     api.UserRole(req.Role),
	})

	if err != nil {
		return nil, err
	}

	return &alexandraProto.Empty{}, err
}

// RemoveUser removes permissions to access a keygroup for a particular user from FReD.
func (s *Server) RemoveUser(ctx context.Context, req *alexandraProto.UserRequest) (*alexandraProto.Empty, error) {
	_, err := s.clientsMgr.getClientTo(s.lighthouse).Client.RemoveUser(ctx, &api.RemoveUserRequest{
		User:     req.User,
		Keygroup: req.Keygroup,
		Role:     api.UserRole(req.Role),
	})

	if err != nil {
		return nil, err
	}

	return &alexandraProto.Empty{}, err
}
