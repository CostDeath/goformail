package interfaces

import (
	"embed"
	"fmt"
	"gitlab.computing.dcu.ie/fonseca3/2025-csc1097-fonseca3-dagohos2/internal/lists"
	"io/fs"
	"log"
	"net/http"
	"strconv"
)

//go:embed out/*
var embedFS embed.FS

func ServeWeb(configs map[string]string) {
	lists.AddListHandlers(http.DefaultServeMux, configs)
	if err := addWebUiHandler(); err != nil {
		log.Fatalf("Error adding UI handler server: %s\n", err)
	}

	// Start the server on port 8080
	port, _ := strconv.Atoi(configs["HTTP_PORT"])
	fmt.Printf("Starting server at http://localhost:%d\n", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
		log.Fatalf("Error starting http server: %s\n", err)
	}
}

func addWebUiHandler() error {
	// Create a file server handler to serve the directory's contents
	uiFS, err := fs.Sub(embedFS, "out")
	fileServer := http.FileServerFS(uiFS)

	// Create a new HTTP handler to serve the UI
	http.Handle("/ui/", http.StripPrefix("/ui/", fileServer))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		http.Redirect(w, r, "/ui/", http.StatusFound)
	})

	return err
}
