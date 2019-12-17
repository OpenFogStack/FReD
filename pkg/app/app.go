package app

import (
	"errors"

	"github.com/mmcloughlin/geohash"

	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/leveldbsd"
	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/memorykg"
	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/memorysd"
)

// Storage is an interface that abstracts the component that stores actual Keygroups data.
type Storage interface {
	Read(kg string, id string) (string, error)
	Update(kg string, id string, data string) error
	Delete(kg string, id string) error
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
	ID string
}

// CreateKeygroup creates a new keygroup with the specified name in Storage.
func (a *App) CreateKeygroup(kgname string) error {
	if err := checkParameters(kgname); err != nil {
		return err
	}

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
	if err := checkParameters(kgname); err != nil {
		return err
	}

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

// Read returns an item with the specified id from the specified keygroup.
func (a *App) Read(kgname string, id string) (string, error) {
	if err := checkParameters(kgname, id); err != nil {
		return "", err
	}

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
func (a *App) Update(kgname string, id string, data string) error {
	if err := checkParameters(kgname, id, data); err != nil {
		return err
	}

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
func (a *App) Delete(kgname string, id string) error {
	if err := checkParameters(kgname, id); err != nil {
		return err
	}

	if !a.kg.Exists(kgname) {
		return errors.New("keygroup not found")
	}

	err := a.sd.Delete(kgname, id)

	if err != nil {
		return err
	}

	return nil
}

func checkParameters(params ...string) error{
	for _, p := range params {
		if p == "" {
			return errors.New("empty parameter")
		}
	}

	return nil
}

// New create a new App.
func New(lat float64, lng float64, storageRuntime string, dbPath string) (a *App, err error) {
	var sd Storage = nil

	switch storageRuntime {
	case "leveldb":
		sd = leveldbsd.New(dbPath)
	case "memory":
		sd = memorysd.New()
	default:
		return nil, errors.New("unknown storage backend")
	}

	a = &App{
		kg: memorykg.New(),
		sd: sd,
		ID: geohash.Encode(lng, lng),
	}

	return
}
