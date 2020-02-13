#!/usr/bin/env bash

HOST=0.nodes.tp.mcc-f.red
PORT=443
APIVERSION=v0
KEYGROUP_NAME=testgroup
ID=1
PRTCL=https
# PRTCL=http

printf "\n"
printf "Deleting the Keygroup...\n"
printf "Calling DELETE %s://%s:%s/%s/keygroupsdsd/%s\n" $PRTCL $HOST $PORT $APIVERSION $KEYGROUP_NAME

curl --request DELETE -sL \
     --url $PRTCL://$HOST:$PORT/$APIVERSION/keygroup/$KEYGROUP_NAME \
     -i

printf "\n"
printf "Creating a Keygroup...\n"
printf "Calling %s://%s:%s/%s/keygroup/%s\n" $PRTCL $HOST $PORT $APIVERSION $KEYGROUP_NAME

curl --request POST -sL \
     --url $PRTCL://$HOST:$PORT/$APIVERSION/keygroup/$KEYGROUP_NAME \
     -i

printf "\n"
printf "Creating a Data Item in Keygroup...\n"
printf "Calling PUT %s://%s:%s/%s/keygroup/%s/data/%s\n" $PRTCL $HOST $PORT $APIVERSION $KEYGROUP_NAME $ID

curl --request PUT -sL \
     --url $PRTCL://$HOST:$PORT/$APIVERSION/keygroup/"$KEYGROUP_NAME"/data/"$ID" \
     -i \
     --data "{\"id\":\"$ID\",\"value\":\"hello world!\",\"keygroup\":\"$KEYGROUP_NAME\"}" \

printf "\n"
printf "Reading Data Item from Keygroup...\n"
printf "Calling GET %s://%s:%s/%s/keygroup/%s/data/%s\n" $PRTCL $HOST $PORT $APIVERSION $KEYGROUP_NAME $ID

curl --request GET -sL \
     --url $PRTCL://$HOST:$PORT/$APIVERSION/keygroup/"$KEYGROUP_NAME"/data/"$ID" \
     -i

printf "\n"
printf "Updating Data Item in Keygroup...\n"
printf "Calling PUT %s://%s:%s/%s/keygroup/%s/data/%s\n" $PRTCL $HOST $PORT $APIVERSION $KEYGROUP_NAME $ID

curl --request PUT -sL \
     --url $PRTCL://$HOST:$PORT/$APIVERSION/keygroup/"$KEYGROUP_NAME"/data/"$ID" \
     -i \
     --data "{\"id\":\"$ID\",\"value\":\"hello other world!\",\"keygroup\":\"$KEYGROUP_NAME\"}" \

printf "\n"
printf "Reading Data Item from Keygroup...\n"
printf "Calling GET %s://%s:%s/%s/keygroup/%s/data/%s\n" $PRTCL $HOST $PORT $APIVERSION $KEYGROUP_NAME $ID

curl --request GET -sL \
     --url $PRTCL://$HOST:$PORT/$APIVERSION/keygroup/"$KEYGROUP_NAME"/data/"$ID" \
     -i

printf "\n"
printf "Deleting Data Item from Keygroup...\n"
printf "Calling DELETE %s://%s:%s/%s/keygroup/%s/data/%s\n" $PRTCL $HOST $PORT $APIVERSION $KEYGROUP_NAME $ID

curl --request DELETE -sL \
     --url $PRTCL://$HOST:$PORT/$APIVERSION/keygroup/"$KEYGROUP_NAME"/data/"$ID" \
     -i

printf "\n"
printf "Reading Data Item from Keygroup...\n"
printf "Calling GET %s://%s:%s/%s/keygroup/%s/data/%s\n" $PRTCL $HOST $PORT $APIVERSION $KEYGROUP_NAME $ID

curl --request GET -sL \
     --url $PRTCL://$HOST:$PORT/$APIVERSION/keygroup/"$KEYGROUP_NAME"/data/"$ID" \
     -i

printf "\n"
printf "Deleting the Keygroup...\n"
printf "Calling DELETE %s://%s:%s/%s/keygroup/%s\n" $PRTCL $HOST $PORT $APIVERSION $KEYGROUP_NAME

curl --request DELETE -sL \
     --url $PRTCL://$HOST:$PORT/$APIVERSION/keygroup/$KEYGROUP_NAME \
     -i