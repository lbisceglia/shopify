package models

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
	"unicode"

	"github.com/rs/xid"
)

const (
	SKU_MIN_LEN = 4
	SKU_MAX_LEN = 12
	ID_LEN      = 20 // tied to xid specification
)

// An ID is a globally-unique identifier for an Item.
// It is allocated for indexing purposes and for use with a database.
// IDs are immutable. An Item maintains the same ID throughout its life.
// It must be 20 characters long and contain only the lowercase letters a-v and digits 0-9.
type ID string

// NewID creates a new, globally-unique ID.
func NewID() ID {
	return ID(xid.New().String())
}

// A SKU is a unique identifier for an Item.
// It is more human-friendly than ID and is allocated for internal use.
// An Item's SKU may be updated over its life, but must always remain unique.
// It may be 4 to 12 characters in length and contain only alphanumeric characters, hyphens, or underscores.
type SKU string

// An Item holds data about an inventory item.
type Item struct {
	ID          ID         `json:"id"`
	SKU         SKU        `json:"sku"`
	Name        string     `json:"name"`
	Description string     `json:"description,omitempty"`
	PriceInCAD  *float64   `json:"price_CAD,omitempty"`
	Quantity    *int       `json:"quantity"`
	DateAdded   *time.Time `json:"-"`
	LastUpdated *time.Time `json:"-"`
}

// GetID returns an item's id field.
func (item *Item) GetID() ID {
	return item.ID
}

// SetID set an item's id field if it has not yet been set.
// Returns an error if the id has already been set or the given id is invalid.
func (item *Item) SetID(id ID) error {
	if !item.IdIsPresent() {
		if _, err := id.isValid(); err != nil {
			return err
		}
		item.ID = id
		return nil
	}
	return errors.New("item id has already been set")
}

// ValidateID checks that the ID is present and formatted according to the API specifcations.
// Returns a 400 Bad Request if the ID is invalid.
func (item *Item) ValidateID() (int, error) {
	return item.ID.isValid()
}

// ValidateSKU checks that the SKU is present and formatted according to the API specifcations.
// Returns a 400 Bad Request if the SKU is invalid.
func (item *Item) ValidateSKU() (int, error) {
	return item.SKU.isValid()
}

// ValidateName checks that the Name is present and formatted according to the API specifications.
// Names are properly formatted if they contain at least 1 non-whitespace character.
// Returns a 400 Bad Request if the SKU is invalid.
func (item *Item) ValidateName() (int, error) {
	item.Name = strings.TrimSpace(item.Name)
	if len(item.Name) == 0 {
		return http.StatusBadRequest, errors.New("name cannot be whitespace or empty")
	}
	return 0, nil
}

// ValidateDescription formats the Description according to the API specification.
// Descriptions are properly formatted if any leading or trailing whitespace is trimmed.
// Returns nil as there are no restrictions on Descriptions.
func (item *Item) ValidateDescription() (int, error) {
	item.Description = strings.TrimSpace(item.Description)
	return 0, nil
}

// ValidatePrice checks that the PriceInCAD is formatted according to the API specifications, if it is present.
// PriceInCAD is an optional field.
// If PriceInCAD is present, it is properly formatted if it is non-negative.
// Returns a 400 Bad Request if the PriceInCAD is invalid.
func (item *Item) ValidatePrice() (int, error) {
	if price := item.PriceInCAD; price != nil && *price < 0 {
		return http.StatusBadRequest, errors.New("price_CAD cannot be negative")
	}
	return 0, nil
}

// ValidateQuantity checks that the Quantity is formatted according to the API specifications, if it is present.
// Quantity is an optional field and will take on a default value of 0 if it is not provided.
// If Quantity is present, it is properly formatted if it is non-negative.
// Returns a 400 Bad Request if the Quantity is invalid.
func (item *Item) ValidateQuantity() (int, error) {
	if qty := item.Quantity; qty != nil && *qty < 0 {
		return http.StatusBadRequest, errors.New("quantity cannot be negative")
	} else if qty == nil {
		q := 0
		item.Quantity = &q
	}
	return 0, nil
}

// isValid checks that the ID is present and formatted according to the API specifcations.
// IDs are properly formatted if they are 20 characters long and contain only lowercase letters a-v and numerical digits 0-9.
// Returns a 400 Bad Request if the ID is invalid.
func (id ID) isValid() (int, error) {
	if len(id) != ID_LEN {
		return http.StatusBadRequest, fmt.Errorf("id must be %d characters in length", ID_LEN)
	}
	for _, c := range id {
		if !(('a' <= c && c <= 'v') || ('0' <= c && c <= '9')) {
			return http.StatusBadRequest, fmt.Errorf("id may only contain [a-v 0-9]")
		}
	}
	return 0, nil
}

// isValid checks that the SKU is present and formatted according to the API specifcations.
// SKUs are properly formatted if they are between 4 and 12 characters long and contain only alphanumeric characters, hyphens, or underscores.
// Returns a 400 Bad Request if the SKU is invalid.
func (sku SKU) isValid() (int, error) {
	if len := len(sku); len < SKU_MIN_LEN || len > SKU_MAX_LEN {
		return http.StatusBadRequest, fmt.Errorf("SKU must be between %d and %d characters in length", SKU_MIN_LEN, SKU_MAX_LEN)
	}
	for _, c := range sku {
		if !(unicode.IsLetter(c) || unicode.IsDigit(c) || c == '-' || c == '_') {
			return http.StatusBadRequest, fmt.Errorf("SKU may only contain [a-z A-Z 0-9 _ -]")
		}
	}
	return 0, nil
}

// ValidateItem ensures that all properties needed to write the Item to database are present and properly formatted.
// SKU and Name are mandatory as they can never be empty.
// Description, PriceInCAD and Quantity may be empty, but will be overwritten to their default values:
// empty string, nil, 0, respectively.
// Returns a 400 Bad Request for invalid Items.
func (item *Item) ValidateItem() (int, error) {
	if code, err := item.ValidateSKU(); err != nil {
		return code, err
	} else if code, err = item.ValidateName(); err != nil {
		return code, err
	} else if code, err = item.ValidateDescription(); err != nil {
		return code, err
	} else if code, err = item.ValidatePrice(); err != nil {
		return code, err
	} else if code, err = item.ValidateQuantity(); err != nil {
		return code, err
	}
	return 0, nil
}

// IdIsPresent returns true if the ID property is present in the Item, false otherwise.
func (item *Item) IdIsPresent() bool {
	return len(item.ID) == ID_LEN
}
