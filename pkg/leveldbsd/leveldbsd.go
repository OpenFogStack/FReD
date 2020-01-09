package leveldbsd

import (
	"github.com/rs/zerolog/log"
	"github.com/syndtr/goleveldb/leveldb"

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
func (s *Storage) Read(i data.Item) (data.Item, error) {
	key := makeKeyName(i.Keygroup, i.ID)

	value, err := s.db.Get([]byte(key), nil)

	i.Data = string(value)

	log.Debug().Err(err).Msgf("Read from levedbsd: in %v, out %s", i, string(value))

	return i, err
}

// Update updates the item with the specified id in the specified keygroup.
func (s *Storage) Update(i data.Item) error {
	key := makeKeyName(i.Keygroup, i.ID)

	err := s.db.Put([]byte(key), []byte(i.Data), nil)

	log.Debug().Err(err).Msgf("Update from levedbsd: in %v", i)

	return err
}

// Delete deletes the item with the specified id from the specified keygroup.
func (s *Storage) Delete(i data.Item) error {
	key := makeKeyName(i.Keygroup, i.ID)

	err := s.db.Delete([]byte(key), nil)

	log.Debug().Err(err).Msgf("Delete from levedbsd: in %v", i)

	return err
}

// Exists checks if the given data item exists in the leveldb database.
func (s *Storage) Exists(i data.Item) bool {
	key := makeKeyName(i.Keygroup, i.ID)

	has, _ := s.db.Has([]byte(key), nil)

	log.Debug().Msgf("Exists from levedbsd: in %v, out: %t", i, has)

	return has
}

// ExistsKeygroup checks if the given keygroup exists in the leveldb database.
func (s *Storage) ExistsKeygroup(i data.Item) bool {
	key := makeKeyName(i.Keygroup, i.ID)

	has, _ := s.db.Has([]byte(key), nil)

	log.Debug().Msgf("ExistsKeygroup from levedbsd: in %v, out: %t", i, has)

	return has
}

// CreateKeygroup creates the given keygroup in the leveldb database.
func (s *Storage) CreateKeygroup(i data.Item) error {
	key := makeKeyName(i.Keygroup, i.ID)

	err := s.db.Put([]byte(key), []byte(i.Data), nil)

	log.Debug().Err(err).Msgf("CreateKeygroup from levedbsd: in %v", i)

	return err
}

// DeleteKeygroup deletes the given keygroup from the leveldb database.
func (s *Storage) DeleteKeygroup(i data.Item) error {
	key := makeKeyName(i.Keygroup, i.ID)

	err := s.db.Delete([]byte(key), nil)

	log.Debug().Err(err).Msgf("DeleteKeygroup from levedbsd: in %v", i)

	return err
}
