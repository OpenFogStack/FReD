FROM golang:1.22-bullseye

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
    python3 -m pip install --no-cache --upgrade grpcio==1.51.1 grpcio-tools==1.51.1 mypy-protobuf==3.4.0 mypy==1.0.0

RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.1
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2

WORKDIR /proto

ENTRYPOINT [ "make" ]