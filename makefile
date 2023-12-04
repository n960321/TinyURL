.PHONY: run build tidy

run:
	go run cmd/tiny-url/tinyurl.go

build:
	go build -v -o bin/tiny-url ./cmd/tiny-url

tidy:
	go mod tidy