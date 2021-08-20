#!/bin/bash

# usage: client.sh <node_addr> <node_id> <client_id>
# check that we got the parameter we needed or exit the script with a usage message
[ $# -ne 3 ] && { echo "Usage: $0 node_addr node_id client_id"; exit 1; }

NODE_ADDR=$1
NODE_ID=$2
CLIENT_ID=$3

printf "Got node address %s\n" "$NODE_ADDR"

# start alexandra

./alexandra --address :10000 \
      --lighthouse "$NODE_ADDR" \
      --ca-cert /cert/ca.crt \
      --alexandra-key /cert/client.key \
      --alexandra-cert /cert/client.crt \
      --clients-key /cert/client.key \
      --clients-cert /cert/client.crt \
      --log-level info \
      --log-handler dev \
      &

sleep 10

# start local client and connect to alexandra, run tests
python3 client.py \
    --id "$CLIENT_ID" \
    --node-id "$NODE_ID" \
    --ops 100 \
    --host localhost:10000 \
    --ca /cert/ca.crt \
    --cert /cert/client.crt \
    --key /cert/client.key

# exit
exit $?