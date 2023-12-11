.PHONY: run build tidy db-run db-remove db-migrate-up db-migrate-down

db-container-id := $(shell docker ps | grep tiny-url-db | awk '{print $$1}')

redis-container-id := $(shell docker ps | grep tiny-url-redis | awk '{print $$1}')

postgresql_url := postgres://postgres:admin@localhost:5432/postgres?sslmode=disable

run:
	@go run cmd/tiny-url/tinyurl.go -config configs/dev.yaml -local true

build:
	go build -v -o bin/tiny-url ./cmd/tiny-url

tidy:
	go mod tidy

db-remove:
	docker rm -f $(db-container-id)

db-run:
	docker run -d --name tiny-url-db -p 5432:5432 -e POSTGRES_PASSWORD=admin postgres:15.5-alpine
	
db-migrate-up:
	migrate -database $(postgresql_url) -path deployment/database/migrations up

db-migrate-down:
	migrate -database $(postgresql_url) -path deployment/database/migrations down

cache-run:
	docker run -d --name tiny-url-redis -p 6379:6379 redis:7.2.3-alpine

cache-remove:
	docker rm -f $(redis-container-id)