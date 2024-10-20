build:
	@go build cmd/api/main.go

run:
	@go run cmd/api/main.go

docker-start:
	@docker-compose up -d

docker-stop:
	@docker-compose down
