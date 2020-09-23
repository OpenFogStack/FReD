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

RUN CGO_ENABLED=0 go install ./cmd/frednode/

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
