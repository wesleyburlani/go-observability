test:
	docker compose exec api go test -v ./...
build:
	go build -o bin/api cmd/api/main.go
clean:
	rm -Rf bin/*
migrations-up:
	./scripts/migrations-up.sh
migrations-down:
	./scripts/migrations-down.sh
generate-db-client:
	docker run --rm -v $$(pwd):/src -w /src kjconroy/sqlc generate
generate-proto:
	docker run --rm -u $(id -u):$(id -g) -v${PWD}:${PWD} \
		-w${PWD} jaegertracing/protobuf:latest \
		--proto_path=${PWD} \
		--go_out=plugins=grpc:${PWD}/internal/transport/grpc \
		-I/usr/include/github.com/gogo/protobuf \
		${PWD}/internal/transport/grpc/proto/*.proto
