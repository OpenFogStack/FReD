package app

import (
	"errors"

	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/memorykg"
	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/memorysd"
)

// Storage is an interface that abstracts the component that stores actual Keygroups data.
type Storage interface {
	Create(kg string, data string) (uint64, error)
	Read(kg string, id uint64) (string, error)
	Update(kg string, id uint64, data string) error
	Delete(kg string, id uint64) error
	CreateKeygroup(kg string) error
	DeleteKeygroup(kg string) error
}

// Keygroups is an interface that abstracts the component that stores Keygroups metadata.
type Keygroups interface {
	Create(kg string) error
	Delete(kg string) error
	Exists(kg string) bool
}

// App has both a store for Keygroups and for Storage.
type App struct {
	kg Keygroups
	sd Storage
}

// CreateKeygroup creates a new keygroup with the specified name in Storage.
func (a *App) CreateKeygroup(kgname string) error {
	err := a.kg.Create(kgname)

	if err != nil {
		return err
	}

	err = a.sd.CreateKeygroup(kgname)

	if err != nil {
		return err
	}

	return nil
}

// DeleteKeygroup removes the keygroup with the specified name from Storage.
func (a *App) DeleteKeygroup(kgname string) error {
	err := a.kg.Delete(kgname)

	if err != nil {
		return err
	}

	err = a.sd.DeleteKeygroup(kgname)

	if err != nil {
		return err
	}

	return nil
}

// Create creates a new item in the specified keygroup.
func (a *App) Create(kgname string, data string) (uint64, error) {

	if !a.kg.Exists(kgname) {
		return 0, errors.New("not found")
	}

	id, err := a.sd.Create(kgname, data)

	if err != nil {
		return 0, err
	}

	return id, nil
}

// Read returns an item with the specified id from the specified keygroup.
func (a *App) Read(kgname string, id uint64) (string, error) {
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
func (a *App) Update(kgname string, id uint64, data string) error {
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
func (a *App) Delete(kgname string, id uint64) error {
	if !a.kg.Exists(kgname) {
		return errors.New("keygroup not found")
	}

	err := a.sd.Delete(kgname, uint64(id))

	if err != nil {
		return err
	}

	return nil
}

// New create a new App.
func New() (a *App) {
	a = &App{
		kg: memorykg.New(),
		sd: memorysd.New(),
	}

	return
}
