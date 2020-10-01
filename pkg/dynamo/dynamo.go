package dynamo

import (
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/rs/zerolog/log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/go-errors/errors"
)

const (
	keyName = "Key"
	sep     = "|"
)

// Storage is a struct that saves all necessary information to access the database, in this case the session for DynamoDB and the table name.
type Storage struct {
	dynamotable string
	svc         dynamodbiface.DynamoDBAPI
}

// makeKeyName creates the internal DynamoDB key given a keygroup name and an id.
func makeKeyName(kgname string, id string) string {
	return kgname + sep + id
}

// makeKeygroupKeyName creates the internal DynamoDB key given a keygroup name.
func makeKeygroupKeyName(kgname string) string {
	return kgname + sep
}

func makeTriggerConfigKeyName(kgname string, tid string) string {
	return sep + "fred" + sep + "triggers" + sep + kgname + sep + tid
}

// getTriggerConfigKey returns the keygroup and id of a key.
func getTriggerConfigKey(key string) (kg, tid string) {
	s := strings.Split(key, sep)

	if len(s) == 5 {
		kg = s[3]
		tid = s[4]
	}

	return
}

// getKey returns the keygroup and id of a key.
func getKey(key string) (kg, id string) {
	s := strings.Split(key, sep)
	kg = s[0]
	if len(s) > 1 {
		id = s[1]
	}
	return
}

// New creates a new Session for DynamoDB.
func New(table, region string) (s *Storage, err error) {
	log.Debug().Msgf("creating a new dynamodb connection to table %s in region %s", table, region)

	log.Debug().Msg("Checked creds - OK!")

	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(region),
	}))

	log.Debug().Msg("Created session - OK!")

	svc := dynamodb.New(sess)

	log.Debug().Msg("Created service - OK!")

	log.Debug().Msgf("Loading table %s...", table)

	// check if the table with that name even exists
	// if not: error out
	desc, err := svc.DescribeTable(&dynamodb.DescribeTableInput{
		TableName: aws.String(table),
	})

	if err != nil {
		return nil, errors.New(err)
	}

	log.Debug().Msgf("Checking table %s...", table)

	// check that the table has the correct fields (i.e. a primary hash key with name "key") for our use
	if len(desc.Table.KeySchema) != 1 {
		return nil, errors.Errorf("expected a single primary range key with name \"%s\" but go %d keys", keyName, len(desc.Table.KeySchema))
	}

	log.Debug().Msg("Checked table fields - OK!")

	if *(desc.Table.KeySchema[0].AttributeName) != keyName && *(desc.Table.KeySchema[0].KeyType) != dynamodb.KeyTypeHash {
		return nil, errors.Errorf("expected the primary key to be named \"%s\" and be of type range but got %s with type %s", keyName, *(desc.Table.KeySchema[0].AttributeName), *(desc.Table.KeySchema[0].KeyType))
	}

	log.Debug().Msg("Checked table keys - OK!")

	return &Storage{
		dynamotable: table,
		svc:         dynamodbiface.DynamoDBAPI(svc),
	}, nil
}

// Close closes the underlying DynamoDB connection (no cleanup needed at the moment).
func (s *Storage) Close() error {
	return nil
}

// Read returns an item with the specified id from the specified keygroup.
func (s *Storage) Read(kg string, id string) (string, error) {

	key := makeKeyName(kg, id)

	result, err := s.svc.GetItem(&dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			keyName: {
				S: aws.String(key),
			},
		},
		TableName: &s.dynamotable,
	})

	if err != nil {
		return "", errors.New(err)
	}

	if result.Item == nil {
		return "", errors.Errorf("could not find item %s in keygroup %s", id, kg)
	}

	Item := struct {
		Key   string
		Value string
	}{}

	err = dynamodbattribute.UnmarshalMap(result.Item, &Item)
	if err != nil {
		return "", errors.New(err)
	}

	return Item.Value, nil

}

// ReadAll returns all items in the specified keygroup.
func (s *Storage) ReadAll(kg string) (map[string]string, error) {
	items := make(map[string]string)

	key := makeKeygroupKeyName(kg)

	filt := expression.Name(keyName).BeginsWith(key)

	expr, err := expression.NewBuilder().WithFilter(filt).Build()
	if err != nil {
		return items, errors.New(err)
	}

	params := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		ProjectionExpression:      expr.Projection(),
		TableName:                 aws.String(s.dynamotable),
	}

	// Make the DynamoDB Query API call
	result, err := s.svc.Scan(params)
	if err != nil {
		return items, errors.New(err)
	}

	for _, i := range result.Items {

		item := struct {
			Key   string
			Value string
		}{}

		err = dynamodbattribute.UnmarshalMap(i, &item)

		if err != nil {
			return items, errors.New(err)
		}

		if item.Key == key {
			continue
		}

		_, id := getKey(item.Key)

		items[id] = item.Value

	}

	return items, nil
}

// Append appends the item to the specified keygroup by incrementing the latest key by one.
func (s *Storage) Append(kg, val string, expiry int) (string, error) {
	// first, get the latest key
	// maximum of 18446744073709551615, though!
	// if you reach this maximum, please send me a letter
	var newest uint64

	key := makeKeygroupKeyName(kg)

	filt := expression.Name(keyName).BeginsWith(key)

	expr, err := expression.NewBuilder().WithFilter(filt).Build()
	if err != nil {
		return "", errors.New(err)
	}

	params := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		ProjectionExpression:      expr.Projection(),
		TableName:                 aws.String(s.dynamotable),
	}

	// Make the DynamoDB Query API call
	result, err := s.svc.Scan(params)
	if err != nil {
		return "", errors.New(err)
	}

	if len(result.Items) == 0 {
		// nothing in existence yet, return MaxUInt (so we can increment it to 0)
		newest = ^uint64(0)
	} else {
		for _, i := range result.Items {

			item := struct {
				Key string
			}{}

			err = dynamodbattribute.UnmarshalMap(i, &item)

			if err != nil {
				return "", errors.New(err)
			}

			if item.Key == key {
				continue
			}

			_, id := getKey(item.Key)

			parsed, err := strconv.ParseUint(id, 10, 64)

			if err != nil {
				return "", errors.New(err)
			}

			if parsed > newest {
				newest = parsed
			}

		}
	}

	// increment by one
	// conveniently, if we reach MaxUint64, we can still increment by 1 to get back to 0
	id := strconv.FormatUint(newest+1, 10)

	Item := struct {
		Key    string
		Value  string
		Expiry int64
	}{
		Key:    id,
		Value:  val,
		Expiry: time.Now().Unix() + int64(expiry),
	}

	av, err := dynamodbattribute.MarshalMap(Item)

	if err != nil {
		return "", errors.New(err)
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(s.dynamotable),
	}

	_, err = s.svc.PutItem(input)
	if err != nil {
		return "", errors.New(err)
	}

	return id, nil
}

// IDs returns the keys of all items in the specified keygroup.
func (s *Storage) IDs(kg string) ([]string, error) {
	var ids []string

	key := makeKeygroupKeyName(kg)

	filt := expression.Name(keyName).BeginsWith(key)

	expr, err := expression.NewBuilder().WithFilter(filt).Build()
	if err != nil {
		return ids, errors.New(err)
	}

	params := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		ProjectionExpression:      expr.Projection(),
		TableName:                 aws.String(s.dynamotable),
	}

	// Make the DynamoDB Query API call
	result, err := s.svc.Scan(params)
	if err != nil {
		return ids, errors.New(err)
	}

	for _, i := range result.Items {

		item := struct {
			Key string
		}{}

		err = dynamodbattribute.UnmarshalMap(i, &item)

		if err != nil {
			return ids, errors.New(err)
		}

		if item.Key == key {
			continue
		}

		_, id := getKey(item.Key)

		ids = append(ids, id)

	}

	return ids, nil
}

// Update updates the item with the specified id in the specified keygroup.
func (s *Storage) Update(kg, id, val string, expiry int) error {

	key := makeKeyName(kg, id)

	Item := struct {
		Key    string
		Value  string
		Expiry int64
	}{
		Key:    key,
		Value:  val,
		Expiry: time.Now().Unix() + int64(expiry),
	}

	av, err := dynamodbattribute.MarshalMap(Item)

	if err != nil {
		return errors.New(err)
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(s.dynamotable),
	}

	_, err = s.svc.PutItem(input)
	if err != nil {
		return errors.New(err)
	}

	return nil
}

// Delete deletes the item with the specified id from the specified keygroup.
func (s *Storage) Delete(kg string, id string) error {
	key := makeKeyName(kg, id)

	input := &dynamodb.DeleteItemInput{
		TableName: aws.String(s.dynamotable),
		Key: map[string]*dynamodb.AttributeValue{
			keyName: {
				S: aws.String(key),
			},
		},
	}

	_, err := s.svc.DeleteItem(input)
	if err != nil {
		return errors.New(err)
	}

	return nil
}

// Exists checks if the given data item exists in the dynamodb database.
func (s *Storage) Exists(kg string, id string) bool {
	key := makeKeyName(kg, id)

	result, err := s.svc.GetItem(&dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			keyName: {
				S: aws.String(key),
			},
		},
		TableName: &s.dynamotable,
	})

	if err != nil {
		return false
	}

	if result.Item == nil {
		return false
	}

	return true
}

// ExistsKeygroup checks if the given keygroup exists in the DynamoDB database.
func (s *Storage) ExistsKeygroup(kg string) bool {
	key := makeKeygroupKeyName(kg)

	result, err := s.svc.GetItem(&dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			keyName: {
				S: aws.String(key),
			},
		},
		TableName: &s.dynamotable,
	})

	if err != nil {
		return false
	}

	if result.Item == nil {
		return false
	}

	return true
}

// CreateKeygroup creates the given keygroup in the DynamoDB database.
func (s *Storage) CreateKeygroup(kg string) error {
	key := makeKeygroupKeyName(kg)

	Item := struct {
		Key   string
		Value string
	}{
		Key:   key,
		Value: key,
	}

	av, err := dynamodbattribute.MarshalMap(Item)

	if err != nil {
		return errors.New(err)
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(s.dynamotable),
	}

	_, err = s.svc.PutItem(input)
	if err != nil {
		return errors.New(err)
	}

	return nil
}

// DeleteKeygroup deletes the given keygroup from the DynamoDB database.
func (s *Storage) DeleteKeygroup(kg string) error {

	// delete all entries for that keygroup
	key := makeKeygroupKeyName(kg)

	filt := expression.Name(keyName).BeginsWith(key)

	expr, err := expression.NewBuilder().WithFilter(filt).Build()
	if err != nil {
		return errors.New(err)
	}

	params := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		ProjectionExpression:      expr.Projection(),
		TableName:                 aws.String(s.dynamotable),
	}

	// Make the DynamoDB Query API call
	result, err := s.svc.Scan(params)
	if err != nil {
		return errors.New(err)
	}

	for _, i := range result.Items {

		item := struct {
			Key string
		}{}

		err = dynamodbattribute.UnmarshalMap(i, &item)

		if err != nil {
			return errors.New(err)
		}

		input := &dynamodb.DeleteItemInput{

			Key: map[string]*dynamodb.AttributeValue{
				keyName: {
					S: aws.String(item.Key),
				},
			},
			TableName: aws.String(s.dynamotable),
		}

		_, err := s.svc.DeleteItem(input)
		if err != nil {
			return errors.New(err)
		}
	}

	// delete the keygroup triggers

	key = makeTriggerConfigKeyName(kg, "")

	filt = expression.Name(keyName).BeginsWith(key)

	expr, err = expression.NewBuilder().WithFilter(filt).Build()

	if err != nil {
		return errors.New(err)
	}

	params = &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		ProjectionExpression:      expr.Projection(),
		TableName:                 aws.String(s.dynamotable),
	}

	// Make the DynamoDB Query API call
	result, err = s.svc.Scan(params)
	if err != nil {
		return errors.New(err)
	}

	for _, i := range result.Items {

		item := struct {
			Key string
		}{}

		err = dynamodbattribute.UnmarshalMap(i, &item)

		if err != nil {
			return errors.New(err)
		}

		input := &dynamodb.DeleteItemInput{

			Key: map[string]*dynamodb.AttributeValue{
				keyName: {
					S: aws.String(item.Key),
				},
			},
			TableName: aws.String(s.dynamotable),
		}

		_, err := s.svc.DeleteItem(input)
		if err != nil {
			return errors.New(err)
		}
	}

	return nil
}

// AddKeygroupTrigger stores a trigger node in the dynamodb database.
func (s *Storage) AddKeygroupTrigger(kg string, id string, host string) error {
	key := makeTriggerConfigKeyName(kg, id)

	Item := struct {
		Key   string
		Value string
	}{
		Key:   key,
		Value: host,
	}

	av, err := dynamodbattribute.MarshalMap(Item)

	if err != nil {
		return errors.New(err)
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(s.dynamotable),
	}

	_, err = s.svc.PutItem(input)
	if err != nil {
		return errors.New(err)
	}

	return nil
}

// DeleteKeygroupTrigger removes a trigger node from the dynamodb database.
func (s *Storage) DeleteKeygroupTrigger(kg string, id string) error {
	key := makeTriggerConfigKeyName(kg, id)

	input := &dynamodb.DeleteItemInput{
		TableName: aws.String(s.dynamotable),
		Key: map[string]*dynamodb.AttributeValue{
			keyName: {
				S: aws.String(key),
			},
		},
	}

	_, err := s.svc.DeleteItem(input)
	if err != nil {
		return errors.New(err)
	}

	return nil
}

// GetKeygroupTrigger returns a map of all trigger nodes from the dynamodb database.
func (s *Storage) GetKeygroupTrigger(kg string) (map[string]string, error) {
	triggers := make(map[string]string)

	key := makeTriggerConfigKeyName(kg, "")

	filt := expression.Name(keyName).BeginsWith(key)

	expr, err := expression.NewBuilder().WithFilter(filt).Build()
	if err != nil {
		return triggers, errors.New(err)
	}

	params := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		ProjectionExpression:      expr.Projection(),
		TableName:                 aws.String(s.dynamotable),
	}

	// Make the DynamoDB Query API call
	result, err := s.svc.Scan(params)
	if err != nil {
		return triggers, errors.New(err)
	}

	for _, i := range result.Items {

		item := struct {
			Key   string
			Value string
		}{}

		err = dynamodbattribute.UnmarshalMap(i, &item)

		if err != nil {
			return triggers, errors.New(err)
		}

		if item.Key == key {
			continue
		}

		_, id := getTriggerConfigKey(item.Key)

		triggers[id] = item.Value

	}

	return triggers, nil
}
