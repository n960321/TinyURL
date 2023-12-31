.PHONY: run build tidy db-run db-remove db-migrate-up db-migrate-down k6-test docker-build docker-remove

db-container-id := $(shell docker ps | grep postgres | awk '{print $$1}')

redis-container-id := $(shell docker ps | grep redis | awk '{print $$1}')

tiny-url-container-id := $(shell docker ps | grep n960321/tiny-url | awk '{print $$1}')

postgresql_url := postgres://postgres:admin@localhost:5432/postgres?sslmode=disable

cur := $(shell pwd)

run:
	@go run cmd/tiny-url/tinyurl.go -config configs/dev.yaml -local true

build:
	@go build -v -o bin/tiny-url ./cmd/tiny-url

tidy:
	@go mod tidy

docker-build:
	@docker build --tag n960321/tiny-url:latest --file build/dockerfile .

docker-run:
	@docker run --name tiny-url \
	-p 8080:8080 \
	--link postgres:postgres \
	--link redis:redis \
	--volume $(cur)/configs:/app/configs \
	--volume $(cur)/deployment:/app/deployment \
	n960321/tiny-url

db-remove:
	docker rm -f $(db-container-id)

db-run:
	docker run -d --name postgres -p 5432:5432 -e POSTGRES_PASSWORD=admin postgres:15.5-alpine
	
db-migrate-up:
	migrate -database $(postgresql_url) -path deployment/database/migrations up

db-migrate-down:
	migrate -database $(postgresql_url) -path deployment/database/migrations down

cache-run:
	docker run -d --name redis -p 6379:6379 redis:7.2.3-alpine

cache-remove:
	docker rm -f $(redis-container-id)

prometheus-run:
	docker run -d --name prometheus -p 9090:9090 -v /Users/h_xian/Documents/playground/TinyURL/deployment/prometheus/prometheus.yml:/etc/prometheus/prometheus.yml prom/prometheus:v2.48.1