# building the binary
FROM golang:1.13-alpine as golang

MAINTAINER Tobias Pfandzelter <tp@mcc.tu-berlin.de>

RUN apk add --no-cache libzmq-static czmq-dev libsodium-static build-base util-linux-dev

RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2

WORKDIR /go/src/gitlab.tu-berlin.de/mcc-fred/fred/

COPY . .

# Static build required so that we can safely copy the binary over.
RUN ls
RUN touch ./cmd/frednode/dummy.cc
RUN go install -a -ldflags '-linkmode external -w -s -extldflags "-static -luuid" ' ./cmd/frednode/

# actual Docker image
FROM scratch

COPY --from=golang /go/bin/frednode frednode

EXPOSE 9001

EXPOSE 5555

ENTRYPOINT ["./frednode"]
