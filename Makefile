gen:
	protoc ./products.proto --go-grpc_out=./Gateway --go_out=./Gateway
	protoc ./products.proto --go-grpc_out=./Product --go_out=./Product