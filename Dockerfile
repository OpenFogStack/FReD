# building the binary
FROM golang:1.15-alpine as golang

MAINTAINER Tobias Pfandzelter <tp@mcc.tu-berlin.de>

# stolen from https://github.com/drone/ca-certs/blob/master/Dockerfile
RUN apk add -U --no-cache ca-certificates
WORKDIR /go/src/gitlab.tu-berlin.de/mcc-fred/fred/

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
