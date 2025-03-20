package interfaces

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"
)

//go:embed out/*
var embedFS embed.FS

func ServeHttp() {
	err := addWebUiHandler()

	if err != nil {
		log.Fatalf("Error adding UI handler server: %s\n", err)
	} else {
		// Start the server on port 8080
		port := 8080
		fmt.Printf("Starting server at http://localhost:%d\n", port)
		err = http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	}

	if err != nil {
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
