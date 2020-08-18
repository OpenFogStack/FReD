package badgerdb

import (
	"strings"
	"time"

	"github.com/dgraph-io/badger/v2"
	"github.com/go-errors/errors"
	"github.com/rs/zerolog/log"
)

// Storage is a struct that saves all necessary information to access the database, in this case just a pointer to the BadgerDB database.
type Storage struct {
	db *badger.DB
}

// makeKeyName creates the internal BadgerDB key given a keygroup name and an id.
func makeKeyName(kgname string, id string) []byte {
	return []byte(kgname + "/" + id)
}

// makeKeygroupKeyName creates the internal BadgerDB key given a keygroup name.
func makeKeygroupKeyName(kgname string) []byte {
	return []byte(kgname + "/")
}

// getKey returns the keygroup and id of a key.
func getKey(key string) (kg, id string) {
	s := strings.Split(key, "/")
	kg = s[0]
	if len(s) > 1 {
		id = s[1]
	}
	return
}

// garbageCollection manages triggering garbage collection for the BadgerDB database. For now, we stick to a schedule.
func garbageCollection(db *badger.DB) {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()
	for range ticker.C {
		log.Debug().Msgf("BadgerDB: triggering garbage collection...")
		for err := db.RunValueLogGC(0.7); err == nil; err = db.RunValueLogGC(0.7) {
			log.Debug().Msgf("BadgerDB: garbage collected!")
		}
		log.Debug().Msgf("BadgerDB: garbage collection done")
	}
}

// New creates a new BadgerDB Storage on disk.
func New(dbPath string) (s *Storage) {
	db, err := badger.Open(badger.DefaultOptions(dbPath))
	if err != nil {
		panic(err)
	}

	s = &Storage{
		db: db,
	}

	go garbageCollection(db)

	return
}

// NewMemory create a new BadgerDB Storage in memory.
func NewMemory() (s *Storage) {
	db, err := badger.Open(badger.DefaultOptions("").WithInMemory(true))
	if err != nil {
		panic(err)
	}

	s = &Storage{
		db: db,
	}

	go garbageCollection(db)

	return
}

// Close closes the underlying BadgerDB.
func (s *Storage) Close() error {
	return s.db.Close()
}

// Read returns an item with the specified id from the specified keygroup.
func (s *Storage) Read(kg string, id string) (string, error) {
	var value string

	err := s.db.View(func(txn *badger.Txn) error {
		key := makeKeyName(kg, id)

		log.Debug().Msgf("our key is %s (%#v)", string(key), key)

		item, err := txn.Get(key)

		if err != nil {
			return err
		}

		val, err := item.ValueCopy(nil)

		if err != nil {
			return err
		}

		value = string(val)

		return nil
	})

	if err != nil {
		// if the error is a "KeyNotFound", there is no need to add a stacktrace
		if errors.Is(err, badger.ErrKeyNotFound) {
			return "", errors.Errorf("key not found in database: %s in keygroup %s", id, kg)
		}
		// if we have a different error, debug with full stacktrace
		return "", errors.New(err)
	}

	return value, nil

}

// ReadAll returns all items in the specified keygroup.
func (s *Storage) ReadAll(kg string) (map[string]string, error) {
	items := make(map[string]string)

	err := s.db.View(func(txn *badger.Txn) error {
		prefix := makeKeygroupKeyName(kg)

		opts := badger.DefaultIteratorOptions
		opts.Prefix = prefix

		it := txn.NewIterator(opts)
		defer it.Close()

		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			item := it.Item()
			_, key := getKey(string(item.Key()))

			if key == "" {
				continue
			}

			v, err := item.ValueCopy(nil)

			if err != nil {
				return err
			}

			items[key] = string(v)
		}

		return nil
	})

	if err != nil {
		return nil, errors.New(err)
	}

	return items, nil
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
			_, key := getKey(string(item.Key()))

			if key == "" {
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
func (s *Storage) Update(kg, id, val string) error {
	err := s.db.Update(func(txn *badger.Txn) error {
		key := makeKeyName(kg, id)

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

// Delete deletes the item with the specified id from the specified keygroup.
func (s *Storage) Delete(kg string, id string) error {
	err := s.db.Update(func(txn *badger.Txn) error {
		err := txn.Delete(makeKeyName(kg, id))
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

// Exists checks if the given data item exists in the leveldb database.
func (s *Storage) Exists(kg string, id string) bool {
	err := s.db.View(func(txn *badger.Txn) error {
		_, err := txn.Get(makeKeyName(kg, id))

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

// ExistsKeygroup checks if the given keygroup exists in the leveldb database.
func (s *Storage) ExistsKeygroup(kg string) bool {
	err := s.db.View(func(txn *badger.Txn) error {
		_, err := txn.Get(makeKeygroupKeyName(kg))

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

// CreateKeygroup creates the given keygroup in the leveldb database.
func (s *Storage) CreateKeygroup(kg string) error {
	err := s.db.Update(func(txn *badger.Txn) error {
		err := txn.Set(makeKeygroupKeyName(kg), []byte(kg))
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

// DeleteKeygroup deletes the given keygroup from the leveldb database.
func (s *Storage) DeleteKeygroup(kg string) error {

	// first, get all the ids in this keygroup (including the marker that the keygroup exists)
	// why is this a string? well, it was a [][]byte before, but for some reason the underlying byte array would be
	// overwritten with older keys in the unit tests, that was really strange
	// so not it's an array of strings, which are immutable
	var ids []string

	err := s.db.View(func(txn *badger.Txn) error {
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

	return nil
}
