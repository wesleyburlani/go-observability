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
