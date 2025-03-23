all: start-env migrate

start-env:
	docker-compose up -d

migrate:
	goose -dir migrations postgres "postgres://postgres:postgrespw@localhost:5432/wb_nats?sslmode=disable" up