package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/lbisceglia/shopify/db"
	"github.com/lbisceglia/shopify/server"
)

const (
	GET    = http.MethodGet
	PUT    = http.MethodPut
	POST   = http.MethodPost
	DELETE = http.MethodDelete
)

func main() {
	// Initialize Router
	r := mux.NewRouter().StrictSlash(true)

	// Initialize Database
	db, err := db.NewSQLDB()
	if err != nil {
		log.Fatal(err)
		return
	}
	defer db.Close()

	// Initialize Server
	s := server.NewServer(db)

	// Routes and Handlers
	r.HandleFunc("/api/items", s.CreateItem).Methods(POST)
	r.HandleFunc("/api/items/{id}", s.UpdateItem).Methods(PUT)
	r.HandleFunc("/api/items/{id}", s.DeleteItem).Methods(DELETE)
	r.HandleFunc("/api/items", s.GetItems).Methods(GET)
	r.HandleFunc("/api/items/{id}", s.GetItem).Methods(GET)

	// TODO: move port to environment var
	log.Fatal(http.ListenAndServe(":8081", r))
}
