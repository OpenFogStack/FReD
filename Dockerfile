# building the binary
FROM golang:1.13-alpine as golang

MAINTAINER Tobias Pfandzelter <tp@mcc.tu-berlin.de>

RUN apk add --no-cache libzmq-static czmq-dev libsodium-static build-base util-linux-dev

# stolen from https://github.com/drone/ca-certs/blob/master/Dockerfile
RUN apk add -U --no-cache ca-certificates

RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2

WORKDIR /go/src/gitlab.tu-berlin.de/mcc-fred/fred/

COPY . .

# Static build required so that we can safely copy the binary over.
RUN ls
RUN touch ./cmd/frednode/dummy.cc
RUN go install -a -ldflags '-linkmode external -w -s -extldflags "-static -luuid" ' ./cmd/frednode/

# actual Docker image
FROM scratch

WORKDIR /

COPY --from=golang /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=golang /go/bin/frednode frednode

# webserver ports
# if use-tls=false, only port 80 will be used
# if use-tls=true, port 80 will be used for ACME and HTTP, port 443 will be used for HTTPS
EXPOSE 80
EXPOSE 443

EXPOSE 5555

ENTRYPOINT ["./frednode"]
