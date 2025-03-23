package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"wb_nats/internal/handlers"
	"wb_nats/internal/listeners/nats"
	"wb_nats/internal/service"
	"wb_nats/internal/storage/postgres"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/xlab/closer"
	"wb_nats/internal/config"
)

func main() {
	cfg := config.MustLoad()

	dbPool, err := createDatabasePool(cfg)
	if err != nil {
		log.Fatal("failed to create database pool: %w", err)
	}

	storage := postgres.NewStorage(dbPool)
	services := service.NewService(storage)
	listeners := nats.NewListeners(services)
	clients := nats.NewClients(*cfg)
	handler := handlers.NewHandlers(services, clients)

	sc := listeners.InitListener(*cfg)
	defer sc.Close()

	subscription, _ := listeners.ListenMessageFromNats(sc)
	defer subscription.Close()

	if cfg.Env == "dev" {
		http.HandleFunc("POST /sendMessage", handler.SendTestMessage)
	}

	http.HandleFunc("GET /getMessage/{order_id}", handler.GetMessage)

	log.Println("Service started")

	host := "127.0.0.1:" + cfg.Port
	if err := http.ListenAndServe(host, nil); err != nil {
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
