package fred

import (
	"sync"

	"git.tu-berlin.de/mcc-fred/fred/pkg/vector"
	"github.com/DistributedClocks/GoVector/govec/vclock"
	"github.com/go-errors/errors"
	"github.com/rs/zerolog/log"
)

// Store is an interface for the storage medium that the key-value val items are persisted on.
type Store interface {
	// Update Needs: keygroup, id, val
	Update(kg string, id string, val string, expiry int, vvector vclock.VClock) error
	// Delete Needs: keygroup, id
	Delete(kg string, id string, vvector vclock.VClock) error
	// Append Needs: keygroup, val, Returns: key
	Append(kg string, id string, val string, expiry int) error
	// Read Needs: keygroup, id; Returns: val, version vector, found
	Read(kg string, id string) ([]string, []vclock.VClock, bool, error)
	// ReadSome Needs: keygroup, id, range; Returns: ids, values, versions
	ReadSome(kg string, id string, count uint64) (map[string][]string, map[string][]vclock.VClock, error)
	// ReadAll Needs: keygroup; Returns: ids, values, versions
	ReadAll(kg string) (map[string][]string, map[string][]vclock.VClock, error)
	// IDs Needs: keygroup, Returns:[] keygroup, id
	IDs(kg string) ([]string, error)
	// Exists Needs: keygroup, id
	Exists(kg string, id string) bool
	// CreateKeygroup Needs: keygroup, Returns: err
	// Doesnt really need to store the KG itself, that is keygroup/store.go's job.
	// But might be useful for databases using it
	CreateKeygroup(kg string) error
	// DeleteKeygroup Same as with CreateKeygroup
	DeleteKeygroup(kg string) error
	// ExistsKeygroup Needs: keygroup
	ExistsKeygroup(kg string) bool
	// AddKeygroupTrigger Needs: keygroup, trigger node id, trigger node host
	AddKeygroupTrigger(kg string, id string, host string) error
	// DeleteKeygroupTrigger Needs: keygroup, trigger node id
	DeleteKeygroupTrigger(kg string, id string) error
	// GetKeygroupTrigger Needs: keygroup; Returns map: trigger node id -> trigger node host
	GetKeygroupTrigger(kg string) (map[string]string, error)
	// Close indicates that the underlying store should be closed as it is no longer needed.
	Close() error
}

type cachedClock struct {
	clocks []vclock.VClock
	*sync.Mutex
}

type keygroupCache struct {
	clocks map[string]*cachedClock
	*sync.RWMutex
}

type storeService struct {
	iS         Store
	id         string
	vCache     map[KeygroupName]keygroupCache
	vCacheLock sync.RWMutex
}

// NewStoreService creates a new val manipulation service.
func newStoreService(iS Store, id NodeID) *storeService {
	return &storeService{
		iS:         iS,
		id:         string(id),
		vCache:     make(map[KeygroupName]keygroupCache),
		vCacheLock: sync.RWMutex{},
	}
}

// Read returns an item from the key-value store.
func (s *storeService) read(kg KeygroupName, id string) ([]Item, error) {
	// TODO: make this part of the single writer thing?
	err := checkKGandID(kg, id)

	if err != nil {
		return nil, err
	}

	if !s.iS.ExistsKeygroup(string(kg)) {
		return nil, errors.Errorf("no such keygroup in store: %+v", kg)
	}

	data, vvectors, found, err := s.iS.Read(string(kg), id)

	if err != nil {
		return nil, err
	}

	if !found {
		return []Item{}, nil
	}

	items := make([]Item, len(data))

	for i := range data {
		items[i] = Item{
			Keygroup:   kg,
			ID:         id,
			Val:        data[i],
			Version:    vvectors[i],
			Tombstoned: data[i] == "",
		}
	}

	//return s.cleanOlderVersions(items)
	return items, nil
}

// readVersion reads data from the data store but also accepts a list of version vectors.
// TODO: what for? Currently it only returns values that are equal or newer than any of the given versions, i.e, it
// TODO: filters out concurrent versions.
func (s *storeService) readVersion(kg KeygroupName, id string, versions []vclock.VClock) ([]Item, error) {
	//TODO: make this part of the single writer thing?
	err := checkKGandID(kg, id)

	if err != nil {
		return nil, err
	}

	if !s.iS.ExistsKeygroup(string(kg)) {
		return nil, errors.Errorf("no such keygroup in store: %+v", kg)
	}

	data, vvectors, found, err := s.iS.Read(string(kg), id)

	if err != nil {
		return nil, err
	}

	if !found {
		return []Item{}, nil
	}

	var items []Item

	for i := range data {
		for _, v := range versions {
			if v.Compare(vvectors[i], vclock.Equal) || v.Compare(vvectors[i], vclock.Ancestor) {
				items = append(items, Item{
					Keygroup:   kg,
					ID:         id,
					Val:        data[i],
					Version:    vvectors[i],
					Tombstoned: data[i] == "",
				})
				continue
			}
		}
	}

	//return s.cleanOlderVersions(items)
	return items, nil
}

// Scan returns a list of count items starting with id from the key-value store.
func (s *storeService) scan(kg KeygroupName, id string, count uint64) ([]Item, error) {
	err := checkKGandID(kg, id)

	if err != nil {
		return nil, err
	}

	if !s.iS.ExistsKeygroup(string(kg)) {
		return nil, errors.Errorf("no such keygroup in store: %+v", kg)
	}

	if !s.exists(Item{
		Keygroup: kg,
		ID:       id,
	}) {
		return nil, errors.Errorf("no such item %s in keygroup %+v", id, kg)
	}

	data, vvectors, err := s.iS.ReadSome(string(kg), id, count)

	if err != nil {
		return nil, err
	}

	items := make([]Item, 0)

	for key, item := range data {
		for i := range item {
			items = append(items, Item{
				Keygroup:   kg,
				ID:         key,
				Val:        item[i],
				Version:    vvectors[key][i],
				Tombstoned: item[i] == "",
			})
		}
	}

	//return s.cleanOlderVersions(items)
	return items, nil
}

// ReadAll returns all items of a particular keygroup from the key-value store.
func (s *storeService) readAll(kg KeygroupName) ([]Item, error) {
	err := checkKeygroup(kg)

	if err != nil {
		return nil, err
	}

	if !s.iS.ExistsKeygroup(string(kg)) {
		return nil, errors.Errorf("no such keygroup in store: %+v", kg)
	}

	data, vvectors, err := s.iS.ReadAll(string(kg))

	if err != nil {
		return nil, err
	}

	items := make([]Item, 0)

	for key, item := range data {
		for i := range item {
			items = append(items, Item{
				Keygroup:   kg,
				ID:         key,
				Val:        item[i],
				Version:    vvectors[key][i],
				Tombstoned: item[i] == "",
			})
		}
	}

	//return s.cleanOlderVersions(items)
	return items, nil
}

// exists checks if an item exists in the key-value store.
func (s *storeService) exists(i Item) bool {
	return s.iS.Exists(string(i.Keygroup), i.ID)
}

// Update updates an item in the key-value store.
func (s *storeService) update(i Item, expiry int) (vclock.VClock, error) {
	// no version given means:
	// increment the locally stored version
	// check the cache for the local version number
	// nothing there? start at 1
	// TODO: in case of node failure, should first check in store if nothing is in cache
	// TODO: question is how do we find out something doesn't exist vs it is not in cache because we failed?
	err := checkItem(i)

	if err != nil {
		return nil, err
	}

	if !s.iS.ExistsKeygroup(string(i.Keygroup)) {
		return nil, errors.Errorf("no such keygroup in store: %+v", i.Keygroup)
	}

	s.vCacheLock.RLock()
	defer s.vCacheLock.RUnlock()

	s.vCache[i.Keygroup].RLock()
	defer s.vCache[i.Keygroup].RUnlock()

	if _, ok := s.vCache[i.Keygroup].clocks[i.ID]; !ok {
		s.vCache[i.Keygroup].RUnlock()
		s.vCache[i.Keygroup].Lock()
		if _, ok := s.vCache[i.Keygroup].clocks[i.ID]; !ok {
			s.vCache[i.Keygroup].clocks[i.ID] = &cachedClock{
				clocks: []vclock.VClock{},
				Mutex:  &sync.Mutex{},
			}
		}
		s.vCache[i.Keygroup].Unlock()
		s.vCache[i.Keygroup].RLock()
	}

	s.vCache[i.Keygroup].clocks[i.ID].Lock()
	defer s.vCache[i.Keygroup].clocks[i.ID].Unlock()

	log.Debug().Msgf("update: before known %+v", s.vCache[i.Keygroup].clocks[i.ID].clocks)

	oldVersion := vclock.VClock{}
	toPrune := make([]vclock.VClock, len(s.vCache[i.Keygroup].clocks[i.ID].clocks))
	for j, c := range s.vCache[i.Keygroup].clocks[i.ID].clocks {
		oldVersion.Merge(c)
		toPrune[j] = c.Copy()
	}

	newVersion := oldVersion
	newVersion.Tick(s.id)

	s.vCache[i.Keygroup].clocks[i.ID].clocks = []vclock.VClock{newVersion}

	err = s.iS.Update(string(i.Keygroup), i.ID, i.Val, expiry, newVersion.GetMap())

	if err != nil {
		return nil, err
	}

	s.prune(string(i.Keygroup), i.ID, toPrune)

	log.Debug().Msgf("update: after known %+v", s.vCache[i.Keygroup].clocks[i.ID].clocks)

	return newVersion.Copy(), nil
}

// updateVersions updates an item in the key-value store given a list of known versions.
func (s *storeService) updateVersions(i Item, versions []vclock.VClock, expiry int) (vclock.VClock, error) {
	// in this case, all existing versions will be removed and there will be a new item with a new version
	// for that to work, we first need to check whether for each of the version we have locally, there is a version in
	// the given versions that is greater or equal
	// the new version will be a merge of all given versions + incrementing our id
	// if any of the checks fail, we will reject the update
	// obviously rejections aren't pretty, but we just assume that two clients accessing the same data won't happen a lot
	// if they were, we could do siblings with dotted version vectors

	err := checkItem(i)

	if err != nil {
		return nil, err
	}

	if !s.iS.ExistsKeygroup(string(i.Keygroup)) {
		return nil, errors.Errorf("no such keygroup in store: %+v", i.Keygroup)
	}

	s.vCacheLock.RLock()
	defer s.vCacheLock.RUnlock()

	s.vCache[i.Keygroup].RLock()
	defer s.vCache[i.Keygroup].RUnlock()

	if _, ok := s.vCache[i.Keygroup].clocks[i.ID]; !ok {
		s.vCache[i.Keygroup].RUnlock()
		s.vCache[i.Keygroup].Lock()
		if _, ok := s.vCache[i.Keygroup].clocks[i.ID]; !ok {
			s.vCache[i.Keygroup].clocks[i.ID] = &cachedClock{
				clocks: []vclock.VClock{},
				Mutex:  &sync.Mutex{},
			}
		}
		s.vCache[i.Keygroup].Unlock()
		s.vCache[i.Keygroup].RLock()
	}

	s.vCache[i.Keygroup].clocks[i.ID].Lock()
	defer s.vCache[i.Keygroup].clocks[i.ID].Unlock()

	// ok, let's check if we want to reject this update!
	// go through all of our known versions and check that the client has seen it in some form by checking that the
	// given versions have one that is greater or equal than that
	for _, v1 := range s.vCache[i.Keygroup].clocks[i.ID].clocks {
		reject := true
		for _, v2 := range versions {
			if v1.Compare(v2, vclock.Equal) || v1.Compare(v2, vclock.Descendant) {
				reject = false
				break
			}
		}

		if reject {
			return nil, errors.Errorf("update rejected as seen version %+v is not covered by given versions %+v", v1, versions)
		}
	}

	// looks good! let's make up a new version for this update
	newVersion := vclock.VClock{}
	for _, v := range versions {
		newVersion.Merge(v)
	}

	newVersion.Tick(s.id)

	err = s.iS.Update(string(i.Keygroup), i.ID, i.Val, expiry, newVersion)

	if err != nil {
		return nil, err
	}

	// let us finally update our local cache by removing all the versions we had before
	s.prune(string(i.Keygroup), i.ID, s.vCache[i.Keygroup].clocks[i.ID].clocks)

	s.vCache[i.Keygroup].clocks[i.ID].clocks = []vclock.VClock{newVersion}

	return newVersion.Copy(), nil
}

func (s *storeService) addVersion(i Item, remoteVersion vclock.VClock, expiry int) error {
	// TODO
	err := checkItem(i)

	if err != nil {
		return err
	}

	if !s.iS.ExistsKeygroup(string(i.Keygroup)) {
		return errors.Errorf("no such keygroup in store: %+v", i.Keygroup)
	}

	s.vCacheLock.RLock()
	defer s.vCacheLock.RUnlock()

	s.vCache[i.Keygroup].RLock()
	defer s.vCache[i.Keygroup].RUnlock()

	if _, ok := s.vCache[i.Keygroup].clocks[i.ID]; !ok {
		s.vCache[i.Keygroup].RUnlock()
		s.vCache[i.Keygroup].Lock()
		if _, ok := s.vCache[i.Keygroup].clocks[i.ID]; !ok {
			s.vCache[i.Keygroup].clocks[i.ID] = &cachedClock{
				clocks: []vclock.VClock{},
				Mutex:  &sync.Mutex{},
			}
		}
		s.vCache[i.Keygroup].Unlock()
		s.vCache[i.Keygroup].RLock()
	}

	s.vCache[i.Keygroup].clocks[i.ID].Lock()
	defer s.vCache[i.Keygroup].clocks[i.ID].Unlock()

	// let's go through all the versions we have and compare to the version we got
	// if we already have a newer version than this one, just discard this
	// if we have an older one, replace it with the newer one
	// if none of that is the case, this is a concurrent version and we just store it

	log.Debug().Msgf("addVersion: before known %+v", s.vCache[i.Keygroup].clocks[i.ID].clocks)

	newClocks := make([]vclock.VClock, 0, len(s.vCache[i.Keygroup].clocks[i.ID].clocks))

	for _, local := range s.vCache[i.Keygroup].clocks[i.ID].clocks {
		// if we already have a newer version, there is nothing to do for us
		if local.Compare(remoteVersion, vclock.Ancestor) {
			log.Debug().Msgf("%s is an ancestor of %s: discarding version", vector.SortedVCString(remoteVersion), vector.SortedVCString(local))
			return nil
		}

		// if this is an equal version, we better hope that we already know about this and the contents aren't any different
		if local.Compare(remoteVersion, vclock.Equal) {
			log.Debug().Msgf("%s is equal to %s: discarding version", vector.SortedVCString(remoteVersion), vector.SortedVCString(local))

			// ok this should actually never happen
			// it can happen when a replica is added while an update is in progress
			// in that case the values should be the same so we can safely ignore it
			// but I'm not a fan of that tbh
			// TODO: figure out why updates can arrive twice while adding a replica

			return nil
		}

		// if this is a newer version than we have, we can actually just overwrite our local version and we're good
		if local.Compare(remoteVersion, vclock.Descendant) {
			log.Debug().Msgf("%s is a descendant of %s", vector.SortedVCString(remoteVersion), vector.SortedVCString(local))
			// to do that, we just need to set our local cache to this new version
			// and then store that version
			// we won't lose any data or anything: a merge would just lead to the same remoteVersion since it is larger
			log.Debug().Msgf("removing version %s", vector.SortedVCString(local))
			s.prune(string(i.Keygroup), i.ID, []vclock.VClock{local})
			continue
		}

		// not ancestor, descendant, equal => concurrent
		log.Debug().Msgf("%s is concurrent to %s", vector.SortedVCString(remoteVersion), vector.SortedVCString(local))
		newClocks = append(newClocks, local)
	}

	log.Debug().Msgf("storing version %s", vector.SortedVCString(remoteVersion))

	err = s.iS.Update(string(i.Keygroup), i.ID, i.Val, expiry, remoteVersion.GetMap())

	if err != nil {
		return err
	}

	// and then since we have seen a remote version, we will update our local cache with this new version we saw
	newClocks = append(newClocks, remoteVersion.Copy())

	s.vCache[i.Keygroup].clocks[i.ID].clocks = newClocks

	log.Debug().Msgf("addVersion: after known %+v", s.vCache[i.Keygroup].clocks[i.ID].clocks)

	return nil
}

// prune removes old versions of an Item
func (s *storeService) prune(kg string, id string, versions []vclock.VClock) {
	// although prune writes to the local database, it doesn't need to be synced
	// because it doesn't change anything on our version number cache
	// errors might happen if an item is removed twice but who cares

	for _, v := range versions {
		if v == nil || len(v) == 0 {
			return
		}

		log.Debug().Msgf("storeservice: pruning version %+v of %s in keygroup %s", v, id, kg)

		err := s.iS.Delete(kg, id, v)

		if err != nil {
			log.Err(err).Msgf("error pruning version %+v of %s in keygroup %s", v, id, kg)
		}
	}

	_, v, _, err := s.iS.Read(kg, id)

	if err != nil {
		log.Err(err).Msgf("error pruning version %+v of %s in keygroup %s", v, id, kg)
	}

	log.Debug().Msgf("pruning: versions %+v pruned, have versions %+v", versions, v)
}

func (s *storeService) tombstone(i Item) (vclock.VClock, error) {
	//TODO
	err := checkKGandID(i.Keygroup, i.ID)

	if err != nil {
		return nil, err
	}

	if !s.iS.ExistsKeygroup(string(i.Keygroup)) {
		return nil, errors.Errorf("no such keygroup in store: %+v", i.Keygroup)
	}

	i.Tombstoned = true
	i.Val = ""

	s.vCacheLock.RLock()
	defer s.vCacheLock.RUnlock()

	s.vCache[i.Keygroup].RLock()
	defer s.vCache[i.Keygroup].RUnlock()

	if _, ok := s.vCache[i.Keygroup].clocks[i.ID]; !ok {
		s.vCache[i.Keygroup].RUnlock()
		s.vCache[i.Keygroup].Lock()
		if _, ok := s.vCache[i.Keygroup].clocks[i.ID]; !ok {
			s.vCache[i.Keygroup].clocks[i.ID] = &cachedClock{
				clocks: []vclock.VClock{},
				Mutex:  &sync.Mutex{},
			}
		}
		s.vCache[i.Keygroup].Unlock()
		s.vCache[i.Keygroup].RLock()
	}

	s.vCache[i.Keygroup].clocks[i.ID].Lock()
	defer s.vCache[i.Keygroup].clocks[i.ID].Unlock()

	oldVersion := vclock.VClock{}
	toPrune := make([]vclock.VClock, len(s.vCache[i.Keygroup].clocks[i.ID].clocks))
	for j, c := range s.vCache[i.Keygroup].clocks[i.ID].clocks {
		oldVersion.Merge(c)
		toPrune[j] = c.Copy()
	}

	newVersion := oldVersion.Copy()
	newVersion.Tick(s.id)

	s.vCache[i.Keygroup].clocks[i.ID].clocks = []vclock.VClock{newVersion}

	err = s.iS.Update(string(i.Keygroup), i.ID, i.Val, 0, newVersion.GetMap())

	if err != nil {
		return nil, err
	}

	s.prune(string(i.Keygroup), i.ID, toPrune)

	return newVersion.Copy(), nil
}

func (s *storeService) tombstoneVersions(i Item, versions []vclock.VClock) (vclock.VClock, error) {
	// TODO
	// this works like updateVersions
	err := checkKGandID(i.Keygroup, i.ID)

	if err != nil {
		return nil, err
	}

	if !s.iS.ExistsKeygroup(string(i.Keygroup)) {
		return nil, errors.Errorf("no such keygroup in store: %+v", i.Keygroup)
	}

	i.Tombstoned = true
	i.Val = ""

	s.vCacheLock.RLock()
	defer s.vCacheLock.RUnlock()

	s.vCache[i.Keygroup].RLock()
	defer s.vCache[i.Keygroup].RUnlock()

	if _, ok := s.vCache[i.Keygroup].clocks[i.ID]; !ok {
		s.vCache[i.Keygroup].RUnlock()
		s.vCache[i.Keygroup].Lock()
		if _, ok := s.vCache[i.Keygroup].clocks[i.ID]; !ok {
			s.vCache[i.Keygroup].clocks[i.ID] = &cachedClock{
				clocks: []vclock.VClock{},
				Mutex:  &sync.Mutex{},
			}
		}
		s.vCache[i.Keygroup].Unlock()
		s.vCache[i.Keygroup].RLock()
	}

	s.vCache[i.Keygroup].clocks[i.ID].Lock()
	defer s.vCache[i.Keygroup].clocks[i.ID].Unlock()

	// ok, let's check if we want to reject this update!
	// look, if the client sends (B:0, C:1) AND (B:0, C:2), it's their fault
	// the client should make sure that it only sends concurrent versions
	for _, v1 := range versions {
		for _, v2 := range s.vCache[i.Keygroup].clocks[i.ID].clocks {
			// the client sent us an older version? reject!
			if v2.Compare(v1, vclock.Ancestor) {
				return nil, errors.Errorf("update rejected as given version %+v is older than seen version %+v", v1, v2)
			}
		}
	}

	// looks good! let's make up a new version for this update
	newVersion := vclock.VClock{}
	for _, v := range versions {
		newVersion.Merge(v)
	}

	newVersion.Tick(s.id)

	err = s.iS.Update(string(i.Keygroup), i.ID, i.Val, 0, newVersion.GetMap())

	if err != nil {
		return nil, err
	}

	// let us finally update our local cache by removing all the older versions
	newClocks := []vclock.VClock{newVersion}
	for _, v := range s.vCache[i.Keygroup].clocks[i.ID].clocks {
		if v.Compare(newVersion, vclock.Concurrent) {
			newClocks = append(newClocks, v.Copy())
		} else {
			s.prune(string(i.Keygroup), i.ID, []vclock.VClock{v})
		}
	}
	s.vCache[i.Keygroup].clocks[i.ID].clocks = newClocks

	return newVersion.Copy(), nil
}

// append appends an item in the key-value store.
func (s *storeService) append(i Item, expiry int) error {
	if !s.iS.ExistsKeygroup(string(i.Keygroup)) {
		return errors.Errorf("no such keygroup in store: %+v", i.Keygroup)
	}

	s.vCacheLock.RLock()
	defer s.vCacheLock.RUnlock()

	s.vCache[i.Keygroup].RLock()
	defer s.vCache[i.Keygroup].RUnlock()

	if _, ok := s.vCache[i.Keygroup].clocks[i.ID]; !ok {
		s.vCache[i.Keygroup].RUnlock()
		s.vCache[i.Keygroup].Lock()
		if _, ok := s.vCache[i.Keygroup].clocks[i.ID]; !ok {
			s.vCache[i.Keygroup].clocks[i.ID] = &cachedClock{
				clocks: []vclock.VClock{},
				Mutex:  &sync.Mutex{},
			}
		}
		s.vCache[i.Keygroup].Unlock()
		s.vCache[i.Keygroup].RLock()
	}

	s.vCache[i.Keygroup].clocks[i.ID].Lock()
	defer s.vCache[i.Keygroup].clocks[i.ID].Unlock()

	err := s.iS.Append(string(i.Keygroup), i.ID, i.Val, expiry)

	if err != nil {
		return err
	}

	return nil
}

// Should be in exthandler and inthandler
// CreateKeygroup creates a new keygroup in the store.
func (s *storeService) createKeygroup(kg KeygroupName) error {
	err := checkKeygroup(kg)

	if err != nil {
		return err
	}

	if s.iS.ExistsKeygroup(string(kg)) {
		return errors.Errorf("keygroup already in store: %+v", kg)
	}

	err = s.iS.CreateKeygroup(string(kg))

	if err != nil {
		return err
	}

	s.vCacheLock.Lock()
	defer s.vCacheLock.Unlock()
	//s.lockGroups[kg] = &singleflight.Group{}
	s.vCache[kg] = keygroupCache{make(map[string]*cachedClock), &sync.RWMutex{}}

	return nil
}

// DeleteKeygroup deletes a keygroup from the store.
func (s *storeService) deleteKeygroup(kg KeygroupName) error {
	err := checkKeygroup(kg)

	if err != nil {
		return err
	}

	if !s.iS.ExistsKeygroup(string(kg)) {
		return errors.Errorf("no such keygroup in store: %+v", kg)
	}

	err = s.iS.DeleteKeygroup(string(kg))

	if err != nil {
		return err
	}

	s.vCacheLock.Lock()
	defer s.vCacheLock.Unlock()
	//delete(s.lockGroups, kg)
	delete(s.vCache, kg)

	return nil
}

func (s *storeService) addKeygroupTrigger(kg KeygroupName, t Trigger) error {
	err := checkKeygroup(kg)

	if err != nil {
		return err
	}

	if !s.iS.ExistsKeygroup(string(kg)) {
		return errors.Errorf("no such keygroup in store: %+v", kg)
	}

	err = s.iS.AddKeygroupTrigger(string(kg), t.ID, t.Host)

	if err != nil {
		return err
	}

	return nil
}

func (s *storeService) deleteKeygroupTrigger(kg KeygroupName, t Trigger) error {
	err := checkKeygroup(kg)

	if err != nil {
		return err
	}

	if !s.iS.ExistsKeygroup(string(kg)) {
		return errors.Errorf("no such keygroup in store: %+v", kg)
	}

	err = s.iS.DeleteKeygroupTrigger(string(kg), t.ID)

	if err != nil {
		return err
	}

	return nil
}

func (s *storeService) getKeygroupTrigger(kg KeygroupName) ([]Trigger, error) {
	err := checkKeygroup(kg)

	if err != nil {
		return nil, err
	}

	if !s.iS.ExistsKeygroup(string(kg)) {
		return nil, errors.Errorf("no such keygroup in store: %+v", kg)
	}

	t, err := s.iS.GetKeygroupTrigger(string(kg))

	if err != nil {
		return nil, err
	}

	tn := make([]Trigger, 0, len(t))

	for id, host := range t {
		tn = append(tn, Trigger{
			ID:   id,
			Host: host,
		})
	}

	log.Debug().Msgf("getKeygroupTrigger: %d items, %+v %+v", len(t), t, tn)

	return tn, nil
}
