package memorysd

import (
	"sync"

	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/commons"
	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/data"
	errors "gitlab.tu-berlin.de/mcc-fred/fred/pkg/errors"
)

// Storage stores a map of keygroup by name.
type Storage struct {
	keygroups map[commons.KeygroupName]Keygroup
	sync.RWMutex
}

// Keygroup stores a map of data by id, has a maxkey to keep track of the unused ids.
type Keygroup struct {
	items map[string]string
	sync.RWMutex
}

// New create a new Storage.
func New() (s *Storage) {
	s = &Storage{
		keygroups: make(map[commons.KeygroupName]Keygroup),
	}

	return
}

// Read returns an item with the specified id from the specified keygroup.
func (s *Storage) Read(kg commons.KeygroupName, id string) (string, error) {
	s.RLock()
	keygroup, ok := s.keygroups[kg]
	s.RUnlock()

	if !ok {
		return "", errors.New(errors.StatusNotFound, "memorysd: no such keygroup")
	}

	keygroup.RLock()
	var value string
	value, ok = keygroup.items[id]
	keygroup.RUnlock()

	if !ok {
		return "", errors.New(errors.StatusNotFound, "memorysd: no such item")
	}

	return value, nil
}

// ReadAll returns all items in the specified keygroup.
func (s *Storage) ReadAll(kg commons.KeygroupName) ([]data.Item, error) {
	s.RLock()
	keygroup, ok := s.keygroups[kg]
	s.RUnlock()

	if !ok {
		return nil, errors.New(errors.StatusNotFound, "memorysd: no such keygroup")
	}

	keygroup.RLock()
	items := make([]data.Item, len(keygroup.items))

	x := 0

	for k := range keygroup.items {
		items[x] = data.Item{
			Keygroup: kg,
			ID:       k,
			Data:     keygroup.items[k],
		}
		x++
	}

	keygroup.RUnlock()

	return items, nil
}

// IDs returns the keys of all items in the specified keygroup.
func (s *Storage) IDs(kg commons.KeygroupName) ([]data.Item, error) {
	s.RLock()
	keygroup, ok := s.keygroups[kg]
	s.RUnlock()

	if !ok {
		return nil, errors.New(errors.StatusNotFound, "memorysd: no such keygroup")
	}

	keygroup.RLock()
	keys := make([]data.Item, len(keygroup.items))

	x := 0

	for k := range keygroup.items {
		keys[x] = data.Item{
			Keygroup: kg,
			ID:       k,
		}
		x++
	}
	keygroup.RUnlock()

	return keys, nil
}

// Update updates the item with the specified id in the specified keygroup.
func (s *Storage) Update(i data.Item) error {
	s.RLock()
	kg, ok := s.keygroups[i.Keygroup]

	if !ok {
		s.RUnlock()
		return errors.New(errors.StatusNotFound, "memorysd: no such keygroup")
	}

	s.RUnlock()

	kg.Lock()

	kg.items[i.ID] = i.Data

	kg.Unlock()

	return nil
}

// Delete deletes the item with the specified id from the specified keygroup.
func (s *Storage) Delete(kg commons.KeygroupName, id string) error {
	s.RLock()
	keygroup, ok := s.keygroups[kg]

	if !ok {
		s.RUnlock()
		return errors.New(errors.StatusNotFound, "memorysd: no such keygroup")
	}

	s.RUnlock()

	keygroup.RLock()
	_, ok = keygroup.items[id]
	keygroup.RUnlock()

	if !ok {
		return errors.New(errors.StatusNotFound, "memorysd: no such item")
	}

	keygroup.Lock()
	delete(keygroup.items, id)
	keygroup.Unlock()

	return nil

}

// Exists checks if the given item exists in the given keygroups map.
func (s *Storage) Exists(kg commons.KeygroupName, id string) bool {
	s.RLock()
	keygroup, ok := s.keygroups[kg]
	s.RUnlock()

	if !ok {
		return false
	}

	keygroup.RLock()
	_, ok = keygroup.items[id]
	keygroup.RUnlock()

	return ok
}

// ExistsKeygroup checks if the keygroup exists in the map.
func (s *Storage) ExistsKeygroup(i data.Item) bool {
	s.RLock()
	_, ok := s.keygroups[i.Keygroup]
	s.RUnlock()

	return ok

}

// CreateKeygroup creates a new keygroup with the specified name in Storage.
func (s *Storage) CreateKeygroup(kg commons.KeygroupName) error {
	s.RLock()
	keygroup, exists := s.keygroups[kg]

	if exists {
		s.RUnlock()
		return errors.New(errors.StatusConflict, "memorysd: keygroup exists")
	}

	s.RUnlock()

	keygroup = Keygroup{
		items: make(map[string]string),
	}

	s.Lock()
	s.keygroups[kg] = keygroup
	s.Unlock()

	return nil
}

// DeleteKeygroup removes the keygroup with the specified name from Storage.
func (s *Storage) DeleteKeygroup(kg commons.KeygroupName) error {
	s.RLock()
	_, ok := s.keygroups[kg]
	s.RUnlock()

	if !ok {
		return errors.New(errors.StatusNotFound, "kmemorysd: keygroup does not exist")
	}

	s.Lock()
	delete(s.keygroups, kg)
	s.Unlock()

	return nil
}
