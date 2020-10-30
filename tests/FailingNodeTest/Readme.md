Currently this test is expected to fail, because nodes cannot recover from missing updates

To run this test, go to the 3NodeTest and run `make failtest`

Command used to generate the pb files:
```bash
# ...Enter some virtualenv...
pip install grpcio-tools
cd ../../proto/client/
python -m  grpc_tools.protoc -I../client  --python_out=. --grpc_python_out=. ../client/client.proto
```
