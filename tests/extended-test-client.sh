#!/usr/bin/env bash

# notice that whatever identifier you put in will also be used to find the right docker container
# terraform plan -var="instance_count=2" -var="identifier=tp" --var-file="secret.tfvars"

IDENTIFIER=tp
DEFAULT_WEB_PORT=80

ID_1=0
HOST_1=$ID_1.nodes.$IDENTIFIER.mcc-f.red
PORT_1=$DEFAULT_WEB_PORT

ID_2=1
HOST_2=$ID_2.nodes.$IDENTIFIER.mcc-f.red
PORT_2=$DEFAULT_WEB_PORT

APIVERSION=v0
KEYGROUP_NAME=testgroup
ID=1


# Create a Keygroup with Host 1
printf "\n"
printf "Create a Keygroup with Host 1\n"
printf "Calling http://%s:%s/%s/keygroup/%s\n" $HOST_1 $PORT_1 $APIVERSION $KEYGROUP_NAME

curl --request POST -sL \
     --url http://$HOST_1:$PORT_1/$APIVERSION/keygroup/$KEYGROUP_NAME \
     -i

# Register Node 2 at Node 1
printf "\n"
printf "Register Node 2 at Node 1\n"
printf "Calling http://%s:%s/%s/replica\n" $HOST_1 $PORT_1 $APIVERSION

curl --request POST -sL \
     --url http://$HOST_1:$PORT_1/$APIVERSION/replica \
     --data "{\"nodes\":[{\"id\":\"$ID_2\",\"addr\":\"$HOST_2\",\"port\":$PORT_2}]}" \
     -i

# Add Node 2 as a Replica Node for the Keygroup at Node 1
printf "\n"
printf "Add Node 2 as a Replica Node for the Keygroup at Node 1\n"
printf "Calling http://%s:%s/%s/keygroup/%s/replica/%s\n" $HOST_1 $PORT_1 $APIVERSION $KEYGROUP_NAME ID_2

curl --request POST -sL \
     --url http://$HOST_1:$PORT_1/$APIVERSION/keygroup/$KEYGROUP_NAME/replica/$ID_2 \
     -i

# Write an item to Node 1
printf "\n"
printf "Write an item to Node 1\n"
printf "Calling PUT http://%s:%s/%s/keygroup/%s/items/%s\n" $HOST_1 $PORT_1 $APIVERSION $KEYGROUP_NAME $ID

curl --request PUT -sL \
     --url http://$HOST_1:$PORT_1/$APIVERSION/keygroup/"$KEYGROUP_NAME"/data/"$ID" \
     --data '{"data":"hello other world!"}' \
     -i

# Read the item at Node 1
printf "\n"
printf "Read this item at Node 1\n"
printf "Calling GET http://%s:%s/%s/keygroup/%s/items/%s\n" $HOST_1 $PORT_1 $APIVERSION $KEYGROUP_NAME $ID

curl --request GET -sL \
     --url http://$HOST_1:$PORT_1/$APIVERSION/keygroup/"$KEYGROUP_NAME"/data/"$ID" \
     -i

# Read this item at Node 2
printf "\n"
printf "Read this item at Node 2\n"
printf "Calling GET http://%s:%s/%s/keygroup/%s/items/%s\n" $HOST_2 $PORT_2 $APIVERSION $KEYGROUP_NAME $ID

curl --request GET -sL \
     --url http://$HOST_2:$PORT_2/$APIVERSION/keygroup/"$KEYGROUP_NAME"/data/"$ID" \
     -i

# Delete the item at Node 2
printf "\n"
printf "Delete the item at Node 2\n"
printf "Calling DELETE http://%s:%s/%s/keygroup/%s/items/%s\n" $HOST_2 $PORT_2 $APIVERSION $KEYGROUP_NAME $ID

curl --request DELETE -sL \
     --url http://$HOST_2:$PORT_2/$APIVERSION/keygroup/"$KEYGROUP_NAME"/data/"$ID" \
     -i

# Read the item at Node 1
printf "\n"
printf "Read this item at Node 1\n"
printf "Calling GET http://%s:%s/%s/keygroup/%s/items/%s\n" $HOST_1 $PORT_1 $APIVERSION $KEYGROUP_NAME $ID

curl --request GET -sL \
     --url http://$HOST_1:$PORT_1/$APIVERSION/keygroup/"$KEYGROUP_NAME"/data/"$ID" \
     -i

# Delete the Keygroup at Node 1
printf "\n"
printf "Delete the Keygroup at Node 1\n"
printf "Calling DELETE http://%s:%s/%s/keygroup/%s\n" $HOST_1 $PORT_1 $APIVERSION $KEYGROUP_NAME

curl --request DELETE -sL \
     --url http://$HOST_1:$PORT_1/$APIVERSION/keygroup/$KEYGROUP_NAME \
     -i