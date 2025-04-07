package rest

import (
	"fmt"
	"gitlab.computing.dcu.ie/fonseca3/2025-csc1097-fonseca3-dagohos2/internal/db"
	"log"
	"net/http"
	"strconv"
)

type Controller struct {
	configs map[string]string
	db      *db.Db
	mux     *http.ServeMux
}

func NewController(configs map[string]string, db *db.Db) *Controller {
	return &Controller{
		configs: configs,
		db:      db,
		mux:     http.DefaultServeMux,
	}
}

func (ctrl Controller) Serve() {
	ctrl.addListHandlers()
	ctrl.addUiHandler()

	// Start the server on port 8080
	port, _ := strconv.Atoi(ctrl.configs["HTTP_PORT"])
	fmt.Printf("Starting server at http://localhost:%d\n", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), ctrl.mux); err != nil {
		log.Fatalf("Error starting http server: %s\n", err)
	}
}
