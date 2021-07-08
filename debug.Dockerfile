# This is a compy from the dockerfile that adds support for debugging. It is used in the 3NodeTest to debug nodeB
# building the binary
FROM golang:1.15-buster as golang

LABEL maintainer="tp@mcc.tu-berlin.de"

WORKDIR /go/src/git.tu-berlin.de/mcc-fred/fred/

RUN apt update && apt install -y ca-certificates git && rm -rf /var/cache/apk/*
COPY nase/tls/ca.crt /usr/local/share/ca-certificates/ca.crt
RUN update-ca-certificates

RUN go get github.com/go-delve/delve/cmd/dlv

# Make an extra layer for the installed packages so that they dont have to be downloaded everytime
COPY go.mod .
COPY go.sum .

RUN go mod download

COPY cmd/frednode cmd/frednode
COPY pkg pkg
COPY proto proto

RUN CGO_ENABLED=0 go install -gcflags="all=-N -l" ./cmd/frednode/

# actual Docker image
FROM debian:buster

WORKDIR /

COPY --from=golang /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=golang /go/bin/frednode frednode
COPY --from=golang /go/bin/dlv dlv

EXPOSE 443
EXPOSE 5555
EXPOSE 40000

#ENTRYPOINT ["/dlv"]
ENTRYPOINT ["./dlv", "--listen=:40000", "--headless=true", "--api-version=2", "--accept-multiclient", "exec", "./frednode", "--"]