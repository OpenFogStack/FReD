package fred

import "github.com/go-errors/errors"

// Store is an interface for the storage medium that the key-value val items are persisted on.
type Store interface {
	// Needs: keygroup, id, val
	Update(kg, id, val string, expiry int) error
	// Needs: keygroup, id
	Delete(kg, id string) error
	// Needs: keygroup, val, Returns: key
	Append(kg, val string, expiry int) (string, error)
	// Needs: keygroup, id; Returns: val
	Read(kg, id string) (string, error)
	// Needs: keygroup, id, range; Returns: ids and values
	ReadSome(kg, id string, count uint64) (map[string]string, error)
	// Needs: keygroup; Returns: ids and values
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
	// Needs: keygroup
	ExistsKeygroup(kg string) bool
	// Needs: keygroup, trigger node id, trigger node host
	AddKeygroupTrigger(kg string, id string, host string) error
	// Needs: keygroup, trigger node id
	DeleteKeygroupTrigger(kg string, id string) error
	// Needs: keygroup; Returns map: trigger node id -> trigger node host
	GetKeygroupTrigger(kg string) (map[string]string, error)
	// Close indicates that the underlying store should be closed as it is no longer needed.
	Close() error
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

	if !s.iS.ExistsKeygroup(string(kg)) {
		return "", errors.Errorf("no such keygroup in store: %#v", kg)
	}

	data, err := s.iS.Read(string(kg), id)

	if err != nil {
		return "", err
	}

	return data, nil
}

// Scan returns a list of count items starting with id from the key-value store.
func (s *storeService) scan(kg KeygroupName, id string, count uint64) ([]Item, error) {
	err := checkKGandID(kg, id)

	if err != nil {
		return nil, err
	}

	if !s.iS.ExistsKeygroup(string(kg)) {
		return nil, errors.Errorf("no such keygroup in store: %#v", kg)
	}

	if !s.exists(Item{
		Keygroup: kg,
		ID:       id,
	}) {
		return nil, errors.Errorf("no such item %s in keygroup %#v", id, kg)
	}

	data, err := s.iS.ReadSome(string(kg), id, count)

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

	return i, nil
}

// ReadAll returns all items of a particular keygroup from the key-value store.
func (s *storeService) readAll(kg KeygroupName) ([]Item, error) {
	err := checkKeygroup(kg)

	if err != nil {
		return nil, err
	}

	if !s.iS.ExistsKeygroup(string(kg)) {
		return nil, errors.Errorf("no such keygroup in store: %#v", kg)
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

	return i, nil
}

// exists checks if an item exists in the key-value store.
func (s *storeService) exists(i Item) bool {
	return s.iS.Exists(string(i.Keygroup), i.ID)
}

// Update updates an item in the key-value store.
func (s *storeService) update(i Item, expiry int) error {
	err := checkItem(i)

	if err != nil {
		return err
	}

	if !s.iS.ExistsKeygroup(string(i.Keygroup)) {
		return errors.Errorf("no such keygroup in store: %#v", i.Keygroup)
	}

	err = s.iS.Update(string(i.Keygroup), i.ID, i.Val, expiry)

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

	if !s.iS.ExistsKeygroup(string(kg)) {
		return errors.Errorf("no such keygroup in store: %#v", kg)
	}

	err = s.iS.Delete(string(kg), id)

	if err != nil {
		return err
	}

	return nil
}

// append appends an item in the key-value store.
func (s *storeService) append(i Item, expiry int) (Item, error) {
	if !s.iS.ExistsKeygroup(string(i.Keygroup)) {
		return i, errors.Errorf("no such keygroup in store: %#v", i.Keygroup)
	}

	k, err := s.iS.Append(string(i.Keygroup), i.Val, expiry)

	if err != nil {
		return i, err
	}

	i.ID = k

	return i, nil
}

// TODO Tobias check whether the CreateKeygroup function was populated the right way
// Should be in exthandler and inthandler
// CreateKeygroup creates a new keygroup in the store.
func (s *storeService) createKeygroup(kg KeygroupName) error {
	err := checkKeygroup(kg)

	if err != nil {
		return err
	}

	if s.iS.ExistsKeygroup(string(kg)) {
		return errors.Errorf("keygroup already in store: %#v", kg)
	}

	err = s.iS.CreateKeygroup(string(kg))

	if err != nil {
		return err
	}

	return nil
}

// DeleteKeygroup deletes a keygroup from the store.
func (s *storeService) deleteKeygroup(kg KeygroupName) error {
	err := checkKeygroup(kg)

	if err != nil {
		return err
	}

	if !s.iS.ExistsKeygroup(string(kg)) {
		return errors.Errorf("no such keygroup in store: %#v", kg)
	}

	err = s.iS.DeleteKeygroup(string(kg))

	if err != nil {
		return err
	}

	return nil
}

func (s *storeService) addKeygroupTrigger(kg KeygroupName, t Trigger) error {
	err := checkKeygroup(kg)

	if err != nil {
		return err
	}

	if !s.iS.ExistsKeygroup(string(kg)) {
		return errors.Errorf("no such keygroup in store: %#v", kg)
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
		return errors.Errorf("no such keygroup in store: %#v", kg)
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
		return nil, errors.Errorf("no such keygroup in store: %#v", kg)
	}

	t, err := s.iS.GetKeygroupTrigger(string(kg))

	if err != nil {
		return nil, err
	}

	var tn []Trigger

	for id, host := range t {
		tn = append(tn, Trigger{
			ID:   id,
			Host: host,
		})
	}

	return tn, nil
}
