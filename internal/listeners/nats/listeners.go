package nats

import (
	"encoding/json"
	"github.com/nats-io/stan.go"
	"log"
	"wb_nats/internal/config"
	"wb_nats/internal/schema"
)

type Service interface {
	SaveMessage(message schema.Request) error
}

type Listeners struct {
	service Service
}

func NewListeners(service Service) *Listeners {
	return &Listeners{service: service}
}

func (l *Listeners) InitListener(cfg config.Config) stan.Conn {
	sc, err := stan.Connect(cfg.Nats.ClusterId, cfg.Nats.ClientId)
	if err != nil {
		log.Fatal(err)
	}

	return sc
}

func (l *Listeners) ListenMessageFromNats(sc stan.Conn) (stan.Subscription, error) {
	return sc.Subscribe("wb", func(m *stan.Msg) {
		log.Printf("Received a message: %s\n", string(m.Data))
		var message schema.Request
		json.Unmarshal(m.Data, &message)

		log.Println(message)

		err := l.service.SaveMessage(message)
		if err != nil {
			log.Println(err)
		}

		m.Ack()
		log.Println("Message successfully processed and acknowledged")

	}, stan.SetManualAckMode())

}
