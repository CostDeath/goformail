package rest

import (
	"encoding/json"
	"gitlab.computing.dcu.ie/fonseca3/2025-csc1097-fonseca3-dagohos2/internal/db"
	"gitlab.computing.dcu.ie/fonseca3/2025-csc1097-fonseca3-dagohos2/internal/rest/model"
	"gitlab.computing.dcu.ie/fonseca3/2025-csc1097-fonseca3-dagohos2/internal/rest/util"
	"net/http"
	"strconv"
)

func (ctrl Controller) addListHandlers() {
	ctrl.mux.HandleFunc("/api/list/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/list/" {
			handleError(w, r, InvalidMethod)
			return
		}
		switch r.Method {
		case "GET":
			ctrl.getList(w, r)
		case "POST":
			ctrl.postList(w, r)
		case "PATCH":
			ctrl.patchList(w, r)
		case "DELETE":
			ctrl.deleteList(w, r)
		default:
			handleError(w, r, InvalidMethod)
		}
	})

	ctrl.mux.HandleFunc("/api/lists/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/lists/" {
			handleError(w, r, InvalidMethod)
			return
		}
		switch r.Method {
		case "GET":
			ctrl.getLists(w, r)
		default:
			handleError(w, r, InvalidMethod)
		}
	})
}

func (ctrl Controller) getList(w http.ResponseWriter, r *http.Request) {
	var id int
	id, err := strconv.Atoi(r.FormValue("id"))
	if err != nil {
		handleError(w, r, ListNotFound)
		return
	}

	list, dbErr := ctrl.db.GetList(id)
	if dbErr != nil {
		switch dbErr.Code {
		case db.ErrNoRows:
			handleError(w, r, ListNotFound)
		default:
			handleError(w, r, DatabaseError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	util.SetResponse("Successfully fetched list!", list, w, r)
}

func (ctrl Controller) postList(w http.ResponseWriter, r *http.Request) {
	var list model.List
	if err := json.NewDecoder(r.Body).Decode(&list); err != nil || !util.ValidateAllSet(list) {
		handleError(w, r, InvalidPayload)
		return
	}

	id, err := ctrl.db.CreateList(list.Name, list.Recipients)
	if err != nil {
		switch err.Code {
		case db.ErrDuplicate:
			handleError(w, r, ListAlreadyExists)
		default:
			handleError(w, r, DatabaseError)
		}
		return
	}

	w.WriteHeader(http.StatusCreated)
	util.SetResponse("Successfully created list!", model.IdObject{Id: id}, w, r)
}

func (ctrl Controller) patchList(w http.ResponseWriter, r *http.Request) {
	var id int
	id, err := strconv.Atoi(r.FormValue("id"))
	if err != nil {
		handleError(w, r, ListNotFound)
		return
	}
	var list model.List
	if err := json.NewDecoder(r.Body).Decode(&list); err != nil {
		handleError(w, r, InvalidPayload)
		return
	}

	if err := ctrl.db.PatchList(id, list.Name, list.Recipients); err != nil {
		switch err.Code {
		case db.ErrDuplicate:
			handleError(w, r, ListAlreadyExists)
		case db.ErrNoRows:
			handleError(w, r, ListNotFound)
		default:
			handleError(w, r, DatabaseError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	util.SetResponse("Successfully patched list!", model.IdObject{Id: id}, w, r)
}

func (ctrl Controller) deleteList(w http.ResponseWriter, r *http.Request) {
	var id int
	id, err := strconv.Atoi(r.FormValue("id"))
	if err != nil {
		handleError(w, r, ListNotFound)
		return
	}

	if err := ctrl.db.DeleteList(id); err != nil {
		switch err.Code {
		case db.ErrNoRows:
			handleError(w, r, ListNotFound)
		default:
			handleError(w, r, DatabaseError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	util.SetResponse("Successfully deleted list!", model.IdObject{Id: id}, w, r)
}

func (ctrl Controller) getLists(w http.ResponseWriter, r *http.Request) {
	lists, dbErr := ctrl.db.GetAllLists()
	if dbErr != nil {
		switch dbErr.Code {
		case db.ErrNoRows:
			handleError(w, r, ListNotFound)
		default:
			handleError(w, r, DatabaseError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	util.SetResponse("Successfully fetched lists!", lists, w, r)
}
