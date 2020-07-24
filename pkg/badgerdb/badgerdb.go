package badgerdb

import (
	"time"

	badger "github.com/dgraph-io/badger/v2"
	"github.com/go-errors/errors"
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

// garbageCollection manages triggering garbage collection for the BadgerDB database. For now, we stick to a schedule.
func garbageCollection(db *badger.DB) {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()
	for range ticker.C {
		for err := db.RunValueLogGC(0.7); err == nil; err = db.RunValueLogGC(0.7) {
		}
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
		item, err := txn.Get(makeKeyName(kg, id))

		if err != nil {
			return errors.New(err)
		}

		val, err := item.ValueCopy(nil)

		if err != nil {
			return errors.New(err)
		}

		value = string(val)

		return err
	})

	return value, err

}

// ReadAll returns all items in the specified keygroup.
func (s *Storage) ReadAll(kg string) (map[string]string, error) {
	items := make(map[string]string)

	err := s.db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()
		prefix := makeKeygroupKeyName(kg)

		// this will remove the entry for our keygroup anchor point which is just "kgname-" without a value in the database
		// iterator returns keys in lexicographical order, so just remove the first one
		it.Seek(prefix)
		it.Next()

		for ; it.ValidForPrefix(prefix); it.Next() {
			item := it.Item()
			v, err := item.ValueCopy(nil)

			if err != nil {
				return errors.New(err)
			}

			items[string(item.Key())] = string(v)
		}
		return nil
	})

	delete(items, string(makeKeygroupKeyName(kg)))

	return items, err
}

// IDs returns the keys of all items in the specified keygroup.
func (s *Storage) IDs(kg string) ([]string, error) {
	var items []string

	err := s.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = false
		it := txn.NewIterator(opts)
		defer it.Close()
		prefix := makeKeygroupKeyName(kg)

		// this will remove the entry for our keygroup anchor point which is just "kgname-" without a value in the database
		// iterator returns keys in lexicographical order, so just remove the first one
		it.Seek(prefix)
		it.Next()

		for ; it.ValidForPrefix(prefix); it.Next() {
			items = append(items, string(it.Item().Key()))
		}
		return nil
	})

	return items, err
}

// Update updates the item with the specified id in the specified keygroup.
func (s *Storage) Update(kg, id, val string) error {
	return s.db.Update(func(txn *badger.Txn) error {
		err := txn.Set(makeKeyName(kg, id), []byte(val))
		if err != nil {
			return errors.New(err)
		}
		return nil
	})
}

// Delete deletes the item with the specified id from the specified keygroup.
func (s *Storage) Delete(kg string, id string) error {
	return s.db.Update(func(txn *badger.Txn) error {
		err := txn.Delete(makeKeyName(kg, id))
		if err != nil {
			return errors.New(err)
		}
		return nil
	})
}

// Exists checks if the given data item exists in the leveldb database.
func (s *Storage) Exists(kg string, id string) bool {
	err := s.db.View(func(txn *badger.Txn) error {
		_, err := txn.Get(makeKeyName(kg, id))

		return err
	})

	return err == nil
}

// ExistsKeygroup checks if the given keygroup exists in the leveldb database.
func (s *Storage) ExistsKeygroup(kg string) bool {
	err := s.db.View(func(txn *badger.Txn) error {
		_, err := txn.Get(makeKeygroupKeyName(kg))

		return err
	})

	return err == nil
}

// CreateKeygroup creates the given keygroup in the leveldb database.
func (s *Storage) CreateKeygroup(kg string) error {
	return s.db.Update(func(txn *badger.Txn) error {
		err := txn.Set(makeKeygroupKeyName(kg), []byte(nil))
		if err != nil {
			return errors.New(err)
		}
		return nil
	})
}

// DeleteKeygroup deletes the given keygroup from the leveldb database.
func (s *Storage) DeleteKeygroup(kg string) error {

	var ids [][]byte

	err := s.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = false
		it := txn.NewIterator(opts)
		defer it.Close()
		prefix := makeKeygroupKeyName(kg)

		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			ids = append(ids, it.Item().Key())
		}
		return nil
	})

	if err != nil {
		return errors.New(err)
	}

	wb := s.db.NewWriteBatch()
	defer wb.Cancel()

	for _, id := range ids {
		err := wb.Delete(id) // Will create txns as needed.

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
