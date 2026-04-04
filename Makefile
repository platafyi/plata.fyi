.PHONY: dev-up dev-down dev-migrate dev-psql build-api build-frontend lint

# Load .env if it exists
-include .env
export

dev-up:
	docker compose up -d postgres
	@echo "Waiting for Postgres to be ready..."
	@until docker compose exec postgres pg_isready -U platafyi -d platafyi > /dev/null 2>&1; do sleep 1; done
	@echo "Postgres is ready."

dev-down:
	docker compose down

dev-migrate:
	cd backend && go run ./cmd/migrate

dev-psql:
	docker compose exec postgres psql -U platafyi -d platafyi

dev-api: dev-migrate
	cd backend && go run -mod=vendor ./cmd/api

dev-frontend:
	cd frontend && PORT=3000 npm run dev

build-api:
	cd backend && go build -o bin/api ./cmd/api

build-migrate:
	cd backend && go build -o bin/migrate ./cmd/migrate

lint:
	cd backend && go vet ./...
	cd frontend && npm run lint

tidy:
	cd backend && go mod tidy
	cd frontend && npm install
