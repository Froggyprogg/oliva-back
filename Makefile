PHONY: generate-structs
generate-structs:
	protoc --go_out=internal/gen --go_opt=paths=import \
			--go-grpc_out=internal/gen --go-grpc_opt=paths=import \
			api/proto/user/auth_srv.proto

PHONY: migrate-db
migrate-db:
	goose -dir migrations postgres "postgres://admin:admin@172.18.0.2:5432/oliva_auth?sslmode=disable" up

PHONY: new-migrate
new-migrate:
	goose -dir migrations create user_table sql

PHONY: conn-psql
conn-psql:
	sudo docker exec -it oliva-sso-db psql -U admin
