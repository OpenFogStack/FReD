#!/usr/bin/env bash

# notice that whatever identifier you put in will also be used to find the right docker container

# 1. terraform plan -var="instance_count=2" -var="identifier=$IDENTIFIER" --var-file="terraform/small-test/secret.tfvars" terraform/small-test
# 2. docker build -t gitlab-registry.tubit.tu-berlin.de/mcc-fred/fred/fred:$IDENTIFIER .
# 3. docker push gitlab-registry.tubit.tu-berlin.de/mcc-fred/fred/fred:$IDENTIFIER
# 4. terraform apply -var="instance_count=2" -var="identifier=$IDENTIFIER" --var-file="terraform/small-test/secret.tfvars" terraform/small-test
# 5. do the test (sh ./tests/extended-test-client.sh $IDENTIFIER)
# 6. terraform destroy -var="instance_count=2" -var="identifier=$IDENTIFIER" --var-file="terraform/small-test/secret.tfvars" terraform/small-test

# for debugging, connect to the instance with `ssh -i "terraform/small-test/terraform.key" ec2-user@0.nodes.$IDENTIFIER.mcc-f.red`
# enter `sudo docker logs fred --follow` there to see the logs

# PRTCL=https
PRTCL=http

IDENTIFIER=$1
DEFAULT_WEB_PORT=80
#DEFAULT_WEB_PORT=9001
DEFAULT_ZMQ_PORT=5555

ID_1=0
#HOST_1=$ID_1.nodes.$IDENTIFIER.mcc-f.red
HOST_1=18.197.166.114
#HOST_1=localhost
PORT_1=$DEFAULT_WEB_PORT
Z_PORT_1=$DEFAULT_ZMQ_PORT

ID_2=1
#HOST_2=$ID_2.nodes.$IDENTIFIER.mcc-f.red
HOST_2=3.123.134.250
PORT_2=$DEFAULT_WEB_PORT
Z_PORT_2=$DEFAULT_ZMQ_PORT

APIVERSION=v0
KEYGROUP_NAME=testgroup
ITEM_ID_1=1
ITEM_ID_2=2

wait () {
  printf "\n"
  printf "\n"
  read -p "Press any key to continue... " -n1 -s
  printf "\n"
  printf "\n"
}

# Seed Host 1
printf "\n"
printf "Seed Host 1\n"
printf "Calling %s://%s:%s/%s/seed\n" $PRTCL $HOST_1 $PORT_1 $APIVERSION

curl --request POST -sL \
     --url $PRTCL://$HOST_1:$PORT_1/$APIVERSION/seed \
     --data "{\"id\":\"$ID_1\",\"addr\":\"$HOST_1\"}" \
     -i \
EOF

wait

# Create a Keygroup with Host 1
printf "\n"
printf "Create a Keygroup with Host 1\n"
printf "Calling %s://%s:%s/%s/keygroup/%s\n" $PRTCL $HOST_1 $PORT_1 $APIVERSION $KEYGROUP_NAME

curl --request POST -sL \
     --url $PRTCL://$HOST_1:$PORT_1/$APIVERSION/keygroup/$KEYGROUP_NAME \
     -i

wait

# Write an item to Node 1
printf "\n"
printf "Write an item to Node 1\n"
printf "Calling PUT %s://%s:%s/%s/keygroup/%s/data/%s\n" $PRTCL $HOST_1 $PORT_1 $APIVERSION $KEYGROUP_NAME $ITEM_ID_1

curl --request PUT -sL \
     --url $PRTCL://$HOST_1:$PORT_1/$APIVERSION/keygroup/"$KEYGROUP_NAME"/data/"$ITEM_ID_1" \
     --data "{\"id\":\"$ITEM_ID_1\",\"value\":\"hello other world!\",\"keygroup\":\"$KEYGROUP_NAME\"}" \
     -i

wait

# Register Node 2 at Node 1
printf "\n"
printf "Register Node 2 at Node 1\n"
printf "Calling %s://%s:%s/%s/replica\n" $PRTCL $HOST_1 $PORT_1 $APIVERSION

curl --request POST -sL \
     --url $PRTCL://$HOST_1:$PORT_1/$APIVERSION/replica \
     --data "{\"nodes\":[{\"id\":\"$ID_2\",\"addr\":\"$HOST_2\",\"zmqPort\":$Z_PORT_2}]}" \
     -i

wait

# Add Node 2 as a Replica Node for the Keygroup at Node 1
printf "\n"
printf "Add Node 2 as a Replica Node for the Keygroup at Node 1\n"
printf "Calling %s://%s:%s/%s/keygroup/%s/replica/%s\n" $PRTCL $HOST_1 $PORT_1 $APIVERSION $KEYGROUP_NAME ID_2

curl --request POST -sL \
     --url $PRTCL://$HOST_1:$PORT_1/$APIVERSION/keygroup/$KEYGROUP_NAME/replica/$ID_2 \
     -i

wait

# Read the item at Node 1
printf "\n"
printf "Read this item at Node 1\n"
printf "Calling GET %s://%s:%s/%s/keygroup/%s/data/%s\n" $PRTCL $HOST_1 $PORT_1 $APIVERSION $KEYGROUP_NAME $ITEM_ID_1

curl --request GET -sL \
     --url $PRTCL://$HOST_1:$PORT_1/$APIVERSION/keygroup/"$KEYGROUP_NAME"/data/"$ITEM_ID_1" \
     -i

wait

# Read this item at Node 2
printf "\n"
printf "Read this item at Node 2\n"
printf "Calling GET %s://%s:%s/%s/keygroup/%s/data/%s\n" $PRTCL $HOST_2 $PORT_2 $APIVERSION $KEYGROUP_NAME $ITEM_ID_1

curl --request GET -sL \
     --url $PRTCL://$HOST_2:$PORT_2/$APIVERSION/keygroup/"$KEYGROUP_NAME"/data/"$ITEM_ID_1" \
     -i

wait

# Write another item to Node 1
printf "\n"
printf "Write another item to Node 1\n"
printf "Calling PUT %s://%s:%s/%s/keygroup/%s/data/%s\n" $PRTCL $HOST_1 $PORT_1 $APIVERSION $KEYGROUP_NAME $ITEM_ID_2

curl --request PUT -sL \
     --url http://$HOST_1:$PORT_1/$APIVERSION/keygroup/"$KEYGROUP_NAME"/data/"$ITEM_ID_2" \
     --data "{\"id\":\"$ITEM_ID_2\",\"value\":\"hello world!\",\"keygroup\":\"$KEYGROUP_NAME\"}" \
     -i

wait

# Read the other item at Node 1
printf "\n"
printf "Read the other item at Node 1\n"
printf "Calling GET %s://%s:%s/%s/keygroup/%s/data/%s\n" $PRTCL $HOST_1 $PORT_1 $APIVERSION $KEYGROUP_NAME $ITEM_ID_2

curl --request GET -sL \
     --url http://$HOST_1:$PORT_1/$APIVERSION/keygroup/"$KEYGROUP_NAME"/data/"$ITEM_ID_2" \
     -i

wait

# Read this other item at Node 2
printf "\n"
printf "Read this other item at Node 2\n"
printf "Calling GET %s://%s:%s/%s/keygroup/%s/data/%s\n" $PRTCL $HOST_2 $PORT_2 $APIVERSION $KEYGROUP_NAME $ITEM_ID_2

curl --request GET -sL \
     --url http://$HOST_2:$PORT_2/$APIVERSION/keygroup/"$KEYGROUP_NAME"/data/"$ITEM_ID_2" \
     -i

wait

# Delete the item at Node 2
printf "\n"
printf "Delete the item at Node 2\n"
printf "Calling DELETE %s://%s:%s/%s/keygroup/%s/data/%s\n" $PRTCL $HOST_2 $PORT_2 $APIVERSION $KEYGROUP_NAME $ITEM_ID_1

curl --request DELETE -sL \
     --url http://$HOST_2:$PORT_2/$APIVERSION/keygroup/"$KEYGROUP_NAME"/data/"$ITEM_ID_1" \
     -i

wait

# Read the item at Node 1
printf "\n"
printf "Read this item at Node 1\n"
printf "Calling GET %s://%s:%s/%s/keygroup/%s/data/%s\n" $PRTCL $HOST_1 $PORT_1 $APIVERSION $KEYGROUP_NAME $ITEM_ID_1

curl --request GET -sL \
     --url http://$HOST_1:$PORT_1/$APIVERSION/keygroup/"$KEYGROUP_NAME"/data/"$ITEM_ID_1" \
     -i

wait

# Delete the Keygroup at Node 1
printf "\n"
printf "Delete the Keygroup at Node 1\n"
printf "Calling DELETE %s://%s:%s/%s/keygroup/%s\n" $PRTCL $HOST_1 $PORT_1 $APIVERSION $KEYGROUP_NAME

curl --request DELETE -sL \
     --url http://$HOST_1:$PORT_1/$APIVERSION/keygroup/$KEYGROUP_NAME \
     -i