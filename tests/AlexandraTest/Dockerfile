# Shamelessly stolen from the original dockerfile
# building the binary
FROM golang:1.17-alpine

LABEL maintainer="tp@mcc.tu-berlin.de"

WORKDIR /go/src/git.tu-berlin.de/mcc-fred/fred/

RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*
RUN update-ca-certificates

# Make an extra layer for the installed packages so that they dont have to be downloaded everytime
COPY go.mod .
COPY go.sum .

RUN go mod download

COPY tests tests
COPY pkg pkg
COPY proto proto
COPY . .

RUN go install ./tests/AlexandraTest/cmd/main/

ENTRYPOINT ["/go/bin/main"]
