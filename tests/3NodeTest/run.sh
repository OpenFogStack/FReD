#!/usr/bin/env bash

cd cmd
go build

# Binary is now in cmd called cmd
docker image build . -t fred:3nodetest