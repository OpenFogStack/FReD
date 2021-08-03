package dynamo

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

const table = "fred"

var db *Storage

func TestMain(m *testing.M) {
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	log.Logger = log.Output(
		zerolog.ConsoleWriter{
			Out:     os.Stderr,
			NoColor: false,
		},
	)

	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Image:        "amazon/dynamodb-local@sha256:bdd26570dc0e0ae49e1ea9d49ff662a6a1afe9121dd25793dc40d02802e7e806",
		Cmd:          []string{"-jar", "DynamoDBLocal.jar", "-inMemory"},
		ExposedPorts: []string{"8000/tcp"},
		//BindMounts:   map[string]string{"8000": "8000"},
		WaitingFor: wait.NewHostPortStrategy("8000"),
	}

	d, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})

	if err != nil {
		log.Fatal().Msg(err.Error())
		return
	}

	defer func(d testcontainers.Container, ctx context.Context) {
		err := d.Terminate(ctx)
		if err != nil {
			log.Fatal().Msgf(err.Error())
		}
	}(d, ctx)

	ip, err := d.Host(ctx)

	if err != nil {
		log.Fatal().Msg(err.Error())
		return
	}

	port, err := d.MappedPort(ctx, "8000")

	if err != nil {
		log.Fatal().Msg(err.Error())
	}

	sess := session.Must(
		session.NewSession(
			aws.NewConfig().
				WithEndpoint(fmt.Sprintf("http://%s:%s", ip, port)).
				WithRegion("eu-central-1").
				WithCredentialsChainVerboseErrors(true)))

	svc := dynamodb.New(sess)

	log.Debug().Msg("Created session - OK!")

	// aws dynamodb create-table --table-name fred --attribute-definitions "AttributeName=Key,AttributeType=S AttributeName=Value,AttributeType=S AttributeName=Expiry,AttributeType=N" --key-schema "AttributeName=Key,KeyType=HASH" --provisioned-throughput "ReadCapacityUnits=1,WriteCapacityUnits=1"
	out1, err := svc.CreateTable(&dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String(keyName),
				AttributeType: aws.String(dynamodb.ScalarAttributeTypeS),
			},
		},
		BillingMode:            nil,
		GlobalSecondaryIndexes: nil,
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String(keyName),
				KeyType:       aws.String(dynamodb.KeyTypeHash),
			},
		},
		LocalSecondaryIndexes: nil,
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(1),
			WriteCapacityUnits: aws.Int64(1),
		},
		SSESpecification:    nil,
		StreamSpecification: nil,
		TableName:           aws.String(table),
		Tags:                nil,
	})

	if err != nil {
		log.Fatal().Msg(err.Error())
		return
	}

	if out1 != nil {
		log.Debug().Msgf("Creating table, output: %#v", out1)
	}

	log.Debug().Msg("Created table - OK!")

	// aws dynamodb update-time-to-live --table-name fred --time-to-live-specification "Enabled=true, AttributeName=Expiry"
	out2, err := svc.UpdateTimeToLive(&dynamodb.UpdateTimeToLiveInput{
		TableName: aws.String(table),
		TimeToLiveSpecification: &dynamodb.TimeToLiveSpecification{
			AttributeName: aws.String(expiryKey),
			Enabled:       aws.Bool(true),
		},
	})

	if err != nil {
		log.Fatal().Msg(err.Error())
		return
	}

	if out2 != nil {
		log.Debug().Msgf("Updating TTL on table, output: %#v", out2)
	}

	log.Debug().Msg("Configured TTL - OK!")

	db, err = NewFromExisting(table, svc)

	if err != nil {
		log.Fatal().Msg(err.Error())
		return
	}

	log.Debug().Msg("Created DB - OK!")

	os.Exit(m.Run())
}

func TestCreateKeygroup(t *testing.T) {
	kg := "test-kg"
	err := db.CreateKeygroup(kg)

	assert.NoError(t, err)
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

		err = db.Update(kg, ids[i], vals[i], false, 0)

		assert.NoError(t, err)

	}

	res, err := db.ReadSome(kg, "id"+strconv.Itoa(scanStart), uint64(scanRange))

	assert.NoError(t, err)

	assert.Len(t, res, scanRange)

	for i := scanStart; i < scanStart+scanRange; i++ {
		assert.Contains(t, res, ids[i])
		assert.Equal(t, vals[i], res[ids[i]])
	}
}

func TestReadAll(t *testing.T) {
	kg := "test-read-all"
	err := db.CreateKeygroup(kg)

	assert.NoError(t, err)

	err = db.Update(kg, "id-1", "data-1", false, 0)

	assert.NoError(t, err)

	err = db.Update(kg, "id-2", "data-2", false, 0)

	assert.NoError(t, err)

	err = db.Update(kg, "id-3", "data-3", false, 0)

	assert.NoError(t, err)

	kg2 := "test-read-all-2"

	err = db.CreateKeygroup(kg2)

	assert.NoError(t, err)

	err = db.Update(kg2, "id-1", "data-1", false, 0)

	assert.NoError(t, err)

	err = db.Update(kg2, "id-2", "data-2", false, 0)

	assert.NoError(t, err)

	err = db.Update(kg2, "id-3", "data-3", false, 0)

	assert.NoError(t, err)

	res, err := db.ReadAll(kg)

	assert.NoError(t, err)

	assert.Equal(t, "data-1", res["id-1"])
	assert.Equal(t, "data-2", res["id-2"])
	assert.Equal(t, "data-3", res["id-3"])

}

func TestIDs(t *testing.T) {
	kg := "test-ids"
	err := db.CreateKeygroup(kg)

	assert.NoError(t, err)

	err = db.Update(kg, "id-1", "data-1", false, 0)

	assert.NoError(t, err)

	err = db.Update(kg, "id-2", "data-2", false, 0)

	assert.NoError(t, err)

	err = db.Update(kg, "id-3", "data-3", false, 0)

	assert.NoError(t, err)

	kg2 := "test-read-all-2"

	err = db.CreateKeygroup(kg2)

	assert.NoError(t, err)

	err = db.Update(kg2, "id-1", "data-1", false, 0)

	assert.NoError(t, err)

	err = db.Update(kg2, "id-2", "data-2", false, 0)

	assert.NoError(t, err)

	err = db.Update(kg2, "id-3", "data-3", false, 0)

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

	err = db.Update(kg, id, value, false, 0)

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

	err = db.Update(kg, id, value, false, 0)

	assert.NoError(t, err)

	retr, err := db.Read(kg, id)

	assert.NoError(t, err)

	assert.Equal(t, value, retr)
}

func TestItemDelete(t *testing.T) {
	kg := "test-kg-item-delete"
	id := "name"
	id2 := "name2"
	value := "value"

	err := db.CreateKeygroup(kg)

	assert.NoError(t, err)

	err = db.Update(kg, id, value, false, 0)

	assert.NoError(t, err)
	retr, err := db.Read(kg, id)
	assert.NoError(t, err)

	assert.Equal(t, value, retr)

	err = db.Delete(kg, id)

	assert.NoError(t, err)

	_, err = db.Read(kg, id)

	assert.Error(t, err, "read a deleted item")

	err = db.Delete(kg, id2)

	assert.NoError(t, err, "deleting non-existent keys should be allowed")

}

func TestItemAfterDeleteKeygroup(t *testing.T) {
	kg := "test-kg-item-delete"
	id := "ndel"
	value := "vdel"

	err := db.CreateKeygroup(kg)

	assert.NoError(t, err)

	err = db.Update(kg, id, value, false, 0)

	assert.NoError(t, err)

	err = db.DeleteKeygroup(kg)

	assert.NoError(t, err)

	_, err = db.Read(kg, id)

	assert.Error(t, err)
}

func TestExpiry(t *testing.T) {
	kg := "test-kg-item"
	id := "name"
	value := "value"

	err := db.CreateKeygroup(kg)

	assert.NoError(t, err)

	err = db.Update(kg, id, value, false, 10)

	assert.NoError(t, err)

	retr, err := db.Read(kg, id)

	assert.NoError(t, err)

	assert.Equal(t, value, retr)

	exists := db.Exists(kg, id)
	assert.True(t, exists)

	time.Sleep(11 * time.Second)

	_, err = db.Read(kg, id)

	assert.Error(t, err)

	exists = db.Exists(kg, id)
	assert.False(t, exists)
}

func TestAppend(t *testing.T) {
	kg := "log"
	v1 := "value-1"
	v2 := "value-2"

	err := db.CreateKeygroup(kg)

	assert.NoError(t, err)

	key1, err := db.Append(kg, v1, 0)

	assert.NoError(t, err)

	assert.Equal(t, "0", key1)

	key2, err := db.Append(kg, v2, 0)

	assert.NoError(t, err)

	assert.Equal(t, "1", key2)

	for i := 2; i < 100; i++ {
		v := "value-" + strconv.Itoa(i)
		key, err := db.Append(kg, v, 0)

		assert.NoError(t, err)

		assert.Equal(t, strconv.Itoa(i), key)
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
				key, err := db.Append(kg, v, 0)

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

	log.Debug().Msgf("List of keygroup triggers: %#v", tList)

	assert.Len(t, tList, 2)
	assert.Contains(t, tList, t1)
	assert.Equal(t, tList[t1], t1host, "t1host not correct")
	assert.Contains(t, tList, t2)
	assert.Equal(t, tList[t2], t2host, "t2host not correct")

	err = db.DeleteKeygroupTrigger(kg, t1)

	assert.NoError(t, err)

	tList, err = db.GetKeygroupTrigger(kg)

	assert.NoError(t, err)

	log.Debug().Msgf("List of keygroup triggers: %#v", tList)

	assert.Len(t, tList, 1)
	assert.Contains(t, tList, t2)
	assert.Equal(t, tList[t2], t2host, "t2host not correct")

	err = db.AddKeygroupTrigger(kg, t3, t3host)

	assert.NoError(t, err)

	err = db.DeleteKeygroup(kg)

	assert.NoError(t, err)

	tList, _ = db.GetKeygroupTrigger(kg)

	log.Debug().Msgf("List of keygroup triggers: %#v", tList)

	assert.Len(t, tList, 0, "got keygroup triggers for nonexistent keygroup")

}

func TestClose(t *testing.T) {
	kg := "test-kg-item"
	id := "name"
	value := "value"

	err := db.CreateKeygroup(kg)

	assert.NoError(t, err)
	err = db.Update(kg, id, value, false, 0)

	assert.NoError(t, err)

	retr, err := db.Read(kg, id)

	assert.NoError(t, err)

	assert.Equal(t, value, retr)

	err = db.Close()

	assert.NoError(t, err)
	// currently close is not implemented in DynamoDB
}
