package rest

import (
	"encoding/json"
	"gitlab.computing.dcu.ie/fonseca3/2025-csc1097-fonseca3-dagohos2/internal/model"
	"gitlab.computing.dcu.ie/fonseca3/2025-csc1097-fonseca3-dagohos2/internal/util"
	"net/http"
	"strconv"
)

func (ctrl Controller) addEmailHandlers() {
	ctrl.mux.HandleFunc("/api/emails/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/emails/" {
			handleUnknownMethod(w, r)
			return
		}
		switch r.Method {
		case "POST":
			ctrl.getEmails(w, r)
		default:
			handleUnknownMethod(w, r)
		}
	})

	ctrl.mux.HandleFunc("/api/emails/approve/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/emails/approve/" {
			handleUnknownMethod(w, r)
			return
		}
		switch r.Method {
		case "POST":
			ctrl.approveEmail(w, r)
		default:
			handleUnknownMethod(w, r)
		}
	})
}

func (ctrl Controller) getEmails(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	_, e := ctrl.auth.CheckTokenValidity(token)
	if e != nil {
		setErrorResponse(w, r, e)
		return
	}

	var reqs model.EmailReqs
	if err := json.NewDecoder(r.Body).Decode(&reqs); err != nil {
		setErrorResponse(w, r, util.NewInvalidObjectError("Invalid json provided", nil))
		return
	}

	resp, err := ctrl.db.GetAllEmails(&reqs)
	if err != nil {
		setErrorResponse(w, r, util.NewGenericError(err.Err))
		return
	}

	w.WriteHeader(http.StatusOK)
	setResponse("Successfully fetched emails!", resp, w, r)
}

func (ctrl Controller) approveEmail(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	reqId, e := ctrl.auth.CheckTokenValidity(token)
	if e != nil {
		setErrorResponse(w, r, e)
		return
	}

	id, err := strconv.Atoi(r.FormValue("id"))
	if err != nil {
		setErrorResponse(w, r, util.NewInvalidObjectError("Invalid id provided", err))
		return
	}

	listId, dbErr := ctrl.db.GetEmailList(id)
	if dbErr != nil {
		setErrorResponse(w, r, util.NewGenericError(dbErr.Err))
		return
	}

	_, e = ctrl.auth.CheckListMods(reqId, listId)
	if e != nil {
		setErrorResponse(w, r, e)
		return
	}

	dbErr = ctrl.db.SetEmailAsApproved(id)
	if dbErr != nil {
		setErrorResponse(w, r, util.NewGenericError(dbErr.Err))
		return
	}

	w.WriteHeader(http.StatusOK)
	setResponse("Successfully approved email!", IdObject{Id: id}, w, r)
}
