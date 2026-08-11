include .env
export

service-run:
	go run main.go

migrate-up:
	@migrate -path ./migrations -database ${CONN_STRING} up	

migrate-down:
	@migrate -path ./migrations -database ${CONN_STRING} down	
docker-d:
	docker compose up -d
docker-up:
	docker compose up
docker-down:
	docker compose down