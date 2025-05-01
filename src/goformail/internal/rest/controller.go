package rest

import (
	"fmt"
	"gitlab.computing.dcu.ie/fonseca3/2025-csc1097-fonseca3-dagohos2/internal/db"
	"gitlab.computing.dcu.ie/fonseca3/2025-csc1097-fonseca3-dagohos2/internal/service"
	"log"
	"net/http"
	"strconv"
)

type Controller struct {
	list service.IListManager
	user service.IUserManager
	auth service.IAuthManager
	db   db.IDb
	mux  *http.ServeMux
}

func NewController(
	listMan *service.ListManager,
	userMan *service.UserManager,
	authMan *service.AuthManager,
	db *db.Db,
) *Controller {
	return &Controller{
		list: listMan,
		user: userMan,
		auth: authMan,
		db:   db,
		mux:  http.DefaultServeMux,
	}
}

func (ctrl Controller) Serve(portStr string) {
	ctrl.addListHandlers()
	ctrl.addUserHandlers()
	ctrl.addAuthHandlers()
	ctrl.addEmailHandlers()
	ctrl.addUiHandler()

	// Start the server on port 8080
	port, _ := strconv.Atoi(portStr)
	fmt.Printf("Starting server at http://localhost:%d\n", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), ctrl.mux); err != nil {
		log.Fatalf("Error starting http server: %s\n", err)
	}
}
