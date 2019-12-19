package app

import "errors"

// Storage is an interface that abstracts the component that stores actual Keygroups data.
type Storage interface {
	Read(kg string, id string) (string, error)
	Update(kg string, id string, data string) error
	Delete(kg string, id string) error
	CreateKeygroup(kg string) error
	DeleteKeygroup(kg string) error
}

// Read returns an item with the specified id from the specified keygroup.
func (a *App) read(kgname string, id string) (string, error) {
	if !a.kg.Exists(kgname) {
		return "", errors.New("keygroup not found")
	}

	data, err := a.sd.Read(kgname, id)

	if err != nil {
		return "", err
	}

	return data, nil
}

// Update updates the item with the specified id in the specified keygroup.
func (a *App) update(kgname string, id string, data string) error {
	if !a.kg.Exists(kgname) {
		return errors.New("keygroup not found")
	}

	err := a.sd.Update(kgname, id, data)

	if err != nil {
		return err
	}

	return nil
}

// Delete deletes the item with the specified id from the specified keygroup.
func (a *App) delete(kgname string, id string) error {
	if !a.kg.Exists(kgname) {
		return errors.New("keygroup not found")
	}

	err := a.sd.Delete(kgname, id)

	if err != nil {
		return err
	}

	return nil
}
