package memorykg

import (
	"sync"
)

type KeygroupStorage struct {
	keygroups map[string]struct{}
	sync.RWMutex
}

func New() (k *KeygroupStorage) {
	k = &KeygroupStorage{
		keygroups: make(map[string]struct{}),
	}

	return
}

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

func (k *KeygroupStorage) Exists(kgname string) bool {
	k.RLock()
	_, ok := k.keygroups[kgname]
	k.RUnlock()

	return ok
}
