FROM gcr.io/etcd-development/etcd:v3.5.7 as ETCD_BINS

FROM debian:bullseye-slim

COPY --from=ETCD_BINS /usr/local/bin/etcd /usr/local/bin/
COPY --from=ETCD_BINS /usr/local/bin/etcdctl /usr/local/bin/
COPY --from=ETCD_BINS /usr/local/bin/etcdutl /usr/local/bin/

WORKDIR /var/etcd/
WORKDIR /var/lib/etcd/

EXPOSE 2379 2380

COPY etcd-entrypoint.sh entrypoint.sh
RUN chmod +x entrypoint.sh

ENTRYPOINT ["/bin/sh", "entrypoint.sh", "/usr/local/bin/etcd"]
