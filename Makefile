-include src/.env
export KAFKA_BROKER_ADDRESSES
export DATABASE_URL
export POSTGRES_USER
export POSTGRES_PASSWORD
export POSTGRES_DB
export GRPC_LISTEN_ADDRESS

# Pre-commit
install-pre-commit:
	pre-commit install

# gRPC code generation
proto-tools:
	go install github.com/bufbuild/buf/cmd/buf@latest
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

proto-gen:
	buf generate

# Database migrations (Atlas + Gorm loader)
ATLAS_BIN = atlas

# Atlas needs a throwaway "dev" database to compute the schema diff.
create-atlas-database:
	@docker exec postgres-go sh -c "createdb -U $(POSTGRES_USER) atlas_dev 2>/dev/null" || echo "Database atlas_dev already exists"

# Generate a new migration from the Gorm models: make makemigrations name=add_orders
makemigrations: create-atlas-database
	@go run ./src/core/database/atlas_loader > schema.sql
	@$(ATLAS_BIN) migrate diff $(name) --env local
	@if exist schema.sql del /f /q schema.sql

migrate-up: create-atlas-database
	@$(ATLAS_BIN) migrate apply --env local --allow-dirty

migrate-down:
	@$(ATLAS_BIN) migrate down --env local

migrate-status:
	@$(ATLAS_BIN) migrate status --env local

# Local development
deps:
	docker compose -p test-worker-service-local --env-file src/.env -f compose.dev.local.yml up -d

local-run: deps
	go run -race ./src

local-down:
	docker compose -p test-worker-service-local --env-file src/.env -f compose.dev.local.yml down

# Reset the local DB
db-reset:
	docker compose -p test-worker-service-local --env-file src/.env -f compose.dev.local.yml down -v
	docker compose -p test-worker-service-local --env-file src/.env -f compose.dev.local.yml up -d

# Start the full Docker stack:
docker-build-bin: export CGO_ENABLED=0
docker-build-bin: export GOOS=linux
docker-build-bin: export GOARCH=amd64
docker-build-bin:
	go build -o bin/app ./src

run: docker-build-bin
	docker compose -p test-service --env-file src/.env -f compose.dev.yml up -d --build

down:
	docker compose -p test-service --env-file src/.env -f compose.dev.yml down

# Lint / test
check:
	go mod tidy
	go fmt ./...
	go vet ./...
	go test -race -v ./...
