#!/usr/bin/env bash
HOST=localhost
PORT=9001
KEYGROUP_NAME=testgroup

printf "\n"
printf "Creating a Keygroup...\n"
printf "Calling http://%s:%s/keygroup/%s\n" $HOST $PORT $KEYGROUP_NAME

curl --request POST -sL \
     --url http://$HOST:$PORT/keygroup/$KEYGROUP_NAME \
     -i

printf "\n"
printf "Creating a Data Item in Keygroup...\n"
printf "Calling POST http://%s:%s/keygroup/%s/items\n" $HOST $PORT $KEYGROUP_NAME

curl --request POST -sL \
     --url http://$HOST:$PORT/keygroup/$KEYGROUP_NAME/items \
     --data '{"data":"hello world!"}' \
     -i

ID=1

printf "\n"
printf "Reading Data Item from Keygroup...\n"
printf "Calling GET http://%s:%s/keygroup/%s/items/%s\n" $HOST $PORT $KEYGROUP_NAME $ID

curl --request GET -sL \
     --url http://$HOST:$PORT/keygroup/"$KEYGROUP_NAME"/items/"$ID" \
     -i

printf "\n"
printf "Updating Data Item in Keygroup...\n"
printf "Calling PUT http://%s:%s/keygroup/%s/items/%s\n" $HOST $PORT $KEYGROUP_NAME $ID

curl --request PUT -sL \
     --url http://$HOST:$PORT/keygroup/"$KEYGROUP_NAME"/items/"$ID" \
     --data '{"data":"hello other world!"}' \
     -i

printf "\n"
printf "Reading Data Item from Keygroup...\n"
printf "Calling GET http://%s:%s/keygroup/%s/items/%s\n" $HOST $PORT $KEYGROUP_NAME $ID

curl --request GET -sL \
     --url http://$HOST:$PORT/keygroup/"$KEYGROUP_NAME"/items/"$ID" \
     -i

printf "\n"
printf "Deleting Data Item from Keygroup...\n"
printf "Calling DELETE http://%s:%s/keygroup/%s/items/%s\n" $HOST $PORT $KEYGROUP_NAME $ID

curl --request DELETE -sL \
     --url http://$HOST:$PORT/keygroup/"$KEYGROUP_NAME"/items/"$ID" \
     -i

printf "\n"
printf "Reading Data Item from Keygroup...\n"
printf "Calling GET http://%s:%s/keygroup/%s/items/%s\n" $HOST $PORT $KEYGROUP_NAME $ID

curl --request GET -sL \
     --url http://$HOST:$PORT/keygroup/"$KEYGROUP_NAME"/items/"$ID" \
     -i

printf "\n"
printf "Deleting the Keygroup...\n"
printf "Calling DELETE http://%s:%s/keygroup/%s\n" $HOST $PORT $KEYGROUP_NAME

curl --request DELETE -sL \
     --url http://$HOST:$PORT/keygroup/$KEYGROUP_NAME \
     -i