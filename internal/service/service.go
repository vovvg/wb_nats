package service

import (
	"wb_nats/internal/schema"
)

type Storage interface {
	InsertMessage(message schema.Request) error
	GetMessage(orderId string) (schema.Request, error)
}

type Service struct {
	storage Storage
}

func NewService(storage Storage) *Service {
	return &Service{storage: storage}
}

func (s *Service) SaveMessage(message schema.Request) error {

	return s.storage.InsertMessage(message)
}

func (s *Service) GetMessage(orderUid string) (schema.Request, error) {

	return s.storage.GetMessage(orderUid)
}
