package leveldbsd

import (
	"github.com/syndtr/goleveldb/leveldb"

	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/data"
)

// Storage is a struct that saves all necessary information to access the database, in this case just a pointer to the LevelDB database.
type Storage struct {
	db *leveldb.DB
}

// makeKeyName creates the internal LevelDB key given a keygroup name and an id
func makeKeyName(kgname string, id string) string {
	return kgname + "/" + id
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

	data, err := s.db.Get([]byte(key), nil)

	i.Data = string(data)

	return i, err
}

// Update updates the item with the specified id in the specified keygroup.
func (s *Storage) Update(i data.Item) error {
	key := makeKeyName(i.Keygroup, i.ID)

	err := s.db.Put([]byte(key), []byte(i.Data), nil)

	return err
}

// Delete deletes the item with the specified id from the specified keygroup.
func (s *Storage) Delete(i data.Item) error {
	key := makeKeyName(i.Keygroup, i.ID)

	err := s.db.Delete([]byte(key), nil)

	return err
}

// Exists checks if the given data item exists in the leveldb database.
func (s *Storage) Exists(i data.Item) bool {
	key := makeKeyName(i.Keygroup, i.ID)

	has, _ := s.db.Has([]byte(key), nil)

	return has
}

// ExistsKeygroup checks if the given keygroup exists in the leveldb database.
func (s *Storage) ExistsKeygroup(i data.Item) bool {
	key := makeKeyName(i.Keygroup, i.ID)

	has, _ := s.db.Has([]byte(key), nil)

	return has
}

// CreateKeygroup creates the given keygroup in the leveldb database.
func (s *Storage) CreateKeygroup(i data.Item) error {
	key := makeKeyName(i.Keygroup, i.ID)

	err := s.db.Put([]byte(key), []byte(i.Data), nil)

	return err
}

// DeleteKeygroup deletes the given keygroup from the leveldb database.
func (s *Storage) DeleteKeygroup(i data.Item) error {
	key := makeKeyName(i.Keygroup, i.ID)

	err := s.db.Delete([]byte(key), nil)

	return err
}
