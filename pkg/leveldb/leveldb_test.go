package leveldb

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/syndtr/goleveldb/leveldb"
)

var db *Storage

// TODO: better tests, maybe even for all packages that implement the Store interface?

func init() {
	db = New("./test.db")
}

func TestKeygroups(t *testing.T) {
	kg := "test-kg"
	db.CreateKeygroup(kg)
	exists := db.ExistsKeygroup(kg)
	if !exists {
		t.Fatal("Keygroup does not exist after creation")
	}
	db.DeleteKeygroup(kg)
	exists = db.ExistsKeygroup(kg)
	if exists {
		t.Fatal("Keygroup does still exist after deletion")
	}

}

func TestReadAll(t *testing.T) {
	kg := "test-read-all"
	db.CreateKeygroup(kg)

	db.Update(kg, "id-1", "data-1")

	db.Update(kg, "id-2", "data-2")

	db.Update(kg, "id-3", "data-3")

	kg2 := "test-read-all-2"

	db.CreateKeygroup(kg2)

	db.Update(kg2, "id-1", "data-1")

	db.Update(kg2, "id-2", "data-2")

	db.Update(kg2, "id-3", "data-3")

	res, _ := db.ReadAll(kg)

	assert.Equal(t, res["id-1"], "data-1")
	assert.Equal(t, res["id-2"], "data-2")
	assert.Equal(t, res["id-3"], "data-3")

}

func TestItemGet(t *testing.T) {
	kg := "test-kg-item"
	id := "name"
	value := "value"

	db.CreateKeygroup(kg)
	db.Update(kg, id, value)

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

	db.CreateKeygroup(kg)
	db.Update(kg, id, value)
	db.DeleteKeygroup(kg)

	retr, err := db.Read(kg, id)
	if err == nil {
		t.Fatalf("Expected an error, but got %s", retr)
	}
	if err != leveldb.ErrNotFound {
		t.Fatalf("Expected the error to be ErrNotFound, but got %s", err)
	}
}
