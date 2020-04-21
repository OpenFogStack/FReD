package leveldbsd

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/syndtr/goleveldb/leveldb"

	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/commons"
	"gitlab.tu-berlin.de/mcc-fred/fred/pkg/data"
)

var db *Storage

func init(){
	db = New("./test.db")
}

func TestKeygroups(t *testing.T) {
	kg := commons.KeygroupName("test-kg")
	db.CreateKeygroup(kg)
	exists := db.ExistsKeygroup(data.Item{Keygroup: kg})
	if !exists {
		t.Fatal("Keygroup does not exist after creation")
	}
	db.DeleteKeygroup(kg)
	exists = db.ExistsKeygroup(data.Item{Keygroup: kg})
	if exists {
		t.Fatal("Keygroup does still exist after deletion")
	}

}

func TestReadAll(t *testing.T) {
	kg := commons.KeygroupName("test-read-all")
	db.CreateKeygroup(kg)
	item1 := data.Item{
		Keygroup: kg,
		ID: "id-1",
		Data: "data-1",
	}
	db.Update(item1)
	item2 := data.Item{
		Keygroup: kg,
		ID: "id-2",
		Data: "data-2",
	}
	db.Update(item2)
	item3 := data.Item{
		Keygroup: kg,
		ID: "id-3",
		Data: "data-3",
	}
	db.Update(item3)

	kg2 := commons.KeygroupName("test-read-all-2")
	db.CreateKeygroup(kg)
	item21 := data.Item{
		Keygroup: kg2,
		ID: "id-1",
		Data: "data-1",
	}
	db.Update(item21)
	item22 := data.Item{
		Keygroup: kg2,
		ID: "id-2",
		Data: "data-2",
	}
	db.Update(item22)
	item23 := data.Item{
		Keygroup: kg2,
		ID: "id-3",
		Data: "data-3",
	}
	db.Update(item23)

	res, _ := db.ReadAll(kg)

	assert.Contains(t, res, item1)
	assert.Contains(t, res, item2)
	assert.Contains(t, res, item3)
}

func TestItemGet(t *testing.T)  {
	kg := commons.KeygroupName("test-kg-item")
	id := "name"
	value := "value"
	item := data.Item{
		Keygroup: kg,
		ID:       id,
		Data:     value,
	}
	db.CreateKeygroup(kg)
	db.Update(item)
	retr, err := db.Read(kg, id)
	if err != nil {
		t.Fatal(err)
	}
	if retr != value {
		t.Fatalf("Expected to get %s but got %s", value, retr)
	}
}

func TestItemAfterDeleteKeygroup(t *testing.T)  {
	kg := commons.KeygroupName("test-kg-item-delete")
	id := "ndel"
	value := "vdel"
	item := data.Item{
		Keygroup: kg,
		ID:       id,
		Data:     value,
	}
	db.CreateKeygroup(kg)
	db.Update(item)
	db.DeleteKeygroup(kg)

	retr, err := db.Read(kg, id)
	if err == nil {
		t.Fatalf("Expected an error, but got %s", retr)
	}
	if err != leveldb.ErrNotFound {
		t.Fatalf("Expected the error to be ErrNotFound, but got %s", err)
	}
}