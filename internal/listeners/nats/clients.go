package nats

import (
	"github.com/nats-io/stan.go"
	"log"
	"wb_nats/internal/config"
)

type Clients struct {
	cfg config.Config
}

func NewClients(cfg config.Config) *Clients {
	return &Clients{cfg: cfg}
}

func (c *Clients) SendMessage(message []byte) {

	sc, err := stan.Connect(c.cfg.Nats.ClusterId, c.cfg.Nats.ProducerId)
	if err != nil {
		log.Println(err)
		return
	}
	defer sc.Close()

	err = sc.Publish(c.cfg.Nats.Subject, message)
	if err != nil {
		log.Println(err)
		return
	}

	log.Println("Message sent")

}
