package db

import (
	"fmt"
	"net/http"
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
type MockDB struct {
	dbBySKU map[models.SKU]*models.Item
	dbByID  map[models.ID]*models.Item
}

// CreateItem writes a brand new Item to the database.
// Returns a 201 Created if successful or a 400 Bad Request if the Item's SKU is not unique.
func (db *MockDB) CreateItem(item *models.Item) (int, error) {
	if _, ok := db.dbBySKU[item.SKU]; ok {
		return http.StatusConflict, fmt.Errorf("there is already an item with SKU %v", item.SKU)
	}

	// Complete item creation
	item.SetID(models.NewID())
	// Mock creation occurs at Jan 1, 2000
	t := db.CreationTime()
	item.DateAdded = t
	item.LastUpdated = t

	// Save item
	db.dbBySKU[item.SKU] = item
	db.dbByID[item.GetID()] = item
	return http.StatusCreated, nil
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
	if v, ok := db.dbByID[*id]; !ok {
		return http.StatusNotFound, fmt.Errorf("there is no item with id %v", item.GetID())
	} else {
		// Update the item with the new values
		if v.SKU != item.SKU {
			// SKU is to be updated, check for uniqueness
			if _, ok := db.dbBySKU[item.SKU]; ok {
				return http.StatusConflict, fmt.Errorf("there is already an item with SKU %v", item.SKU)
			}
			delete(db.dbBySKU, v.SKU)
			v.SKU = item.SKU
			db.dbBySKU[v.SKU] = v
		}

		v.Name = item.Name
		v.Description = item.Description
		v.PriceInCAD = item.PriceInCAD
		v.Quantity = item.Quantity

		db.UpdateTime(v)
		return http.StatusNoContent, nil
	}
}

// DeleteItem performs a 'hard delete' and permanently removes an item from the database.
// Returns a 204 No Content if successful.
// Returns a 404 Not Found if there is no Item with the given ID in the database.
func (db *MockDB) DeleteItem(id *models.ID) (int, error) {
	var sku *models.SKU
	if v, ok := db.dbByID[*id]; !ok {
		return http.StatusNotFound, fmt.Errorf("there is no item with ID %v", *id)
	} else {
		sku = &v.SKU
	}

	// Delete item
	delete(db.dbBySKU, *sku)
	delete(db.dbByID, *id)
	return http.StatusNoContent, nil
}

// GetItems returns a collection of all Items in the database.
// The mock implementation of GetItems never fails.
// Returns all items and a 200 OK.
func (db *MockDB) GetItems() ([]*models.Item, int, error) {
	items := make([]*models.Item, len(db.dbBySKU))
	i := 0
	for _, v := range db.dbBySKU {
		items[i] = v
		i++
	}
	return items, http.StatusOK, nil
}

// GetItem returns a single Item from the database.
// Returns the Item and a 200 OK if successful.
// Returns nil and a 404 Not Found if there is no Item with the given ID in the database.
func (db *MockDB) GetItem(id *models.ID) (*models.Item, int, error) {
	if v, ok := db.dbByID[*id]; !ok {
		return nil, http.StatusNotFound, fmt.Errorf("there is no item with ID %v", *id)
	} else {
		return v, http.StatusOK, nil
	}
}

// CreationTime returns the time that an object was created.
// Encapsulates time creation logic for the purposes of unit testing.
// The mock implementation hard codes every creation date to 2000-01-01 00:00:00 +0000 UTC
func (db *MockDB) CreationTime() *time.Time {
	t := time.Date(2000, time.January, 01, 00, 00, 00, 000, time.UTC)
	return &t
}

// UpdateTime updates the LastUpdated time to reflect that an Item has just been updated.
// Encapsulates time updating logic for the purposes of unit testing.
// The mock implementation increments the LastUpdated field by one day each time it is called.
func (db *MockDB) UpdateTime(item *models.Item) {
	if item.LastUpdated == nil {
		item.LastUpdated = db.CreationTime()
	} else {
		t := item.LastUpdated.AddDate(0, 0, 1)
		item.LastUpdated = &t
	}
}

// NewMockDB creates an in-memory mock database.
// It is designed for testing purposes and should not be used in production.
func NewMockDB() DB {
	return &MockDB{
		dbBySKU: make(map[models.SKU]*models.Item),
		dbByID:  make(map[models.ID]*models.Item),
	}
}

// LoadTestItems loads the Items directly into the database.
// It assumes that all Items have been validated for correctness.
// This method bypasses CreateItem and should only be called during testing,
// never in production code.
func (db *MockDB) LoadTestItems(items []models.Item) {
	for i := range items {
		db.dbByID[items[i].ID] = &items[i]
		db.dbBySKU[items[i].SKU] = &items[i]
	}
}
