.PHONY: run build test docker-up docker-down clean

run:
	go run cmd/api/main.go

build:
	go build -o bin/api cmd/api/main.go

test:
	go test -v ./...

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

docker-logs:
	docker-compose logs -f app

clean:
	rm -rf bin/
	go clean

install:
	go mod download
	go mod tidy