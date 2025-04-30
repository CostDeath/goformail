package rest

import (
	"encoding/json"
	"gitlab.computing.dcu.ie/fonseca3/2025-csc1097-fonseca3-dagohos2/internal/model"
	"gitlab.computing.dcu.ie/fonseca3/2025-csc1097-fonseca3-dagohos2/internal/util"
	"net/http"
	"strconv"
)

func (ctrl Controller) addListHandlers() {
	ctrl.mux.HandleFunc("/api/list/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/list/" {
			handleUnknownMethod(w, r)
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
			handleUnknownMethod(w, r)
		}
	})

	ctrl.mux.HandleFunc("/api/lists/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/lists/" {
			handleUnknownMethod(w, r)
			return
		}
		switch r.Method {
		case "GET":
			ctrl.getLists(w, r)
		default:
			handleUnknownMethod(w, r)
		}
	})
}

func (ctrl Controller) getList(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	_, e := ctrl.auth.CheckTokenValidity(token)
	if e != nil {
		setErrorResponse(w, r, e)
		return
	}

	id, err := strconv.Atoi(r.FormValue("id"))
	if err != nil {
		setErrorResponse(w, r, util.NewInvalidObjectError("Invalid id provided", nil))
		return
	}

	list, e := ctrl.list.GetList(id)
	if e != nil {
		setErrorResponse(w, r, e)
		return
	}

	w.WriteHeader(http.StatusOK)
	setResponse("Successfully fetched list!", list, w, r)
}

func (ctrl Controller) postList(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	id, e := ctrl.auth.CheckTokenValidity(token)
	if e != nil {
		setErrorResponse(w, r, e)
		return
	}

	var list model.ListRequest
	if err := json.NewDecoder(r.Body).Decode(&list); err != nil {
		setErrorResponse(w, r, util.NewInvalidObjectError("Invalid json provided", nil))
		return
	}

	_, e = ctrl.auth.CheckPerms(id, "CRT_LIST")
	if e != nil {
		setErrorResponse(w, r, e)
		return
	}

	list.Mods = append(list.Mods, int64(id))
	id, e = ctrl.list.CreateList(&list)
	if e != nil {
		setErrorResponse(w, r, e)
		return
	}

	w.WriteHeader(http.StatusCreated)
	setResponse("Successfully created list!", IdObject{Id: id}, w, r)
}

func (ctrl Controller) patchList(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	reqId, e := ctrl.auth.CheckTokenValidity(token)
	if e != nil {
		setErrorResponse(w, r, e)
		return
	}

	id, err := strconv.Atoi(r.FormValue("id"))
	if err != nil {
		setErrorResponse(w, r, util.NewInvalidObjectError("Invalid id provided", nil))
		return
	}
	var list model.ListRequest
	if err := json.NewDecoder(r.Body).Decode(&list); err != nil {
		setErrorResponse(w, r, util.NewInvalidObjectError("Invalid json provided", nil))
		return
	}

	_, e = ctrl.auth.CheckListMods(reqId, id)
	if e != nil {
		setErrorResponse(w, r, e)
		return
	}

	e = ctrl.list.UpdateList(id, &list)
	if e != nil {
		setErrorResponse(w, r, e)
		return
	}

	w.WriteHeader(http.StatusOK)
	setResponse("Successfully patched list!", IdObject{Id: id}, w, r)
}

func (ctrl Controller) deleteList(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	reqId, e := ctrl.auth.CheckTokenValidity(token)
	if e != nil {
		setErrorResponse(w, r, e)
		return
	}

	id, err := strconv.Atoi(r.FormValue("id"))
	if err != nil {
		setErrorResponse(w, r, util.NewInvalidObjectError("Invalid id provided", nil))
		return
	}

	_, e = ctrl.auth.CheckListMods(reqId, id)
	if e != nil {
		setErrorResponse(w, r, e)
		return
	}

	e = ctrl.list.DeleteList(id)
	if e != nil {
		setErrorResponse(w, r, e)
		return
	}

	w.WriteHeader(http.StatusOK)
	setResponse("Successfully deleted list!", IdObject{Id: id}, w, r)
}

func (ctrl Controller) getLists(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	_, e := ctrl.auth.CheckTokenValidity(token)
	if e != nil {
		setErrorResponse(w, r, e)
		return
	}

	lists, e := ctrl.list.GetAllLists()
	if e != nil {
		setErrorResponse(w, r, e)
		return
	}

	w.WriteHeader(http.StatusOK)
	setResponse("Successfully fetched lists!", lists, w, r)
}
