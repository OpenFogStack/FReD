# Shamelessly stolen from the original dockerfile
# building the binary
FROM golang:1.14-alpine as golang

MAINTAINER Tobias Pfandzelter <tp@mcc.tu-berlin.de>

WORKDIR /go/src/gitlab.tu-berlin.de/mcc-fred/fred/

# Make an extra layer for the installed packages so that they dont have to be downloaded everytime
COPY go.mod .
COPY go.sum .

RUN go mod download

COPY pkg pkg
COPY cmd cmd

# Static build required so that we can safely copy the binary over.
RUN CGO_ENABLED=0 go install ./cmd/leveldbServer/

# actual Docker image
FROM scratch

WORKDIR /
COPY --from=golang /go/bin/leveldbServer leveldbServer

# webserver ports
# if use-tls=false, only port 80 will be used
# if use-tls=true, port 80 will be used for ACME and HTTP, port 443 will be used for HTTPS

EXPOSE 1337

ENTRYPOINT ["./leveldbServer"]