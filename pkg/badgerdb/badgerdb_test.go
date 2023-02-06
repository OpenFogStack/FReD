package badgerdb

import (
	"fmt"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/DistributedClocks/GoVector/govec/vclock"
	"github.com/go-errors/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

const badgerDBPath = "./test.db"

var db *Storage

// TODO: better tests, maybe even for all packages that implement the Store interface?

func TestMain(m *testing.M) {
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	log.Logger = log.Output(
		zerolog.ConsoleWriter{
			Out:     os.Stderr,
			NoColor: false,
		},
	)

	fInfo, err := os.Stat(badgerDBPath)

	if err == nil {
		if !fInfo.IsDir() {
			panic(errors.Errorf("%s is not a directory", badgerDBPath))
		}

		err = os.RemoveAll(badgerDBPath)
		if err != nil {
			panic(err)
		}
	}

	db = New(badgerDBPath)

	stat := m.Run()

	fInfo, err = os.Stat(badgerDBPath)

	if err == nil {
		if !fInfo.IsDir() {
			panic(errors.Errorf("%s is not a directory", badgerDBPath))
		}

		err = os.RemoveAll(badgerDBPath)
		if err != nil {
			panic(err)
		}
	}

	os.Exit(stat)
}

func TestKeygroups(t *testing.T) {
	kg := "test-kg"
	err := db.CreateKeygroup(kg)

	assert.NoError(t, err)

	exists := db.ExistsKeygroup(kg)

	assert.True(t, exists, "Keygroup does not exist after creation")

	err = db.DeleteKeygroup(kg)

	assert.NoError(t, err)

	exists = db.ExistsKeygroup(kg)

	assert.False(t, exists, "Keygroup still exists after deletion")

}

func TestReadSome(t *testing.T) {
	kg := "test-kg-scan"
	updates := 10
	scanStart := 3
	scanRange := 5

	err := db.CreateKeygroup(kg)

	assert.NoError(t, err)
	// 2. put in a bunch of items
	ids := make([]string, updates)
	vals := make([]string, updates)

	for i := 0; i < updates; i++ {
		ids[i] = "id" + strconv.Itoa(i)
		vals[i] = "val" + strconv.Itoa(i)

		err = db.Update(kg, ids[i], vals[i], 0, vclock.VClock{})

		assert.NoError(t, err)

	}

	keys, values, _, err := db.ReadSome(kg, "id"+strconv.Itoa(scanStart), uint64(scanRange))

	assert.NoError(t, err)

	assert.Len(t, keys, scanRange)
	assert.Len(t, values, scanRange)

	for i := scanStart; i < scanStart+scanRange; i++ {
		assert.Contains(t, keys, ids[i])
		assert.Contains(t, values, vals[i])
		assert.Equal(t, ids[i], keys[i-scanStart])
		assert.Equal(t, vals[i], values[i-scanStart])
	}
}

func TestReadAll(t *testing.T) {
	kg := "test-read-all"
	err := db.CreateKeygroup(kg)

	assert.NoError(t, err)

	err = db.Update(kg, "id-1", "data-1", 0, vclock.VClock{})

	assert.NoError(t, err)

	err = db.Update(kg, "id-2", "data-2", 0, vclock.VClock{})

	assert.NoError(t, err)

	err = db.Update(kg, "id-3", "data-3", 0, vclock.VClock{})

	assert.NoError(t, err)

	kg2 := "test-read-all-2"

	err = db.CreateKeygroup(kg2)

	assert.NoError(t, err)

	err = db.Update(kg2, "id-4", "data-4", 0, vclock.VClock{})

	assert.NoError(t, err)

	err = db.Update(kg2, "id-5", "data-5", 0, vclock.VClock{})

	assert.NoError(t, err)

	err = db.Update(kg2, "id-6", "data-6", 0, vclock.VClock{})

	assert.NoError(t, err)

	keys, values, _, err := db.ReadAll(kg)

	assert.NoError(t, err)

	assert.Len(t, keys, 3)
	assert.Len(t, values, 3)
	assert.Equal(t, "id-1", keys[0])
	assert.Equal(t, "data-1", values[0])
	assert.Equal(t, "id-2", keys[1])
	assert.Equal(t, "data-2", values[1])
	assert.Equal(t, "id-3", keys[2])
	assert.Equal(t, "data-3", values[2])
}

func TestIDs(t *testing.T) {
	kg := "test-ids"
	err := db.CreateKeygroup(kg)

	assert.NoError(t, err)

	err = db.Update(kg, "id-1", "data-1", 0, vclock.VClock{})

	assert.NoError(t, err)

	err = db.Update(kg, "id-2", "data-2", 0, vclock.VClock{})

	assert.NoError(t, err)

	err = db.Update(kg, "id-3", "data-3", 0, vclock.VClock{})

	assert.NoError(t, err)

	kg2 := "test-read-all-2"

	err = db.CreateKeygroup(kg2)

	assert.NoError(t, err)

	err = db.Update(kg2, "id-1", "data-1", 0, vclock.VClock{})

	assert.NoError(t, err)

	err = db.Update(kg2, "id-2", "data-2", 0, vclock.VClock{})

	assert.NoError(t, err)

	err = db.Update(kg2, "id-3", "data-3", 0, vclock.VClock{})

	assert.NoError(t, err)

	res, err := db.IDs(kg)

	assert.NoError(t, err)

	assert.Len(t, res, 3)
	assert.Contains(t, res, "id-1")
	assert.Contains(t, res, "id-2")
	assert.Contains(t, res, "id-3")
}

func TestItemExists(t *testing.T) {
	kg := "test-kg-item"
	id := "name"
	id2 := "name2"
	value := "value"

	err := db.CreateKeygroup(kg)

	assert.NoError(t, err)

	err = db.Update(kg, id, value, 0, vclock.VClock{})

	assert.NoError(t, err)
	ex := db.Exists(kg, id)

	assert.True(t, ex, "exists says existing item doesn't exist")

	ex = db.Exists(kg, id2)
	assert.False(t, ex, "exists says non-existent item exists")

}

func TestItemGet(t *testing.T) {
	kg := "test-kg-item"
	id := "name"
	value := "value"

	err := db.CreateKeygroup(kg)

	assert.NoError(t, err)

	err = db.Update(kg, id, value, 0, vclock.VClock{})

	assert.NoError(t, err)

	retr, _, _, err := db.Read(kg, id)
	assert.NoError(t, err)

	assert.Len(t, retr, 1)
	assert.Equal(t, retr[0], value)
}

func TestItemDelete(t *testing.T) {
	kg := "test-kg-item-delete"
	id := "name"
	id2 := "name2"
	value := "value"

	err := db.CreateKeygroup(kg)

	assert.NoError(t, err)
	err = db.Update(kg, id, value, 0, vclock.VClock{})

	assert.NoError(t, err)

	retr, _, _, err := db.Read(kg, id)

	assert.NoError(t, err)

	assert.Len(t, retr, 1)
	assert.Equal(t, retr[0], value)

	err = db.Delete(kg, id, nil)

	assert.NoError(t, err, "read a deleted item")

	_, _, found, err := db.Read(kg, id)

	assert.NoError(t, err)
	assert.False(t, found)

	err = db.Delete(kg, id2, vclock.VClock{})

	assert.NoError(t, err, "deleting non-existent keys should be allowed")

}

func TestItemAfterDeleteKeygroup(t *testing.T) {
	kg := "test-kg-item-delete"
	id := "ndel"
	value := "vdel"

	err := db.CreateKeygroup(kg)

	assert.NoError(t, err)

	err = db.Update(kg, id, value, 0, vclock.VClock{})

	assert.NoError(t, err)

	err = db.DeleteKeygroup(kg)

	assert.NoError(t, err)

	_, _, found, err := db.Read(kg, id)

	assert.NoError(t, err)
	assert.False(t, found)
}

func TestExpiry(t *testing.T) {
	kg := "test-kg-expiry"
	id := "name"
	value := "value"

	err := db.CreateKeygroup(kg)

	assert.NoError(t, err)

	err = db.Update(kg, id, value, 10, vclock.VClock{})

	assert.NoError(t, err)

	retr, _, _, err := db.Read(kg, id)

	assert.NoError(t, err)

	assert.Len(t, retr, 1)
	assert.Equal(t, retr[0], value)

	exists := db.Exists(kg, id)
	assert.True(t, exists)

	time.Sleep(11 * time.Second)

	_, _, found, err := db.Read(kg, id)

	assert.NoError(t, err)
	assert.False(t, found)

	exists = db.Exists(kg, id)
	assert.False(t, exists)
}

func TestAppend(t *testing.T) {
	kg := "log"
	v1 := "value-1"
	v2 := "value-2"

	err := db.CreateKeygroup(kg)

	assert.NoError(t, err)

	err = db.Append(kg, "0", v1, 0)

	assert.NoError(t, err)

	err = db.Append(kg, "1", v2, 0)

	assert.NoError(t, err)

	for i := 2; i < 100; i++ {
		v := "value-" + strconv.Itoa(i)
		id := strconv.Itoa(i)
		err := db.Append(kg, id, v, 0)

		assert.NoError(t, err)
	}
}

func TestConcurrentAppend(t *testing.T) {
	kg := "logconcurrent"
	concurrent := 4
	items := 100

	err := db.CreateKeygroup(kg)

	assert.NoError(t, err)

	keys := make([]map[string]struct{}, concurrent)
	done := make(chan struct{})

	for i := 0; i < concurrent; i++ {
		keys[i] = make(map[string]struct{})
		go func(id int, keys *map[string]struct{}) {
			for j := 2; j < items; j++ {
				v := fmt.Sprintf("value-%d-%d", id, j)
				key := strconv.Itoa(items*id + j)
				err := db.Append(kg, key, v, 0)

				assert.NoError(t, err)

				(*keys)[key] = struct{}{}
			}
			done <- struct{}{}
		}(i, &keys[i])
	}

	for i := 0; i < concurrent; i++ {
		<-done
	}

	for i, k := range keys {
		for key := range k {
			found := false
			for j := i + 1; j < concurrent; j++ {
				_, ok := keys[j][key]
				found = found || ok
			}
			assert.False(t, found, "key given out multiple times")
		}
	}

}

func TestDualAppend(t *testing.T) {
	kg := "logdual"
	v1 := "value-1"
	v2 := "value-2"

	err := db.CreateKeygroup(kg)

	assert.NoError(t, err)

	err = db.Append(kg, "0", v1, 0)

	assert.NoError(t, err)

	err = db.Append(kg, "0", v2, 0)

	assert.Error(t, err)
}

func TestTriggerNodes(t *testing.T) {
	kg := "kg1"

	err := db.CreateKeygroup(kg)

	assert.NoError(t, err)

	t1 := "t1"
	t1host := "1.1.1.1:3000"

	t2 := "t2"
	t2host := "2.2.2.2:3000"

	t3 := "t3"
	t3host := "3.3.3.3:3000"

	err = db.AddKeygroupTrigger(kg, t1, t1host)

	assert.NoError(t, err)

	err = db.AddKeygroupTrigger(kg, t1, t1host)

	assert.NoError(t, err)

	err = db.AddKeygroupTrigger(kg, t2, t2host)

	assert.NoError(t, err)

	tList, err := db.GetKeygroupTrigger(kg)

	assert.NoError(t, err)

	log.Debug().Msgf("List of keygroup triggers: %+v", tList)

	assert.Len(t, tList, 2)

	assert.Contains(t, tList, t1)

	assert.Contains(t, tList, t2)

	assert.Equal(t, t1host, tList[t1], "t1host not correct")

	assert.Equal(t, t2host, tList[t2], "t2host not correct")

	err = db.DeleteKeygroupTrigger(kg, t1)

	assert.NoError(t, err)

	tList, err = db.GetKeygroupTrigger(kg)

	assert.NoError(t, err)

	log.Debug().Msgf("List of keygroup triggers: %+v", tList)

	assert.Len(t, tList, 1)

	assert.Contains(t, tList, t2)

	assert.Equal(t, t2host, tList[t2], "t2host not correct")

	err = db.AddKeygroupTrigger(kg, t3, t3host)

	assert.NoError(t, err)

	err = db.DeleteKeygroup(kg)

	assert.NoError(t, err)

	tList, _ = db.GetKeygroupTrigger(kg)

	log.Debug().Msgf("List of keygroup triggers: %+v", tList)

	assert.Len(t, tList, 0, "got keygroup triggers for nonexistent keygroup")

}

func TestPutItemVersion(t *testing.T) {
	kg := "test-kg-item-put-version"
	id := "name"
	value := "value"

	err := db.CreateKeygroup(kg)

	assert.NoError(t, err)

	v1 := vclock.VClock{}
	v1.Set("X", 10)
	v1.Set("Y", 3)
	v1.Set("Z", 7)
	err = db.Update(kg, id, value, 0, v1)

	assert.NoError(t, err)

	retr, clock, found, err := db.Read(kg, id)

	assert.NoError(t, err)

	assert.True(t, found)
	assert.Len(t, retr, 1)
	assert.Equal(t, retr[0], value)
	assert.Len(t, retr, 1)
	assert.True(t, clock[0].Compare(v1, vclock.Equal))
}

func TestDeleteVersion(t *testing.T) {
	kg := "test-kg-item-delete-version"
	id := "name"
	id2 := "name2"
	value := "value"

	err := db.CreateKeygroup(kg)

	assert.NoError(t, err)

	v1 := vclock.VClock{}
	v1.Set("X", 10)
	v1.Set("Y", 3)
	v1.Set("Z", 7)
	err = db.Update(kg, id, value, 0, v1)

	assert.NoError(t, err)

	err = db.Delete(kg, id, v1.Copy())

	assert.NoError(t, err)

	_, _, found, err := db.Read(kg, id)

	assert.NoError(t, err)
	assert.False(t, found)

	err = db.Delete(kg, id2, v1.Copy())

	assert.NoError(t, err)

}

func TestClose(t *testing.T) {
	kg := "test-kg-item"
	id := "name"
	value := "value"

	err := db.CreateKeygroup(kg)

	assert.NoError(t, err)

	err = db.Update(kg, id, value, 0, vclock.VClock{})

	assert.NoError(t, err)

	retr, _, _, err := db.Read(kg, id)
	assert.NoError(t, err)

	assert.Len(t, retr, 1)
	assert.Equal(t, retr[0], value)

	err = db.Close()

	assert.NoError(t, err)

	_, _, _, err = db.Read(kg, id)
	assert.Error(t, err)

}
