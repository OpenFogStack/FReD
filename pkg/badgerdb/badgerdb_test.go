package badgerdb

import (
	"testing"
	"time"

	"github.com/go-errors/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

var db *Storage

// TODO: better tests, maybe even for all packages that implement the Store interface?

func init() {
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	db = New("./test.db")
}

func TestKeygroups(t *testing.T) {
	kg := "test-kg"
	err := db.CreateKeygroup(kg)

	if err != nil {
		log.Err(err).Msg(err.(*errors.Error).ErrorStack())
	}

	exists := db.ExistsKeygroup(kg)
	if !exists {
		t.Fatal("Keygroup does not exist after creation")
	}

	err = db.DeleteKeygroup(kg)

	if err != nil {
		log.Err(err).Msg(err.(*errors.Error).ErrorStack())
	}

	exists = db.ExistsKeygroup(kg)
	if exists {
		t.Fatal("Keygroup does still exist after deletion")
	}

}

func TestReadAll(t *testing.T) {
	kg := "test-read-all"
	err := db.CreateKeygroup(kg)

	if err != nil {
		log.Err(err).Msg(err.(*errors.Error).ErrorStack())
	}

	err = db.Update(kg, "id-1", "data-1", 0)

	if err != nil {
		log.Err(err).Msg(err.(*errors.Error).ErrorStack())
	}

	err = db.Update(kg, "id-2", "data-2", 0)

	if err != nil {
		log.Err(err).Msg(err.(*errors.Error).ErrorStack())
	}

	err = db.Update(kg, "id-3", "data-3", 0)

	if err != nil {
		log.Err(err).Msg(err.(*errors.Error).ErrorStack())
	}

	kg2 := "test-read-all-2"

	err = db.CreateKeygroup(kg2)

	if err != nil {
		log.Err(err).Msg(err.(*errors.Error).ErrorStack())
	}

	err = db.Update(kg2, "id-1", "data-1", 0)

	if err != nil {
		log.Err(err).Msg(err.(*errors.Error).ErrorStack())
	}

	err = db.Update(kg2, "id-2", "data-2", 0)

	if err != nil {
		log.Err(err).Msg(err.(*errors.Error).ErrorStack())
	}

	err = db.Update(kg2, "id-3", "data-3", 0)

	if err != nil {
		log.Err(err).Msg(err.(*errors.Error).ErrorStack())
	}

	res, err := db.ReadAll(kg)

	if err != nil {
		log.Err(err).Msg(err.(*errors.Error).ErrorStack())
	}

	assert.Equal(t, "data-1", res["id-1"])
	assert.Equal(t, "data-2", res["id-2"])
	assert.Equal(t, "data-3", res["id-3"])

}

func TestItemGet(t *testing.T) {
	kg := "test-kg-item"
	id := "name"
	value := "value"

	err := db.CreateKeygroup(kg)

	if err != nil {
		log.Err(err).Msg(err.(*errors.Error).ErrorStack())
	}

	err = db.Update(kg, id, value, 0)

	if err != nil {
		log.Err(err).Msg(err.(*errors.Error).ErrorStack())
	}

	retr, err := db.Read(kg, id)
	if err != nil {
		t.Fatal(err)
	}
	if retr != value {
		t.Fatalf("Expected to get %s but got %s", value, retr)
	}
}

func TestItemAfterDeleteKeygroup(t *testing.T) {
	kg := "test-kg-item-delete"
	id := "ndel"
	value := "vdel"

	err := db.CreateKeygroup(kg)

	if err != nil {
		log.Err(err).Msg(err.(*errors.Error).ErrorStack())
	}

	err = db.Update(kg, id, value, 0)

	if err != nil {
		log.Err(err).Msg(err.(*errors.Error).ErrorStack())
	}

	err = db.DeleteKeygroup(kg)

	if err != nil {
		log.Err(err).Msg(err.(*errors.Error).ErrorStack())
	}

	retr, err := db.Read(kg, id)
	if err == nil {
		t.Fatalf("Expected an error, but got %s", retr)
	}
}

func TestExpiry(t *testing.T) {
	kg := "test-kg-item"
	id := "name"
	value := "value"

	err := db.CreateKeygroup(kg)

	if err != nil {
		log.Err(err).Msg(err.(*errors.Error).ErrorStack())
	}

	err = db.Update(kg, id, value, 10)

	if err != nil {
		log.Err(err).Msg(err.(*errors.Error).ErrorStack())
	}

	retr, err := db.Read(kg, id)
	if err != nil {
		t.Fatal(err)
	}
	if retr != value {
		t.Fatalf("Expected to get %s but got %s", value, retr)
	}

	time.Sleep(10 * time.Second)

	_, err = db.Read(kg, id)
	if err == nil {
		t.Fatal(err)
	}
}
