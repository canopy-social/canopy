.PHONY: dev run build test lint migrate-up migrate-down docker-up docker-down sqlc clean

# Development
dev:
	go run ./cmd/server

worker:
	go run ./cmd/worker

# Build
build:
	go build -ldflags="-w -s" -o canopy-server ./cmd/server
	go build -ldflags="-w -s" -o canopy-worker ./cmd/worker

# Test
test:
	go test ./... -v -count=1

test-cover:
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html

# Lint
lint:
	golangci-lint run ./...

# Database migrations
migrate-up:
	migrate -path migrations -database "$(DATABASE_URL)" up

migrate-down:
	migrate -path migrations -database "$(DATABASE_URL)" down 1

migrate-create:
	@read -p "Migration name: " name; \
	migrate create -ext sql -dir migrations -seq $$name

# Docker (development infrastructure)
docker-up:
	docker compose up -d

docker-down:
	docker compose down

docker-logs:
	docker compose logs -f

# sqlc code generation
sqlc:
	sqlc generate

# Clean build artifacts
clean:
	rm -f canopy-server canopy-worker coverage.out coverage.html
