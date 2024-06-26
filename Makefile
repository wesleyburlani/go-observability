test:
	make migrations-up && go test -v ./...
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
		--proto_path=./protofiles \
		--go_out=./internal/ports/grpc \
		--go-grpc_out=./internal/ports/grpc \
		./protofiles/*.proto
grpc-client:
	evans ./protofiles/*.proto --host api --port 4000
