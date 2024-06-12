package dynamo

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	"git.tu-berlin.de/mcc-fred/fred/pkg/vector"
	"git.tu-berlin.de/mcc-fred/vclock"
	"github.com/aws/aws-sdk-go-v2/aws"
	shttp "github.com/aws/aws-sdk-go-v2/aws/transport/http"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	dynamoDBTypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/smithy-go"
	"github.com/go-errors/errors"
	"github.com/rs/zerolog/log"
)

const (
	keygroupName = "Keygroup"
	keyName      = "Key"
	valName      = "Value"
	triggerName  = "Trigger"
	expiryKey    = "Expiry"
	NULLValue    = "%NULL%"
)

// Storage is a struct that saves all necessary information to access the database, in this case the session for
// DynamoDB and the table name. The DynamoDB table is set up with the following attributes:
// Keygroup (S) | Key (S) | Value (Document) | Expiry (N) | Trigger (Document)
// where "Keygroup" is the partition key and "Key" is the sort key (both together form the primary key).
// Set this up with the following aws-cli command (table name in this case is "fred"):
//
//	aws dynamodb create-table --table-name fred \
//		--attribute-definitions "AttributeName=Keygroup,AttributeType=S AttributeName=Key,AttributeType=S" \
//		--key-schema "AttributeName=Keygroup,KeyType=HASH AttributeName=Key,KeyType=RANGE" \
//		--provisioned-throughput "ReadCapacityUnits=1,WriteCapacityUnits=1"
//
// The "Expiry" attribute is used to expire data items automatically in DynamoDB. Set this up with this command:
//
//	aws dynamodb update-time-to-live --table-name fred \
//		--time-to-live-specification "Enabled=true, AttributeName=Expiry"
//
// Two types of items are stored here:
//   - Keygroup configuration is stored with the NULL "Key" and the keygroup name: this has  the "Trigger" attribute that
//     stores a map of trigger nodes for that keygroup
//   - Keys are stored with a "Keygroup" and unique "Key", where the Value is a list of version vectors and values - the
//     additional "Expiry" attribute can be set to let the keys expire, and it is updated with each update to the data item
//     (note that this means that in DynamoDB, all versions of an item expire at the same time, not necessarily in the
//     order in which they appeared)
type Storage struct {
	dynamotable string
	svc         *dynamodb.Client
}

func vectorFromString(s string) (vclock.VClock, error) {
	raw, err := url.QueryUnescape(s)

	if err != nil {
		return nil, err
	}

	var b map[string]uint64
	err = json.Unmarshal([]byte(raw), &b)

	if err != nil {
		return nil, err
	}

	return b, nil
}

func vectorToString(v vclock.VClock) string {
	return url.QueryEscape(vector.SortedVCString(v))
}

func New(table string, region string, endpoint string, create bool) (*Storage, error) {
	log.Trace().Msgf("creating a new dynamodb connection to table %s in region %s", table, region)

	log.Trace().Msg("Checked creds - OK!")

	opts := dynamodb.Options{
		Region: region,
	}

	if endpoint != "" {
		log.Debug().Msgf("creating a new local dynamodb connection (endpoint: http://%s) to table %s in region %s", endpoint, table, region)

		creds := credentials.NewStaticCredentialsProvider("TEST_KEY", "TEST_SECRET", "")

		opts = dynamodb.Options{
			Region:           region,
			Credentials:      creds,
			EndpointResolver: dynamodb.EndpointResolverFromURL(fmt.Sprintf("http://%s", endpoint)),
			// this is a custom endpoint, we might as well ddos it
			RetryMaxAttempts: 100,
		}
	}

	log.Trace().Msg("Created session - OK!")

	svc := dynamodb.New(opts)

	log.Trace().Msg("Created service - OK!")

	if create {
		// aws dynamodb create-table --table-name fred --attribute-definitions "AttributeName=Keygroup,AttributeType=S AttributeName=Key,AttributeType=S" --key-schema "AttributeName=Keygroup,KeyType=HASH AttributeName=Key,KeyType=RANGE" --provisioned-throughput "ReadCapacityUnits=1,WriteCapacityUnits=1"
		_, err := svc.CreateTable(context.TODO(), &dynamodb.CreateTableInput{
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
			return nil, err
		}

		log.Trace().Msg("DynamoDB: Created table - OK!")

		// aws dynamodb update-time-to-live --table-name fred --time-to-live-specification "Enabled=true, AttributeName=Expiry"
		_, err = svc.UpdateTimeToLive(context.TODO(), &dynamodb.UpdateTimeToLiveInput{
			TableName: aws.String(table),
			TimeToLiveSpecification: &dynamoDBTypes.TimeToLiveSpecification{
				AttributeName: aws.String(expiryKey),
				Enabled:       aws.Bool(true),
			},
		})

		if err != nil {
			return nil, err
		}

		log.Trace().Msg("DynamoDB: Configured TTL - OK!")

		return &Storage{
			dynamotable: table,
			svc:         svc,
		}, nil
	}

	log.Trace().Msgf("Loading table %s...", table)
	// check if the table with that name even exists
	// if not: error out
	desc, err := svc.DescribeTable(context.TODO(), &dynamodb.DescribeTableInput{
		TableName: aws.String(table),
	})

	if err != nil {
		log.Error().Msg(errors.New(err).ErrorStack())
		return nil, errors.New(err)
	}

	log.Trace().Msgf("Checking table %s...", table)

	// check that the table has the correct fields (i.e. a primary hash key with name "Keygroup" and secondary range key
	// "Key") for our use
	if len(desc.Table.KeySchema) != 2 {
		return nil, errors.Errorf("expected a composite primary key with hash key \"%s\" and range key \"%s\" but got %d keys", keygroupName, keyName, len(desc.Table.KeySchema))
	}

	log.Trace().Msg("Checked table fields - OK!")

	if *(desc.Table.KeySchema[0].AttributeName) != keygroupName || desc.Table.KeySchema[0].KeyType != dynamoDBTypes.KeyTypeHash {
		return nil, errors.Errorf("expected the first part of primary key to be named \"%s\" and be of type hash but got %s with type %s", keygroupName, *(desc.Table.KeySchema[0].AttributeName), desc.Table.KeySchema[0].KeyType)
	}

	if *(desc.Table.KeySchema[1].AttributeName) != keyName || desc.Table.KeySchema[1].KeyType != dynamoDBTypes.KeyTypeRange {
		return nil, errors.Errorf("expected the second part of primary key to be named \"%s\" and be of type range but got %s with type %s", keyName, *(desc.Table.KeySchema[1].AttributeName), desc.Table.KeySchema[1].KeyType)
	}

	log.Trace().Msg("Checked table keys - OK!")

	return &Storage{
		dynamotable: table,
		svc:         svc,
	}, nil
}

// Close closes the underlying DynamoDB connection (no cleanup needed at the moment).
func (s *Storage) Close() error {
	return nil
}

// Read returns an item with the specified id from the specified keygroup.
func (s *Storage) Read(kg string, id string) ([]string, []vclock.VClock, bool, error) {
	log.Trace().Msgf("Reading key %s from keygroup %s", id, kg)

	// To read, we need to get the item with the "Keygroup" kg and "Key" id and convert the returned "Value" to a list
	// of strings and vclocks.
	proj := expression.NamesList(expression.Name(keyName), expression.Name(valName), expression.Name(expiryKey))

	expr, err := expression.NewBuilder().WithProjection(proj).Build()

	if err != nil {
		log.Error().Msg(errors.New(err).ErrorStack())
		return nil, nil, false, errors.New(err)
	}

	result, err := s.svc.GetItem(context.TODO(), &dynamodb.GetItemInput{
		Key: map[string]dynamoDBTypes.AttributeValue{
			keygroupName: &dynamoDBTypes.AttributeValueMemberS{
				Value: kg,
			},
			keyName: &dynamoDBTypes.AttributeValueMemberS{
				Value: id,
			},
		},
		TableName:                &s.dynamotable,
		ExpressionAttributeNames: expr.Names(),
		ProjectionExpression:     expr.Projection(),
	})

	if err != nil {
		log.Error().Msg(errors.New(err).ErrorStack())
		return nil, nil, false, errors.New(err)
	}

	// check if the item was found at all
	if result.Item == nil {
		return nil, nil, false, nil
	}

	// let's check the value of the item we got back
	val, ok := result.Item[valName]

	if !ok {
		return nil, nil, false, nil
	}

	// check that the item hasn't actually expired but wasn't cleaned up
	if e, ok := result.Item[expiryKey]; ok {
		expiration, ok := e.(*dynamoDBTypes.AttributeValueMemberN)

		if ok {
			expiry, err := strconv.Atoi(expiration.Value)
			if err != nil {
				log.Error().Msg(errors.New(err).ErrorStack())
				return nil, nil, false, errors.New(err)
			}

			log.Trace().Msgf("Read found key expiring at %d, it is %d now", expiry, time.Now().UTC().Unix())

			// oops, item has expired - we treat this as "not found"
			if int64(expiry) < time.Now().UTC().Unix() {
				return nil, nil, false, nil
			}
		}
	}

	// ok, now we have the item in "val"
	// since "Value" is a map of vector clocks to data, we convert that now
	values, ok := val.(*dynamoDBTypes.AttributeValueMemberM)

	if !ok || values == nil || len(values.Value) == 0 {
		return nil, nil, false, nil
	}

	items := make([]string, 0, len(values.Value))
	vvectors := make([]vclock.VClock, 0, len(values.Value))

	for v, data := range values.Value {
		version, err := vectorFromString(v)

		if err != nil {
			log.Error().Msg(errors.New(err).ErrorStack())
			return nil, nil, false, errors.New(err)
		}

		i, ok := data.(*dynamoDBTypes.AttributeValueMemberS)

		if !ok {
			return nil, nil, false, errors.Errorf("Read: can't parse item value")
		}

		vvectors = append(vvectors, version)

		items = append(items, i.Value)
	}

	return items, vvectors, true, nil

}

func (s *Storage) ReadSome(kg string, id string, count uint64) ([]string, []string, []vclock.VClock, error) {
	log.Trace().Msgf("Reading %d keys from keygroup %s starting at %s", count, kg, id)

	// in this case we need to get all items with "Keygroup" kg and then sort them
	filt := expression.Name(keygroupName).Equal(expression.Value(kg)).And(expression.Name(keyName).GreaterThanEqual(expression.Value(id)))
	proj := expression.NamesList(expression.Name(keyName), expression.Name(valName), expression.Name(expiryKey))

	expr, err := expression.NewBuilder().WithFilter(filt).WithProjection(proj).Build()

	if err != nil {
		log.Error().Msg(errors.New(err).ErrorStack())
		return nil, nil, nil, errors.New(err)
	}

	params := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		ProjectionExpression:      expr.Projection(),
		TableName:                 aws.String(s.dynamotable),
	}

	// Make the DynamoDB Query API call
	result, err := s.svc.Scan(context.TODO(), params)
	if err != nil {
		log.Error().Msg(errors.New(err).ErrorStack())
		// continue: propably a projection issue
		//return nil, nil, errors.New(err)
	}

	type item struct {
		key     string
		val     string
		version vclock.VClock
	}

	items := make([]item, 0, len(result.Items))

	for _, i := range result.Items {
		key, ok := i[keyName]

		if !ok {
			return nil, nil, nil, errors.Errorf("ReadSome: internal error, can't find key")
		}

		k, ok := key.(*dynamoDBTypes.AttributeValueMemberS)

		if !ok {
			return nil, nil, nil, errors.Errorf("ReadSome: malformed key")
		}

		if k.Value == NULLValue {
			continue
		}

		val, ok := i[valName]

		if !ok {
			return nil, nil, nil, errors.Errorf("ReadSome: malformed value")
		}

		if e, ok := i[expiryKey]; ok {
			expiration, ok := e.(*dynamoDBTypes.AttributeValueMemberN)

			if ok {
				expiry, err := strconv.Atoi(expiration.Value)
				if err != nil {
					log.Error().Msg(errors.New(err).ErrorStack())
					return nil, nil, nil, errors.New(err)
				}

				log.Trace().Msgf("ReadSome found key expiring at %d, it is %d now", expiry, time.Now().UTC().Unix())

				// oops, item has expired - we treat this as "not found"
				if int64(expiry) < time.Now().UTC().Unix() {
					continue
				}
			}
		}

		values, ok := val.(*dynamoDBTypes.AttributeValueMemberM)

		if !ok || values == nil || len(values.Value) == 0 {
			continue
		}

		for v, data := range values.Value {
			version, err := vectorFromString(v)

			if err != nil {
				log.Error().Msg(errors.New(err).ErrorStack())
				return nil, nil, nil, errors.New(err)
			}

			it, ok := data.(*dynamoDBTypes.AttributeValueMemberS)

			if !ok {
				return nil, nil, nil, errors.Errorf("ReadSome: malformed item")
			}

			items = append(items, item{
				k.Value,
				it.Value,
				version,
			})

		}
	}

	// now we have a list of items, sort them and convert to lists
	keys := make([]string, 0, len(items))
	values := make([]string, 0, len(items))
	versions := make([]vclock.VClock, 0, len(items))

	// now we have lists of keys and items, we need to sort them by the key attribute alphabetically
	sort.Slice(items, func(i, j int) bool {
		return strings.ToLower(items[i].key) < strings.ToLower(items[j].key)
	})

	i := 0
	curr := ""
	for _, x := range items {
		if x.key != curr {
			curr = x.key
			i++
			if i > int(count) {
				break
			}
		}

		keys = append(keys, x.key)
		values = append(values, x.val)
		versions = append(versions, x.version)
	}

	return keys, values, versions, nil
}

// ReadAll returns all items in the specified keygroup.
func (s *Storage) ReadAll(kg string) ([]string, []string, []vclock.VClock, error) {
	log.Trace().Msgf("Reading all keys from keygroup %s", kg)

	// in this case we need to get all items with "Keygroup" kg and then sort them
	filt := expression.Name(keygroupName).Equal(expression.Value(kg))
	proj := expression.NamesList(expression.Name(keyName), expression.Name(valName), expression.Name(expiryKey))

	expr, err := expression.NewBuilder().WithFilter(filt).WithProjection(proj).Build()

	if err != nil {
		log.Error().Msg(errors.New(err).ErrorStack())
		return nil, nil, nil, errors.New(err)
	}

	params := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		ProjectionExpression:      expr.Projection(),
		TableName:                 aws.String(s.dynamotable),
	}

	// Make the DynamoDB Query API call
	result, err := s.svc.Scan(context.TODO(), params)
	if err != nil {
		log.Error().Msg(errors.New(err).ErrorStack())
		return nil, nil, nil, errors.New(err)
	}

	log.Trace().Msgf("ReadAll: got %d items", len(result.Items))

	type item struct {
		key     string
		val     string
		version vclock.VClock
	}

	items := make([]item, 0, len(result.Items))

	for _, i := range result.Items {
		key, ok := i[keyName]

		if !ok {
			return nil, nil, nil, nil
		}

		k, ok := key.(*dynamoDBTypes.AttributeValueMemberS)

		if !ok {
			return nil, nil, nil, errors.Errorf("ReadAll: malformed key")
		}

		if k.Value == NULLValue {
			continue
		}

		val, ok := i[valName]

		if !ok {
			return nil, nil, nil, errors.Errorf("ReadAll: malformed value")
		}

		if e, ok := i[expiryKey]; ok {
			expiration, ok := e.(*dynamoDBTypes.AttributeValueMemberN)

			if ok {
				expiry, err := strconv.Atoi(expiration.Value)
				if err != nil {
					log.Error().Msg(errors.New(err).ErrorStack())
					return nil, nil, nil, errors.New(err)
				}

				log.Trace().Msgf("ReadAll found key expiring at %d, it is %d now", expiry, time.Now().UTC().Unix())

				// oops, item has expired - we treat this as "not found"
				if int64(expiry) < time.Now().UTC().Unix() {
					continue
				}
			}
		}
		values, ok := val.(*dynamoDBTypes.AttributeValueMemberM)

		if !ok || values == nil || len(values.Value) == 0 {
			continue
		}

		for v, data := range values.Value {
			version, err := vectorFromString(v)

			if err != nil {
				log.Error().Msg(errors.New(err).ErrorStack())
				return nil, nil, nil, errors.New(err)
			}

			it, ok := data.(*dynamoDBTypes.AttributeValueMemberS)

			if !ok {
				return nil, nil, nil, errors.Errorf("ReadAll: malformed item")
			}

			items = append(items, item{
				k.Value,
				it.Value,
				version,
			})
		}
	}

	// now we have a list of items, sort them and convert to lists
	keys := make([]string, 0, len(items))
	values := make([]string, 0, len(items))
	versions := make([]vclock.VClock, 0, len(items))

	// now we have lists of keys and items, we need to sort them by the key attribute alphabetically
	sort.Slice(items, func(i, j int) bool {
		return strings.ToLower(items[i].key) < strings.ToLower(items[j].key)
	})

	for _, x := range items {
		keys = append(keys, x.key)
		values = append(values, x.val)
		versions = append(versions, x.version)
	}

	return keys, values, versions, nil
}

// Append appends the item to the specified keygroup by incrementing the latest key by one.
func (s *Storage) Append(kg string, id string, val string, expiry int) error {
	log.Trace().Msgf("Appending key %s to keygroup %s", id, kg)

	cond := expression.AttributeNotExists(expression.Name(keygroupName))

	expr, err := expression.NewBuilder().WithCondition(cond).Build()

	if err != nil {
		log.Error().Msg(errors.New(err).ErrorStack())
		return errors.New(err)
	}

	in := &dynamodb.PutItemInput{
		Item: map[string]dynamoDBTypes.AttributeValue{
			keygroupName: &dynamoDBTypes.AttributeValueMemberS{
				Value: kg,
			},
			keyName: &dynamoDBTypes.AttributeValueMemberS{
				Value: id,
			},
			valName: &dynamoDBTypes.AttributeValueMemberM{
				Value: map[string]dynamoDBTypes.AttributeValue{
					vectorToString(vclock.VClock{}): &dynamoDBTypes.AttributeValueMemberS{
						Value: val,
					},
				},
			},
		},
		TableName:                 aws.String(s.dynamotable),
		ConditionExpression:       expr.Condition(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
	}

	if expiry > 0 {
		in.Item[expiryKey] = &dynamoDBTypes.AttributeValueMemberN{
			Value: strconv.FormatInt(time.Now().UTC().Unix()+int64(expiry), 10),
		}
	}

	_, err = s.svc.PutItem(context.TODO(), in)
	if err != nil {
		log.Error().Msg(errors.New(err).ErrorStack())
		return errors.New(err)
	}

	return nil
}

// IDs returns the keys of all items in the specified keygroup.
func (s *Storage) IDs(kg string) ([]string, error) {
	log.Trace().Msgf("Getting all keys from keygroup %s", kg)

	// in this case we need to get all items with "Keygroup" kg and then sort them
	filt := expression.Name(keygroupName).Equal(expression.Value(kg))
	proj := expression.NamesList(expression.Name(keyName), expression.Name(expiryKey))

	expr, err := expression.NewBuilder().WithFilter(filt).WithProjection(proj).Build()
	if err != nil {
		log.Error().Msg(errors.New(err).ErrorStack())
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

	result, err := s.svc.Scan(context.TODO(), params)
	if err != nil {
		log.Error().Msg(errors.New(err).ErrorStack())
		// continue: propably a projection issue
		// return nil, errors.New(err)
	}

	ids := make(map[string]struct{})

	for _, i := range result.Items {
		key, ok := i[keyName]

		if !ok {
			return nil, errors.Errorf("malformed key")
		}

		k, ok := key.(*dynamoDBTypes.AttributeValueMemberS)

		if !ok {
			return nil, errors.Errorf("malformed key")
		}

		if k.Value == NULLValue {
			continue
		}

		if e, ok := i[expiryKey]; ok {
			expiration, ok := e.(*dynamoDBTypes.AttributeValueMemberN)

			if ok {
				expiry, err := strconv.Atoi(expiration.Value)
				if err != nil {
					log.Error().Msg(errors.New(err).ErrorStack())
					return nil, errors.New(err)
				}

				log.Trace().Msgf("Read found key expiring at %d, it is %d now", expiry, time.Now().UTC().Unix())

				// oops, item has expired - we treat this as "not found"
				if int64(expiry) < time.Now().UTC().Unix() {
					continue
				}
			}
		}

		ids[k.Value] = struct{}{}
	}

	idList := make([]string, 0, len(ids))

	for i := range ids {
		idList = append(idList, i)
	}

	return idList, nil
}

// Update updates the item with the specified id in the specified keygroup.
func (s *Storage) Update(kg string, id string, val string, expiry int, vvector vclock.VClock) error {
	log.Trace().Msgf("Updating key %s in keygroup %s", id, kg)

	version := vectorToString(vvector)

	input := &dynamodb.UpdateItemInput{
		TableName: aws.String(s.dynamotable),
		Key: map[string]dynamoDBTypes.AttributeValue{
			keygroupName: &dynamoDBTypes.AttributeValueMemberS{
				Value: kg,
			},
			keyName: &dynamoDBTypes.AttributeValueMemberS{
				Value: id,
			},
		},
		ExpressionAttributeNames: map[string]string{
			"#value":   valName,
			"#version": version,
		},
		ExpressionAttributeValues: map[string]dynamoDBTypes.AttributeValue{
			":data": &dynamoDBTypes.AttributeValueMemberS{
				Value: val,
			},
		},
		UpdateExpression: aws.String("SET #value.#version = :data"),
	}

	if expiry > 0 {
		input.ExpressionAttributeNames["#expiry"] = expiryKey
		input.ExpressionAttributeValues[":expiry"] = &dynamoDBTypes.AttributeValueMemberN{
			Value: strconv.FormatInt(time.Now().UTC().Unix()+int64(expiry), 10),
		}
		input.UpdateExpression = aws.String(*input.UpdateExpression + " SET  #expiry = :expiry")
	}

	log.Trace().Msgf("Update: setting item %s (keygroup %s) with version %s to %s with expiry %d", id, kg, version, val, expiry)

	_, err := s.svc.UpdateItem(context.TODO(), input)

	if err != nil {
		// yes this is terrible
		// so apparently when you update a map value, like we do here, before the first time you add something to the
		// map you have to actually create the map. there is no way to do it on the fly.
		// so if we get an error that says "hey, there is no map here!", we have to create the data item with a
		// pre-filled map with our values first
		//
		// you know what's best about this (except that we have to do the same thing on other parts of the code
		// as well [smelly])?
		// there's not even a way to find out if we have the right error in a good way
		// the folks at AWS have forgotten to put the ValidationException error in the SDK
		// https://github.com/aws/aws-sdk/issues/47
		if e, ok := err.(*smithy.OperationError); ok {
			if e2, ok := e.Unwrap().(*shttp.ResponseError); ok {
				if e2.HTTPStatusCode() == http.StatusBadRequest {
					input := &dynamodb.UpdateItemInput{
						TableName: aws.String(s.dynamotable),
						Key: map[string]dynamoDBTypes.AttributeValue{
							keygroupName: &dynamoDBTypes.AttributeValueMemberS{
								Value: kg,
							},
							keyName: &dynamoDBTypes.AttributeValueMemberS{
								Value: id,
							},
						},
						ExpressionAttributeNames: map[string]string{
							"#value": valName,
						},
						ExpressionAttributeValues: map[string]dynamoDBTypes.AttributeValue{
							":data": &dynamoDBTypes.AttributeValueMemberM{
								Value: map[string]dynamoDBTypes.AttributeValue{
									version: &dynamoDBTypes.AttributeValueMemberS{
										Value: val,
									},
								},
							},
						},
						UpdateExpression: aws.String("SET #value = :data"),
					}

					if expiry > 0 {
						input.ExpressionAttributeNames["#expiry"] = expiryKey
						input.ExpressionAttributeValues[":expiry"] = &dynamoDBTypes.AttributeValueMemberN{
							Value: strconv.FormatInt(time.Now().UTC().Unix()+int64(expiry), 10),
						}
						input.UpdateExpression = aws.String(*input.UpdateExpression + ", #expiry = :expiry")
					}

					_, err := s.svc.UpdateItem(context.TODO(), input)
					if err != nil {
						log.Error().Msg(errors.New(err).ErrorStack())
						return errors.New(err)
					}
					return nil
				}
			}
		}

		log.Error().Msg(errors.New(err).ErrorStack())
		return errors.New(err)
	}

	return nil
}

// Delete deletes the item with the specified id from the specified keygroup.
func (s *Storage) Delete(kg string, id string, vvector vclock.VClock) error {
	log.Trace().Msgf("Deleting key %s from keygroup %s", id, kg)

	version := vectorToString(vvector)

	input := &dynamodb.UpdateItemInput{
		TableName: aws.String(s.dynamotable),
		Key: map[string]dynamoDBTypes.AttributeValue{
			keygroupName: &dynamoDBTypes.AttributeValueMemberS{
				Value: kg,
			},
			keyName: &dynamoDBTypes.AttributeValueMemberS{
				Value: id,
			},
		},
		ExpressionAttributeNames: map[string]string{
			"#version": version,
			"#value":   valName,
		},
		UpdateExpression: aws.String("REMOVE #value.#version"),
	}

	log.Trace().Msgf("Delete: removing item %s (keygroup %s) with version %s", id, kg, version)

	_, err := s.svc.UpdateItem(context.TODO(), input)
	if err != nil {
		if e, ok := err.(*smithy.OperationError); ok {
			// ignore error if it says that the key does not exist
			// want to allow idempotent deletes
			if e2, ok := e.Unwrap().(*shttp.ResponseError); ok {
				if e2.HTTPStatusCode() == http.StatusBadRequest {
					return nil
				}
			}

		}
		log.Error().Msg(errors.New(err).ErrorStack())
		return errors.New(err)

	}

	return nil
}

// Exists checks if the given data item exists in the dynamodb database.
func (s *Storage) Exists(kg string, id string) bool {
	log.Trace().Msgf("Checking if key %s exists in keygroup %s", id, kg)

	proj := expression.NamesList(expression.Name(expiryKey))
	expr, err := expression.NewBuilder().WithProjection(proj).Build()

	if err != nil {
		log.Error().Msg(errors.New(err).ErrorStack())
		log.Error().Msgf("Exists: %+v", err)
		return false
	}

	result, err := s.svc.GetItem(context.TODO(), &dynamodb.GetItemInput{
		Key: map[string]dynamoDBTypes.AttributeValue{
			keygroupName: &dynamoDBTypes.AttributeValueMemberS{
				Value: kg,
			},
			keyName: &dynamoDBTypes.AttributeValueMemberS{
				Value: id,
			},
		},
		TableName:                &s.dynamotable,
		ExpressionAttributeNames: expr.Names(),
		ProjectionExpression:     expr.Projection(),
	})

	if err != nil {
		log.Error().Msg(errors.New(err).ErrorStack())
		log.Error().Msgf("Exists: %s", err.Error())
		return false
	}

	// check if the item was found at all
	if result.Item == nil {
		return false
	}

	// check that the item hasn't actually expired but wasn't cleaned up
	if e, ok := result.Item[expiryKey]; ok {
		expiration, ok := e.(*dynamoDBTypes.AttributeValueMemberN)

		if ok {
			expiry, err := strconv.Atoi(expiration.Value)
			if err != nil {
				log.Error().Msg(errors.New(err).ErrorStack())
				return false
			}

			log.Trace().Msgf("Exists found key expiring at %d, it is %d now", expiry, time.Now().UTC().Unix())

			// oops, item has expired - we treat this as "not found"
			if int64(expiry) < time.Now().UTC().Unix() {
				return false
			}
		}
	}

	return true
}

// ExistsKeygroup checks if the given keygroup exists in the DynamoDB database.
func (s *Storage) ExistsKeygroup(kg string) bool {
	log.Trace().Msgf("Checking if keygroup %s exists", kg)

	result, err := s.svc.GetItem(context.TODO(), &dynamodb.GetItemInput{
		Key: map[string]dynamoDBTypes.AttributeValue{
			keygroupName: &dynamoDBTypes.AttributeValueMemberS{
				Value: kg,
			},
			keyName: &dynamoDBTypes.AttributeValueMemberS{
				Value: NULLValue,
			},
		},
		TableName: &s.dynamotable,
	})

	if err != nil {
		log.Error().Msg(errors.New(err).ErrorStack())
		log.Error().Msgf("ExistsKeygroup: %s", err.Error())
		return false
	}

	// check if the item was found at all
	if result.Item == nil {
		return false
	}

	return true
}

// CreateKeygroup creates the given keygroup in the DynamoDB database.
func (s *Storage) CreateKeygroup(kg string) error {
	log.Trace().Msgf("Creating keygroup %s", kg)

	input := &dynamodb.PutItemInput{
		Item: map[string]dynamoDBTypes.AttributeValue{
			keygroupName: &dynamoDBTypes.AttributeValueMemberS{
				Value: kg,
			},
			keyName: &dynamoDBTypes.AttributeValueMemberS{
				Value: NULLValue,
			},
			triggerName: &dynamoDBTypes.AttributeValueMemberM{
				Value: map[string]dynamoDBTypes.AttributeValue{},
			},
		},
		TableName: aws.String(s.dynamotable),
	}

	_, err := s.svc.PutItem(context.TODO(), input)
	if err != nil {
		log.Error().Msg(errors.New(err).ErrorStack())
		return errors.New(err)
	}

	return nil
}

// DeleteKeygroup deletes the given keygroup from the DynamoDB database.
func (s *Storage) DeleteKeygroup(kg string) error {
	log.Trace().Msgf("Deleting keygroup %s", kg)

	// delete all entries for that keygroup
	filt := expression.Name(keygroupName).Equal(expression.Value(kg))
	proj := expression.NamesList(expression.Name(keyName))

	expr, err := expression.NewBuilder().WithFilter(filt).WithProjection(proj).Build()
	if err != nil {
		log.Error().Msg(errors.New(err).ErrorStack())
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
	result, err := s.svc.Scan(context.TODO(), params)
	if err != nil {
		log.Error().Msg(errors.New(err).ErrorStack())
		return errors.New(err)
	}

	for _, i := range result.Items {

		k, ok := i[keyName]

		if !ok {
			return errors.Errorf("wrong format")
		}

		key, ok := k.(*dynamoDBTypes.AttributeValueMemberS)

		if !ok {
			return errors.Errorf("wrong format")
		}

		input := &dynamodb.DeleteItemInput{

			Key: map[string]dynamoDBTypes.AttributeValue{
				keygroupName: &dynamoDBTypes.AttributeValueMemberS{
					Value: kg,
				},
				keyName: &dynamoDBTypes.AttributeValueMemberS{
					Value: key.Value,
				},
			},
			TableName: aws.String(s.dynamotable),
		}

		_, err := s.svc.DeleteItem(context.TODO(), input)
		if err != nil {
			log.Error().Msg(errors.New(err).ErrorStack())
			return errors.New(err)
		}
	}

	// delete the keygroup configuration

	input := &dynamodb.DeleteItemInput{
		Key: map[string]dynamoDBTypes.AttributeValue{
			keygroupName: &dynamoDBTypes.AttributeValueMemberS{
				Value: kg,
			},
			keyName: &dynamoDBTypes.AttributeValueMemberS{
				Value: NULLValue,
			},
		},
		TableName: aws.String(s.dynamotable),
	}

	_, err = s.svc.DeleteItem(context.TODO(), input)
	if err != nil {
		log.Error().Msg(errors.New(err).ErrorStack())
		return errors.New(err)
	}

	return nil
}

// AddKeygroupTrigger stores a trigger node in the dynamodb database.
func (s *Storage) AddKeygroupTrigger(kg string, tid string, host string) error {
	log.Trace().Msgf("Adding trigger %s to keygroup %s", tid, kg)

	input := &dynamodb.UpdateItemInput{
		TableName: aws.String(s.dynamotable),
		Key: map[string]dynamoDBTypes.AttributeValue{
			keygroupName: &dynamoDBTypes.AttributeValueMemberS{
				Value: kg,
			},
			keyName: &dynamoDBTypes.AttributeValueMemberS{
				Value: NULLValue,
			},
		},
		ExpressionAttributeNames: map[string]string{
			"#tid":      tid,
			"#triggers": triggerName,
		},
		ExpressionAttributeValues: map[string]dynamoDBTypes.AttributeValue{
			":thost": &dynamoDBTypes.AttributeValueMemberS{
				Value: host,
			},
		},
		UpdateExpression: aws.String("SET #triggers.#tid = :thost"),
	}

	_, err := s.svc.UpdateItem(context.TODO(), input)

	if err != nil {
		// see Update method
		if e, ok := err.(*smithy.OperationError); ok {
			if e2, ok := e.Unwrap().(*shttp.ResponseError); ok {
				if e2.HTTPStatusCode() == http.StatusBadRequest {
					input := &dynamodb.UpdateItemInput{
						TableName: aws.String(s.dynamotable),
						Key: map[string]dynamoDBTypes.AttributeValue{
							keygroupName: &dynamoDBTypes.AttributeValueMemberS{
								Value: kg,
							},
							keyName: &dynamoDBTypes.AttributeValueMemberS{
								Value: NULLValue,
							},
						},
						ExpressionAttributeNames: map[string]string{
							"#triggers": triggerName,
						},
						ExpressionAttributeValues: map[string]dynamoDBTypes.AttributeValue{
							":triggers": &dynamoDBTypes.AttributeValueMemberM{
								Value: map[string]dynamoDBTypes.AttributeValue{
									tid: &dynamoDBTypes.AttributeValueMemberS{
										Value: host,
									},
								},
							},
						},
						UpdateExpression: aws.String("SET #triggers = :triggers"),
					}

					_, err := s.svc.UpdateItem(context.TODO(), input)
					if err != nil {
						log.Error().Msg(errors.New(err).ErrorStack())
						return errors.New(err)
					}
					return nil
				}
			}
		}

		log.Error().Msg(errors.New(err).ErrorStack())
		return errors.New(err)
	}

	return nil
}

// DeleteKeygroupTrigger removes a trigger node from the dynamodb database.
func (s *Storage) DeleteKeygroupTrigger(kg string, tid string) error {
	log.Trace().Msgf("Deleting trigger %s from keygroup %s", tid, kg)

	input := &dynamodb.UpdateItemInput{
		TableName: aws.String(s.dynamotable),
		Key: map[string]dynamoDBTypes.AttributeValue{
			keygroupName: &dynamoDBTypes.AttributeValueMemberS{
				Value: kg,
			},
			keyName: &dynamoDBTypes.AttributeValueMemberS{
				Value: NULLValue,
			},
		},
		ExpressionAttributeNames: map[string]string{
			"#tid":      tid,
			"#triggers": triggerName,
		},
		UpdateExpression: aws.String("REMOVE #triggers.#tid"),
	}

	_, err := s.svc.UpdateItem(context.TODO(), input)

	if err != nil {
		log.Error().Msg(errors.New(err).ErrorStack())
		return errors.New(err)
	}

	return nil
}

// GetKeygroupTrigger returns a map of all trigger nodes from the dynamodb database.
func (s *Storage) GetKeygroupTrigger(kg string) (map[string]string, error) {
	log.Trace().Msgf("Getting triggers for keygroup %s", kg)

	proj := expression.NamesList(expression.Name(triggerName))
	expr, err := expression.NewBuilder().WithProjection(proj).Build()

	if err != nil {
		log.Error().Msg(errors.New(err).ErrorStack())
		return nil, errors.New(err)
	}

	result, err := s.svc.GetItem(context.TODO(), &dynamodb.GetItemInput{
		Key: map[string]dynamoDBTypes.AttributeValue{
			keygroupName: &dynamoDBTypes.AttributeValueMemberS{
				Value: kg,
			},
			keyName: &dynamoDBTypes.AttributeValueMemberS{
				Value: NULLValue,
			},
		},
		TableName:                &s.dynamotable,
		ExpressionAttributeNames: expr.Names(),
		ProjectionExpression:     expr.Projection(),
	})

	if err != nil {
		log.Error().Msg(errors.New(err).ErrorStack())
		// continue: propably a projection issue
		// return nil, errors.New(err)
	}

	// check if the item was found at all
	if result.Item == nil {
		return nil, errors.Errorf("GetKeygroupTriggers: keygroup %s not found", kg)
	}

	// let's check the value of the item we got back
	triggers, ok := result.Item[triggerName]

	if !ok {
		return nil, errors.Errorf("GetKeygroupTriggers: keygroup %s not configured correctly", kg)
	}

	triggerMap, ok := triggers.(*dynamoDBTypes.AttributeValueMemberM)

	if !ok {
		return nil, errors.Errorf("GetKeygroupTriggers: keygroup %s not configured correctly", kg)
	}

	t := make(map[string]string)

	for tid, v := range triggerMap.Value {
		host, ok := v.(*dynamoDBTypes.AttributeValueMemberS)

		if !ok {
			return nil, errors.Errorf("GetKeygroupTriggers: keygroup %s not configured correctly", kg)
		}

		t[tid] = host.Value
	}

	return t, nil
}
