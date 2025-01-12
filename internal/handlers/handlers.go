package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/nats-io/stan.go"
	"io"
	"log"
	"net/http"
	"wb_nats/internal/schema"
)

type Service interface {
	SaveMessage(message schema.Request) error
	GetMessage(orderUid string) error
}

type Handlers struct {
	service Service
}

func NewHandlers(service Service) *Handlers {
	return &Handlers{service: service}
}

func (h *Handlers) SendMessage(w http.ResponseWriter, r *http.Request) {

	sc, err := stan.Connect("nats_wb", "wb_req")
	if err != nil {
		log.Println(err)
		return
	}
	defer sc.Close()

	var request schema.Request

	body, _ := io.ReadAll(r.Body)

	if err := json.Unmarshal(body, &request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		_, errOut := fmt.Fprintf(w, "{\"message\":\"%s\"}", err)
		if errOut != nil {
			log.Printf("POST /sendMessage out failed: %s", errOut.Error())
			return
		}
		return
	}
	log.Println(request)

	err = sc.Publish("wb", body)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Сообщение отправлено")

}

func (h *Handlers) GetMessageFromNats(sc stan.Conn) (stan.Subscription, error) {
	return sc.Subscribe("wb", func(m *stan.Msg) {
		log.Printf("Received a message: %s\n", string(m.Data))
		var message schema.Request
		json.Unmarshal(m.Data, &message)

		log.Println(message)

		err := h.service.SaveMessage(message)
		if err != nil {
			log.Fatal(err)
		}

	}, stan.StartWithLastReceived())

}

// TODO: check that message was processed
func (h *Handlers) isMessageProcessed(sequence uint64) bool {
	return false
}
