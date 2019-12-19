package app

import (
	"errors"

	"github.com/mmcloughlin/geohash"

	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/leveldbsd"
	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/memorykg"
	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/memorysd"
)

// App has both a store for Keygroups and for Storage.
type App struct {
	kg Keygroups
	sd Storage
	ID string
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
		ID: geohash.Encode(lat, lng),
	}

	return
}