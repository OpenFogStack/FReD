package data

import (
	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/commons"
)

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
func (s *service) Read(kg commons.KeygroupName, id string) (string, error) {
	err := checkKGandID(kg, id)

	if err != nil {
		return "", err
	}

	data, err := s.iS.Read(kg, id)

	if err != nil {
		return "", err
	}

	return data, nil
}

// ReadAll returns all items of a particular keygroup from the key-value store.
func (s *service) ReadAll(kg commons.KeygroupName) ([]Item, error) {
	err := checkKeygroup(kg)

	if err != nil {
		return nil, err
	}

	data, err := s.iS.ReadAll(kg)

	return data, err
}

// Update updates an item in the key-value store.
func (s *service) Update(i Item) error {
	err := checkItem(i)

	if err != nil {
		return err
	}

	err = s.iS.Update(i)

	if err != nil {
		return err
	}

	return nil
}

// Delete removes an item from the key-value store.
func (s *service) Delete(kg commons.KeygroupName, id string) error {
	err := checkKGandID(kg, id)

	if err != nil {
		return err
	}

	err = s.iS.Delete(kg, id)

	if err != nil {
		return err
	}

	return nil
}

// TODO Tobias check whether the CreateKeygroup function was populated the right way
// Should be in exthandler and inthandler
// CreateKeygroup creates a new keygroup in the store.
func (s *service) CreateKeygroup(kg commons.KeygroupName) error {
	err := checkKeygroup(kg)

	if err != nil {
		return err
	}

	return nil
}

// TODO implement me!
// DeleteKeygroup deletes a keygroup from the store.
func (s *service) DeleteKeygroup(kg commons.KeygroupName) error {
	err := checkKeygroup(kg)

	if err != nil {
		return err
	}

	return nil
}
