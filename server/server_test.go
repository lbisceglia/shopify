package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/lbisceglia/shopify/db"
	"github.com/lbisceglia/shopify/models"
)

const (
	GET     = http.MethodGet
	PUT     = http.MethodPut
	POST    = http.MethodPost
	DELETE  = http.MethodDelete
	rootURL = "/api/items"
)

func Router(s InventoryServer) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/api/items", s.CreateItem).Methods(POST)
	r.HandleFunc("/api/items/{id}", s.UpdateItem).Methods(PUT)
	r.HandleFunc("/api/items/{id}", s.DeleteItem).Methods(DELETE)
	r.HandleFunc("/api/items", s.GetItems).Methods(GET)
	r.HandleFunc("/api/items/{id}", s.GetItem).Methods(GET)
	return r
}

func Setup() *mux.Router {
	s := NewServer(db.NewMockDB())
	return Router(s)
}

func InitHTTP(method string, url string, bodyMap map[string]interface{}) (*http.Request, *httptest.ResponseRecorder) {
	body, _ := json.Marshal(bodyMap)
	req, _ := http.NewRequest(method, url, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	res := httptest.NewRecorder()
	return req, res
}
func TestGetItemsEmpty(t *testing.T) {
	r := Setup()

	// Get no items
	req, res := InitHTTP(GET, rootURL, nil)
	r.ServeHTTP(res, req)

	var items []models.Item
	if err := json.Unmarshal(res.Body.Bytes(), &items); err != nil {
		t.Fatal("Parse JSON Data Error")
	}

	// Check there are no items
	if got, want := res.Code, http.StatusOK; got != want {
		t.Errorf("got %v; want %v", got, want)
	}
	if len(items) != 0 {
		t.Error("expected an empty list of items")
	}
}

func TestGetItems(t *testing.T) {
	r := Setup()

	// Create the item
	bodyMap := map[string]interface{}{
		"sku":  "AAAAAAAA",
		"name": "Thing1",
	}

	req, res := InitHTTP(POST, rootURL, bodyMap)
	r.ServeHTTP(res, req)

	resp := res.Result()

	// Check the item was created successfully
	if got, want := resp.StatusCode, http.StatusCreated; got != want {
		t.Errorf("got %v; want %v", got, want)
	}

	header := res.Result().Header
	location := header.Values("Location")

	if location == nil || len(location) != 1 {
		t.Fatalf("got %v; want %v", len(location), 1)
	}

	// Get the item
	req, res = InitHTTP(GET, rootURL, nil)
	r.ServeHTTP(res, req)

	var items []models.Item
	if err := json.Unmarshal(res.Body.Bytes(), &items); err != nil {
		t.Fatal("Parse JSON Data Error")
	}

	// Check an item is returned and matches what is expected
	if got, want := res.Code, http.StatusOK; got != want {
		t.Errorf("got %v; want %v", got, want)
	}
	if len(items) != 1 {
		t.Fatal("expected one item to be returned")
	}

	item := items[0]
	id := models.ID(location[0][1:])
	if item.ID != id {
		t.Errorf(`expected item to have id "%s" matching its location`, id)
	}
	if item.SKU != "AAAAAAAA" {
		t.Errorf(`expected item to have sku "AAAAAAAA"; got %s`, item.SKU)
	}
	if item.Name != "Thing1" {
		t.Errorf(`expected item to have name "Thing1"; got %s`, item.Name)
	}
	if *item.Quantity != 0 {
		t.Errorf(`expected item to have quantity 0; got %d`, *item.Quantity)
	}
}

func TestCreateAndGetItem(t *testing.T) {
	r := Setup()

	// Create the item
	bodyMap := map[string]interface{}{
		"sku":  "AAAAAAAA",
		"name": "Thing1",
	}

	req, res := InitHTTP(POST, rootURL, bodyMap)
	r.ServeHTTP(res, req)

	// Check the item was created successfully
	if got, want := res.Code, http.StatusCreated; got != want {
		t.Errorf("got %v; want %v", got, want)
	}

	header := res.Result().Header
	location := header.Values("Location")

	if location == nil || len(location) != 1 {
		t.Fatalf("got %v; want %v", len(location), 1)
	}

	// Get the item
	req, res = InitHTTP(GET, rootURL+location[0], nil)
	r.ServeHTTP(res, req)

	var item models.Item
	bytes := res.Body.Bytes()
	if err := json.Unmarshal(bytes, &item); err != nil {
		t.Fatal("Parse JSON Data Error")
	}
	if got, want := res.Code, http.StatusOK; got != want {
		t.Errorf("got %v; want %v", got, want)
	}

	id := models.ID(location[0][1:])
	if item.ID != id {
		t.Errorf(`expected item to have id "%s" matching its location`, id)
	}
	if item.SKU != "AAAAAAAA" {
		t.Errorf(`expected item to have sku "AAAAAAAA"; got %s`, item.SKU)
	}
	if item.Name != "Thing1" {
		t.Errorf(`expected item to have name "Thing1"; got %s`, item.Name)
	}
	if *item.Quantity != 0 {
		t.Errorf(`expected item to have quantity 0; got %d`, *item.Quantity)
	}
}

func TestGetItemNotFound(t *testing.T) {
	// Get non-existent item at /api/items/00000000000000000000
	r := Setup()

	req, res := InitHTTP(GET, rootURL+"/00000000000000000000", nil)
	r.ServeHTTP(res, req)

	if got, want := res.Code, http.StatusNotFound; got != want {
		t.Errorf("got %v; want %v", got, want)
	}
}

func TestDeleteExistingItem(t *testing.T) {
	r := Setup()

	// Create the item
	bodyMap := map[string]interface{}{
		"sku":  "AAAAAAAA",
		"name": "Thing1",
	}

	req, res := InitHTTP(POST, rootURL, bodyMap)
	r.ServeHTTP(res, req)

	// Check the item was created successfully
	if got, want := res.Code, http.StatusCreated; got != want {
		t.Errorf("got %v; want %v", got, want)
	}

	header := res.Result().Header
	location := header.Values("Location")

	if location == nil || len(location) != 1 {
		t.Fatalf("got %v; want %v", len(location), 1)
	}

	// Delete the item
	req, res = InitHTTP(DELETE, rootURL+location[0], nil)
	r.ServeHTTP(res, req)

	// Check that the item was deleted successfully
	if got, want := res.Code, http.StatusNoContent; got != want {
		t.Errorf("got %v; want %v", got, want)
	}
}

func TestDeleteItemNotFound(t *testing.T) {
	r := Setup()

	// Delete the non-existent item at /api/items/00000000000000000000
	req, res := InitHTTP(DELETE, rootURL+"/00000000000000000000", nil)
	r.ServeHTTP(res, req)

	// Check that the item was deleted successfully
	if got, want := res.Code, http.StatusNotFound; got != want {
		t.Errorf("got %v; want %v", got, want)
	}
}

func TestCreateItemInvalid(t *testing.T) {
	r := Setup()

	// Attempt to create malformed items
	tests := map[string]map[string]interface{}{
		"no sku": {
			"name": "Thing1",
		},
		"short sku": {
			"sku":  "ABC",
			"name": "Thing1",
		},
		"long sku": {
			"sku":  "ZZZZZZZZZZZZZZZZZZZZ",
			"name": "Thing1",
		},
		"invalid character in sku": {
			"sku":  "AAAAAAA?",
			"name": "Thing1",
		},
		"no name": {
			"sku": "AAAAAAAA",
		},
		"empty name": {
			"sku":  "AAAAAAAA",
			"name": "",
		},
		"whitespace name": {
			"sku":  "AAAAAAAA",
			"name": "      ",
		},
		"negative price": {
			"sku":       "AAAAAAAA",
			"name":      "Thing1",
			"price_CAD": -0.01,
		},
		"negative quantity": {
			"sku":      "AAAAAAAA",
			"name":     "Thing1",
			"quantity": -1,
		},
		"float quantity": {
			"sku":      "AAAAAAAA",
			"name":     "Thing1",
			"quantity": 1.5,
		},
	}

	for name, bodyMap := range tests {
		t.Run(name, func(t *testing.T) {
			req, res := InitHTTP(POST, rootURL, bodyMap)
			r.ServeHTTP(res, req)

			// Check the item was rejected
			if got, want := res.Code, http.StatusBadRequest; got != want {
				t.Errorf("got %v; want %v", got, want)
			}
		})
	}
}

func TestCreateItemDuplicateSKU(t *testing.T) {
	r := Setup()

	// Create the item
	bodyMap := map[string]interface{}{
		"sku":  "AAAAAAAA",
		"name": "Thing1",
	}

	req, res := InitHTTP(POST, rootURL, bodyMap)
	r.ServeHTTP(res, req)

	// Check the item was created successfully
	if got, want := res.Code, http.StatusCreated; got != want {
		t.Errorf("got %v; want %v", got, want)
	}

	header := res.Result().Header
	location := header.Values("Location")

	if location == nil || len(location) != 1 {
		t.Fatalf("got %v; want %v", len(location), 1)
	}

	// Create the item again
	req, res = InitHTTP(POST, rootURL, bodyMap)
	r.ServeHTTP(res, req)

	// Check the item was rejected for being a duplicate
	if got, want := res.Code, http.StatusConflict; got != want {
		t.Errorf("got %v; want %v", got, want)
	}
}

func TestUpdateItem(t *testing.T) {
	r := Setup()

	// STEP 1
	// Create the item
	bodyMap := map[string]interface{}{
		"sku":         "AAAAAAAA",
		"name":        "Thing1",
		"description": "First thing's first",
		"price_CAD":   15.00,
		"quantity":    9,
	}

	req, res := InitHTTP(POST, rootURL, bodyMap)
	r.ServeHTTP(res, req)

	// Check the item was created successfully
	if got, want := res.Code, http.StatusCreated; got != want {
		t.Errorf("got %v; want %v", got, want)
	}

	header := res.Result().Header
	location := header.Values("Location")

	if location == nil || len(location) != 1 {
		t.Fatalf("got %v; want %v", len(location), 1)
	}

	// STEP 2
	// Get the item
	req, res = InitHTTP(GET, rootURL+location[0], nil)
	r.ServeHTTP(res, req)

	var item models.Item
	bytes := res.Body.Bytes()
	if err := json.Unmarshal(bytes, &item); err != nil {
		t.Fatal("Parse JSON Data Error")
	}
	if got, want := res.Code, http.StatusOK; got != want {
		t.Errorf("got %v; want %v", got, want)
	}

	// Ensure fields were successfully set prior to overwriting
	id := models.ID(location[0][1:])
	if item.ID != id {
		t.Errorf(`expected item to have id "%s" matching its location`, id)
	}
	if item.SKU != "AAAAAAAA" {
		t.Errorf(`expected item to have sku "AAAAAAAA"; got %s`, item.SKU)
	}
	if item.Name != "Thing1" {
		t.Errorf(`expected item to have name "Thing1"; got %s`, item.Name)
	}
	if item.Description != "First thing's first" {
		t.Errorf(`expected item to have description "First thing's first"; got %s`, item.Description)
	}
	if *item.PriceInCAD != 15.00 {
		t.Errorf(`expected item to have price 15.00; got %f`, *item.PriceInCAD)
	}
	if *item.Quantity != 9 {
		t.Errorf(`expected item to have quantity 9; got %d`, *item.Quantity)
	}

	// STEP 3
	// Update the item
	bodyMap = map[string]interface{}{
		"sku":  "BBBBBBBB",
		"name": "ThingOne",
	}

	req, res = InitHTTP(PUT, rootURL+location[0], bodyMap)
	r.ServeHTTP(res, req)

	// Check the item was updated successfully
	if got, want := res.Code, http.StatusNoContent; got != want {
		t.Errorf("got %v; want %v", got, want)
	}

	// Get the updated item
	req, res = InitHTTP(GET, rootURL+location[0], nil)
	r.ServeHTTP(res, req)

	item = models.Item{}
	bytes = res.Body.Bytes()
	if err := json.Unmarshal(bytes, &item); err != nil {
		t.Fatal("Parse JSON Data Error")
	}
	if got, want := res.Code, http.StatusOK; got != want {
		t.Errorf("got %v; want %v", got, want)
	}

	// Ensure fields were successfully updated
	if item.ID != id {
		t.Errorf(`expected item to have id "%s" matching its location`, id)
	}
	if item.SKU != "BBBBBBBB" {
		t.Errorf(`expected item to have sku "BBBBBBBB"; got %s`, item.SKU)
	}
	if item.Name != "ThingOne" {
		t.Errorf(`expected item to have name "ThingOne"; got %s`, item.Name)
	}
	if item.Description != "" {
		t.Errorf(`expected item to have no description"; got %s`, item.Description)
	}
	if item.PriceInCAD != nil {
		t.Errorf(`expected item to have no price; got %f`, *item.PriceInCAD)
	}
	if *item.Quantity != 0 {
		t.Errorf(`expected item to have quantity 0; got %d`, *item.Quantity)
	}
}

func TestUpdateItemNotFound(t *testing.T) {
	r := Setup()

	// Create the item
	bodyMap := map[string]interface{}{
		"sku":  "AAAAAAAA",
		"name": "Thing1",
	}

	// Update non-existent item at /api/items/00000000000000000000
	req, res := InitHTTP(PUT, rootURL+"/00000000000000000000", bodyMap)
	r.ServeHTTP(res, req)

	// Check the item was updated successfully
	if got, want := res.Code, http.StatusNotFound; got != want {
		t.Errorf("got %v; want %v", got, want)
	}
}

func TestUpdateItemSameSKU(t *testing.T) {
	r := Setup()

	// STEP 1
	// Create the item
	bodyMap := map[string]interface{}{
		"sku":  "AAAAAAAA",
		"name": "Thing1",
	}

	req, res := InitHTTP(POST, rootURL, bodyMap)
	r.ServeHTTP(res, req)

	// Check the item was created successfully
	if got, want := res.Code, http.StatusCreated; got != want {
		t.Errorf("got %v; want %v", got, want)
	}

	header := res.Result().Header
	location := header.Values("Location")

	if location == nil || len(location) != 1 {
		t.Fatalf("got %v; want %v", len(location), 1)
	}

	bodyMap = map[string]interface{}{
		"sku":  "AAAAAAAA",
		"name": "Same SKU, new Name",
	}

	// Make an idempotent update
	req, res = InitHTTP(PUT, rootURL+location[0], bodyMap)
	r.ServeHTTP(res, req)

	// Check the item was created successfully
	if got, want := res.Code, http.StatusNoContent; got != want {
		t.Errorf("got %v; want %v", got, want)
	}

	// Get the updated item
	req, res = InitHTTP(GET, rootURL+location[0], nil)
	r.ServeHTTP(res, req)

	item := models.Item{}
	bytes := res.Body.Bytes()
	if err := json.Unmarshal(bytes, &item); err != nil {
		t.Fatal("Parse JSON Data Error")
	}
	if got, want := res.Code, http.StatusOK; got != want {
		t.Errorf("got %v; want %v", got, want)
	}

	// Ensure fields were successfully updated
	if item.SKU != "AAAAAAAA" {
		t.Errorf(`expected item to have sku "AAAAAAAA"; got %s`, item.SKU)
	}
	if item.Name != "Same SKU, new Name" {
		t.Errorf(`expected item to have name "Same SKU, new Name"; got %s`, item.Name)
	}
}

func TestUpdateItemDuplicateSKU(t *testing.T) {
	r := Setup()

	// STEP 1
	// Create the first item
	bodyMap1 := map[string]interface{}{
		"sku":  "AAAAAAAA",
		"name": "Thing1",
	}

	req, res := InitHTTP(POST, rootURL, bodyMap1)
	r.ServeHTTP(res, req)

	// Check the item was created successfully
	if got, want := res.Code, http.StatusCreated; got != want {
		t.Errorf("got %v; want %v", got, want)
	}

	header1 := res.Result().Header
	location1 := header1.Values("Location")

	if location1 == nil || len(location1) != 1 {
		t.Fatalf("got %v; want %v", len(location1), 1)
	}

	// STEP 2
	// Create the second item
	bodyMap2 := map[string]interface{}{
		"sku":  "BBBBBBBB",
		"name": "Thing2",
	}

	req, res = InitHTTP(POST, rootURL, bodyMap2)
	r.ServeHTTP(res, req)

	// Check the item was created successfully
	if got, want := res.Code, http.StatusCreated; got != want {
		t.Errorf("got %v; want %v", got, want)
	}

	header2 := res.Result().Header
	location2 := header2.Values("Location")

	if location2 == nil || len(location2) != 1 {
		t.Fatalf("got %v; want %v", len(location2), 1)
	}

	// Update item 1 SKU to item 2's SKU
	req, res = InitHTTP(PUT, rootURL+location1[0], bodyMap2)
	r.ServeHTTP(res, req)

	// Check the item was created successfully
	if got, want := res.Code, http.StatusConflict; got != want {
		t.Errorf("got %v; want %v", got, want)
	}
}
