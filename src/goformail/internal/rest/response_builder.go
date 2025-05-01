package rest

import (
	"encoding/json"
	"fmt"
	"gitlab.computing.dcu.ie/fonseca3/2025-csc1097-fonseca3-dagohos2/internal/util"
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

type LoginObject struct {
	Token string `json:"token"`
	User  int    `json:"user"`
}

type TokenObject struct {
	Token string `json:"token"`
}

func setResponse(msg string, data interface{}, w http.ResponseWriter, r *http.Request) {
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

func setErrorResponse(w http.ResponseWriter, r *http.Request, err *util.Error) {
	switch err.Code {
	case util.ErrInvalidObject:
		http.Error(w, err.Message, http.StatusBadRequest)
	case util.ErrUserAlreadyExists:
		http.Error(w, err.Message, http.StatusConflict)
	case util.ErrNoUser:
		http.Error(w, err.Message, http.StatusNotFound)
	case util.ErrListAlreadyExists:
		http.Error(w, err.Message, http.StatusConflict)
	case util.ErrNoList:
		http.Error(w, err.Message, http.StatusNotFound)
	case util.ErrNoEmail:
		http.Error(w, err.Message, http.StatusNotFound)
	case util.ErrIncorrectPassword:
		http.Error(w, err.Message, http.StatusUnauthorized)
	case util.ErrInvalidToken:
		http.Error(w, err.Message, http.StatusUnauthorized)
	case util.ErrNoPermission:
		http.Error(w, err.Message, http.StatusForbidden)
	default:
		http.Error(w, err.Message, http.StatusInternalServerError)
	}
}

func handleUnknownMethod(w http.ResponseWriter, r *http.Request) {
	msg := fmt.Sprintf("No known method %s for endpoint %s", r.Method, r.RequestURI)
	http.Error(w, msg, http.StatusNotFound)
}
