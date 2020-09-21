ETCD_VER=v3.4.7

# choose either URL
GOOGLE_URL=https://storage.googleapis.com/etcd
DOWNLOAD_URL=${GOOGLE_URL}

curl -L ${DOWNLOAD_URL}/${ETCD_VER}/etcd-${ETCD_VER}-linux-amd64.tar.gz -o etcd-${ETCD_VER}-linux-amd64.tar.gz