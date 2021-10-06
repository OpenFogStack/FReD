FROM golang:1.17-alpine

ENV PYTHONUNBUFFERED=1
RUN apk add --update --no-cache grpc make protoc python3

RUN python3 -m ensurepip && \
    python3 -m pip install --no-cache --upgrade wheel pip setuptools && \
    apk add --update --no-cache g++ && \
    python3 -m pip install --no-cache --upgrade grpcio grpcio-tools && \
    apk del g++ \
    python3 -m pip uninstall --no-cache wheel setuptool

RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

WORKDIR /proto

ENTRYPOINT [ "make" ]