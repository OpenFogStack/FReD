package leveldbsd

import (
	"github.com/syndtr/goleveldb/leveldb"
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
func (s *Storage) Read(kgname string, id string) (string, error) {
	key := makeKeyName(kgname, id)

	data, err := s.db.Get([]byte(key), nil)

	return string(data), err
}

// Update updates the item with the specified id in the specified keygroup.
func (s *Storage) Update(kgname string, id string, data string) error {
	key := makeKeyName(kgname, id)

	err := s.db.Put([]byte(key), []byte(data), nil)

	return err
}

// Delete deletes the item with the specified id from the specified keygroup.
func (s *Storage) Delete(kgname string, id string) error {
	key := makeKeyName(kgname, id)

	err := s.db.Delete([]byte(key), nil)

	return err
}

// CreateKeygroup does nothing.
func (s *Storage) CreateKeygroup(kgname string) error {
	return nil
}

// DeleteKeygroup does nothing.
func (s *Storage) DeleteKeygroup(kgname string) error {
	return nil
}
