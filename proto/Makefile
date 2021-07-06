.PHONY: all client peering storage trigger middleware docker

all: client peering storage trigger middleware

client peering storage trigger middleware: ## Compile all proto files
	@protoc -I $@/ $@.proto --go_out=$@ --go_opt=paths=source_relative --go-grpc_out=$@ --go-grpc_opt=require_unimplemented_servers=false,paths=source_relative

docker: ## Compile proto files in Docker
	@docker build -t protoc .
	@docker run --rm -v $(CURDIR)/.:/proto protoc client