PHONY: generate-structs
generate-structs:
	protoc --go_out=pkg/user_v1 --go_opt=paths=import \
			--go-grpc_out=pkg/user_v1 --go-grpc_opt=paths=import \
			api/proto/user_v1/service.proto