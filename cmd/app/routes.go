package main

import (
	"github.com/gorilla/mux"
)

func (app *application) routes() *mux.Router {
	// Register handler functions.
	r := mux.NewRouter()
	r.HandleFunc("/api/bpcalc/", app.all).Methods("GET")
	r.HandleFunc("/api/bpcalc/{id}", app.findByID).Methods("GET")
	r.HandleFunc("/api/bpcalc/", app.insert).Methods("POST")
	r.HandleFunc("/api/bpcalc/{id}", app.delete).Methods("DELETE")
	// r.HandleFunc("/", app.all).Methods("GET")
	// r.HandleFunc("/", app.insert).Methods("POST")

	return r
}
