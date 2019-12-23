package memorykg

import (
	"errors"
	"sync"

	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/keygroup"
)

// KeygroupStorage saves a set of all available keygroup.
type KeygroupStorage struct {
	keygroups map[string]struct{}
	sync.RWMutex
}

// New creates a new KeygroupStorage.
func New() (kS *KeygroupStorage) {
	kS = &KeygroupStorage{
		keygroups: make(map[string]struct{}),
	}

	return
}

// Create adds a keygroup to the KeygroupStorage.
func (kS *KeygroupStorage) Create(k keygroup.Keygroup) error {
	kS.RLock()
	_, ok := kS.keygroups[k.Name]
	kS.RUnlock()

	if ok {
		return nil
	}

	kS.Lock()
	kS.keygroups[k.Name] = struct{}{}
	kS.Unlock()

	return nil
}

// Delete removes a keygroup from the KeygroupStorage.
func (kS *KeygroupStorage) Delete(k keygroup.Keygroup) error {
	kS.RLock()
	_, ok := kS.keygroups[k.Name]
	kS.RUnlock()

	if !ok {
		return errors.New("memorykg: no such keygroup")
	}

	kS.Lock()
	delete(kS.keygroups, k.Name)
	kS.Unlock()

	return nil
}

// Exists checks if a keygroup exists in the KeygroupStorage.
func (kS *KeygroupStorage) Exists(k keygroup.Keygroup) bool {
	kS.RLock()
	_, ok := kS.keygroups[k.Name]
	kS.RUnlock()

	return ok
}
