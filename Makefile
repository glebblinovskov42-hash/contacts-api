include .env
export

service-run:
	go run main.go

migrate-up:
	migrate -path ./migrations -database "postgres://postgres:sos11982@localhost:5432/postgres?sslmode=disable" up	

migrate-down:
	migrate -path ./migrations -database "postgres://postgres:sos11982@localhost:5432/postgres?sslmode=disable" down	