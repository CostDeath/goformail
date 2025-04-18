package util

import (
	"encoding/json"
	"log"
	"net/http"
)

type response struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type IdObject struct {
	Id int `json:"id"`
}

func SetResponse(msg string, data interface{}, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	resp, err := json.Marshal(response{Message: msg, Data: data})
	if err != nil {
		log.Printf("Error marshalling response on request for %s. Error: %s\n", r.RequestURI, err)
		return
	}

	if _, err := w.Write(resp); err != nil {
		log.Printf("Error setting message %s on request for %s. Error: %s\n", msg, r.RequestURI, err)
		return
	}
}
