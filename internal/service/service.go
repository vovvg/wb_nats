package service

type Storage interface {
}

type Listener interface {
}

type Service struct {
	storage  Storage
	listener Listener
}

func NewService(storage Storage, listener Listener) *Service {
	return &Service{storage: storage, listener: listener}
}
