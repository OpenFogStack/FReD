HOST=localhost
PORT=9000
KEYGROUP_NAME=testgroup

echo "Creating a Keygroup..."

curl --request POST -sL \
     --url http://$HOST:$PORT/keygroup/$KEYGROUP_NAME \
     -i

echo "Creating a Data Item in Keygroup..."

ID=$(curl --request POST -sL \
     --url http://$HOST:$PORT/keygroup/$KEYGROUP_NAME/items \
     --data "{'data':'hello world!'}" \
     )

echo "Reading Data Item from Keygroup..."

curl --request GET -sL \
     --url http://$HOST:$PORT/keygroup/"$KEYGROUP_NAME"/items/"$ID" \
     -i

echo "Updating Data Item in Keygroup..."

curl --request PUT -sL \
     --url http://$HOST:$PORT/keygroup/"$KEYGROUP_NAME"/items/"$ID" \
     --data "{'data':'hello other world!'}" \
     -i

echo "Reading Data Item from Keygroup..."

curl --request GET -sL \
     --url http://$HOST:$PORT/keygroup/"$KEYGROUP_NAME"/items/"$ID" \
     -i

echo "Deleting Data Item from Keygroup..."

curl --request DELETE -sL \
     --url http://$HOST:$PORT/keygroup/"$KEYGROUP_NAME"/items/"$ID" \
     -i

echo "Reading Data Item from Keygroup..."

curl --request GET -sL \
     --url http://$HOST:$PORT/keygroup/"$KEYGROUP_NAME"/items/"$ID" \
     -i

echo "Deleting the Keygroup..."

curl --request DELETE -sL \
     --url http://$HOST:$PORT/keygroup/$KEYGROUP_NAME \
     -i