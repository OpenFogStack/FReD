package storageserver

import (
	"context"

	"git.tu-berlin.de/mcc-fred/fred/pkg/fred"
	"git.tu-berlin.de/mcc-fred/fred/proto/storage"
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
func (s *Server) Update(_ context.Context, item *storage.UpdateItem) (*storage.Response, error) {
	log.Debug().Msgf("GRPCServer: Update in=%#v", item)

	err := s.store.Update(item.Keygroup, item.Id, item.Val, item.Append, int(item.Expiry))
	if err != nil {
		log.Err(err).Msgf("GRPCServer has encountered an error while updating item %#v", item)
		return &storage.Response{Success: false}, err
	}
	return &storage.Response{Success: true}, nil
}

// Append calls specific method of the storage interface
func (s *Server) Append(_ context.Context, item *storage.AppendItem) (*storage.Key, error) {
	log.Debug().Msgf("GRPCServer: Append in=%#v", item)

	res, err := s.store.Append(item.Keygroup, item.Val, int(item.Expiry))

	if err != nil {
		log.Err(err).Msgf("GRPCServer has encountered an error while appending item %#v", item)
		return &storage.Key{}, err
	}

	return &storage.Key{
		Keygroup: item.Keygroup,
		Id:       res,
	}, nil
}

// Delete calls specific method of the storage interface
func (s Server) Delete(_ context.Context, key *storage.Key) (*storage.Response, error) {
	log.Debug().Msgf("GRPCServer: Delete in=%#v", key)
	err := s.store.Delete(key.Keygroup, key.Id)
	if err != nil {
		log.Err(err).Msgf("GRPCServer has encountered an error while deleting item %#v", key)
		return &storage.Response{Success: false, Message: "Server has encountered an error while deleting an item"}, err
	}
	return &storage.Response{Success: true}, nil
}

// Read calls specific method of the storage interface
func (s Server) Read(_ context.Context, key *storage.Key) (*storage.Val, error) {
	log.Debug().Msgf("GRPCServer: Read in=%#v", key)
	res, err := s.store.Read(key.Keygroup, key.Id)
	if err != nil {
		log.Err(err).Msgf("GRPCServer has encountered an error while reading item %#v", key)
		return &storage.Val{}, err
	}
	return &storage.Val{Val: res}, nil
}

// Scan calls specific method of the storage interface
func (s Server) Scan(req *storage.ScanRequest, server storage.Database_ScanServer) error {
	// Stream: call server.send for every item, return if none left.
	log.Debug().Msgf("GRPCServer: Scan in=%#v", req)
	res, err := s.store.ReadSome(req.Key.Keygroup, req.Key.Id, req.Count)
	if err != nil {
		log.Err(err).Msgf("GRPCServer has encountered an error while scanning %d items from keygroup %#v", req.Count, req.Key.Keygroup)
		return err
	}
	for id, elem := range res {
		err := server.Send(&storage.Item{
			Id:  id,
			Val: elem,
		})
		if err != nil {
			return err
		}
	}
	// Return nil == successful transfer
	return nil
}

// ReadAll calls specific method of the storage interface
func (s Server) ReadAll(kg *storage.Keygroup, server storage.Database_ReadAllServer) error {
	// Stream: call server.send for every item, return if none left.
	log.Debug().Msgf("GRPCServer: ReadAll in=%#v", kg)
	res, err := s.store.ReadAll(kg.Keygroup)
	if err != nil {
		log.Err(err).Msgf("GRPCServer has encountered an error while reading whole keygroup %#v", kg)
		return err
	}
	for id, elem := range res {
		err := server.Send(&storage.Item{
			Id:  id,
			Val: elem,
		})
		if err != nil {
			return err
		}
	}
	// Return nil == successful transfer
	return nil
}

// IDs calls specific method of the storage interface
func (s Server) IDs(kg *storage.Keygroup, server storage.Database_IDsServer) error {
	log.Debug().Msgf("GRPCServer: IDs in=%#v", kg)
	res, err := s.store.IDs(kg.Keygroup)
	if err != nil {
		log.Err(err).Msgf("GRPCServer has encountered an error while reading IDs %#v", kg)
		return err
	}
	for _, elem := range res {
		err := server.Send(&storage.Key{Id: elem})
		if err != nil {
			return err
		}
	}
	// Return nil == successful transfer
	return nil
}

// Exists calls specific method of the storage interface
func (s Server) Exists(_ context.Context, key *storage.Key) (*storage.Response, error) {
	log.Debug().Msgf("GRPCServer: Exists in=%#v", key)
	exists := s.store.Exists(key.Keygroup, key.Id)
	return &storage.Response{Success: exists}, nil
}

// CreateKeygroup calls specific method of the storage interface
func (s Server) CreateKeygroup(_ context.Context, kg *storage.Keygroup) (*storage.Response, error) {
	log.Debug().Msgf("GRPCServer: CreateKeygroup in=%#v", kg)
	err := s.store.CreateKeygroup(kg.Keygroup)
	if err != nil {
		log.Err(err).Msgf("GRPCServer has encountered an error while creating keygroup %#v", kg)
		return &storage.Response{Success: false}, err
	}
	return &storage.Response{Success: true}, nil
}

// DeleteKeygroup calls specific method of the storage interface
func (s Server) DeleteKeygroup(_ context.Context, kg *storage.Keygroup) (*storage.Response, error) {
	log.Debug().Msgf("GRPCServer: DeleteKeygroup in=%#v", kg)
	err := s.store.DeleteKeygroup(kg.Keygroup)
	if err != nil {
		log.Err(err).Msgf("GRPCServer has encountered an error while deleting keygroup %#v", kg)
		return &storage.Response{Success: false}, err
	}
	return &storage.Response{Success: true}, nil
}

// ExistsKeygroup calls specific method of the storage interface
func (s Server) ExistsKeygroup(ctx context.Context, kg *storage.Keygroup) (*storage.Response, error) {
	log.Debug().Msgf("GRPCServer: ExistsKeygroup in=%#v", kg)
	exists := s.store.ExistsKeygroup(kg.Keygroup)
	return &storage.Response{Success: exists}, nil
}

// AddKeygroupTrigger calls specific method of the storage interface.
func (s *Server) AddKeygroupTrigger(ctx context.Context, t *storage.KeygroupTrigger) (*storage.Response, error) {
	log.Debug().Msgf("GRPCServer: AddKeygroupTrigger in=%#v", t)
	err := s.store.AddKeygroupTrigger(t.Keygroup, t.Trigger.Id, t.Trigger.Host)
	if err != nil {
		log.Err(err).Msgf("GRPCServer has encountered an error while adding trigger %#v", t)
		return &storage.Response{Success: false}, err
	}
	return &storage.Response{Success: true}, nil
}

// DeleteKeygroupTrigger calls specific method of the storage interface.
func (s *Server) DeleteKeygroupTrigger(ctx context.Context, t *storage.KeygroupTrigger) (*storage.Response, error) {
	log.Debug().Msgf("GRPCServer: DeleteKeygroupTrigger in=%#v", t)
	err := s.store.DeleteKeygroupTrigger(t.Keygroup, t.Trigger.Id)
	if err != nil {
		log.Err(err).Msgf("GRPCServer has encountered an error while deleting trigger %#v", t)
		return &storage.Response{Success: false}, err
	}
	return &storage.Response{Success: true}, nil
}

// GetKeygroupTrigger calls specific method of the storage interface.
func (s *Server) GetKeygroupTrigger(kg *storage.Keygroup, server storage.Database_GetKeygroupTriggerServer) error {
	// Steam: call server.send for every trigger, return if none left.
	log.Debug().Msgf("GRPCServer: GetKeygroupTrigger in=%#v", kg)
	res, err := s.store.GetKeygroupTrigger(kg.Keygroup)
	if err != nil {
		log.Err(err).Msgf("GRPCServer has encountered an error while reading all triggers for keygroup %#v", kg)
		return err
	}
	for id, host := range res {
		err := server.Send(&storage.Trigger{
			Id:   id,
			Host: host,
		})
		if err != nil {
			return err
		}
	}
	// Return nil == successful transfer
	return nil
}
