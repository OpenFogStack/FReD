FROM golang:1.19

ENV PYTHONUNBUFFERED=1
RUN apt-get update && apt-get install -y --no-install-recommends \
    build-essential \
    ca-certificates \
    curl \
    git \
    make \
    protobuf-compiler \
    python3 \
    python3-pip \
    python3-dev \
    && rm -rf /var/lib/apt/lists/*
# RUN apk add --update --no-cache grpc make protoc python3 python3-dev

RUN python3 -m pip install --no-cache --upgrade wheel pip setuptools && \
    python3 -m pip install --no-cache --upgrade grpcio grpcio-tools mypy-protobuf mypy

RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

WORKDIR /proto

ENTRYPOINT [ "make" ]