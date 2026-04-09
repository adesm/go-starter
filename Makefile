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

migrate:
	go run cmd/migrate/main.go

# Module Generator
# Usage: make generate-module NAME=order
generate-module:
	@chmod +x scripts/genmodule.sh
	@./scripts/genmodule.sh $(NAME)