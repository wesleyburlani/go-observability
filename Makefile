test:
	go test -v ./...
build:
	go build -o bin/api cmd/api/main.go
clean:
	rm -Rf bin/*
migrations-up:
	./scripts/migrations-up.sh
migrations-down:
	./scripts/migrations-down.sh
generate-db-client:
	sqlc generate
generate-proto:
	protoc \
		--proto_path=./internal/transport/grpc/proto \
		--go_out=./internal/transport/grpc \
		--go-grpc_out=./internal/transport/grpc \
		./internal/transport/grpc/proto/*.proto
grpc-client:
	evans ./internal/transport/grpc/proto/*.proto --host api --port 4000
