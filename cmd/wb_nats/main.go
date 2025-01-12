package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"wb_nats/internal/handlers"
	"wb_nats/internal/service"
	"wb_nats/internal/storage/postgres"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nats-io/stan.go"
	"github.com/xlab/closer"
	"wb_nats/internal/config"
)

type Order struct {
	OrderUid string `json:"order_uid"`
}

func main() {
	cfg := config.MustLoad()

	dbPool, err := createDatabasePool(cfg)
	if err != nil {
		log.Fatal("failed to create database pool: %w", err)
	}

	storage := postgres.NewStorage(dbPool)
	services := service.NewService(storage)
	handler := handlers.NewHandlers(services)

	sc, err := stan.Connect(cfg.Nats.ClusterId, cfg.Nats.ClientId)
	if err != nil {
		log.Fatal(err)
	}
	defer sc.Close()

	subscription, _ := handler.GetMessageFromNats(sc)

	defer subscription.Close()

	if cfg.Env == "dev" {
		http.HandleFunc("POST /sendMessage", handler.SendMessage)
	}

	log.Println("Service started")
	if err := http.ListenAndServe("127.0.0.1:8080", nil); err != nil {
		panic(err)
	}
}

func createDatabasePool(cfg *config.Config) (*pgxpool.Pool, error) {

	dbpool, err := pgxpool.New(context.Background(), cfg.Database.Url)

	if err != nil {
		return nil, fmt.Errorf("error creating database pool: %w", err)
	}

	closer.Bind(dbpool.Close)

	return dbpool, nil
}
