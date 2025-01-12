package postgres

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"wb_nats/internal/schema"
)

type Storage struct {
	pool *pgxpool.Pool
}

func NewStorage(pool *pgxpool.Pool) *Storage {
	return &Storage{pool: pool}
}

func (s *Storage) InsertMessage(message schema.Request) error {
	_, err := s.pool.Exec(context.Background(),
		"INSERT INTO delivery (name, phone, zip, city, address, region, email) VALUES ($1, $2, $3, $4, $5, $6, $7)",
		message.Delivery.Name, message.Delivery.Phone, message.Delivery.Zip, message.Delivery.City, message.Delivery.Address, message.Delivery.Region, message.Delivery.Email)
	return err
}

func (s *Storage) GetMessage(orderUid string) error {
	return nil
}
