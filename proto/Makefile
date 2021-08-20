.PHONY: all client peering storage trigger middleware docker

all: client peering storage trigger middleware

client peering storage trigger: ## Compile all proto files for Go
	@protoc -I $@/ $@.proto --go_out=$@ --go_opt=paths=source_relative --go-grpc_out=$@ --go-grpc_opt=require_unimplemented_servers=false,paths=source_relative

middleware: ## Compile all proto files for Go and Python
	@protoc -I $@/ $@.proto --go_out=$@ --go_opt=paths=source_relative --go-grpc_out=$@ --go-grpc_opt=require_unimplemented_servers=false,paths=source_relative
	@python3 -m grpc_tools.protoc -I $@/ --python_out=$@ --grpc_python_out=$@ --mypy_out=$@ $@.proto

docker: ## Compile proto files in Docker
	@docker build -t protoc .
	@docker run --rm -v $(CURDIR)/.:/proto protoc all
