include .env
export

service-run:
	go run main.go

migrate-up:
	@migrate -path ./migrations -database "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable" up	

migrate-down:
	@migrate -path ./migrations -database "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable" down	
docker-upd:
	docker compose up -d
docker-up:
	docker compose up
docker-stop:
	docker compose stop