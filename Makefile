#!make
include .env

DATABASE_URL = "postgres://${DB_USER}:${DB_PASS}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable"

migration-create:
	migrate create -ext sql -dir ./migrations $(name)

migration-up:
	migrate -path "./migrations" -database ${DATABASE_URL} up

migration-down:
	migrate -path "./migrations" -database ${DATABASE_URL} down 1

run:
	air -c .air.toml