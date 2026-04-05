BINARY=bin/server
MAIN=main.go
DB_URL=postgres://postgres:secret@localhost:5434/tournament_manager?sslmode=disable

.PHONY: run build test lint docker-up docker-down migrate-up migrate-down clean

run:
	go run $(MAIN)

build:
	go build -o $(BINARY) $(MAIN)

test:
	go test ./... -v -cover

lint:
	golangci-lint run ./...

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

migrate-up:
	migrate -path infra/db/migrations -database "$(DB_URL)" up

migrate-down:
	migrate -path infra/db/migrations -database "$(DB_URL)" down

clean:
	rm -rf bin/
