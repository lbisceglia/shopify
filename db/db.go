package db

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/lbisceglia/shopify/models"
	_ "github.com/lib/pq"
)

// A DB is a database for an inventory management CRUD application.
type DB interface {
	InitDB() error
	CreateItem(item *models.Item) (int, error)
	UpdateItem(id *models.ID, item *models.Item) (int, error)
	DeleteItem(id *models.ID) (int, error)
	GetItems() ([]models.Item, int, error)
	GetItem(id *models.ID) (models.Item, int, error)
	CreationTime() *time.Time
	UpdateTime(item *models.Item)
	LoadTestItems(items []models.Item)
	Close() error
}

// SQLDB is an implementation of a DB capable of managing inventory items.
// It uses a PostgreSQL database.
type SQLDB struct {
	db *sql.DB
}

// NewSQLDB creates a new PostgreSQL database with an active connection.
// It assumes that the caller will also call Close to end the connection.
// Returns a reference to the new DB and nil if the connection was successful,
// otherwise returns a reference to an empty DB and an error.
func NewSQLDB() (DB, error) {
	db := &SQLDB{}
	if err := db.InitDB(); err != nil {
		db.db = nil
		return db, err
	}
	return db, nil
}

// newTestDB creates a reference to the PostgreSQL testing database and
// removes all records to prepare it for a fresh test.
// It assumes that the caller will also call Close to end the connection.
// Returns a reference to the new DB and nil if the connection was successful,
// otherwise returns a reference to an empty DB and an error.
func newTestDB() (*SQLDB, error) {
	db := &SQLDB{}
	if err := db.initDB("postgres", "postgres", "localhost", "5432", "inventory_test"); err != nil {
		db.db = nil
		return db, err
	}
	if err := db.clearTestDB(); err != nil {
		db.db = nil
		return db, err
	}
	return db, nil
}

// clearTestDB removes all records from the database.
// It is only designed to be called on the test databse and should NEVER be called on a production database.
func (db *SQLDB) clearTestDB() error {
	if _, err := db.db.Query(`DELETE FROM items`); err != nil {
		return err
	}
	if _, err := db.db.Query(`DELETE FROM deleted_items`); err != nil {
		return err
	}
	return nil
}

// initDB initializes the database connection.
// It assumes that the caller will also call Close to end the connection.
func (db *SQLDB) initDB(user, password, host, port, dbname string) error {
	// connection string
	psqlconn := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable", user, password, host, port, dbname)

	// open database
	sqldb, err := sql.Open("postgres", psqlconn)
	if err != nil {
		return err
	}

	// check db
	if err := sqldb.Ping(); err != nil {
		return err
	}

	db.db = sqldb

	fmt.Println("server successfully connected to database")
	return nil
}

// InitDB connects the server to the database.
func (db *SQLDB) InitDB() error {

	// TODO: use os environment variables

	user := os.Getenv("DB_USERNAME")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	dbname := os.Getenv("DB_NAME")

	// user := "postgres"
	// password := "postgres"
	// host := "localhost"
	// port := "5432"
	// dbname := "inventory"

	return db.initDB(user, password, host, port, dbname)
}

// Close closes the databse connection so no more queries or statements may be sent to it.
func (db *SQLDB) Close() error {
	return db.db.Close()
}

// CreateItem writes a brand new Item to the database.
// Returns a 201 Created if successful or a 409 Conflict if the Item's SKU is not unique.
func (db *SQLDB) CreateItem(item *models.Item) (int, error) {
	sqlStmt := `
	INSERT into items (id, sku, name, description, price_cad, quantity, date_added, last_updated)
	VALUES($1, $2, $3, $4, $5, $6, now(), now());
	`

	var price interface{}
	if item.PriceInCAD == nil {
		price = nil
	} else {
		price = *item.PriceInCAD
	}

	// Complete item creation
	item.SetID(models.NewID())
	t := time.Now()
	item.DateAdded = &t
	item.LastUpdated = &t

	_, err := db.db.Exec(sqlStmt, item.ID, item.SKU, item.Name, item.Description, price, *item.Quantity)
	if err != nil {
		return http.StatusConflict, err
	}
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
func (db *SQLDB) UpdateItem(id *models.ID, item *models.Item) (int, error) {
	sqlStmt := `
	UPDATE items
	SET sku = $1, name = $2, description = $3, price_cad = $4, quantity = $5, last_updated = now()
	WHERE id = $6;
	`

	var price interface{}
	if item.PriceInCAD == nil {
		price = nil
	} else {
		price = *item.PriceInCAD
	}

	db.UpdateTime(item)

	res, err := db.db.Exec(sqlStmt, item.SKU, item.Name, item.Description, price, *item.Quantity, *id)
	if err != nil {
		return http.StatusConflict, err
	}
	if count, err := res.RowsAffected(); count == 0 {
		return http.StatusNotFound, fmt.Errorf("there is no item with ID %v", *id)
	} else if err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusNoContent, nil
}

// DeleteItem performs a 'hard delete' and permanently removes an item from the databse.
// Returns a 204 No Content if successful.
// Returns a 404 Not Found if there is no Item with the given ID in the database.
func (db *SQLDB) DeleteItem(id *models.ID) (int, error) {
	// TODO: change to soft delete
	sqlStmt := `DELETE FROM items WHERE id = $1;`

	if res, err := db.db.Exec(sqlStmt, *id); err == nil {
		if count, err := res.RowsAffected(); err == nil && count == 0 {
			return http.StatusNotFound, fmt.Errorf("there is no item with ID %v", *id)
		}
	}
	return http.StatusNoContent, nil
}

// GetItems returns a collection of all Items in the database.
// Returns all Items, a 200 OK, and nil if successful.
// Returns an empty slice of Items, 500 Internal Server Error, and an error if there is an error fetching the data.
func (db *SQLDB) GetItems() ([]models.Item, int, error) {
	sqlStmt := `SELECT * FROM items;`
	rows, err := db.db.Query(sqlStmt)

	if err != nil {
		return []models.Item{}, http.StatusInternalServerError, err
	}

	items := []models.Item{}
	for rows.Next() {
		item := models.Item{}

		if err := rows.Scan(&item.ID, &item.SKU, &item.Name, &item.Description, &item.PriceInCAD, &item.Quantity, &item.DateAdded, &item.LastUpdated); err != nil {
			return []models.Item{}, http.StatusInternalServerError, err
		}

		items = append(items, item)
	}
	return items, http.StatusOK, nil
}

// GetItem returns a single Item from the database.
// Returns the Item, a 200 OK, and nil if successful.
// Returns an empty Item, 404 Not Found, and an error if there is no Item with the given ID in the database.
// Returns an empty Item, 500 Internal Server Error and an error if there is an error fetching the data.
func (db *SQLDB) GetItem(id *models.ID) (models.Item, int, error) {
	sqlStmt := `SELECT * FROM items where id = $1;`
	rows, err := db.db.Query(sqlStmt, *id)

	if err != nil {
		return models.Item{}, http.StatusInternalServerError, err
	}

	item := models.Item{}
	i := 0
	for rows.Next() {
		if i >= 1 {
			return models.Item{}, http.StatusInternalServerError, fmt.Errorf("items are not unique by id")
		}

		if err := rows.Scan(&item.ID, &item.SKU, &item.Name, &item.Description, &item.PriceInCAD, &item.Quantity, &item.DateAdded, &item.LastUpdated); err != nil {
			return models.Item{}, http.StatusInternalServerError, err
		}
		i++
	}

	if i < 1 {
		return models.Item{}, http.StatusNotFound, fmt.Errorf("there is no item with ID %v", *id)
	}

	return item, http.StatusOK, nil
}

// CreationTime returns the time that an object was created.
// Encapsulates time creation logic for the purposes of unit testing.
// Returns the current time.
func (db *SQLDB) CreationTime() *time.Time {
	t := time.Now()
	return &t
}

// UpdateTime updates the LastUpdated time to reflect that an Item has just been updated.
// Encapsulates time updating logic for the purposes of unit testing.
// Updates the LastUpdated field to the current time.
func (db *SQLDB) UpdateTime(item *models.Item) {
	t := time.Now()
	item.LastUpdated = &t
}

// LoadTestItems loads the Items directly into the database.
// It assumes that all Items have been validated for correctness.
// This method bypasses CreateItem and should only be called during development,
// never in production code.
func (db *SQLDB) LoadTestItems(items []models.Item) {
	for i := range items {
		db.CreateItem(&items[i])
	}
}

/*
Mock Implementation
*/

// A MockDB is an in-memory mock database to be used during unit testing.
type MockDB struct {
	dbBySKU map[models.SKU]*models.Item
	dbByID  map[models.ID]*models.Item
}

// InitDB does nothing for the mock implementation.
func (db *MockDB) InitDB() error {
	return nil
}

// CreateItem writes a brand new Item to the database.
// Returns a 201 Created if successful or a 409 Conflict if the Item's SKU is not unique.
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
func (db *MockDB) GetItems() ([]models.Item, int, error) {
	items := make([]models.Item, len(db.dbBySKU))
	i := 0
	for _, v := range db.dbBySKU {
		items[i] = *v
		i++
	}
	return items, http.StatusOK, nil
}

// GetItem returns a single Item from the database.
// Returns the Item and a 200 OK if successful.
// Returns nil and a 404 Not Found if there is no Item with the given ID in the database.
func (db *MockDB) GetItem(id *models.ID) (models.Item, int, error) {
	if v, ok := db.dbByID[*id]; !ok {
		return models.Item{}, http.StatusNotFound, fmt.Errorf("there is no item with ID %v", *id)
	} else {
		return *v, http.StatusOK, nil
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

// Close closes the database connection. It does nothing in the mock implementation.
func (db *MockDB) Close() error {
	return nil
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
