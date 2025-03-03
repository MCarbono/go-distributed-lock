export POSTGRES_URL="postgres://user:password@localhost:5432/order-postgres?sslmode=disable"

## create a new migration file. Usage `make migration-create name=<migration_name>`
.PHONY: migration-create
migration-create:
	@echo "Creating a new migration..."
	@go run github.com/golang-migrate/migrate/v4/cmd/migrate create -ext sql -dir postgres/migrations -seq $(name)

.ONESHELL:
.PHONY: migrate-up
migrate-up:
	go run --tags postgres github.com/golang-migrate/migrate/v4/cmd/migrate \
		-path=./postgres/migrations \
		-database=$(POSTGRES_URL) \
		-verbose up

.ONESHELL:
.PHONY: migrate-down
migrate-down:
	go run --tags postgres github.com/golang-migrate/migrate/v4/cmd/migrate \
		-path=./postgres/migrations \
		-database=$(POSTGRES_URL) \
		-verbose down