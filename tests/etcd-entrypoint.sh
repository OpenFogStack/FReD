#!/bin/sh
set -e

PARAMS=${@:-""}

# check if INSECURE env variable is set
if [ -n "$INSECURE" ]; then
    echo "Running insecure etcd"
    # remove param client-cert-auth
    PARAMS=$(echo $PARAMS | sed -e 's/--client-cert-auth//g')
fi
/bin/sh $PARAMS
