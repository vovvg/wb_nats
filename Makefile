all: start-env migrate

start-env:
	docker-compose up -d

migrate:
	goose -dir migrations postgres "user=postgres dbname=wb_nats password=postgres sslmode=disable" up