.PHONY: all client peering storage trigger middleware

all: client peering storage trigger middleware

client peering storage trigger middleware: ## Compile all proto files
	protoc -I $@/ $@.proto --go_out=plugins=grpc:$@