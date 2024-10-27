package postgres

import "github.com/jackc/pgx/v5/pgxpool"

type Storage struct {
	dbpool *pgxpool.Pool
}

func NewStorage(dbpool *pgxpool.Pool) *Storage {
	return &Storage{dbpool: dbpool}
}
