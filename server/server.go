package server

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/lbisceglia/shopify/db"
	"github.com/lbisceglia/shopify/models"
)

// An InventoryServer responds to HTTP requests on the inventory.
// It supports to the following RESTful actions:
// - Create a new inventory item;
// - Update the data on an existing inventory item;
// - Permanently delete an existing inventory item;
// - Retrieve all items in inventory; and
// - Retrieve a single inventory item.
type InventoryServer interface {
	CreateItem(w http.ResponseWriter, r *http.Request)
	UpdateItem(w http.ResponseWriter, r *http.Request)
	DeleteItem(w http.ResponseWriter, r *http.Request)
	GetItems(w http.ResponseWriter, r *http.Request)
	GetItem(w http.ResponseWriter, r *http.Request)
}

// A Server is an implementation of an Inventory Server.
type Server struct {
	db db.DB
}

// NewServer creates a new instance of an Inventory Server.
func NewServer() InventoryServer {
	// TODO: change to real database
	db := db.NewMockDB()
	return newServer(db)
}

// newServer creates a new instance of an Inventory server with the specified database.
func newServer(db db.DB) InventoryServer {
	return &Server{
		db: db,
	}
}

// CreateItem creates an inventory Item according to the request.
// It ensures the request Item is well-formed in accordance with the API specification.
//
// Returns a 201 Created and responds with the relative URL of the newly-created resource
// (Header: Location) upon success.
// Returns a 400 Bad Request if the request is malformed.
// Returns a 409 Conflict if a non-unique SKU is provided.
func (s *Server) CreateItem(w http.ResponseWriter, r *http.Request) {
	s.setHeader(w)
	var item models.Item

	// Decode and validate the request
	if !s.decodeRequestItem(w, r.Body, &item) || !s.validateItem(w, &item) {
		return
	}

	// Save item to database
	code, err := s.db.CreateItem(&item)

	if err != nil {
		// Handle database errors
		writeError(w, code, err)
		return
	}

	// Respond with URL of newly-created resource
	relativeURL := fmt.Sprintf("/%s", item.GetID())
	w.Header().Set("Location", relativeURL)
	w.WriteHeader(code)
}

// UpdateItem updates an inventory Item according to the request.
// It ensures the request Item is well-formed in accordance with the API specification.
// It does not perform partial updates; any optional fields will be overwritten with
// their default values if they are missing from the request.
//
// Returns a 204 No Content on success.
// Returns a 400 Bad Request if the request is malformed.
// Returns a 404 Not Found if there is no resource corresponding to the URL endpoint.
// Returns a 409 Conflict if a non-unique SKU is provided as part of the update.
func (s *Server) UpdateItem(w http.ResponseWriter, r *http.Request) {
	s.setHeader(w)
	var item models.Item

	// Decode and validate the request
	if !s.decodeRequestItem(w, r.Body, &item) || !s.validateItem(w, &item) {
		return
	}

	// Update item in database
	id := models.ID(mux.Vars(r)["id"])
	code, err := s.db.UpdateItem(&id, &item)

	if err != nil {
		// Handle database errors
		writeError(w, code, err)
		return
	}

	w.WriteHeader(code)
}

// Delete Item permanently removes an item from inventory.
//
// Returns a 204 No Content on success.
// Returns a 404 Not Found if there is no resource corresponding to the URL endpoint.
func (s *Server) DeleteItem(w http.ResponseWriter, r *http.Request) {
	s.setHeader(w)

	// Delete item from database
	id := models.ID(mux.Vars(r)["id"])
	code, err := s.db.DeleteItem(&id)

	if err != nil {
		// Handle database errors
		writeError(w, code, err)
		return
	}

	w.WriteHeader(code)
}

// GetItems returns a collection of all Items in inventory.
//
// Returns all Items and a 200 OK on success.
func (s *Server) GetItems(w http.ResponseWriter, r *http.Request) {
	// TODO: paginate
	s.setHeader(w)

	// Get items from databse
	items, code, err := s.db.GetItems()

	if err != nil {
		// Handle database errors
		writeError(w, code, err)
		return
	}

	w.WriteHeader(code)

	// Respond with items
	if err := json.NewEncoder(w).Encode(items); err != nil {
		log.Println(err)
	}
}

// GetItem returns a single inventory Item
//
// Returns the Item and a 200 OK on success.
// Returns a 404 Not Found if there is no resource corresponding to the URL endpoint.
func (s *Server) GetItem(w http.ResponseWriter, r *http.Request) {
	s.setHeader(w)

	// Get item from database
	id := models.ID(mux.Vars(r)["id"])
	item, code, err := s.db.GetItem(&id)

	if err != nil {
		// Handle database errors
		writeError(w, code, err)
		return
	}

	w.WriteHeader(code)

	// Respond with items
	if err := json.NewEncoder(w).Encode(item); err != nil {
		log.Println(err)
	}
}

/*
  Helper Methods
*/

// setHeader sets the header's content type to application/json.
func (s *Server) setHeader(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
}

// writeError writes error states to the response.
// It assumes the error is not nil and will panic if passed a nil error.
func writeError(w http.ResponseWriter, code int, err error) {
	msg, _ := json.Marshal(err.Error())
	w.WriteHeader(code)
	w.Write(msg)
}

// decodeRequestItem decodes the json Item embedded in a Request and validates it for type errors.
// Returns true if decoded successfully, false otherwise.
func (s *Server) decodeRequestItem(w http.ResponseWriter, body io.ReadCloser, item *models.Item) bool {
	if err := json.NewDecoder(body).Decode(&item); err != nil {
		// Malformed request
		writeError(w, http.StatusBadRequest, err)
		return false
	}
	return true
}

// validateItem validates an Item embedded in a Request to ensure it adheres to API specification.
// Returns true if the Item is valid, false otherwise.
func (s *Server) validateItem(w http.ResponseWriter, item *models.Item) bool {
	if code, err := item.ValidateItem(); err != nil {
		// Invalid Item in request
		writeError(w, code, err)
		return false
	}
	return true
}
