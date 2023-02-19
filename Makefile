migrate:
	go run ./cmd/cli/... migrate up

migrate-down:
	go run ./cmd/cli/... migrate down

seed:
	go run ./cmd/cli/... seed

test:
	go test ./application/cart/tests/... -v

swagger:
	swag init -g cmd/api/main.go

docker:
	docker-compose up -d

run:
	go run ./cmd/api/...

proto:
	protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative ./application/cart/delivery/rpc/codegen/cart.proto


