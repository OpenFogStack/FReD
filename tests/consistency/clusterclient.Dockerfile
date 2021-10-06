# building the binary
FROM golang:1.17-alpine as golang

LABEL maintainer="tp@mcc.tu-berlin.de"

WORKDIR /go/src/git.tu-berlin.de/mcc-fred/fred/

RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*
RUN update-ca-certificates

# Make an extra layer for the installed packages so that they dont have to be downloaded everytime
COPY go.mod .
COPY go.sum .

RUN go mod download

COPY pkg pkg
COPY proto proto
COPY cmd/alexandra cmd/alexandra

RUN CGO_ENABLED=0 go build -o alexandra ./cmd/alexandra/

# actual Docker image
FROM python:3.9-slim-buster

WORKDIR /

COPY --from=golang /go/src/git.tu-berlin.de/mcc-fred/fred/alexandra alexandra
COPY tests/consistency/requirements.txt requirements.txt
RUN python3 -m pip install -r requirements.txt
COPY proto proto
COPY tests/consistency/client.sh client.sh
COPY tests/consistency/client.py client.py

RUN apt-get clean && \
    apt-get update && \
    apt-get install iproute2 iputils-ping -y && \
    apt-get clean

ENTRYPOINT [ "/bin/bash", "./client.sh" ]
