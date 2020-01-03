# building the binary
FROM golang:alpine as golang

MAINTAINER Tobias Pfandzelter <tp@mcc.tu-berlin.de>

WORKDIR /go/src/gitlab.tu-berlin.de/mcc-fred/fred/

COPY . .

# Static build required so that we can safely copy the binary over.
RUN CGO_ENABLED=0 GOOS=linux go install -a -tags netgo -ldflags '-w' ./cmd/frednode/main.go

# actual Docker image
FROM scratch

COPY --from=golang /go/bin/main frednode

EXPOSE 9001

ENTRYPOINT ["./frednode"]
