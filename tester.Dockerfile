# Shamelessly stolen from the original dockerfile
# building the binary
FROM golang:1.15-alpine

MAINTAINER Tobias Pfandzelter <tp@mcc.tu-berlin.de>

WORKDIR /go/src/gitlab.tu-berlin.de/mcc-fred/fred/

# Make an extra layer for the installed packages so that they dont have to be downloaded everytime
COPY go.mod .
COPY go.sum .

RUN go mod download

COPY tests tests
COPY ext ext
COPY pkg pkg

RUN go install ./tests/3NodeTest/cmd/main/

ENTRYPOINT ["/go/bin/main"]