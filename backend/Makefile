include .env
export

# Database configuration
DB_URL=postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=$(POSTGRES_SSLMODE)

.PHONY: all run build test test-race test-coverage migrate-up migrate-down docker-up docker-down docker-build docker-logs fmt vet clean

all: build

run:
	go run cmd/api/main.go

build:
	go build -o tmp/api cmd/api/main.go

migrate-up:
	migrate -path migrations -database "$(DB_URL)" up

migrate-down:
	migrate -path migrations -database "$(DB_URL)" down

test:
	go test -v ./...

test-race:
	go test -v -race ./...

test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

docker-restart:
	docker-compose restart

docker-build:
	docker-compose build --no-cache

docker-logs:
	docker-compose logs -f

fmt:
	go fmt ./...

vet:
	go vet ./...

clean:
	rm -rf tmp/
	rm -f coverage.out
