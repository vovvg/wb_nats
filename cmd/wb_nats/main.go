package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"wb_nats/internal/service"
	"wb_nats/internal/storage/postgres"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/stan.go"
	"github.com/xlab/closer"
	"wb_nats/internal/config"
)

type Delivery struct {
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Zip     string `json:"zip"`
	City    string `json:"city"`
	Address string `json:"address"`
	Region  string `json:"region"`
	Email   string `json:"email"`
}

func main() {
	cfg := config.MustLoad()

	dbPool, err := createDatabasePool(cfg)

	if err != nil {
		log.Fatal("failed to create database pool: %w", err)
	}

	//createNatsConnection(cfg)

	sc, err := stan.Connect("nats_wb", "wb_resp")
	if err != nil {
		log.Fatal(err)
	}
	defer sc.Close()

	if cfg.Env == "local" {
		http.HandleFunc("POST /sendMessage", service.SendMessage)
	}
	//mux.HandleFunc("GET /products/{id}/reviews", service.GetReviews)

	// Создаем подписку на канал
	_, err = sc.Subscribe("wb", func(m *stan.Msg) {
		log.Printf("Received a message: %s\n", string(m.Data))
		var message Delivery
		json.Unmarshal(m.Data, &message)

		err := insertDelivery(dbPool, message)
		if err != nil {
			log.Fatal(err)
		}

	}, stan.StartWithLastReceived())

	if err := http.ListenAndServe("127.0.0.1:8080", nil); err != nil {
		panic(err)
	}

	storage := postgres.NewStorage(dbPool)
	log.Println(storage)
}

func insertDelivery(conn *pgxpool.Pool, deliveryMessage Delivery) error {
	_, err := conn.Exec(context.Background(),
		"INSERT INTO public.delivery (name, phone, zip, city, address, region, email) VALUES ($1, $2, $3, $4, $5, $6, $7)",
		deliveryMessage.Name, deliveryMessage.Phone, deliveryMessage.Zip, deliveryMessage.City, deliveryMessage.Address, deliveryMessage.Region, deliveryMessage.Email)
	return err
}

func createDatabasePool(cfg *config.Config) (*pgxpool.Pool, error) {

	dbpool, err := pgxpool.New(context.Background(), cfg.Database.Url)

	if err != nil {
		return nil, fmt.Errorf("error creating database pool: %w", err)
	}

	closer.Bind(dbpool.Close)

	return dbpool, nil
}

func createNatsConnection(cfg *config.Config) {

	nc, err := nats.Connect(cfg.Nats.Url)
	if err != nil {
		log.Fatal("failed to create nats connection: %w", err)
	}
	defer nc.Drain()

}
