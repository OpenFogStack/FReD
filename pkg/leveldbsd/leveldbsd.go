package leveldbsd

import (
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"

	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/commons"
	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/data"
)

// Storage is a struct that saves all necessary information to access the database, in this case just a pointer to the LevelDB database.
type Storage struct {
	db *leveldb.DB
}

// makeKeyName creates the internal LevelDB key given a keygroup name and an id
func makeKeyName(kgname commons.KeygroupName, id string) string {
	return string(kgname) + "/" + id
}

// makeKeygroupKeyName creates the internal LevelDB key given a keygroup name
func makeKeygroupKeyName(kgname commons.KeygroupName) string {
	return string(kgname) + "/"
}

// New create a new Storage.
func New(dbPath string) (s *Storage) {
	db, err := leveldb.OpenFile(dbPath, nil)

	if err != nil {
		panic(err)
	}

	s = &Storage{
		db: db,
	}

	return
}

// Read returns an item with the specified id from the specified keygroup.
func (s *Storage) Read(kg commons.KeygroupName, id string) (string, error) {
	key := makeKeyName(kg, id)

	value, err := s.db.Get([]byte(key), nil)

	log.Debug().Err(err).Msgf("Read from levedbsd: in %#v, %#v, out %s", kg, id, string(value))

	return string(value), err
}

// ReadAll returns all items in the specified keygroup.
func (s *Storage) ReadAll(kg commons.KeygroupName) ([]data.Item, error) {
	key := makeKeygroupKeyName(kg)

	iter := s.db.NewIterator(util.BytesPrefix([]byte(key)), nil)
	defer iter.Release()

	var items []data.Item
	l := 0

	for iter.Next() {
		l++
	}
	log.Debug().Msgf("ReadAll from levekdbsd: Fount %d items with this prefix", l)

	next := iter.First()
	for next {
		if string(iter.Key()) == key {
			next = iter.Next()
			continue
		}

		items = append(items, data.Item{
			Keygroup: kg,
			// The Key is kg-name/key, so only take the key part
			ID:       strings.Split(string(iter.Key()), "/")[1],
			Data:     string(iter.Value()),
		})

		next = iter.Next()
	}

	err := iter.Error()

	log.Debug().Err(err).Msgf("ReadAll from levedbsd: in %#v, out %#v", kg, items)

	return items, err
}

// IDs returns the keys of all items in the specified keygroup.
func (s *Storage) IDs(kg commons.KeygroupName) ([]data.Item, error) {
	key := makeKeygroupKeyName(kg)

	iter := s.db.NewIterator(util.BytesPrefix([]byte(key)), nil)

	defer iter.Release()

	l := 0

	for iter.Next() {
		l++
	}

	keys := make([]data.Item, l)

	next := iter.First()

	for next {
		if string(iter.Key()) == key {
			next = iter.Next()
			continue
		}

		keys[len(keys)-l] = data.Item{
			Keygroup: kg,
			ID:       string(iter.Key()),
		}

		next = iter.Next()

		l--
	}

	err := iter.Error()

	log.Debug().Err(err).Msgf("IDs from levedbsd: in %#v, out %#v", kg, keys)

	return keys, err
}

// Update updates the item with the specified id in the specified keygroup.
func (s *Storage) Update(i data.Item) error {
	key := makeKeyName(i.Keygroup, i.ID)

	err := s.db.Put([]byte(key), []byte(i.Data), nil)

	log.Debug().Err(err).Msgf("Update from levedbsd: in %#v", i)

	return err
}

// Delete deletes the item with the specified id from the specified keygroup.
func (s *Storage) Delete(kg commons.KeygroupName, id string) error {
	key := makeKeyName(kg, id)

	err := s.db.Delete([]byte(key), nil)

	log.Debug().Err(err).Msgf("Delete from levedbsd: in %#v, %#v", kg, id)

	return err
}

// Exists checks if the given data item exists in the leveldb database.
func (s *Storage) Exists(kg commons.KeygroupName, id string) bool {
	key := makeKeyName(kg, id)

	has, _ := s.db.Has([]byte(key), nil)

	log.Debug().Msgf("Exists from levedbsd: in %#v, %#v, out: %t", kg, id, has)

	return has
}

// ExistsKeygroup checks if the given keygroup exists in the leveldb database.
func (s *Storage) ExistsKeygroup(i data.Item) bool {
	key := makeKeyName(i.Keygroup, i.ID)

	has, _ := s.db.Has([]byte(key), nil)

	log.Debug().Msgf("ExistsKeygroup from levedbsd: in %#v, out: %t", i, has)

	return has
}

// CreateKeygroup creates the given keygroup in the leveldb database.
func (s *Storage) CreateKeygroup(kg commons.KeygroupName) error {
	// TODO Tobias this used to also use the data of the item passed, was this necessary?
	key := makeKeygroupKeyName(kg)

	err := s.db.Put([]byte(key), nil, nil)

	log.Debug().Err(err).Msgf("CreateKeygroup from levedbsd: in %#v", kg)

	return err
}

// DeleteKeygroup deletes the given keygroup from the leveldb database.
func (s *Storage) DeleteKeygroup(kg commons.KeygroupName) (err error) {
	key := makeKeygroupKeyName(kg)
	log.Debug().Msgf("DeleteKeygroup called for KG %s", kg)

	iter := s.db.NewIterator(util.BytesPrefix([]byte(key)), nil)
	defer iter.Release()

	for iter.Next() {
		log.Debug().Msgf("DeleteKeygroup deleting item %s", iter.Key())
		err = s.db.Delete(iter.Key(), nil)
	}

	log.Debug().Err(err).Msgf("DeleteKeygroup from levedbsd: in %#v", kg)

	return err
}
