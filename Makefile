include config/.env

.SILENT:

DB_URL=postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DATABASE)?sslmode=disable

tidy:
	go mod tidy
	go mod vendor

run:
	go run cmd/main.go

migrate:
	migrate create -ext sql -dir ./migration -seq $(name)

migrateup:
	@migrate -path ./migration -database "$(DB_URL)" -verbose up

migratedown:
	@migrate -path ./migration -database "$(DB_URL)" -verbose down