.PHONY: deps tidy run build test docker-up docker-down docker-seed docker-reset-db

deps:
	brew install dbmate
	go mod download

tidy:
	go mod tidy

run:
	go run main.go

build:
	go build -o bin/api main.go

test:
	go test ./...

docker-up:
	docker compose up -d --build

docker-down:
	docker compose down

docker-seed:
	docker compose exec -T db psql -U postgres -d machin8 < db/fixtures/01-sampledata.sql

docker-reset-db:
	docker compose down -v
