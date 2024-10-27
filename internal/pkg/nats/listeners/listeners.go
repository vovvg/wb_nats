package listeners

import "wb_nats/internal/config"

type Listener struct {
	cfg *config.Config
}

func NewListener(cfg *config.Config) *Listener {
	return &Listener{cfg: cfg}
}

func (l *Listener) getMessage() error {
	return nil
}
