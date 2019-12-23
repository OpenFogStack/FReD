package data

import "errors"

type service struct {
	iS Store
}

// New creates a new data manipulation service.
func New(iS Store) Service {
	return &service{
		iS: iS,
	}
}

// Read returns an item from the key-value store.
func (s *service) Read(i Item) (Item, error) {
	err := checkItem(i)

	if err != nil {
		return Item{}, err
	}

	if !s.iS.ExistsKeygroup(i) {
		return Item{}, errors.New("data: keygroup not found")
	}

	data, err := s.iS.Read(i)

	if err != nil {
		return Item{}, err
	}

	return data, nil
}

// Update updates an item in the key-value store.
func (s *service) Update(i Item) error {
	err := checkItem(i)

	if err != nil {
		return err
	}

	if !s.iS.ExistsKeygroup(i) {
		return errors.New("data: keygroup not found")
	}

	err = s.iS.Update(i)

	if err != nil {
		return err
	}

	return nil
}

// Delete removes an item from the key-value store.
func (s *service) Delete(i Item) error {
	err := checkItem(i)

	if err != nil {
		return err
	}

	if !s.iS.ExistsKeygroup(i) {
		return errors.New("data: keygroup not found")
	}

	err = s.iS.Delete(i)

	if err != nil {
		return err
	}

	return nil
}

// CreateKeygroup creates a new keygroup in the store.
func (s *service) CreateKeygroup(i Item) error {
	err := checkKeygroup(i)

	if err != nil {
		return err
	}

	return nil
}

// DeleteKeygroup deletes a keygroup from the store.
func (s *service) DeleteKeygroup(i Item) error {
	err := checkKeygroup(i)

	if err != nil {
		return err
	}

	return nil
}
