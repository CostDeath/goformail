package rest

import (
	"encoding/json"
	"gitlab.computing.dcu.ie/fonseca3/2025-csc1097-fonseca3-dagohos2/internal/model"
	"gitlab.computing.dcu.ie/fonseca3/2025-csc1097-fonseca3-dagohos2/internal/util"
	"net/http"
	"strings"
)

func (ctrl Controller) addAuthHandlers() {
	ctrl.mux.HandleFunc("/api/login/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/login/" {
			handleUnknownMethod(w, r)
			return
		}
		switch r.Method {
		case "POST":
			ctrl.login(w, r)
		default:
			handleUnknownMethod(w, r)
		}
	})

	ctrl.mux.HandleFunc("/api/validateToken/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/validateToken/" {
			handleUnknownMethod(w, r)
			return
		}
		switch r.Method {
		case "POST":
			ctrl.validate(w, r)
		default:
			handleUnknownMethod(w, r)
		}
	})
}

func (ctrl Controller) login(w http.ResponseWriter, r *http.Request) {
	var creds model.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		setErrorResponse(w, r, util.NewInvalidObjectError("Invalid json provided", nil))
		return
	}

	token, id, err := ctrl.auth.Login(&creds)
	if err != nil {
		setErrorResponse(w, r, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	setResponse("Successfully logged in!", LoginObject{Token: token, User: id}, w, r)
}

func (ctrl Controller) validate(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	_, err := ctrl.auth.CheckTokenValidity(token)
	if err != nil {
		setErrorResponse(w, r, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	token = strings.TrimPrefix(token, "Bearer ")
	setResponse("Token is valid!", TokenObject{Token: token}, w, r)
}
