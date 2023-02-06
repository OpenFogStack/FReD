# building the binary
FROM golang:1.19-alpine as golang

LABEL maintainer="tp@mcc.tu-berlin.de"

WORKDIR /go/src/git.tu-berlin.de/mcc-fred/fred/

RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*
RUN update-ca-certificates

# Make an extra layer for the installed packages so that they dont have to be downloaded everytime
COPY go.mod .
COPY go.sum .

RUN go mod download

COPY cmd/frednode cmd/frednode
COPY pkg pkg
COPY proto proto

RUN CGO_ENABLED=0 go build -o fred ./cmd/frednode/

# actual Docker image
FROM ubuntu:18.04

WORKDIR /

COPY --from=golang /go/src/git.tu-berlin.de/mcc-fred/fred/fred fred

RUN apt-get update && \
    apt-get install iproute2 iputils-ping -y

ENTRYPOINT [ "./fred" ]

