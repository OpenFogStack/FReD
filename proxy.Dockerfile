# Shamelessly stolen from the original dockerfile
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

COPY cmd/fredproxy cmd/fredproxy
COPY pkg pkg
COPY proto proto

# Static build required so that we can safely copy the binary over.
RUN CGO_ENABLED=0 go install ./cmd/fredproxy/

# actual Docker image
FROM scratch

WORKDIR /
COPY --from=golang /go/bin/fredproxy fredproxy

EXPOSE 5555
EXPOSE 9001

ENV PATH=.
ENTRYPOINT ["fredproxy"]
