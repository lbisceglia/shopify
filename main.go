package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/lbisceglia/shopify/server"
)

const (
	GET    = http.MethodGet
	PUT    = http.MethodPut
	POST   = http.MethodPost
	DELETE = http.MethodDelete
)

func main() {
	r := mux.NewRouter().StrictSlash(true)
	s := server.NewServer()

	// Routes and Handlers
	r.HandleFunc("/api/items", s.CreateItem).Methods(POST)
	r.HandleFunc("/api/items/{id}", s.UpdateItem).Methods(PUT)
	r.HandleFunc("/api/items/{id}", s.DeleteItem).Methods(DELETE)
	r.HandleFunc("/api/items", s.GetItems).Methods(GET)
	r.HandleFunc("/api/items/{id}", s.GetItem).Methods(GET)

	// TODO: move port to environment var
	log.Fatal(http.ListenAndServe(":8081", r))
}
