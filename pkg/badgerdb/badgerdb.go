package badgerdb

import (
	"strings"
	"time"

	"git.tu-berlin.de/mcc-fred/fred/pkg/vector"
	"github.com/DistributedClocks/GoVector/govec/vclock"
	"github.com/dgraph-io/badger/v3"
	"github.com/go-errors/errors"
	"github.com/rs/zerolog/log"
)

const sep = "|"
const gcInterval = 5 * time.Minute
const gcDiscardRatio = 0.7

// Storage is a struct that saves all necessary information to access the database, in this case just a pointer to the BadgerDB database.
type Storage struct {
	db  *badger.DB
	seq map[string]*badger.Sequence
}

// makeKeyNamePrefix creates the internal BadgerDB key given a keygroup name and an id but without a version vector
func makeKeyNamePrefix(kgname string, id string) []byte {
	return []byte(kgname + sep + id + sep)
}

// makeKeyName creates the internal BadgerDB key given a keygroup name, an id, and a version vector.
func makeKeyName(kgname string, id string, vvector vclock.VClock) []byte {
	return append([]byte(kgname+sep+id+sep), vector.Bytes(vvector)...)
}

// makeKeygroupKeyName creates the internal BadgerDB key given a keygroup name.
func makeKeygroupKeyName(kgname string) []byte {
	return []byte(kgname + sep)
}

func makeKeygroupConfigKeyName(kgname string) []byte {
	return []byte(sep + "fred" + sep + "keygroup" + sep + kgname)
}

func makeTriggerConfigKeyName(kgname string, tid string) []byte {
	return []byte(sep + "fred" + sep + "triggers" + sep + kgname + sep + tid)
}

func makeLogConfigKeyName(kgname string) []byte {
	return []byte(sep + "fred" + sep + "rolling" + sep + kgname)
}

// getTriggerConfigKey returns the keygroup and id of a key.
func getTriggerConfigKey(key string) (kg, tid string) {
	s := strings.Split(key, sep)

	if len(s) == len([]string{"nil", "fred", "triggers", "keygroup", "trigger id"}) {
		kg = s[3]
		tid = s[4]
	}

	return
}

// getKey returns the keygroup and id of a key.
func getKey(key string) (kg, id string, vvector vclock.VClock) {
	s := strings.SplitN(key, sep, 3)
	kg = s[0]

	if len(s) == len([]string{"keygroup", "identifier", "version vector"}) {
		id = s[1]
		vvector, _ = vector.FromBytes([]byte(s[2]))
	}

	return
}

// garbageCollection manages triggering garbage collection for the BadgerDB database. For now, we stick to a schedule.
func garbageCollection(db *badger.DB) {
	ticker := time.NewTicker(gcInterval)
	defer ticker.Stop()
	for range ticker.C {
		log.Debug().Msgf("BadgerDB: triggering garbage collection...")
		for err := db.RunValueLogGC(gcDiscardRatio); err == nil; err = db.RunValueLogGC(gcDiscardRatio) {
			log.Debug().Msgf("BadgerDB: garbage collected!")
		}
		log.Debug().Msgf("BadgerDB: garbage collection done")
	}
}

// New creates a new BadgerDB Storage on disk.
func New(dbPath string) (s *Storage) {
	db, err := badger.Open(badger.DefaultOptions(dbPath).WithLoggingLevel(badger.ERROR))
	if err != nil {
		panic(err)
	}

	s = &Storage{
		db:  db,
		seq: make(map[string]*badger.Sequence),
	}

	go garbageCollection(db)

	return
}

// NewMemory create a new BadgerDB Storage in memory.
func NewMemory() (s *Storage) {
	db, err := badger.Open(badger.DefaultOptions("").WithInMemory(true).WithLoggingLevel(badger.ERROR))
	if err != nil {
		panic(err)
	}

	s = &Storage{
		db:  db,
		seq: make(map[string]*badger.Sequence),
	}

	go garbageCollection(db)

	return
}

// Close closes the underlying BadgerDB.
func (s *Storage) Close() error {
	return s.db.Close()
}

// Read returns an item with the specified id from the specified keygroup.
func (s *Storage) Read(kg string, id string) ([]string, []vclock.VClock, bool, error) {
	values := make([]string, 0)
	vvectors := make([]vclock.VClock, 0)

	err := s.db.View(func(txn *badger.Txn) error {
		prefix := makeKeyNamePrefix(kg, id)

		opts := badger.DefaultIteratorOptions
		opts.Prefix = prefix

		it := txn.NewIterator(opts)
		defer it.Close()

		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			item := it.Item()
			_, _, vvector := getKey(string(item.Key()))

			v, err := item.ValueCopy(nil)

			if err != nil {
				return err
			}

			values = append(values, string(v))
			vvectors = append(vvectors, vvector)
		}

		return nil
	})

	if err != nil {
		// if the error is a "KeyNotFound", there is no need to add a stacktrace
		if errors.Is(err, badger.ErrKeyNotFound) {
			return nil, nil, false, nil
		}
		// if we have a different error, debug with full stacktrace
		return nil, nil, false, errors.New(err)
	}

	if len(values) == 0 {
		return nil, nil, false, nil
	}

	return values, vvectors, true, nil

}

// ReadSome returns count number of items in the specified keygroup starting at id.
func (s *Storage) ReadSome(kg, id string, count uint64) (map[string][]string, map[string][]vclock.VClock, error) {
	items := make(map[string][]string)
	vvectors := make(map[string][]vclock.VClock)

	err := s.db.View(func(txn *badger.Txn) error {
		prefix := makeKeygroupKeyName(kg)
		start := makeKeyNamePrefix(kg, id)

		opts := badger.DefaultIteratorOptions
		opts.Prefix = prefix

		it := txn.NewIterator(opts)
		defer it.Close()

		var i uint64
		for it.Seek(start); it.ValidForPrefix(prefix) && i < count; it.Next() {
			item := it.Item()
			_, key, vvector := getKey(string(item.Key()))

			v, err := item.ValueCopy(nil)

			if err != nil {
				return err
			}

			if _, ok := items[key]; ok {
				items[key] = append(items[key], string(v))
				vvectors[key] = append(vvectors[key], vvector)
				continue
			}

			items[key] = []string{string(v)}
			vvectors[key] = []vclock.VClock{vvector}
			i++
		}

		return nil
	})

	if err != nil {
		return nil, nil, errors.New(err)
	}

	return items, vvectors, nil
}

// ReadAll returns all items in the specified keygroup.
func (s *Storage) ReadAll(kg string) (map[string][]string, map[string][]vclock.VClock, error) {
	items := make(map[string][]string)
	vvectors := make(map[string][]vclock.VClock)

	err := s.db.View(func(txn *badger.Txn) error {
		prefix := makeKeygroupKeyName(kg)

		opts := badger.DefaultIteratorOptions
		opts.Prefix = prefix

		it := txn.NewIterator(opts)
		defer it.Close()

		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			item := it.Item()
			_, key, vvector := getKey(string(item.Key()))

			v, err := item.ValueCopy(nil)

			if err != nil {
				return err
			}

			if _, ok := items[key]; ok {
				items[key] = append(items[key], string(v))
				vvectors[key] = append(vvectors[key], vvector)
				continue
			}

			items[key] = []string{string(v)}
			vvectors[key] = []vclock.VClock{vvector}
		}

		return nil
	})

	if err != nil {
		return nil, nil, errors.New(err)
	}

	return items, vvectors, nil
}

// IDs returns the keys of all items in the specified keygroup.
func (s *Storage) IDs(kg string) ([]string, error) {
	var items []string

	err := s.db.View(func(txn *badger.Txn) error {
		prefix := makeKeygroupKeyName(kg)

		opts := badger.DefaultIteratorOptions
		opts.Prefix = prefix
		opts.PrefetchValues = false

		it := txn.NewIterator(opts)
		defer it.Close()

		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			item := it.Item()
			_, key, _ := getKey(string(item.Key()))

			if len(items) > 0 && key == items[len(items)-1] {
				continue
			}

			items = append(items, key)
		}
		return nil
	})

	if err != nil {
		return nil, errors.New(err)
	}

	return items, nil
}

// Update updates the item with the specified id in the specified keygroup.
func (s *Storage) Update(kg, id, val string, expiry int, vvector vclock.VClock) error {
	err := s.db.Update(func(txn *badger.Txn) error {
		key := makeKeyName(kg, id, vvector)

		if expiry > 0 {
			err := txn.SetEntry(&badger.Entry{
				Key:       key,
				Value:     []byte(val),
				ExpiresAt: uint64(time.Now().Unix()) + uint64(expiry),
			})
			if err != nil {
				return err
			}
			return nil
		}

		err := txn.Set(key, []byte(val))

		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return errors.New(err)
	}

	return nil
}

// Delete deletes the item with the specified id from the specified keygroup. Will delete a specific version or all versions if vvector is nil.
func (s *Storage) Delete(kg string, id string, vvector vclock.VClock) error {
	if vvector != nil {
		err := s.db.Update(func(txn *badger.Txn) error {
			err := txn.Delete(makeKeyName(kg, id, vvector))
			if err != nil {
				return err
			}
			return nil

		})

		if err != nil {
			// if the error is a "KeyNotFound", there is no need to add a stacktrace
			if errors.Is(err, badger.ErrKeyNotFound) {
				return errors.Errorf("key not found in database: %s in keygroup %s", id, kg)
			}
			// if we have a different error, debug with full stacktrace
			return errors.New(err)
		}

		return nil
	}

	var ids []string

	err := s.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = false
		it := txn.NewIterator(opts)
		defer it.Close()

		prefix := makeKeyNamePrefix(kg, id)

		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			ids = append(ids, string(it.Item().Key()))
		}
		return nil
	})

	if err != nil {
		return errors.New(err)
	}

	// then, delete all the keys
	wb := s.db.NewWriteBatch()
	defer wb.Cancel()

	for _, d := range ids {
		err := wb.Delete([]byte(d)) // Will create txns as needed.

		if err != nil {
			return errors.New(err)
		}
	}
	err = wb.Flush() // Wait for all txns to finish.

	if err != nil {
		return errors.New(err)
	}

	return nil
}

// Append appends the item to the specified keygroup by incrementing the latest key by one.
func (s *Storage) Append(kg string, id string, val string, expiry int) error {
	if s.Exists(kg, id) {
		return errors.Errorf("key %s for keygroup %s exists in database already and may not be changed", id, kg)
	}

	err := s.db.Update(func(txn *badger.Txn) error {
		key := makeKeyName(kg, id, vclock.VClock{})

		if expiry > 0 {
			err := txn.SetEntry(&badger.Entry{
				Key:       key,
				Value:     []byte(val),
				ExpiresAt: uint64(time.Now().Unix()) + uint64(expiry),
			})
			if err != nil {
				return errors.New(err)
			}
			return nil
		}

		err := txn.Set(key, []byte(val))

		if err != nil {
			return errors.New(err)
		}
		return nil
	})

	if err != nil {
		return errors.New(err)
	}

	return nil
}

// Exists checks if the given data item exists in the badgerdb database.
func (s *Storage) Exists(kg string, id string) bool {
	err := s.db.View(func(txn *badger.Txn) error {
		prefix := makeKeyNamePrefix(kg, id)

		opts := badger.DefaultIteratorOptions
		opts.Prefix = prefix

		it := txn.NewIterator(opts)
		defer it.Close()
		it.Seek(prefix)
		if it.ValidForPrefix(prefix) {
			return nil
		}

		return badger.ErrKeyNotFound
	})

	if err == nil {
		return true
	}
	// if the error is not a "KeyNotFound", there is something else going on
	if !errors.Is(err, badger.ErrKeyNotFound) {
		log.Err(err).Msg("")
	}

	return false
}

// ExistsKeygroup checks if the given keygroup exists in the badgerdb database.
func (s *Storage) ExistsKeygroup(kg string) bool {
	err := s.db.View(func(txn *badger.Txn) error {
		_, err := txn.Get(makeKeygroupConfigKeyName(kg))

		return err
	})

	if err == nil {
		return true
	}
	// if the error is not a "KeyNotFound", there is something else going on
	if !errors.Is(err, badger.ErrKeyNotFound) {
		log.Err(err).Msg("")
	}

	return false
}

// CreateKeygroup creates the given keygroup in the badgerdb database.
func (s *Storage) CreateKeygroup(kg string) error {
	err := s.db.Update(func(txn *badger.Txn) error {
		err := txn.Set(makeKeygroupConfigKeyName(kg), []byte(kg))
		if err != nil {
			return err
		}

		s.seq[kg], err = s.db.GetSequence(makeLogConfigKeyName(kg), 100)

		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return errors.New(err)
	}

	return nil
}

// DeleteKeygroup deletes the given keygroup from the badgerdb database.
func (s *Storage) DeleteKeygroup(kg string) error {
	err := s.db.Update(func(txn *badger.Txn) error {
		err := txn.Delete(makeKeygroupConfigKeyName(kg))
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		// if the error is a "KeyNotFound", there is no need to add a stacktrace
		if errors.Is(err, badger.ErrKeyNotFound) {
			return errors.Errorf("keygroup not found in database: %s", kg)
		}
		// if we have a different error, debug with full stacktrace
		return errors.New(err)
	}

	// first, get all the ids in this keygroup (including the marker that the keygroup exists)
	// why is this a string? well, it was a [][]byte before, but for some reason the underlying byte array would be
	// overwritten with older keys in the unit tests, that was really strange
	// so now it's an array of strings, which are immutable
	var ids []string

	err = s.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = false
		it := txn.NewIterator(opts)
		defer it.Close()
		prefix := makeKeygroupKeyName(kg)

		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			ids = append(ids, string(it.Item().Key()))
		}
		return nil
	})

	if err != nil {
		return errors.New(err)
	}

	// then, delete all the keys
	wb := s.db.NewWriteBatch()
	defer wb.Cancel()

	for _, id := range ids {
		err := wb.Delete([]byte(id)) // Will create txns as needed.

		if err != nil {
			return errors.New(err)
		}
	}
	err = wb.Flush() // Wait for all txns to finish.

	if err != nil {
		return errors.New(err)
	}

	// we also need to remove all keygroup triggers
	var triggers []string

	err = s.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = false
		it := txn.NewIterator(opts)
		defer it.Close()
		prefix := makeTriggerConfigKeyName(kg, "")

		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			triggers = append(triggers, string(it.Item().Key()))
		}
		return nil
	})

	if err != nil {
		return errors.New(err)
	}

	// then, delete all the keys for the trigger nodes
	wb = s.db.NewWriteBatch()
	defer wb.Cancel()

	for _, t := range triggers {
		err := wb.Delete([]byte(t)) // Will create txns as needed.

		if err != nil {
			return errors.New(err)
		}
	}
	err = wb.Flush() // Wait for all txns to finish.

	if err != nil {
		return errors.New(err)
	}

	if seq, ok := s.seq[kg]; ok {
		err = seq.Release()

		if err != nil {
			return err
		}

		delete(s.seq, kg)
	}

	return nil
}

// AddKeygroupTrigger adds a trigger node to the given keygroup in the badgerdb database.
func (s *Storage) AddKeygroupTrigger(kg string, id string, host string) error {
	err := s.db.Update(func(txn *badger.Txn) error {
		key := makeTriggerConfigKeyName(kg, id)

		err := txn.Set(key, []byte(host))
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return errors.New(err)
	}

	return nil
}

// DeleteKeygroupTrigger removes a trigger node from the given keygroup in the badgerdb database.
func (s *Storage) DeleteKeygroupTrigger(kg string, id string) error {
	err := s.db.Update(func(txn *badger.Txn) error {
		err := txn.Delete(makeTriggerConfigKeyName(kg, id))
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return errors.New(err)
	}

	return nil
}

// GetKeygroupTrigger returns a list of all trigger nodes for the given keygroup in the badgerdb database.
func (s *Storage) GetKeygroupTrigger(kg string) (map[string]string, error) {
	items := make(map[string]string)

	err := s.db.View(func(txn *badger.Txn) error {
		prefix := makeTriggerConfigKeyName(kg, "")

		opts := badger.DefaultIteratorOptions
		opts.Prefix = prefix

		it := txn.NewIterator(opts)
		defer it.Close()

		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			item := it.Item()
			_, tid := getTriggerConfigKey(string(item.Key()))

			if tid == "" {
				continue
			}

			v, err := item.ValueCopy(nil)

			if err != nil {
				return err
			}

			items[tid] = string(v)
		}

		return nil
	})

	if err != nil {
		return nil, errors.New(err)
	}

	return items, nil
}
