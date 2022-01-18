package server

import (
	"net/http"
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
type Server struct{}

// NewServer creates a new instance of an Inventory Server.
func NewServer() InventoryServer {
	return &Server{}
}

// CreateItem creates an inventory Item according to the request.
// It ensures the request Item is well-formed in accordance with the API specification.
//
// Returns a 201 Created and responds with the relative URL of the newly-created resource
// (Header: Location) upon success.
// Returns a 400 Bad Request if the request is malformed.
// Returns a 409 Conflict if a non-unique SKU is provided.
func (s *Server) CreateItem(w http.ResponseWriter, r *http.Request) {
	// TODO
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
	// TODO
}

// Delete Item permanently removes an item from inventory.
//
// Returns a 204 No Content on success.
// Returns a 404 Not Found if there is no resource corresponding to the URL endpoint.
func (s *Server) DeleteItem(w http.ResponseWriter, r *http.Request) {
	// TODO
}

// GetItems returns a collection of all Items in inventory.
//
// Returns all Items and a 200 OK on success.
func (s *Server) GetItems(w http.ResponseWriter, r *http.Request) {
	// TODO
}

// GetItem returns a single inventory Item
//
// Returns the Item and a 200 OK on success.
// Returns a 404 Not Found if there is no resource corresponding to the URL endpoint.
func (s *Server) GetItem(w http.ResponseWriter, r *http.Request) {
	// TODO
}
