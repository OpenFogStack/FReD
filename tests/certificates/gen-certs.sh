#!/bin/bash

# usage: gen-certs.sh

openssl genrsa -out ca.key 2048

openssl req -x509 -new -nodes \
     -key ca.key -sha512 \
     -days 1825 -out ca.crt

function generate {
  # give better names to parameter variables
  NAME=$1

  # generate a key
  openssl genrsa -out "${NAME}".key 2048

  # write the config file
  # shellcheck disable=SC2086
  cat > ${NAME}.conf <<EOF

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
  echo "CN = ${NAME}" >> "${NAME}".conf

  # shellcheck disable=SC2086
  cat >> ${NAME}.conf <<EOF

[v3_req]
keyUsage = keyEncipherment, dataEncipherment, digitalSignature
extendedKeyUsage = serverAuth, clientAuth
subjectAltName = @alt_names

[alt_names]
DNS.1 = localhost
IP.1 = 127.0.0.1
EOF

  # write the IP SAN into the config file
  IPNUM=1

  for i in "$@"; do
    # skip the first parameter
    if [ "$IPNUM" -eq 1 ]; then
      IPNUM=$((IPNUM+1))
      continue
    fi

    echo "IP.${IPNUM} = ${i}" >> "${NAME}".conf
    IPNUM=$((IPNUM+1))
  done

  # generate the CSR
  openssl req -new -key "${NAME}".key -out "${NAME}".csr -config "${NAME}".conf

  # build the certificate
  openssl x509 -req -in "${NAME}".csr -CA ca.crt -CAkey ca.key \
  -CAcreateserial -out "${NAME}".crt -days 1825 \
  -extfile "${NAME}".conf -extensions v3_req

  # delete the config file and csr
  rm "${NAME}".conf
  rm "${NAME}".csr
}

generate "alexandra" "172.26.4.1"
generate "alexandraTester" "172.26.7.1"
generate "client" "172.26.4.1"
generate "etcd" "172.26.6.1"
generate "littleclient" "172.26.4.1"
generate "nodeA" "172.26.1.1" "172.26.1.101" "172.26.1.102" "172.26.1.103"
generate "nodeB" "172.26.2.1"
generate "nodeC" "172.26.3.1"
generate "storeA" "172.26.1.104"
generate "storeB" "172.26.2.2"
generate "storeC" "172.26.3.2"
generate "trigger" "172.26.5.1"