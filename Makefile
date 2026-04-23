-include .env

DB_URL=postgres://myuser:mypassword@localhost:5432/sr-restored?sslmode=disable

.PHONY: db-up db-down migrate-up migrate-down run

db-up:
	docker-compose up -d

db-down:
	docker-compose down

migrate-up:
	migrate -path migrations -database "$(DB_URL)" -verbose up

migrate-down:
	migrate -path migrations -database "$(DB_URL)" -verbose down

run:
	DATABASE_URL=$(DB_URL) go run main.go
