package alexandra

import (
	"fmt"
	"sync"

	"git.tu-berlin.de/mcc-fred/vclock"
	"github.com/rs/zerolog/log"
)

type vcache struct {
	keygroups map[string]struct {
		items map[string]*struct {
			clocks []vclock.VClock
			*sync.RWMutex
		}
		*sync.RWMutex
	}
	*sync.RWMutex
}

func newVCache() *vcache {
	return &vcache{
		keygroups: make(map[string]struct {
			items map[string]*struct {
				clocks []vclock.VClock
				*sync.RWMutex
			}
			*sync.RWMutex
		}),
		RWMutex: &sync.RWMutex{},
	}
}

func (vc *vcache) ensureExists(kg string, id string) {
	// first check if the keygroup exists
	// this requires rlocking the cache
	vc.RLock()
	if _, ok := vc.keygroups[kg]; !ok {
		// if it doesn't exist, we need to upgrade to a write lock
		vc.RUnlock()
		vc.Lock()
		// and check again, because it might have been created in the meantime
		if _, ok := vc.keygroups[kg]; !ok {
			// if it still doesn't exist, we need to create it
			vc.keygroups[kg] = struct {
				items map[string]*struct {
					clocks []vclock.VClock
					*sync.RWMutex
				}
				*sync.RWMutex
			}{
				items: make(map[string]*struct {
					clocks []vclock.VClock
					*sync.RWMutex
				}),
				RWMutex: &sync.RWMutex{},
			}
		}
		// then go back to a read lock
		vc.Unlock()
		vc.RLock()
	}

	// now we can check if the item exists in the keygroup cache
	// this requires rlocking the keygroup cache
	vc.keygroups[kg].RLock()
	if _, ok := vc.keygroups[kg].items[id]; !ok {
		// if it doesn't exist, we need to upgrade to a write lock
		vc.keygroups[kg].RUnlock()
		vc.keygroups[kg].Lock()
		// and check again, because it might have been created in the meantime
		if _, ok := vc.keygroups[kg].items[id]; !ok {
			// if it still doesn't exist, we need to create it
			vc.keygroups[kg].items[id] = &struct {
				clocks []vclock.VClock
				*sync.RWMutex
			}{
				clocks: make([]vclock.VClock, 0),
				RWMutex:  &sync.RWMutex{},
			}
		}
		// then go back to a read lock
		vc.keygroups[kg].Unlock()
		vc.keygroups[kg].RLock()
	}
}

func (vc *vcache) cRLock(kg string, id string) {
	vc.ensureExists(kg, id)

	// now we can read/write-lock the item
	vc.keygroups[kg].items[id].RLock()
}
func (vc *vcache) cRUnlock(kg string, id string) {
	// need to assume that it exists and was already locked
	vc.keygroups[kg].items[id].RUnlock()
	vc.keygroups[kg].RUnlock()
	vc.RUnlock()
}

func (vc *vcache) cLock(kg string, id string) {
	vc.ensureExists(kg, id)

	// now we can read/write-lock the item
	vc.keygroups[kg].items[id].Lock()
}

func (vc *vcache) cUnlock(kg string, id string) {
	// need to assume that it exists and was already locked
	vc.keygroups[kg].items[id].Unlock()
	vc.keygroups[kg].RUnlock()
	vc.RUnlock()
}

func (vc *vcache) isLocked(_ string, _ string) bool {
	// TODO: the sync.RWMutex does not provide a way to check if it is locked
	// so we need to find a way to do this ourselves
	return true
}

func (vc *vcache) add(kg string, id string, version vclock.VClock) error {

	if !vc.isLocked(kg, id) {
		return fmt.Errorf("add failed because item %s in keygroup %s is not locked in cache", id, kg)
	}

	newClocks := make([]vclock.VClock, 0, len(vc.keygroups[kg].items[id].clocks))

	// let's first go to the existing versions
	for _, v := range vc.keygroups[kg].items[id].clocks {
		switch v.Order(version) {
		// seems like we have that version in store already!
		// just return, nothing to do here
		case vclock.Equal:
			{
				return nil
			}
		// if there is a concurrent version to our new version, we need to keep that
		case vclock.Concurrent:
			{
				newClocks = append(newClocks, v)
				continue
			}
		// if by chance we come upon a newer version, this has all been pointless
		// actually, this is a bad sign: we seem to have read outdated data!
		case vclock.Ancestor:
			{
				return fmt.Errorf("add failed because your version %v of %s in keygroup %s is older than cached version %v", version, id, kg, v)
			}
		}
	}

	// now actually add the new version
	newClocks = append(newClocks, version.Copy())

	log.Debug().Msgf("Cache Add (kg %s item %s): new: %+v old: %+v", kg, id, newClocks, vc.keygroups[kg].items[id].clocks)

	// and store our new clocks
	vc.keygroups[kg].items[id].clocks = newClocks
	return nil
}

func (vc *vcache) supersede(kg string, id string, known []vclock.VClock, version vclock.VClock) error {

	if !vc.isLocked(kg, id) {
		return fmt.Errorf("supersede failed because item %s in keygroup %s is not locked in cache", id, kg)
	}

	newClocks := []vclock.VClock{
		version.Copy(),
	}

	// add all clocks from cache as well, unless they are in the "known" array
	for _, v := range vc.keygroups[kg].items[id].clocks {
		discard := false
		for _, k := range known {
			if v.Compare(k, vclock.Equal) {
				discard = true
				break
			}
		}
		if discard {
			continue
		}

		newClocks = append(newClocks, v)
	}

	vc.keygroups[kg].items[id].clocks = newClocks

	return nil
}

func (vc *vcache) get(kg string, id string) ([]vclock.VClock, error) {
	if !vc.isLocked(kg, id) {
		return nil, fmt.Errorf("get failed because item %s in keygroup %s is not locked in cache", id, kg)
	}

	return vc.keygroups[kg].items[id].clocks, nil
}
