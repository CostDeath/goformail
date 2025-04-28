package rest

import (
	"encoding/json"
	"gitlab.computing.dcu.ie/fonseca3/2025-csc1097-fonseca3-dagohos2/internal/db"
	"gitlab.computing.dcu.ie/fonseca3/2025-csc1097-fonseca3-dagohos2/internal/model"
	"gitlab.computing.dcu.ie/fonseca3/2025-csc1097-fonseca3-dagohos2/internal/rest/util"
	"net/http"
	"strconv"
)

func (ctrl Controller) addListHandlers() {
	ctrl.mux.HandleFunc("/api/list/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/list/" {
			util.HandleError(w, r, util.InvalidMethod)
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
			util.HandleError(w, r, util.InvalidMethod)
		}
	})

	ctrl.mux.HandleFunc("/api/lists/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/lists/" {
			util.HandleError(w, r, util.InvalidMethod)
			return
		}
		switch r.Method {
		case "GET":
			ctrl.getLists(w, r)
		default:
			util.HandleError(w, r, util.InvalidMethod)
		}
	})
}

func (ctrl Controller) getList(w http.ResponseWriter, r *http.Request) {
	var id int
	id, err := strconv.Atoi(r.FormValue("id"))
	if err != nil {
		util.HandleError(w, r, util.ListNotFound)
		return
	}

	list, dbErr := ctrl.db.GetList(id)
	if dbErr != nil {
		switch dbErr.Code {
		case db.ErrNoRows:
			util.HandleError(w, r, util.ListNotFound)
		default:
			util.HandleError(w, r, util.DatabaseError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	setResponse("Successfully fetched list!", list, w, r)
}

func (ctrl Controller) postList(w http.ResponseWriter, r *http.Request) {
	var list model.List
	if err := json.NewDecoder(r.Body).Decode(&list); err != nil || !util.ValidateAllSet(list) {
		util.HandleError(w, r, util.InvalidPayload)
		return
	}

	id, err := ctrl.db.CreateList(&list)
	if err != nil {
		switch err.Code {
		case db.ErrDuplicate:
			util.HandleError(w, r, util.ListAlreadyExists)
		default:
			util.HandleError(w, r, util.DatabaseError)
		}
		return
	}

	w.WriteHeader(http.StatusCreated)
	setResponse("Successfully created list!", IdObject{Id: id}, w, r)
}

func (ctrl Controller) patchList(w http.ResponseWriter, r *http.Request) {
	var id int
	id, err := strconv.Atoi(r.FormValue("id"))
	if err != nil {
		util.HandleError(w, r, util.ListNotFound)
		return
	}
	var list model.List
	if err := json.NewDecoder(r.Body).Decode(&list); err != nil {
		util.HandleError(w, r, util.InvalidPayload)
		return
	}

	if err := ctrl.db.PatchList(id, &list); err != nil || r.Body == http.NoBody {
		switch err.Code {
		case db.ErrDuplicate:
			util.HandleError(w, r, util.ListAlreadyExists)
		case db.ErrNoRows:
			util.HandleError(w, r, util.ListNotFound)
		default:
			util.HandleError(w, r, util.DatabaseError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	setResponse("Successfully patched list!", IdObject{Id: id}, w, r)
}

func (ctrl Controller) deleteList(w http.ResponseWriter, r *http.Request) {
	var id int
	id, err := strconv.Atoi(r.FormValue("id"))
	if err != nil {
		util.HandleError(w, r, util.ListNotFound)
		return
	}

	if err := ctrl.db.DeleteList(id); err != nil {
		switch err.Code {
		case db.ErrNoRows:
			util.HandleError(w, r, util.ListNotFound)
		default:
			util.HandleError(w, r, util.DatabaseError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	setResponse("Successfully deleted list!", IdObject{Id: id}, w, r)
}

func (ctrl Controller) getLists(w http.ResponseWriter, r *http.Request) {
	lists, dbErr := ctrl.db.GetAllLists()
	if dbErr != nil {
		switch dbErr.Code {
		case db.ErrNoRows:
			util.HandleError(w, r, util.ListNotFound)
		default:
			util.HandleError(w, r, util.DatabaseError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	setResponse("Successfully fetched lists!", lists, w, r)
}
