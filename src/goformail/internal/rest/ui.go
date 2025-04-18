package rest

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
)

//go:embed out/*
var embedFS embed.FS

func (ctrl Controller) addUiHandler() {
	// Create a file server handler to serve the directory's contents
	uiFS, err := fs.Sub(embedFS, "out")
	if err != nil {
		log.Fatal(err)
	}
	fileServer := http.FileServerFS(uiFS)

	// Create a new HTTP handler to serve the file server
	ctrl.mux.Handle("/ui/", http.StripPrefix("/ui/", fileServer))
	ctrl.mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		http.Redirect(w, r, "/ui/", http.StatusFound)
	})
}
