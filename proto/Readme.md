# Installing GRPC
Theoretically all the generated files (`*.pb.go`) are already commited to this repo, so there is no need to install grpc.

If the files need to be changed / a new language binding needs to be created here are the install guidelines: https://grpc.io/docs/languages/go/quickstart/

On Linux, this should work:

```sh
sudo [YOUR_PACKAGE_MANAGER] install grpc protoc make
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

make
```