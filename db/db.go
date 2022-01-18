package db

import (
	"time"

	"github.com/lbisceglia/shopify/models"
)

// A DB is a database for an inventory management CRUD application.
type DB interface {
	CreateItem(item *models.Item) (int, error)
	UpdateItem(id *models.ID, item *models.Item) (int, error)
	DeleteItem(id *models.ID) (int, error)
	GetItems() ([]*models.Item, int, error)
	GetItem(id *models.ID) (*models.Item, int, error)
	LoadTestItems(items []models.Item)
	CreationTime() *time.Time
	UpdateTime(item *models.Item)
}

// A MockDB is an in-memory mock database to be used during unit testing.
type MockDB struct{}

// CreateItem writes a brand new Item to the database.
// Returns a 201 Created if successful or a 400 Bad Request if the Item's SKU is not unique.
func (db *MockDB) CreateItem(item *models.Item) (int, error) {
	// TODO
	return 0, nil
}

// UpdateItem updates editable properties of an existing Item in the database.
// Editable properties are properties managed by the user;
// specifically, all properties aside from ID, DateAdded, and LastUpdated.
//
// SKUs may only be updated to a unique SKU that does not already exist in the database.
// Returns a 204 No Content if successful.
// Returns a 404 Not Found if there is no Item with the given ID in the database.
// Returns a 409 Conflict if the user attempts to change the SKU to something non-unique.
func (db *MockDB) UpdateItem(id *models.ID, item *models.Item) (int, error) {
	// TODO
	return 0, nil
}

// DeleteItem performs a 'hard delete' and permanently removes an item from the database.
// Returns a 204 No Content if successful.
// Returns a 404 Not Found if there is no Item with the given ID in the database.
func (db *MockDB) DeleteItem(id *models.ID) (int, error) {
	// TODO
	return 0, nil
}

// GetItems returns a collection of all Items in the database.
// The mock implementation of GetItems never fails.
// Returns all items and a 200 OK.
func (db *MockDB) GetItems() ([]*models.Item, int, error) {
	// TODO
	return []*models.Item{}, 0, nil
}

// GetItem returns a single Item from the database.
// Returns the Item and a 200 OK if successful.
// Returns nil and a 404 Not Found if there is no Item with the given ID in the database.
func (db *MockDB) GetItem(id *models.ID) (*models.Item, int, error) {
	// TODO
	return &models.Item{}, 0, nil
}

// CreationTime returns the time that an object was created.
// Encapsulates time creation logic for the purposes of unit testing.
// The mock implementation hard codes every creation date to 2000-01-01 00:00:00 +0000 UTC
func (db *MockDB) CreationTime() *time.Time {
	// TODO
	t := time.Now()
	return &t
}

// UpdateTime updates the LastUpdated time to reflect that an Item has just been updated.
// Encapsulates time updating logic for the purposes of unit testing.
// The mock implementation increments the LastUpdated field by one day each time it is called.
func (db *MockDB) UpdateTime(item *models.Item) {
	// TODO
}

// NewMockDB creates an in-memory mock database.
// It is designed for testing purposes and should not be used in production.
func NewMockDB() DB {
	// TODO
	return &MockDB{}
}

// LoadTestItems loads the Items directly into the database.
// It assumes that all Items have been validated for correctness.
// This method bypasses CreateItem and should only be called during testing,
// never in production code.
func (db *MockDB) LoadTestItems(items []models.Item) {
	// TODO
}
