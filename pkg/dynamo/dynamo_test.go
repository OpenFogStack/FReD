package dynamo

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/DistributedClocks/GoVector/govec/vclock"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	dynamoDBTypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
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

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion("eu-central-1"),
	)

	if err != nil {
		log.Fatal().Msg(err.Error())
		return
	}

	cfg.Credentials = credentials.NewStaticCredentialsProvider("TEST_KEY", "TEST_SECRET", "")

	svc := dynamodb.New(dynamodb.Options{
		Region:           cfg.Region,
		HTTPClient:       cfg.HTTPClient,
		Credentials:      cfg.Credentials,
		APIOptions:       cfg.APIOptions,
		Logger:           cfg.Logger,
		ClientLogMode:    cfg.ClientLogMode,
		EndpointResolver: dynamodb.EndpointResolverFromURL(fmt.Sprintf("http://%s:%s", ip, port)),
	})

	log.Debug().Msg("Created session - OK!")

	// aws dynamodb create-table --table-name fred --attribute-definitions "AttributeName=Keygroup,AttributeType=S AttributeName=Key,AttributeType=S" --key-schema "AttributeName=Keygroup,KeyType=HASH AttributeName=Key,KeyType=RANGE" --provisioned-throughput "ReadCapacityUnits=1,WriteCapacityUnits=1"
	out1, err := svc.CreateTable(context.TODO(), &dynamodb.CreateTableInput{
		AttributeDefinitions: []dynamoDBTypes.AttributeDefinition{
			{
				AttributeName: aws.String(keygroupName),
				AttributeType: dynamoDBTypes.ScalarAttributeTypeS,
			},
			{
				AttributeName: aws.String(keyName),
				AttributeType: dynamoDBTypes.ScalarAttributeTypeS,
			},
		},
		GlobalSecondaryIndexes: nil,
		KeySchema: []dynamoDBTypes.KeySchemaElement{
			{
				AttributeName: aws.String(keygroupName),
				KeyType:       dynamoDBTypes.KeyTypeHash,
			},
			{
				AttributeName: aws.String(keyName),
				KeyType:       dynamoDBTypes.KeyTypeRange,
			},
		},
		ProvisionedThroughput: &dynamoDBTypes.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(1),
			WriteCapacityUnits: aws.Int64(1),
		},
		TableName: aws.String(table),
	})

	if err != nil {
		log.Fatal().Msg(err.Error())
		return
	}

	if out1 != nil {
		log.Debug().Msgf("Creating table, output: %+v", out1)
	}

	log.Debug().Msg("Created table - OK!")

	// aws dynamodb update-time-to-live --table-name fred --time-to-live-specification "Enabled=true, AttributeName=Expiry"
	out2, err := svc.UpdateTimeToLive(context.TODO(), &dynamodb.UpdateTimeToLiveInput{
		TableName: aws.String(table),
		TimeToLiveSpecification: &dynamoDBTypes.TimeToLiveSpecification{
			AttributeName: aws.String(expiryKey),
			Enabled:       aws.Bool(true),
		},
	})

	if err != nil {
		log.Fatal().Msg(err.Error())
		return
	}

	if out2 != nil {
		log.Debug().Msgf("Updating TTL on table, output: %+v", out2)
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

	data, _, err := db.ReadSome(kg, "id"+strconv.Itoa(scanStart), uint64(scanRange))

	assert.NoError(t, err)

	assert.Len(t, data, scanRange)

	for i := scanStart; i < scanStart+scanRange; i++ {
		assert.Contains(t, data, ids[i])
		assert.Len(t, data[ids[i]], 1)
		assert.Equal(t, data[ids[i]][0], vals[i])
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

	data, _, err := db.ReadAll(kg)

	assert.NoError(t, err)

	assert.Len(t, data, 3)
	assert.Len(t, data["id-1"], 1)
	assert.Equal(t, "data-1", data["id-1"][0])
	assert.Len(t, data["id-2"], 1)
	assert.Equal(t, "data-2", data["id-2"][0])
	assert.Len(t, data["id-3"], 1)
	assert.Equal(t, "data-3", data["id-3"][0])

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

	tList, err := db.GetKeygroupTrigger(kg)

	assert.NoError(t, err)

	log.Debug().Msgf("List of keygroup triggers: %+v", tList)

	assert.Len(t, tList, 0)

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

	tList, err = db.GetKeygroupTrigger(kg)

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
	kg := "test-kg-close"
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
	// currently close is not implemented in DynamoDB
}
