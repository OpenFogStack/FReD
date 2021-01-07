package badgerdb

import (
	"strconv"
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

func TestAppend(t *testing.T) {
	kg := "log"
	v1 := "value-1"
	v2 := "value-2"

	err := db.CreateKeygroup(kg)

	if err != nil {
		log.Err(err).Msg(err.(*errors.Error).ErrorStack())
	}

	key1, err := db.Append(kg, v1, 0)

	if err != nil {
		t.Fatal(err)
	}
	if key1 != "0" {
		t.Fatalf("Expected to get %s but got %s", "0", key1)
	}

	key2, err := db.Append(kg, v2, 0)

	if err != nil {
		t.Fatal(err)
	}
	if key2 != "1" {
		t.Fatalf("Expected to get %s but got %s", "0", key2)
	}

	for i := 2; i < 100; i++ {
		v := "value-" + strconv.Itoa(i)
		key, err := db.Append(kg, v, 0)

		if err != nil {
			t.Fatal(err)
		}

		if key != strconv.Itoa(i) {
			t.Fatalf("Expected to get %s but got %s", strconv.Itoa(i), key)
		}
	}
}

func TestTriggerNodes(t *testing.T) {
	kg := "kg1"

	err := db.CreateKeygroup(kg)

	if err != nil {
		t.Fatal(err)
	}

	t1 := "t1"
	t1host := "1.1.1.1:3000"

	t2 := "t2"
	t2host := "2.2.2.2:3000"

	t3 := "t3"
	t3host := "3.3.3.3:3000"

	err = db.AddKeygroupTrigger(kg, t1, t1host)

	if err != nil {
		t.Fatal(err)
	}

	err = db.AddKeygroupTrigger(kg, t1, t1host)

	if err != nil {
		t.Fatal(err)
	}

	err = db.AddKeygroupTrigger(kg, t2, t2host)

	if err != nil {
		t.Fatal(err)
	}

	tList, err := db.GetKeygroupTrigger(kg)

	if err != nil {
		t.Fatal(err)
	}

	if len(tList) != 2 {
		t.Fatal("not the right number of triggers for this keygroup")
	}

	if _, ok := tList[t1]; !ok {
		t.Fatal("t1 not in list of triggers for this keygroup")
	}

	if _, ok := tList[t2]; !ok {
		t.Fatal("t2 not in list of triggers for this keygroup")
	}

	if host := tList[t1]; host != t1host {
		t.Fatal("t1host not correct")
	}

	if host := tList[t2]; host != t2host {
		t.Fatal("t1host not correct")
	}

	err = db.DeleteKeygroupTrigger(kg, t1)

	if err != nil {
		t.Fatal(err)
	}

	if len(tList) != 1 {
		t.Fatal("not the right number of triggers for this keygroup")
	}

	if _, ok := tList[t2]; !ok {
		t.Fatal("t2 not in list of triggers for this keygroup")
	}

	if host := tList[t2]; host != t2host {
		t.Fatal("t1host not correct")
	}

	err = db.DeleteKeygroup(kg)

	if err != nil {
		t.Fatal(err)
	}

	tList, err = db.GetKeygroupTrigger(kg)

	if err == nil {
		t.Fatal("got keygroup triggers for nonexistent keygroup")
	}

	if len(tList) != 0 {
		t.Fatal("got keygroup triggers for nonexistent keygroup")
	}

	err = db.AddKeygroupTrigger(kg, t3, t3host)

	if err == nil {
		t.Fatal("added keygroup trigger to nonexistent keygroup")
	}

}
