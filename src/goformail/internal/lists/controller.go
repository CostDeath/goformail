package lists

import (
	"encoding/json"
	"fmt"
	"gitlab.computing.dcu.ie/fonseca3/2025-csc1097-fonseca3-dagohos2/internal/lists/model"
	"gitlab.computing.dcu.ie/fonseca3/2025-csc1097-fonseca3-dagohos2/internal/util"
	"net/http"
)

type ErrorReason int

var configs map[string]string

const (
	InvalidMethod ErrorReason = iota
	InvalidPayload
)

func AddListHandlers(mux *http.ServeMux, cfg map[string]string) {
	configs = cfg
	mux.HandleFunc("/api/list/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/list/" {
			handleError(w, r, InvalidMethod)
			return
		}
		switch r.Method {
		case "POST":
			postList(w, r)
		case "PATCH":
			patchList(w, r)
		case "DELETE":
			deleteList(w, r)
		default:
			handleError(w, r, InvalidMethod)
		}
	})
}

func postList(w http.ResponseWriter, r *http.Request) {
	var list model.List
	if err := json.NewDecoder(r.Body).Decode(&list); err != nil || !util.ValidateAllSet(list) {
		handleError(w, r, InvalidPayload)
		return
	}

	err := createList()
	if err != nil {
		return
	}

	msg, err := json.Marshal(list)
	if err != nil {
		fmt.Println("2")
		return
	}

	setResponse(msg, http.StatusCreated, w, r)
}

func patchList(w http.ResponseWriter, r *http.Request) {
}

func deleteList(w http.ResponseWriter, r *http.Request) {

}

func handleError(w http.ResponseWriter, r *http.Request, e ErrorReason) {
	var msg string
	var code int
	switch e {
	case InvalidMethod:
		w.WriteHeader(http.StatusNotFound)
		msg = fmt.Sprintf("No known method %s for endpoint %s", r.Method, r.RequestURI)
	case InvalidPayload:
		w.WriteHeader(http.StatusBadRequest)
		msg = "Invalid payload"
	default:
		w.WriteHeader(http.StatusBadRequest)
		msg = "An unknown error occurred"
	}

	setResponse([]byte(fmt.Sprintf("{\"message\": \"%s\"}", msg)), code, w, r)
}

func setResponse(msg []byte, code int, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	_, err := w.Write(msg)
	if err != nil {
		fmt.Printf("Error setting message %s on request for %s\n", msg, r.RequestURI)
		return
	}
}
