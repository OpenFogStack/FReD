package fred

import (
	"sync"

	"github.com/rs/zerolog/log"
)

// keygroupStore saves a set of all available keygroups.
type keygroupStore struct {
	keygroups map[KeygroupName]struct{}
	sync.RWMutex
	nodeID string
}

// newKeygroupStore creates a new KeygroupStorage.
func newKeygroupStore(n string) (kS *keygroupStore) {
	kS = &keygroupStore{
		keygroups: make(map[KeygroupName]struct{}),
		nodeID:    n,
	}

	return
}

// create adds a keygroup to the KeygroupStorage.
func (kS *keygroupStore) create(k Keygroup) error {
	err := checkKeygroup(k.Name)

	if err != nil {
		log.Err(err).Msg("Keygroup service cannot create a new keygroup")
		return err
	}

	log.Debug().Msgf("CreateKeygroup from memorykg: in %#v", k)
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

// delete removes a keygroup from the KeygroupStorage.
func (kS *keygroupStore) delete(k Keygroup) error {
	err := checkKeygroup(k.Name)

	if err != nil {
		log.Err(err).Msg("Keygroup service cannot delete a keygroup")
		return err
	}

	log.Debug().Msgf("DeleteKeygroup from memorykg: in %#v", k)
	kS.RLock()
	_, ok := kS.keygroups[k.Name]
	kS.RUnlock()

	if !ok {
		return newError(StatusNotFound, "no such keygroup")
	}

	kS.Lock()
	delete(kS.keygroups, k.Name)
	kS.Unlock()

	return nil
}

// exists checks if a keygroup exists in the KeygroupStorage.
func (kS *keygroupStore) exists(k Keygroup) bool {
	kS.RLock()
	_, ok := kS.keygroups[k.Name]
	kS.RUnlock()

	log.Debug().Msgf("Exists from memorykg: in %#v, out %#v", k, ok)

	return ok
}
