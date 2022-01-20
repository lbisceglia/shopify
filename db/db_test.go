package db

import (
	"net/http"
	"testing"

	"github.com/lbisceglia/shopify/models"
)

type CreateResult struct {
	item      *models.Item
	toLoad    []models.Item
	code      int
	isError   bool
	itemCount int
}

type UpdateResult struct {
	item      *models.Item
	id        *models.ID
	want      models.Item
	toLoad    []models.Item
	code      int
	isError   bool
	itemCount int
}

type DeleteResult struct {
	id        *models.ID
	toLoad    []models.Item
	code      int
	isError   bool
	itemCount int
}

type GetItemResult struct {
	id        *models.ID
	toLoad    []models.Item
	code      int
	isError   bool
	itemCount int
}

var itemA = models.Item{
	ID:          "00000000000000000001",
	SKU:         "AAAAAAAA",
	Name:        "Thing1",
	Description: "First thing's first",
	PriceInCAD:  price(20.00),
	Quantity:    quantity(3),
}

func TestCreateItem(t *testing.T) {
	tests := map[string]CreateResult{
		"valid minimal": {
			item: &models.Item{
				SKU:      "01234567",
				Name:     "Thing1",
				Quantity: quantity(0),
			},
			toLoad:    nil,
			code:      http.StatusCreated,
			isError:   false,
			itemCount: 1,
		},
		"valid maximal": {
			item: &models.Item{
				ID:          "aaaaaaaaaaaaaaaaaaaa",
				SKU:         "01234567",
				Name:        "Thing1",
				Description: "First thing's first",
				PriceInCAD:  price(10.00),
				Quantity:    quantity(200),
			},
			toLoad:    nil,
			code:      http.StatusCreated,
			isError:   false,
			itemCount: 1,
		},
		"invalid duplicate sku": {
			item: &models.Item{
				SKU:      "01234567",
				Name:     "Thing2",
				Quantity: quantity(7),
			},
			toLoad: []models.Item{
				{
					SKU:      "01234567",
					Name:     "Thing1",
					Quantity: quantity(0),
				},
			},
			code:      http.StatusConflict,
			isError:   true,
			itemCount: 1,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			db, err := newTestDB()
			if err != nil {
				t.Fatalf(err.Error())
			}
			defer db.Close()
			db.LoadTestItems(test.toLoad)

			code, err := db.CreateItem(test.item)
			isError := err != nil
			if isError != test.isError {
				t.Errorf("got %v; want %v", err, test.isError)
			}
			if code != test.code {
				t.Errorf("got %v; want %v", code, test.code)
			}
			if !isError {
				if !test.item.IdIsPresent() {
					t.Fatal("id was not set")
				}
				if *test.item.LastUpdated != *test.item.DateAdded {
					t.Error("LastUpdated time does not match DateAdded time")
				}
			}

			items, _, _ := db.GetItems()
			if got, want := len(items), test.itemCount; got != want {
				t.Errorf("got %v; want %v", got, want)
			}
			db.clearTestDB()
			db.Close()
		})
	}
}

func TestUpdateItem(t *testing.T) {
	tests := map[string]UpdateResult{
		"valid minimal idempotent": {
			item: &models.Item{
				SKU:         "AAAAAAAA",
				Name:        "Thing1",
				Description: "First thing's first",
				PriceInCAD:  price(20.00),
				Quantity:    quantity(3),
			},
			id: id("00000000000000000001"),
			want: models.Item{
				ID:          "00000000000000000001",
				SKU:         "AAAAAAAA",
				Name:        "Thing1",
				Description: "First thing's first",
				PriceInCAD:  price(20.00),
				Quantity:    quantity(3),
			},
			toLoad:    []models.Item{itemA},
			code:      http.StatusNoContent,
			isError:   false,
			itemCount: 1,
		},
		"valid id idempotent": {
			item: &models.Item{
				ID:          "99999999999999999999",
				SKU:         "AAAAAAAA",
				Name:        "Thing1",
				Description: "First thing's first",
				PriceInCAD:  price(20.00),
				Quantity:    quantity(3),
			},
			id: id("00000000000000000001"),
			want: models.Item{
				ID:          "00000000000000000001",
				SKU:         "AAAAAAAA",
				Name:        "Thing1",
				Description: "First thing's first",
				PriceInCAD:  price(20.00),
				Quantity:    quantity(3),
			},
			toLoad:    []models.Item{itemA},
			code:      http.StatusNoContent,
			isError:   false,
			itemCount: 1,
		},
		"valid date_added idempotent": {
			item: &models.Item{
				SKU:         "AAAAAAAA",
				Name:        "Thing1",
				Description: "First thing's first",
				PriceInCAD:  price(20.00),
				Quantity:    quantity(3),
			},
			id: id("00000000000000000001"),
			want: models.Item{
				ID:          "00000000000000000001",
				SKU:         "AAAAAAAA",
				Name:        "Thing1",
				Description: "First thing's first",
				PriceInCAD:  price(20.00),
				Quantity:    quantity(3),
			},
			toLoad:    []models.Item{itemA},
			code:      http.StatusNoContent,
			isError:   false,
			itemCount: 1,
		},
		"valid last_updated idempotent": {
			item: &models.Item{
				SKU:         "AAAAAAAA",
				Name:        "Thing1",
				Description: "First thing's first",
				PriceInCAD:  price(20.00),
				Quantity:    quantity(3),
			},
			id: id("00000000000000000001"),
			want: models.Item{
				ID:          "00000000000000000001",
				SKU:         "AAAAAAAA",
				Name:        "Thing1",
				Description: "First thing's first",
				PriceInCAD:  price(20.00),
				Quantity:    quantity(3),
			},
			toLoad:    []models.Item{itemA},
			code:      http.StatusNoContent,
			isError:   false,
			itemCount: 1,
		},
		"valid sku": {
			item: &models.Item{
				SKU:         "01234567",
				Name:        "Thing1",
				Description: "First thing's first",
				PriceInCAD:  price(20.00),
				Quantity:    quantity(3),
			},
			id: id("00000000000000000001"),
			want: models.Item{
				ID:          "00000000000000000001",
				SKU:         "01234567",
				Name:        "Thing1",
				Description: "First thing's first",
				PriceInCAD:  price(20.00),
				Quantity:    quantity(3),
			},
			toLoad:    []models.Item{itemA},
			code:      http.StatusNoContent,
			isError:   false,
			itemCount: 1,
		},
		"invalid duplicate sku other": {
			item: &models.Item{
				SKU:         "BBBBBBBB",
				Name:        "Thing1",
				Description: "First thing's first",
				PriceInCAD:  price(20.00),
				Quantity:    quantity(3),
			},
			id: id("00000000000000000001"),
			want: models.Item{
				ID:          "00000000000000000001",
				SKU:         "AAAAAAAA",
				Name:        "Thing1",
				Description: "First thing's first",
				PriceInCAD:  price(20.00),
				Quantity:    quantity(3),
			},
			toLoad: []models.Item{
				itemA,
				{
					ID:       "00000000000000000002",
					SKU:      "BBBBBBBB",
					Name:     "Thing2",
					Quantity: quantity(0),
				},
			},
			code:      http.StatusConflict,
			isError:   true,
			itemCount: 2,
		},
		"valid Name": {
			item: &models.Item{
				SKU:         "AAAAAAAA",
				Name:        "Thingamabob",
				Description: "First thing's first",
				PriceInCAD:  price(20.00),
				Quantity:    quantity(3),
			},
			id: id("00000000000000000001"),
			want: models.Item{
				ID:          "00000000000000000001",
				SKU:         "AAAAAAAA",
				Name:        "Thingamabob",
				Description: "First thing's first",
				PriceInCAD:  price(20.00),
				Quantity:    quantity(3),
			},
			toLoad:    []models.Item{itemA},
			code:      http.StatusNoContent,
			isError:   false,
			itemCount: 1,
		},
		"valid Description": {
			item: &models.Item{
				SKU:         "AAAAAAAA",
				Name:        "Thing1",
				Description: "If you're not first you're last",
				PriceInCAD:  price(20.00),
				Quantity:    quantity(3),
			},
			id: id("00000000000000000001"),
			want: models.Item{
				ID:          "00000000000000000001",
				SKU:         "AAAAAAAA",
				Name:        "Thing1",
				Description: "If you're not first you're last",
				PriceInCAD:  price(20.00),
				Quantity:    quantity(3),
			},
			toLoad:    []models.Item{itemA},
			code:      http.StatusNoContent,
			isError:   false,
			itemCount: 1,
		},
		"valid overwrite Description": {
			item: &models.Item{
				SKU:        "AAAAAAAA",
				Name:       "Thing1",
				PriceInCAD: price(20.00),
				Quantity:   quantity(3),
			},
			id: id("00000000000000000001"),
			want: models.Item{
				ID:          "00000000000000000001",
				SKU:         "AAAAAAAA",
				Name:        "Thing1",
				Description: "",
				PriceInCAD:  price(20.00),
				Quantity:    quantity(3),
			},
			toLoad:    []models.Item{itemA},
			code:      http.StatusNoContent,
			isError:   false,
			itemCount: 1,
		},
		"valid Price": {
			item: &models.Item{
				SKU:         "AAAAAAAA",
				Name:        "Thing1",
				Description: "First thing's first",
				PriceInCAD:  price(1.00),
				Quantity:    quantity(3),
			},
			id: id("00000000000000000001"),
			want: models.Item{
				ID:          "00000000000000000001",
				SKU:         "AAAAAAAA",
				Name:        "Thing1",
				Description: "First thing's first",
				PriceInCAD:  price(1.00),
				Quantity:    quantity(3),
			},
			toLoad:    []models.Item{itemA},
			code:      http.StatusNoContent,
			isError:   false,
			itemCount: 1,
		},
		"valid overwrite Price": {
			item: &models.Item{
				SKU:         "AAAAAAAA",
				Name:        "Thing1",
				Description: "First thing's first",
				Quantity:    quantity(3),
			},
			id: id("00000000000000000001"),
			want: models.Item{
				ID:          "00000000000000000001",
				SKU:         "AAAAAAAA",
				Name:        "Thing1",
				Description: "First thing's first",
				PriceInCAD:  nil,
				Quantity:    quantity(3),
			},
			toLoad:    []models.Item{itemA},
			code:      http.StatusNoContent,
			isError:   false,
			itemCount: 1,
		},
		"valid Quantity": {
			item: &models.Item{
				SKU:         "AAAAAAAA",
				Name:        "Thing1",
				Description: "First thing's first",
				PriceInCAD:  price(20.00),
				Quantity:    quantity(9999),
			},
			id: id("00000000000000000001"),
			want: models.Item{
				ID:          "00000000000000000001",
				SKU:         "AAAAAAAA",
				Name:        "Thing1",
				Description: "First thing's first",
				PriceInCAD:  price(20.00),
				Quantity:    quantity(9999),
			},
			toLoad:    []models.Item{itemA},
			code:      http.StatusNoContent,
			isError:   false,
			itemCount: 1,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			db, err := newTestDB()
			if err != nil {
				t.Fatalf(err.Error())
			}
			defer db.Close()
			db.LoadTestItems(test.toLoad)

			code, err := db.UpdateItem(test.id, test.item)
			isError := err != nil
			if isError != test.isError {
				t.Errorf("got %v; want %v", err, test.isError)
			}
			if code != test.code {
				t.Errorf("got %v; want %v", code, test.code)
			}

			if code != http.StatusNotFound {
				got, _, err := db.GetItem(test.id)
				if err != nil {
					t.Fatal("GetItem not working, cannot fetch an item which exists")
				}
				if got, want := got, test.want; !itemsEqual(got, want) {
					t.Errorf("got %v; want %v", got, want)
				}
			}

			items, _, _ := db.GetItems()
			if got, want := len(items), test.itemCount; got != want {
				t.Errorf("got %v; want %v", got, want)
			}
			db.clearTestDB()
		})
	}
}

func TestDeleteItems(t *testing.T) {
	tests := map[string]DeleteResult{
		"valid delete": {
			toLoad:    []models.Item{itemA},
			id:        id("00000000000000000001"),
			code:      http.StatusNoContent,
			isError:   false,
			itemCount: 0,
		},
		"invalid delete": {
			toLoad:    []models.Item{itemA},
			id:        id("00000000000000000002"),
			code:      http.StatusNotFound,
			isError:   true,
			itemCount: 1,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			db, err := newTestDB()
			if err != nil {
				t.Fatalf(err.Error())
			}
			defer db.Close()
			db.LoadTestItems(test.toLoad)

			code, err := db.DeleteItem(test.id)
			isError := err != nil
			if isError != test.isError {
				t.Errorf("got %v; want %v", err, test.isError)
			}
			if code != test.code {
				t.Errorf("got %v; want %v", code, test.code)
			}

			// TODO: re-enable after GetItems implemented
			items, _, _ := db.GetItems()
			if got, want := len(items), test.itemCount; got != want {
				t.Errorf("got %v; want %v", got, want)
			}
			db.clearTestDB()
		})
	}
}

func TestGetItem(t *testing.T) {
	tests := map[string]GetItemResult{
		"valid get": {
			toLoad:    []models.Item{itemA},
			id:        id("00000000000000000001"),
			code:      http.StatusOK,
			isError:   false,
			itemCount: 1,
		},
		"invalid get": {
			toLoad:    []models.Item{itemA},
			id:        id("00000000000000000002"),
			code:      http.StatusNotFound,
			isError:   true,
			itemCount: 1,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			db, err := newTestDB()
			if err != nil {
				t.Fatalf(err.Error())
			}
			defer db.Close()
			db.LoadTestItems(test.toLoad)

			_, code, err := db.GetItem(test.id)
			isError := err != nil
			if isError != test.isError {
				t.Errorf("got %v; want %v", err, test.isError)
			}
			if code != test.code {
				t.Errorf("got %v; want %v", code, test.code)
			}

			items, _, _ := db.GetItems()
			if got, want := len(items), test.itemCount; got != want {
				t.Errorf("got %v; want %v", got, want)
			}
			db.clearTestDB()
		})
	}
}

func TestGetItems(t *testing.T) {
	tests := map[string]GetItemResult{
		"valid get empty": {
			toLoad:    []models.Item{},
			code:      http.StatusOK,
			isError:   false,
			itemCount: 0,
		},
		"valid get": {
			toLoad:    []models.Item{itemA},
			code:      http.StatusOK,
			isError:   false,
			itemCount: 1,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			db, err := newTestDB()
			if err != nil {
				t.Fatalf(err.Error())
			}
			defer db.Close()
			db.LoadTestItems(test.toLoad)

			items, code, err := db.GetItems()
			if isError := err != nil; isError != test.isError {
				t.Errorf("got %v; want %v", err, test.isError)
			}
			if code != test.code {
				t.Errorf("got %v; want %v", code, test.code)
			}
			if got, want := len(items), test.itemCount; got != want {
				t.Errorf("got %v; want %v", got, want)
			}
			db.clearTestDB()
		})
	}
}

func itemsEqual(item1 models.Item, item2 models.Item) bool {
	values := item1.ID == item2.ID &&
		item1.SKU == item2.SKU &&
		item1.Name == item2.Name &&
		item1.Description == item2.Description

	pointers := true
	if item1.PriceInCAD == nil {
		pointers = pointers && item2.PriceInCAD == nil
	} else if item2.PriceInCAD == nil {
		pointers = false
	} else {
		pointers = pointers && *item1.PriceInCAD == *item2.PriceInCAD
	}

	pointers = pointers && *item1.Quantity == *item2.Quantity

	return values && pointers
}

func price(p float64) *float64 {
	return &p
}

func quantity(q int) *int {
	return &q
}

func id(id models.ID) *models.ID {
	return &id
}
