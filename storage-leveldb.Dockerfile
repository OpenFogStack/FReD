# Shamelessly stolen from the original dockerfile
# building the binary
FROM golang:1.13-alpine as golang

MAINTAINER Tobias Pfandzelter <tp@mcc.tu-berlin.de>

RUN apk add --no-cache build-base util-linux-dev

WORKDIR /go/src/gitlab.tu-berlin.de/mcc-fred/fred/

COPY . .

# Static build required so that we can safely copy the binary over.
RUN go install -a -ldflags '-linkmode external -w -s -extldflags "-static -luuid" ' ./cmd/leveldbServer/

# actual Docker image
FROM scratch

WORKDIR /
COPY --from=golang /go/bin/leveldbServer leveldbServer

# webserver ports
# if use-tls=false, only port 80 will be used
# if use-tls=true, port 80 will be used for ACME and HTTP, port 443 will be used for HTTPS

EXPOSE 1337

ENTRYPOINT ["./leveldbServer"]