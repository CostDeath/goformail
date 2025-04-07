package rest

import (
	"fmt"
	"net/http"
)

type ErrorReason int

const (
	InvalidMethod ErrorReason = iota
	InvalidPayload
	ListAlreadyExists
	ListNotFound
	DatabaseError
)

func handleError(w http.ResponseWriter, r *http.Request, e ErrorReason) {
	var msg string
	switch e {
	case InvalidMethod:
		msg = fmt.Sprintf("No known method %s for endpoint %s", r.Method, r.RequestURI)
		http.Error(w, msg, http.StatusNotFound)
	case InvalidPayload:
		http.Error(w, "Invalid payload", http.StatusNotFound)
	case ListAlreadyExists:
		http.Error(w, "A list with this name already exists", http.StatusConflict)
	case ListNotFound:
		http.Error(w, "List not found", http.StatusNotFound)
	case DatabaseError:
		http.Error(w, "An error occurred relating to the database", http.StatusInternalServerError)
	default:
		http.Error(w, "An unknown error occurred", http.StatusNotFound)
	}
}
