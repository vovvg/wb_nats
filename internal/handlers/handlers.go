package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"wb_nats/internal/schema"
)

type Service interface {
	SaveMessage(message schema.Request) error
	GetMessage(orderUid string) (schema.Request, error)
}

type Client interface {
	SendMessage(message []byte)
}

type Handlers struct {
	service Service
	client  Client
}

func NewHandlers(service Service, client Client) *Handlers {
	return &Handlers{service: service, client: client}
}

func (h *Handlers) GetMessage(w http.ResponseWriter, r *http.Request) {
	orderId := r.PathValue("order_id")

	order, err := h.service.GetMessage(orderId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	resp, err := json.Marshal(order)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, errOut := fmt.Fprintf(w, string(resp))
	if errOut != nil {
		log.Printf("GET order/{order_id} out failed: %s", errOut.Error())

		return
	}
}

func (h *Handlers) SendTestMessage(w http.ResponseWriter, r *http.Request) {
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

	h.client.SendMessage(body)
}
