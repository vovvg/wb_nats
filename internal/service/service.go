package service

import "wb_nats/internal/schema"

type Storage interface {
	InsertMessage(message schema.Request) error
	GetMessage(orderUid string) error
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

func (s *Service) GetMessage(orderUid string) error {

	return s.storage.GetMessage(orderUid)
}
