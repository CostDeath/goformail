package util

import (
	"encoding/json"
	"fmt"
	"gitlab.computing.dcu.ie/fonseca3/2025-csc1097-fonseca3-dagohos2/internal/rest/model"
	"log"
	"net/http"
)

func SetResponse(msg string, data interface{}, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	resp, err := json.Marshal(model.Response{Message: msg, Data: data})
	if err != nil {
		log.Printf("Error marshalling response on request for %s. Error: %s\n", r.RequestURI, err)
		return
	}

	if _, err := w.Write(resp); err != nil {
		fmt.Printf("Error setting message %s on request for %s. Error: %s\n", msg, r.RequestURI, err)
		return
	}
}
