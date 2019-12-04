package memorykg

import (
	"sync"
)

// KeygroupStorage saves a set of all available keygroups.
type KeygroupStorage struct {
	keygroups map[string]struct{}
	sync.RWMutex
}

// New create a new KeygroupStorage.
func New() (k *KeygroupStorage) {
	k = &KeygroupStorage{
		keygroups: make(map[string]struct{}),
	}

	return
}

// Create adds a keygroup to the KeygroupStorage.
func (k *KeygroupStorage) Create(kgname string) error {
	k.RLock()
	_, ok := k.keygroups[kgname]
	k.RUnlock()

	if ok {
		return nil
	}

	k.Lock()
	k.keygroups[kgname] = struct{}{}
	k.Unlock()

	return nil
}

// Delete removes a keygroup from the KeygroupStorage.
func (k *KeygroupStorage) Delete(kgname string) error {
	k.RLock()
	_, ok := k.keygroups[kgname]
	k.RUnlock()

	if !ok {
		return nil
	}

	k.Lock()
	delete(k.keygroups, kgname)
	k.Unlock()

	return nil
}

// Exists checks if a keygroup exists in the KeygroupStorage.
func (k *KeygroupStorage) Exists(kgname string) bool {
	k.RLock()
	_, ok := k.keygroups[kgname]
	k.RUnlock()

	return ok
}
