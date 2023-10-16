# execute all tests on the repository
test:
	ENV=test go test -v ./...
build:
	make -B generate-db-client && go build -o bin/api cmd/api/main.go
# deletes the contents of bin/ folder
clean:
	rm -Rf bin/*
migrations-up:
	./scripts/migrations-up.sh
migrations-down:
	./scripts/migrations-down.sh
generate-db-client:
	docker run --rm -v $$(pwd):/src -w /src kjconroy/sqlc generate
