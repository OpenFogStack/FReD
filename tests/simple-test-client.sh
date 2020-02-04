#!/usr/bin/env bash

HOST=localhost
PORT=9002
APIVERSION=v0
KEYGROUP_NAME=testgroup
ID=1

printf "\n"
printf "Deleting the Keygroup...\n"
printf "Calling DELETE http://%s:%s/%s/keygroupsdsd/%s\n" $HOST $PORT $APIVERSION $KEYGROUP_NAME

curl --request DELETE -sL \
     --url http://$HOST:$PORT/$APIVERSION/keygroup/$KEYGROUP_NAME \
     -i

printf "\n"
printf "Creating a Keygroup...\n"
printf "Calling http://%s:%s/%s/keygroup/%s\n" $HOST $PORT $APIVERSION $KEYGROUP_NAME

curl --request POST -sL \
     --url http://$HOST:$PORT/$APIVERSION/keygroup/$KEYGROUP_NAME \
     -i

printf "\n"
printf "Creating a Data Item in Keygroup...\n"
printf "Calling PUT http://%s:%s/%s/keygroup/%s/data/%s\n" $HOST $PORT $APIVERSION $KEYGROUP_NAME $ID

curl \
--request PUT -sL \
--url http://$HOST:$PORT/$APIVERSION/keygroup/"$KEYGROUP_NAME"/data/"$ID" \
-i \
-d \
'
{
  "data": {
    "id": "$ID",
    "type": "item",
    "attributes": {
      "value": "Hello World!"
    }
  }
}'


printf "\n"
printf "Reading Data Item from Keygroup...\n"
printf "Calling GET http://%s:%s/%s/keygroup/%s/data/%s\n" $HOST $PORT $APIVERSION $KEYGROUP_NAME $ID

curl --request GET -sL \
     --url http://$HOST:$PORT/$APIVERSION/keygroup/"$KEYGROUP_NAME"/data/"$ID" \
     -i

printf "\n"
printf "Updating Data Item in Keygroup...\n"
printf "Calling PUT http://%s:%s/%s/keygroup/%s/data/%s\n" $HOST $PORT $APIVERSION $KEYGROUP_NAME $ID

curl --request PUT -sL \
     --url http://$HOST:$PORT/$APIVERSION/keygroup/"$KEYGROUP_NAME"/data/"$ID" \
     -i \
     --data-binary @- << EOF
{
  "data": {
    "id": "1",
    "type": "item",
    "attributes": {
      "value": "Hello Other World!"
    }
  }
}
EOF

printf "\n"
printf "Reading Data Item from Keygroup...\n"
printf "Calling GET http://%s:%s/%s/keygroup/%s/data/%s\n" $HOST $PORT $APIVERSION $KEYGROUP_NAME $ID

curl --request GET -sL \
     --url http://$HOST:$PORT/$APIVERSION/keygroup/"$KEYGROUP_NAME"/data/"$ID" \
     -i

printf "\n"
printf "Deleting Data Item from Keygroup...\n"
printf "Calling DELETE http://%s:%s/%s/keygroup/%s/data/%s\n" $HOST $PORT $APIVERSION $KEYGROUP_NAME $ID

curl --request DELETE -sL \
     --url http://$HOST:$PORT/$APIVERSION/keygroup/"$KEYGROUP_NAME"/data/"$ID" \
     -i

printf "\n"
printf "Reading Data Item from Keygroup...\n"
printf "Calling GET http://%s:%s/%s/keygroup/%s/data/%s\n" $HOST $PORT $APIVERSION $KEYGROUP_NAME $ID

curl --request GET -sL \
     --url http://$HOST:$PORT/$APIVERSION/keygroup/"$KEYGROUP_NAME"/data/"$ID" \
     -i

printf "\n"
printf "Deleting the Keygroup...\n"
printf "Calling DELETE http://%s:%s/%s/keygroup/%s\n" $HOST $PORT $APIVERSION $KEYGROUP_NAME

curl --request DELETE -sL \
     --url http://$HOST:$PORT/$APIVERSION/keygroup/$KEYGROUP_NAME \
     -i