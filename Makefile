include .env
export

service-run:
	go run main.go

migrate-up:
	@migrate -path ./migrations -database "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):5433/$(DB_NAME)?sslmode=disable" up	

migrate-down:
	@migrate -path ./migrations -database "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):5433/$(DB_NAME)?sslmode=disable" down	
docker-upd:
	docker compose up -d
docker-up:
	docker compose up
docker-stop:
	docker compose stop
unit-test:
	go test ./internal/contacts/feature/service -v
e2e-test:
	go test ./internal/contacts/feature/transport -v -run TestCreateContacte2e