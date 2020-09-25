# Shamelessly stolen from the original dockerfile
# building the binary
FROM golang:1.15-alpine as golang

MAINTAINER Tobias Pfandzelter <tp@mcc.tu-berlin.de>

WORKDIR /go/src/gitlab.tu-berlin.de/mcc-fred/fred/

RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*
COPY nase/tls/ca.crt /usr/local/share/ca-certificates/ca.crt
RUN update-ca-certificates

# Make an extra layer for the installed packages so that they dont have to be downloaded everytime
COPY go.mod .
COPY go.sum .

RUN go mod download

COPY pkg pkg
COPY cmd cmd
COPY proto proto

# Static build required so that we can safely copy the binary over.
RUN CGO_ENABLED=0 go install ./cmd/storageserver/

# actual Docker image
FROM scratch

WORKDIR /
COPY --from=golang /go/bin/storageserver storageserver

EXPOSE 1337

ENTRYPOINT ["./storageserver"]