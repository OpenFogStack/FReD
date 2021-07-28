package dynamo

import (
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/go-errors/errors"
	"github.com/rs/zerolog/log"
)

const (
	keyName     = "Key"
	valName     = "Value"
	counterName = "Counter"
	expiryKey   = "Expiry"
	sep         = "|"
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

	if len(s) == len([]string{"nil", "fred", "triggers", "kgname", "triggerID"}) {
		kg = s[3]
		tid = s[4]
	}

	return
}

// getKey returns the keygroup and id of a key.
func getKey(key string) (kg, id string) {
	s := strings.Split(key, sep)
	kg = s[0]
	if len(s) >= len([]string{"kgname", "id"}) {
		id = s[1]
	}
	return
}

func NewFromExisting(table string, svc dynamodbiface.DynamoDBAPI) (s *Storage, err error) {
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
		svc:         svc,
	}, nil
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

	return NewFromExisting(table, svc)
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

	val, ok := result.Item[valName]

	if !ok {
		return "", errors.Errorf("could not find item %s in keygroup %s", id, kg)
	}

	if e, ok := result.Item[expiryKey]; ok {
		expiry, err := strconv.Atoi(*e.N)
		if err != nil {
			return "", errors.Errorf("could not find item %s in keygroup %s", id, kg)
		}

		log.Debug().Msgf("Read found key expiring at %d, it is %d now", expiry, time.Now().Unix())

		if int64(expiry) < time.Now().Unix() {
			return "", errors.Errorf("could not find item %s in keygroup %s", id, kg)
		}
	}

	return *val.S, nil

}

func (s *Storage) ReadSome(kg, id string, count uint64) (map[string]string, error) {
	key := makeKeygroupKeyName(kg)
	start := makeKeyName(kg, id)

	filt := expression.Name(keyName).GreaterThanEqual(expression.Value(start))

	expr, err := expression.NewBuilder().WithFilter(filt).Build()

	if err != nil {
		return nil, errors.New(err)
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
		return nil, errors.New(err)
	}

	items := make(map[string]string)

	for _, i := range result.Items {

		k, ok := i[keyName]

		if !ok {
			return items, errors.Errorf("wrong format: key not in %#v", i)
		}

		if *k.S == key {
			continue
		}

		if e, ok := i[expiryKey]; ok {
			expiry, err := strconv.Atoi(*e.N)
			if err != nil {
				return items, errors.Errorf("wrong format: expiry is not a number: %s", err.Error())
			}

			log.Debug().Msgf("ReadSome found key expiring at %d, it is %d now", expiry, time.Now().Unix())

			if int64(expiry) < time.Now().Unix() {
				continue
			}
		}

		v, ok := i[valName]

		if !ok {
			return items, errors.Errorf("wrong format: value not in %#v", i)
		}

		_, keyID := getKey(*k.S)

		items[keyID] = *v.S
	}

	// sort and filter
	ids := make([]string, len(items))

	i := 0
	for k := range items {
		ids[i] = k
		i++
	}

	sort.Strings(ids)

	for i = len(items) - 1; i >= int(count); i-- {
		delete(items, ids[i])
	}

	return items, nil
}

// ReadAll returns all items in the specified keygroup.
func (s *Storage) ReadAll(kg string) (map[string]string, error) {

	key := makeKeygroupKeyName(kg)

	filt := expression.Name(keyName).BeginsWith(key)

	expr, err := expression.NewBuilder().WithFilter(filt).Build()
	if err != nil {
		return nil, errors.New(err)
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
		return nil, errors.New(err)
	}

	items := make(map[string]string)

	for _, i := range result.Items {

		k, ok := i[keyName]

		if !ok {
			return items, errors.Errorf("wrong format: key not in %#v", i)
		}

		if *k.S == key {
			continue
		}

		if e, ok := i[expiryKey]; ok {
			expiry, err := strconv.Atoi(*e.N)
			if err != nil {
				return items, errors.Errorf("wrong format: expiry is not a number: %s", err.Error())
			}

			log.Debug().Msgf("ReadAll found key expiring at %d, it is %d now", expiry, time.Now().Unix())

			if int64(expiry) < time.Now().Unix() {
				continue
			}
		}

		v, ok := i[valName]

		if !ok {
			return items, errors.Errorf("wrong format: value not in %#v", i)
		}

		_, id := getKey(*k.S)

		items[id] = *v.S

	}

	return items, nil
}

// Append appends the item to the specified keygroup by incrementing the latest key by one.
func (s *Storage) Append(kg, val string, expiry int) (string, error) {
	key := makeKeygroupKeyName(kg)

	input := &dynamodb.UpdateItemInput{
		TableName: aws.String(s.dynamotable),
		Key: map[string]*dynamodb.AttributeValue{
			keyName: {
				S: aws.String(key),
			},
		},
		ExpressionAttributeNames: map[string]*string{
			"#counter": aws.String(counterName),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":increase": {
				N: aws.String("1"),
			},
		},
		UpdateExpression: aws.String("SET #counter = #counter + :increase"),
		ReturnValues:     aws.String("UPDATED_NEW"),
	}

	out, err := s.svc.UpdateItem(input)

	if err != nil {
		log.Debug().Msgf("Append: could not update key %s: %s", key, err.Error())
		return "", errors.New(err)
	}

	c, ok := out.Attributes[counterName]

	if !ok {
		return "", errors.Errorf("could not increase counter")
	}

	newID := c.N

	in := &dynamodb.PutItemInput{
		Item: map[string]*dynamodb.AttributeValue{
			keyName: {
				S: aws.String(makeKeyName(kg, *newID)),
			},
			valName: {
				S: aws.String(val),
			},
		},
		TableName: aws.String(s.dynamotable),
	}

	if expiry > 0 {
		in.Item[expiryKey] = &dynamodb.AttributeValue{
			N: aws.String(strconv.FormatInt(time.Now().Unix()+int64(expiry), 10)),
		}
	}

	_, err = s.svc.PutItem(in)
	if err != nil {
		return "", errors.New(err)
	}

	return *newID, nil
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
		k, ok := i[keyName]

		if !ok {
			return ids, errors.Errorf("wrong format: key not in %#v", i)
		}

		if *k.S == key {
			continue
		}

		if e, ok := i[expiryKey]; ok {
			expiry, err := strconv.Atoi(*e.N)
			if err != nil {
				return ids, errors.Errorf("wrong format: expiry is not a number: %s", err.Error())
			}

			log.Debug().Msgf("IDs found key expiring at %d, it is %d now", expiry, time.Now().Unix())

			if int64(expiry) < time.Now().Unix() {
				continue
			}
		}

		_, id := getKey(*k.S)

		ids = append(ids, id)

	}

	return ids, nil
}

// Update updates the item with the specified id in the specified keygroup.
func (s *Storage) Update(kg, id, val string, _ bool, expiry int) error {

	key := makeKeyName(kg, id)

	input := &dynamodb.PutItemInput{
		Item: map[string]*dynamodb.AttributeValue{
			keyName: {
				S: aws.String(key),
			},
			valName: {
				S: aws.String(val),
			},
		},
		TableName: aws.String(s.dynamotable),
	}

	if expiry > 0 {
		input.Item[expiryKey] = &dynamodb.AttributeValue{
			N: aws.String(strconv.FormatInt(time.Now().Unix()+int64(expiry), 10)),
		}
	}

	_, err := s.svc.PutItem(input)
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

	if e, ok := result.Item[expiryKey]; ok {
		expiry, err := strconv.Atoi(*e.N)
		if err != nil {
			return false
		}

		log.Debug().Msgf("Exists found key expiring at %d, it is %d now", expiry, time.Now().Unix())

		return int64(expiry) >= time.Now().Unix()
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

	input := &dynamodb.PutItemInput{
		Item: map[string]*dynamodb.AttributeValue{
			keyName: {
				S: aws.String(key),
			},
			counterName: {
				N: aws.String("-1"),
			},
		},
		TableName: aws.String(s.dynamotable),
	}

	_, err := s.svc.PutItem(input)
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

		k, ok := i[keyName]

		if !ok {
			return errors.Errorf("wrong format")
		}

		input := &dynamodb.DeleteItemInput{

			Key: map[string]*dynamodb.AttributeValue{
				keyName: {
					S: aws.String(*k.S),
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

		k, ok := i[keyName]

		if !ok {
			return errors.Errorf("wrong format")
		}

		input := &dynamodb.DeleteItemInput{

			Key: map[string]*dynamodb.AttributeValue{
				keyName: {
					S: aws.String(*k.S),
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

	input := &dynamodb.PutItemInput{
		Item: map[string]*dynamodb.AttributeValue{
			keyName: {
				S: aws.String(key),
			},
			valName: {
				S: aws.String(host),
			},
		},
		TableName: aws.String(s.dynamotable),
	}

	_, err := s.svc.PutItem(input)
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
		k, ok := i[keyName]

		if !ok {
			return triggers, errors.Errorf("wrong format")
		}

		v, ok := i[valName]

		if !ok {
			return triggers, errors.Errorf("wrong format")
		}

		if *k.S == key {
			continue
		}

		_, id := getTriggerConfigKey(*k.S)

		triggers[id] = *v.S

	}

	return triggers, nil
}
