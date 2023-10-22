test:
	docker compose exec api go test -v ./...
build:
	docker compose exec tools go build -o bin/api cmd/api/main.go
clean:
	rm -Rf bin/*
migrations-up:
	docker compose exec tools ./scripts/migrations-up.sh
migrations-down:
	docker compose exec tools ./scripts/migrations-down.sh
generate-db-client:
	docker compose exec tools sqlc generate
generate-proto:
	docker compose exec tools protoc \
		--proto_path=./internal/transport/grpc/proto \
		--go_out=./internal/transport/grpc \
		--go-grpc_out=./internal/transport/grpc \
		./internal/transport/grpc/proto/*.proto
