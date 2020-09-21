.PHONY: all client peering storage trigger

all: client peering storage trigger

client peering storage trigger: ## Compile all proto files
	protoc -I $@/ $@.proto --go_out=plugins=grpc:$@