package service

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"wb_nats/internal/schema"
)

func (s *Service) GetMessages(w http.ResponseWriter, r *http.Request) {

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
