package models

import (
	"net/http"
	"testing"
	"time"
)

type GetIDResult struct {
	item Item
	want ID
}

type SetIDResult struct {
	item    Item
	toSet   ID
	isError bool
}

type ValidateResult struct {
	item    Item
	code    int
	isError bool
}

func TestGetID(t *testing.T) {
	tests := map[string]GetIDResult{
		"no id": {
			item: Item{},
			want: "",
		},
		"empty id": {
			item: Item{ID: ""},
			want: "",
		},
		"valid id": {
			item: Item{ID: "abcdefghijklmnopqrst"},
			want: "abcdefghijklmnopqrst",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			if id := test.item.GetID(); id != test.want {
				t.Errorf("got %v; want %v", id, test.want)
			}
		})
	}
}

func TestSetID(t *testing.T) {
	tests := map[string]SetIDResult{
		"valid id; not yet set": {
			item:    Item{},
			toSet:   "abcdefghijklmnopqrst",
			isError: false,
		},
		"invalid id; not yet set": {
			item:    Item{},
			toSet:   "abc",
			isError: true,
		},
		"invalid id; already set": {
			item:    Item{ID: "abcdefghijklmnopqrst"},
			toSet:   "abc",
			isError: true,
		},
		"valid id; already set": {
			item:    Item{ID: "abcdefghijklmnopqrst"},
			toSet:   "abcdefghijklmnopqrst",
			isError: true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			err := test.item.SetID(test.toSet)
			if isError := err != nil; isError != test.isError {
				t.Errorf("got %v; want %v", err, test.isError)
			}
		})
	}

}

func TestValidateID(t *testing.T) {
	tests := map[string]ValidateResult{
		"invalid no id": {
			item:    Item{},
			code:    http.StatusBadRequest,
			isError: true,
		},
		"invalid short id": {
			item:    Item{ID: "abcdefghijklmnopqrs"},
			code:    http.StatusBadRequest,
			isError: true,
		},
		"valid id all letters": {
			item:    Item{ID: "abcdefghijklmnopqrst"},
			code:    0,
			isError: false,
		},
		"valid id all numbers": {
			item:    Item{ID: "01234567890123456789"},
			code:    0,
			isError: false,
		},
		"valid id mixed": {
			item:    Item{ID: "01234abcde56789fghij"},
			code:    0,
			isError: false,
		},
		"invalid lowercase letter": {
			item:    Item{ID: "aaaaaaaaaaaaaaaaaaaw"},
			code:    http.StatusBadRequest,
			isError: true,
		},
		"invalid uppercase letter": {
			item:    Item{ID: "aaaaaaaaaaaaaaaaaaaA"},
			code:    http.StatusBadRequest,
			isError: true,
		},
		"invalid character": {
			item:    Item{ID: "0000000000000000000?"},
			code:    http.StatusBadRequest,
			isError: true,
		},
		"invalid long id": {
			item:    Item{ID: "abcdefghijklmnopqrstu"},
			code:    http.StatusBadRequest,
			isError: true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			code, err := test.item.ValidateID()
			if isError := err != nil; isError != test.isError {
				t.Errorf("got %v; want %v", err, test.isError)
			}
			if code != test.code {
				t.Errorf("got %v; want %v", code, test.code)
			}
		})
	}
}

func TestValidateSKU(t *testing.T) {
	tests := map[string]ValidateResult{
		"invalid no sku": {
			item:    Item{},
			code:    http.StatusBadRequest,
			isError: true,
		},
		"invalid short sku": {
			item:    Item{SKU: "ABC"},
			code:    http.StatusBadRequest,
			isError: true,
		},
		"valid sku minimal": {
			item:    Item{SKU: "A_-0"},
			code:    0,
			isError: false,
		},
		"valid sku only letters": {
			item:    Item{SKU: "ABCDefgh"},
			code:    0,
			isError: false,
		},
		"valid sku only numbers": {
			item:    Item{SKU: "01234567"},
			code:    0,
			isError: false,
		},
		"valid sku only hyphens": {
			item:    Item{SKU: "--------"},
			code:    0,
			isError: false,
		},
		"valid sku only underscores": {
			item:    Item{SKU: "________"},
			code:    0,
			isError: false,
		},
		"valid id maximal": {
			item:    Item{SKU: "Ab_0-12_Cd-3"},
			code:    0,
			isError: false,
		},
		"invalid character": {
			item:    Item{SKU: "AAAAAAA?"},
			code:    http.StatusBadRequest,
			isError: true,
		},
		"invalid long id": {
			item:    Item{SKU: "Ab_0-12345678"},
			code:    http.StatusBadRequest,
			isError: true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			code, err := test.item.ValidateSKU()
			if isError := err != nil; isError != test.isError {
				t.Errorf("got %v; want %v", err, test.isError)
			}
			if code != test.code {
				t.Errorf("got %v; want %v", code, test.code)
			}
		})
	}
}

func TestValidateName(t *testing.T) {
	tests := map[string]ValidateResult{
		"invalid no name": {
			item:    Item{Name: ""},
			code:    http.StatusBadRequest,
			isError: true,
		},
		"invalid whitespace name": {
			item: Item{Name: "    	"},
			code:    http.StatusBadRequest,
			isError: true,
		},
		"valid name": {
			item:    Item{Name: "Thingamajig"},
			code:    0,
			isError: false,
		},
		"valid name with spaces": {
			item: Item{Name: "  Thingamabob	"},
			code:    0,
			isError: false,
		},
		"valid name with internal spaces": {
			item:    Item{Name: "Amazing Crazy Doohickey!"},
			code:    0,
			isError: false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			code, err := test.item.ValidateName()
			if isError := err != nil; isError != test.isError {
				t.Errorf("got %v; want %v", err, test.isError)
			}
			if code != test.code {
				t.Errorf("got %v; want %v", code, test.code)
			}
		})
	}
}

func TestValidatePrice(t *testing.T) {
	testPricePositive := 15.0
	testPriceZero := 0.0
	testPriceNegative := -0.1

	tests := map[string]ValidateResult{
		"valid no price": {
			item:    Item{PriceInCAD: nil},
			code:    0,
			isError: false,
		},
		"valid price positive": {
			item:    Item{PriceInCAD: &testPricePositive},
			code:    0,
			isError: false,
		},
		"valid price zero": {
			item:    Item{PriceInCAD: &testPriceZero},
			code:    0,
			isError: false,
		},
		"invalid price negative": {
			item:    Item{PriceInCAD: &testPriceNegative},
			code:    http.StatusBadRequest,
			isError: true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			code, err := test.item.ValidatePrice()
			if isError := err != nil; isError != test.isError {
				t.Errorf("got %v; want %v", err, test.isError)
			}
			if code != test.code {
				t.Errorf("got %v; want %v", code, test.code)
			}
		})
	}
}

func TestValidateQuantity(t *testing.T) {
	testQuantityPositive := 5
	testQuantityZero := 0
	testQuantityNegative := -1

	tests := map[string]ValidateResult{
		"valid no quantity": {
			item:    Item{Quantity: nil},
			code:    0,
			isError: false,
		},
		"valid quantity positive": {
			item:    Item{Quantity: &testQuantityPositive},
			code:    0,
			isError: false,
		},
		"valid quantity zero": {
			item:    Item{Quantity: &testQuantityZero},
			code:    0,
			isError: false,
		},
		"invalid quantity negative": {
			item:    Item{Quantity: &testQuantityNegative},
			code:    http.StatusBadRequest,
			isError: true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			code, err := test.item.ValidateQuantity()
			if isError := err != nil; isError != test.isError {
				t.Errorf("got %v; want %v", err, test.isError)
			}
			if code != test.code {
				t.Errorf("got %v; want %v", code, test.code)
			}
		})
	}
}

func TestValidateItem(t *testing.T) {
	time := time.Date(2021, time.January, 10, 18, 38, 38, 500, time.UTC)
	testPriceZero := 0.00
	testPriceNegative := -0.01
	testQuantityZero := 0
	testQuantityNegative := -1

	tests := map[string]ValidateResult{
		"valid minimal": {
			item: Item{
				SKU:  "00000001",
				Name: "Thing1",
			},
			code:    0,
			isError: false,
		},
		"valid maximal": {
			item: Item{
				ID:          "abcdefghijklmnopqrst",
				SKU:         "00000001",
				Name:        "Thing1",
				Description: "The first thing",
				PriceInCAD:  &testPriceZero,
				Quantity:    &testQuantityZero,
				DateAdded:   &time,
				LastUpdated: &time,
			},
			code:    0,
			isError: false,
		},
		"missing sku": {
			item:    Item{Name: "Thing1"},
			code:    http.StatusBadRequest,
			isError: true,
		},
		"invalid sku": {
			item:    Item{SKU: "ABC"},
			code:    http.StatusBadRequest,
			isError: true,
		},
		"missing name": {
			item:    Item{SKU: "0123456789"},
			code:    http.StatusBadRequest,
			isError: true,
		},
		"invalid name": {
			item: Item{
				SKU:  "00000001",
				Name: "    ",
			},
			code:    http.StatusBadRequest,
			isError: true,
		},
		"invalid price": {
			item: Item{
				SKU:        "00000001",
				Name:       "Thing1",
				PriceInCAD: &testPriceNegative,
			},
			code:    http.StatusBadRequest,
			isError: true,
		},
		"invalid quantity": {
			item: Item{
				SKU:      "00000001",
				Name:     "Thing1",
				Quantity: &testQuantityNegative,
			},
			code:    http.StatusBadRequest,
			isError: true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			code, err := test.item.ValidateItem()
			if isError := err != nil; isError != test.isError {
				t.Errorf("got %v; want %v", err, test.isError)
			}
			if code != test.code {
				t.Errorf("got %v; want %v", code, test.code)
			}
		})
	}
}
