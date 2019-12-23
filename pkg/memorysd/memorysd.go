package memorysd

import (
	"errors"
	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/data"
	"sync"
)

// Storage stores a map of keygroup by name.
type Storage struct {
	keygroups map[string]Keygroup
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
		keygroups: make(map[string]Keygroup),
	}

	return
}

// Read returns an item with the specified id from the specified keygroup.
func (s *Storage) Read(i data.Item) (data.Item, error) {
	s.RLock()
	kg, ok := s.keygroups[i.Keygroup]
	s.RUnlock()

	if !ok {
		return i, errors.New("memorysd: no such keygroup")
	}

	kg.RLock()
	var value string
	value, ok = kg.items[i.ID]
	kg.RUnlock()

	if !ok {
		return i, errors.New("memorysd: no such item")
	}

	i.Data = value

	return i, nil
}

// Update updates the item with the specified id in the specified keygroup.
func (s *Storage) Update(i data.Item) error {
	s.RLock()
	kg, ok := s.keygroups[i.Keygroup]

	if !ok {
		s.RUnlock()
		return errors.New("memorysd: no such keygroup")
	}

	s.RUnlock()

	kg.Lock()

	kg.items[i.ID] = i.Data

	kg.Unlock()

	return nil
}

// Delete deletes the item with the specified id from the specified keygroup.
func (s *Storage) Delete(i data.Item) error {
	s.RLock()
	kg, ok := s.keygroups[i.Keygroup]

	if !ok {
		s.RUnlock()
		return errors.New("memorysd: no such keygroup")
	}

	s.RUnlock()

	kg.RLock()
	_, ok = kg.items[i.ID]
	kg.RUnlock()

	if !ok {
		return errors.New("memorysd: no such item")
	}

	kg.Lock()
	delete(kg.items, i.ID)
	kg.Unlock()

	return nil

}

// Exists checks if the given item exists in the given keygroups map.
func (s *Storage) Exists(i data.Item) bool {
	s.RLock()
	kg, ok := s.keygroups[i.Keygroup]
	s.RUnlock()

	if !ok {
		return false
	}

	kg.RLock()
	_, ok = kg.items[i.ID]
	kg.RUnlock()

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
func (s *Storage) CreateKeygroup(i data.Item) error {
	s.RLock()
	kg, exists := s.keygroups[i.Keygroup]

	if exists {
		s.RUnlock()
		return errors.New("memorysd: keygroup exists")
	}

	s.RUnlock()

	kg = Keygroup{
		items: make(map[string]string),
	}

	s.Lock()
	s.keygroups[i.Keygroup] = kg
	s.Unlock()

	return nil
}

// DeleteKeygroup removes the keygroup with the specified name from Storage.
func (s *Storage) DeleteKeygroup(i data.Item) error {
	s.RLock()
	_, ok := s.keygroups[i.Keygroup]
	s.RUnlock()

	if !ok {
		return errors.New("kmemorysd: keygroup does not exist")
	}

	s.Lock()
	delete(s.keygroups, i.Keygroup)
	s.Unlock()

	return nil
}
