PHONY: generate-structs
generate-structs:
	protoc --go_out=internal/gen --go_opt=paths=import \
			--go-grpc_out=internal/gen --go-grpc_opt=paths=import \
			api/proto/user/service.proto