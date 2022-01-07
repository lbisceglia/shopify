package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/lbisceglia/shopify/server"
)

const (
	GET    = "GET"
	PUT    = "PUT"
	POST   = "POST"
	DELETE = "DELETE"
)

func main() {
	r := mux.NewRouter()

	// Routes and Handlers
	r.HandleFunc("/api/items", server.CreateItem).Methods(POST)
	r.HandleFunc("/api/items/{id}", server.UpdateItem).Methods(PUT)
	r.HandleFunc("/api/items/{id}", server.DeleteItem).Methods(DELETE)
	r.HandleFunc("/api/items", server.GetItems).Methods(GET)

	// TODO: move port to environment var
	log.Fatal(http.ListenAndServe(":8081", r))
}
