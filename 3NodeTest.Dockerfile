#build stage
FROM golang:alpine AS builder
WORKDIR /go/src/gitlab.tu-berlin.de/mcc-fred/fred/
COPY . .
RUN apk add --no-cache git
RUN go get -d -v ./tests/3NodeTest
RUN go install -v ./tests/3NodeTest
#final stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /go/src/gitlab.tu-berlin.de/mcc-fred/fred/tests/3NodeTest /go/src/gitlab.tu-berlin.de/mcc-fred/fred/tests/3NodeTest
ENTRYPOINT ls -Rla ./go/src/gitlab.tu-berlin.de/mcc-fred/fred/tests/3NodeTest/