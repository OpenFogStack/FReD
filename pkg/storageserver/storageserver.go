package storageserver

import (
	"context"

	"git.tu-berlin.de/mcc-fred/fred/pkg/fred"
	"git.tu-berlin.de/mcc-fred/fred/proto/storage"
	"github.com/DistributedClocks/GoVector/govec/vclock"
	"github.com/rs/zerolog/log"
)

// Server implements the DatabaseServer interface
type Server struct {
	store fred.Store
}

// NewStorageServer creates a new Server to serve GRPC Requests. It answers according to the data/Storage Interface
func NewStorageServer(store *fred.Store) *Server {
	log.Debug().Msgf("Setting up new Server with store %#v", *store)
	return &Server{store: *store}
}

// Update calls specific method of the storage interface
func (s *Server) Update(_ context.Context, req *storage.UpdateRequest) (*storage.UpdateResponse, error) {
	log.Debug().Msgf("GRPCServer: Update in=%#v", req)

	err := s.store.Update(req.Keygroup, req.Id, req.Val, int(req.Expiry), vclock.VClock{}.CopyFromMap(req.Version))

	if err != nil {
		log.Err(err).Msgf("GRPCServer has encountered an error while updating item %#v", req)
		return nil, err
	}
	return &storage.UpdateResponse{}, nil
}

// Append calls specific method of the storage interface
func (s *Server) Append(_ context.Context, req *storage.AppendRequest) (*storage.AppendResponse, error) {
	log.Debug().Msgf("GRPCServer: Append in=%#v", req)

	err := s.store.Append(req.Keygroup, req.Id, req.Val, int(req.Expiry))

	if err != nil {
		log.Err(err).Msgf("GRPCServer has encountered an error while appending item %#v", req)
		return nil, err
	}

	return &storage.AppendResponse{}, nil
}

// Delete calls specific method of the storage interface
func (s Server) Delete(_ context.Context, req *storage.DeleteRequest) (*storage.DeleteResponse, error) {
	log.Debug().Msgf("GRPCServer: Delete in=%#v", req)

	err := s.store.Delete(req.Keygroup, req.Id, vclock.VClock{}.CopyFromMap(req.Version))
	if err != nil {
		log.Err(err).Msgf("GRPCServer has encountered an error while deleting item %#v", req)
		return nil, err
	}
	return &storage.DeleteResponse{}, nil
}

// Read calls specific method of the storage interface
func (s Server) Read(_ context.Context, req *storage.ReadRequest) (*storage.ReadResponse, error) {
	log.Debug().Msgf("GRPCServer: Read in=%#v", req)
	vals, vvectors, found, err := s.store.Read(req.Keygroup, req.Id)

	if err != nil {
		log.Err(err).Msgf("GRPCServer has encountered an error while reading item %#v", req)
		return nil, err
	}

	if !found {
		return &storage.ReadResponse{
			Items: []*storage.Item{},
		}, nil
	}

	items := make([]*storage.Item, len(vals))

	for i := range vals {
		items[i] = &storage.Item{
			Keygroup: req.Keygroup,
			Id:       req.Id,
			Val:      vals[i],
			Version:  vvectors[i],
		}
	}

	return &storage.ReadResponse{
		Items: items,
	}, nil
}

// Scan calls specific method of the storage interface
func (s Server) Scan(_ context.Context, req *storage.ScanRequest) (*storage.ScanResponse, error) {
	// Stream: call server.send for every item, return if none left.
	log.Debug().Msgf("GRPCServer: Scan in=%#v", req)

	vals, vvectors, err := s.store.ReadSome(req.Keygroup, req.Start, req.Count)

	if err != nil {
		log.Err(err).Msgf("GRPCServer has encountered an error while scanning %d items from keygroup %#v", req.Count, req.Keygroup)
		return nil, err
	}
	items := make([]*storage.Item, 0)

	for key, data := range vals {
		for i := range data {
			items = append(items, &storage.Item{
				Keygroup: req.Keygroup,
				Id:       key,
				Val:      vals[key][i],
				Version:  vvectors[key][i],
			})
		}
	}

	return &storage.ScanResponse{
		Items: items,
	}, nil
}

// ReadAll calls specific method of the storage interface
func (s Server) ReadAll(_ context.Context, req *storage.ReadAllRequest) (*storage.ReadAllResponse, error) {
	// Stream: call server.send for every item, return if none left.
	log.Debug().Msgf("GRPCServer: ReadAll in=%#v", req)
	vals, vvectors, err := s.store.ReadAll(req.Keygroup)
	if err != nil {
		log.Err(err).Msgf("GRPCServer has encountered an error while reading whole keygroup %#v", req)
		return nil, err
	}
	items := make([]*storage.Item, 0)

	for key, data := range vals {
		for i := range data {
			items = append(items, &storage.Item{
				Keygroup: req.Keygroup,
				Id:       key,
				Val:      vals[key][i],
				Version:  vvectors[key][i],
			})
		}
	}

	return &storage.ReadAllResponse{
		Items: items,
	}, nil
}

// IDs calls specific method of the storage interface
func (s Server) IDs(_ context.Context, req *storage.IDsRequest) (*storage.IDsResponse, error) {
	log.Debug().Msgf("GRPCServer: IDs in=%#v", req)
	ids, err := s.store.IDs(req.Keygroup)
	if err != nil {
		log.Err(err).Msgf("GRPCServer has encountered an error while reading IDs %#v", req)
		return nil, err
	}

	return &storage.IDsResponse{
		Ids: ids,
	}, nil
}

// Exists calls specific method of the storage interface
func (s Server) Exists(_ context.Context, req *storage.ExistsRequest) (*storage.ExistsResponse, error) {
	log.Debug().Msgf("GRPCServer: Exists in=%#v", req)

	exists := s.store.Exists(req.Keygroup, req.Id)

	return &storage.ExistsResponse{Exists: exists}, nil
}

// CreateKeygroup calls specific method of the storage interface
func (s Server) CreateKeygroup(_ context.Context, req *storage.CreateKeygroupRequest) (*storage.CreateKeygroupResponse, error) {
	log.Debug().Msgf("GRPCServer: CreateKeygroup in=%#v", req)
	err := s.store.CreateKeygroup(req.Keygroup)

	if err != nil {
		log.Err(err).Msgf("GRPCServer has encountered an error while creating keygroup %#v", req)
		return nil, err
	}

	return &storage.CreateKeygroupResponse{}, nil
}

// DeleteKeygroup calls specific method of the storage interface
func (s Server) DeleteKeygroup(_ context.Context, req *storage.DeleteKeygroupRequest) (*storage.DeleteKeygroupResponse, error) {
	log.Debug().Msgf("GRPCServer: DeleteKeygroup in=%#v", req)
	err := s.store.DeleteKeygroup(req.Keygroup)
	if err != nil {
		log.Err(err).Msgf("GRPCServer has encountered an error while deleting keygroup %#v", req)
		return nil, err
	}
	return &storage.DeleteKeygroupResponse{}, nil
}

// ExistsKeygroup calls specific method of the storage interface
func (s Server) ExistsKeygroup(_ context.Context, req *storage.ExistsKeygroupRequest) (*storage.ExistsKeygroupResponse, error) {
	log.Debug().Msgf("GRPCServer: ExistsKeygroup in=%#v", req)

	exists := s.store.ExistsKeygroup(req.Keygroup)

	return &storage.ExistsKeygroupResponse{Exists: exists}, nil
}

// AddKeygroupTrigger calls specific method of the storage interface.
func (s *Server) AddKeygroupTrigger(_ context.Context, req *storage.AddKeygroupTriggerRequest) (*storage.AddKeygroupTriggerResponse, error) {
	log.Debug().Msgf("GRPCServer: AddKeygroupTrigger in=%#v", req)
	err := s.store.AddKeygroupTrigger(req.Keygroup, req.Id, req.Host)

	if err != nil {
		log.Err(err).Msgf("GRPCServer has encountered an error while adding trigger %#v", req)
		return nil, err
	}

	return &storage.AddKeygroupTriggerResponse{}, nil
}

// DeleteKeygroupTrigger calls specific method of the storage interface.
func (s *Server) DeleteKeygroupTrigger(_ context.Context, req *storage.DeleteKeygroupTriggerRequest) (*storage.DeleteKeygroupTriggerResponse, error) {
	log.Debug().Msgf("GRPCServer: DeleteKeygroupTrigger in=%#v", req)

	err := s.store.DeleteKeygroupTrigger(req.Keygroup, req.Id)

	if err != nil {
		log.Err(err).Msgf("GRPCServer has encountered an error while deleting trigger %#v", req)
		return nil, err
	}

	return &storage.DeleteKeygroupTriggerResponse{}, nil
}

// GetKeygroupTrigger calls specific method of the storage interface.
func (s *Server) GetKeygroupTrigger(_ context.Context, req *storage.GetKeygroupTriggerRequest) (*storage.GetKeygroupTriggerResponse, error) {
	// Steam: call server.send for every trigger, return if none left.
	log.Debug().Msgf("GRPCServer: GetKeygroupTrigger in=%#v", req)
	res, err := s.store.GetKeygroupTrigger(req.Keygroup)
	if err != nil {
		log.Err(err).Msgf("GRPCServer has encountered an error while reading all triggers for keygroup %#v", req)
		return nil, err
	}

	triggers := make([]*storage.Trigger, 0, len(res))

	for id, host := range res {
		triggers = append(triggers, &storage.Trigger{
			Id:   id,
			Host: host,
		})
	}

	return &storage.GetKeygroupTriggerResponse{Triggers: triggers}, nil
}
