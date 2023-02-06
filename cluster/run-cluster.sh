#!/bin/bash

#constants
CLUSTER_NAME=cluster-
NET_NAME=clusternetwork
CERT_FOLDER="$(pwd)"/certs
GATEWAY=172.18.20.1
SUBNET=172.18.20.0/24
BASE_IP=172.18.20.
ETCD_IP=172.18.20.2
ETCD_VERSION=v3.5.7

docker network remove "$NET_NAME" 2&> /dev/null || true

gen_cert() {
  NAME=$1
  IP=$2

  rm "$CERT_FOLDER"/"${NAME}".crt || true
  rm "$CERT_FOLDER"/"${NAME}".key || true

  # generate a key
  openssl genrsa -out "$CERT_FOLDER"/"${NAME}".key 2048

  # write the config file
# shellcheck disable=SC2086
  cat > "$CERT_FOLDER"/${NAME}.conf <<EOF

  [ req ]
  default_bits = 2048
  prompt = no
  default_md = sha512
  req_extensions = v3_req
  distinguished_name = dn

  [ dn ]
  C = DE
  ST = Berlin
  L = Berlin
  O = MCC
  OU = FRED
EOF

  # write the CN into the config file
  echo "CN = ${NAME}" >> "$CERT_FOLDER"/"${NAME}".conf

  cat >> "$CERT_FOLDER"/"${NAME}".conf <<EOF
  [v3_req]
  keyUsage = keyEncipherment, dataEncipherment, digitalSignature
  extendedKeyUsage = serverAuth, clientAuth
  subjectAltName = @alt_names

  [alt_names]
  DNS.1 = localhost
  IP.1 = 127.0.0.1
EOF

  # write the IP SAN into the config file
  echo "IP.2 = ${IP}" >> "$CERT_FOLDER"/"${NAME}".conf

  # generate the CSR
  openssl req -new \
    -key "$CERT_FOLDER"/"${NAME}".key \
    -out "$CERT_FOLDER"/"${NAME}".csr \
    -config "$CERT_FOLDER"/"${NAME}".conf

  # build the certificate
  openssl x509 -req -in "$CERT_FOLDER"/"${NAME}".csr \
    -CA "$CERT_FOLDER"/ca.crt \
    -CAkey "$CERT_FOLDER"/ca.key \
    -CAcreateserial \
    -out "$CERT_FOLDER"/"${NAME}".crt \
    -days 1825 \
    -extfile "$CERT_FOLDER"/"${NAME}".conf \
    -extensions v3_req

  # delete the config file and csr
  rm "$CERT_FOLDER"/"${NAME}".conf
  rm "$CERT_FOLDER"/"${NAME}".csr
}

# usage: run-cluster.sh <num_nodes>
# check that we got the parameter we needed or exit the script with a usage message
[ $# -ne 1 ] ||  echo "$1" | grep -E -q -v '^[0-9]+$'  && { echo "Usage: $0 num_nodes"; exit 1; }

# prettier name
NUM_NODES=$1

# create a network
docker network create "$NET_NAME" --gateway "$GATEWAY" --subnet "$SUBNET" || exit 1

# generate certificates
gen_cert etcd "$ETCD_IP" || exit 1

for (( i = 1; i <= NUM_NODES; i=i+1 ))
do
       gen_cert node"$i" "$BASE_IP$(( i+2 ))" || exit 1
done

# start etcd

docker pull gcr.io/etcd-development/etcd:${ETCD_VERSION}
docker run -d \
  --name "$CLUSTER_NAME"etcd \
  -v "$CERT_FOLDER"/etcd.crt:/cert/etcd.crt \
  -v "$CERT_FOLDER"/etcd.key:/cert/etcd.key \
  -v "$CERT_FOLDER"/ca.crt:/cert/ca.crt \
  --network="$NET_NAME" \
  --ip="$ETCD_IP" \
  gcr.io/etcd-development/etcd:${ETCD_VERSION} \
  etcd \
  --name s-1 \
  --data-dir /tmp/etcd/s-1 \
  --listen-client-urls https://"$ETCD_IP":2379 \
  --advertise-client-urls https://"$ETCD_IP":2379 \
  --listen-peer-urls http://"$ETCD_IP":2380 \
  --initial-advertise-peer-urls http://"$ETCD_IP":2380 \
  --initial-cluster s-1=http://"$ETCD_IP":2380 \
  --initial-cluster-token tkn \
  --initial-cluster-state new \
  --cert-file=/cert/etcd.crt \
  --key-file=/cert/etcd.key \
  --client-cert-auth \
  --trusted-ca-file=/cert/ca.crt

# start as many containers as needed
docker build -t fred ../.

for (( i = 1; i <= NUM_NODES; i=i+1 ))
do

  docker run -d  \
    --name "$CLUSTER_NAME"node"$(( i ))" \
    -v "$CERT_FOLDER"/node"$i".crt:/cert/node.crt \
    -v "$CERT_FOLDER"/node"$i".key:/cert/node.key \
    -v "$CERT_FOLDER"/ca.crt:/cert/ca.crt \
    --network="$NET_NAME" \
    --ip="$BASE_IP$(( i+2 ))" \
    fred \
    --log-level debug \
    --handler dev \
    --peer-host "$BASE_IP$(( i+2 ))":5555 \
    --nodeID node"$(( i ))" \
    --host "$BASE_IP$(( i+2 ))":9001 \
    --cert /cert/node.crt \
    --key /cert/node.key \
    --ca-file /cert/ca.crt \
    --peer-cert /cert/node.crt \
    --peer-key /cert/node.key \
    --peer-ca /cert/ca.crt \
    --adaptor memory \
    --nase-host https://"$ETCD_IP":2379 \
    --nase-cert /cert/node.crt \
    --nase-key /cert/node.key \
    --nase-ca /cert/ca.crt \
    --nase-cached \
    --handler dev \
    --remote-storage-cert /cert/node.crt \
    --remote-storage-key /cert/node.key  \
    --remote-storage-ca /cert/ca.crt  \
    --trigger-cert /cert/node.crt \
    --trigger-key /cert/node.key \
    --trigger-ca /cert/ca.crt

done

cleanup() {
  docker stop "$CLUSTER_NAME"etcd
  docker rm "$CLUSTER_NAME"etcd

  for (( i = 1; i <= NUM_NODES; i=i+1 ))
  do
         docker stop "$CLUSTER_NAME"node"$(( i ))"
         docker rm "$CLUSTER_NAME"node"$(( i ))"
  done

  docker network remove "$NET_NAME"

  exit 0
}

trap 'cleanup' INT

echo "press Ctrl-C to stop cluster..."

while true ; do
    true
done