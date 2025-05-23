package rest

import (
	"encoding/json"
	"gitlab.computing.dcu.ie/fonseca3/2025-csc1097-fonseca3-dagohos2/internal/model"
	"gitlab.computing.dcu.ie/fonseca3/2025-csc1097-fonseca3-dagohos2/internal/util"
	"net/http"
	"strconv"
)

func (ctrl Controller) addUserHandlers() {
	ctrl.mux.HandleFunc("/api/user/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/user/" {
			handleUnknownMethod(w, r)
			return
		}
		switch r.Method {
		case "GET":
			ctrl.getUser(w, r)
		case "POST":
			ctrl.postUser(w, r)
		case "PATCH":
			ctrl.patchUser(w, r)
		case "DELETE":
			ctrl.deleteUser(w, r)
		default:
			handleUnknownMethod(w, r)
		}
	})

	ctrl.mux.HandleFunc("/api/users/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/users/" {
			handleUnknownMethod(w, r)
			return
		}
		switch r.Method {
		case "GET":
			ctrl.getUsers(w, r)
		default:
			handleUnknownMethod(w, r)
		}
	})
}

func (ctrl Controller) getUser(w http.ResponseWriter, r *http.Request) {
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

	user, e := ctrl.user.GetUser(id)
	if e != nil {
		setErrorResponse(w, r, e)
		return
	}

	w.WriteHeader(http.StatusOK)
	setResponse("Successfully fetched user!", user, w, r)
}

func (ctrl Controller) postUser(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	id, e := ctrl.auth.CheckTokenValidity(token)
	if e != nil {
		setErrorResponse(w, r, e)
		return
	}

	var user model.UserRequest
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		setErrorResponse(w, r, util.NewInvalidObjectError("Invalid json provided", nil))
		return
	}

	_, e = ctrl.auth.CheckUserPerms(id, "CRT_USER", user.Permissions)
	if e != nil {
		setErrorResponse(w, r, e)
		return
	}

	id, e = ctrl.user.CreateUser(&user)
	if e != nil {
		setErrorResponse(w, r, e)
		return
	}

	w.WriteHeader(http.StatusCreated)
	setResponse("Successfully created user!", IdObject{Id: id}, w, r)
}

func (ctrl Controller) patchUser(w http.ResponseWriter, r *http.Request) {
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
	var user model.UserRequest
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		setErrorResponse(w, r, util.NewInvalidObjectError("Invalid json provided", nil))
		return
	}
	
	if reqId != id || len(user.Permissions) != 0 {
		_, e = ctrl.auth.CheckUserPerms(reqId, "MOD_USER", user.Permissions)
		if e != nil {
			setErrorResponse(w, r, e)
			return
		}
	}

	e = ctrl.user.UpdateUser(id, &user)
	if e != nil {
		setErrorResponse(w, r, e)
		return
	}

	w.WriteHeader(http.StatusOK)
	setResponse("Successfully patched user!", IdObject{Id: id}, w, r)
}

func (ctrl Controller) deleteUser(w http.ResponseWriter, r *http.Request) {
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

	if reqId != id {
		_, e = ctrl.auth.CheckPerms(reqId, "MOD_USER")
		if e != nil {
			setErrorResponse(w, r, e)
			return
		}
	}

	e = ctrl.user.DeleteUser(id)
	if e != nil {
		setErrorResponse(w, r, e)
		return
	}

	w.WriteHeader(http.StatusOK)
	setResponse("Successfully deleted user!", IdObject{Id: id}, w, r)
}

func (ctrl Controller) getUsers(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	_, e := ctrl.auth.CheckTokenValidity(token)
	if e != nil {
		setErrorResponse(w, r, e)
		return
	}

	users, e := ctrl.user.GetAllUsers()
	if e != nil {
		setErrorResponse(w, r, e)
		return
	}

	w.WriteHeader(http.StatusOK)
	setResponse("Successfully fetched users!", users, w, r)
}
