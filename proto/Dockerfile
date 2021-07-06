FROM golang:1.16-alpine

RUN apk add grpc make protoc
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

WORKDIR /proto

ENTRYPOINT [ "make" ]