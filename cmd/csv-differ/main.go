package main

import (
	"log"
	"net/http"
	"os"

	"github.com/a-h/templ"
	"github.com/cybre/csv-differ/internal/handler"
	"github.com/cybre/csv-differ/internal/templates"
	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	r.Handle("/", templ.Handler(templates.Index()))
	r.Handle("/diff", handler.Diff())

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	server := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	log.Fatal(server.ListenAndServe())
}
