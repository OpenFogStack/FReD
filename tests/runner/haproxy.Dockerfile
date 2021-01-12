FROM haproxy:lts-alpine

COPY nodeA-haproxy.cfg /usr/local/etc/haproxy/haproxy.cfg