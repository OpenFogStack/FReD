package fred

// Store is an interface for the storage medium that the key-value val items are persisted on.
type Store interface {
	// Needs: keygroup, id, val
	Update(kg, id, val string) error
	// Needs: keygroup, id
	Delete(kg, id string) error
	// Needs: keygroup, id; Returns: val
	Read(kg, id string) (string, error)
	// Needs: keygroup; Returns: keygroup, id, val
	ReadAll(kg string) (map[string]string, error)
	// Needs: keygroup, Returns:[] keygroup, id
	IDs(kg string) ([]string, error)
	// Needs: keygroup, id
	Exists(kg, id string) bool
	// Needs: keygroup, Returns: err
	// Doesnt really need to store the KG itself, that is keygroup/store.go's job.
	// But might be useful for databases using it
	CreateKeygroup(kg string) error
	// Same as with CreateKeygroup
	DeleteKeygroup(kg string) error
}

type storeService struct {
	iS Store
}

// NewStoreService creates a new val manipulation service.
func newStoreService(iS Store) *storeService {
	return &storeService{
		iS: iS,
	}
}

// Read returns an item from the key-value store.
func (s *storeService) read(kg KeygroupName, id string) (string, error) {
	err := checkKGandID(kg, id)

	if err != nil {
		return "", err
	}

	data, err := s.iS.Read(string(kg), id)

	if err != nil {
		return "", err
	}

	return data, nil
}

// ReadAll returns all items of a particular keygroup from the key-value store.
func (s *storeService) readAll(kg KeygroupName) ([]Item, error) {
	err := checkKeygroup(kg)

	if err != nil {
		return nil, err
	}

	data, err := s.iS.ReadAll(string(kg))

	if err != nil {
		return nil, err
	}

	i := make([]Item, len(data))
	c := 0

	for id, val := range data {
		i[c] = Item{
			Keygroup: kg,
			ID:       id,
			Val:      val,
		}

		c++
	}

	return i, err
}

// Update updates an item in the key-value store.
func (s *storeService) update(i Item) error {
	err := checkItem(i)

	if err != nil {
		return err
	}

	err = s.iS.Update(string(i.Keygroup), i.ID, i.Val)

	if err != nil {
		return err
	}

	return nil
}

// Delete removes an item from the key-value store.
func (s *storeService) delete(kg KeygroupName, id string) error {
	err := checkKGandID(kg, id)

	if err != nil {
		return err
	}

	err = s.iS.Delete(string(kg), id)

	if err != nil {
		return err
	}

	return nil
}

// TODO Tobias check whether the CreateKeygroup function was populated the right way
// Should be in exthandler and inthandler
// CreateKeygroup creates a new keygroup in the store.
func (s *storeService) createKeygroup(kg KeygroupName) error {
	err := checkKeygroup(kg)

	if err != nil {
		return err
	}

	return nil
}

// TODO implement me!
// DeleteKeygroup deletes a keygroup from the store.
func (s *storeService) deleteKeygroup(kg KeygroupName) error {
	err := checkKeygroup(kg)

	if err != nil {
		return err
	}

	return nil
}
