package service

import (
	"encoding/json"
	"fmt"
	"github.com/nats-io/stan.go"
	"io"
	"log"
	"net/http"
	schema "wb_nats/internal/shema"
)

func SendMessage(w http.ResponseWriter, r *http.Request) {

	sc, err := stan.Connect("nats_wb", "wb_req")
	if err != nil {
		log.Fatal(err)
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
	// Публикуем сообщение
	err = sc.Publish("wb", body)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Сообщение отправлено")

}

func GetMessages(w http.ResponseWriter, r *http.Request) {

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

}
