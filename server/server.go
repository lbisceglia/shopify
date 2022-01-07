package server

import (
	"fmt"
	"net/http"
)

// createItem creates an inventory item.
func CreateItem(w http.ResponseWriter, r *http.Request) {
	// TODO
}

// updateItem updates an inventory item's information.
func UpdateItem(w http.ResponseWriter, r *http.Request) {
	// TODO
}

// deleteItem removes an item from inventory.
func DeleteItem(w http.ResponseWriter, r *http.Request) {
	// TODO
}

// getItems lists all items in inventory in a paginated format.
func GetItems(w http.ResponseWriter, r *http.Request) {
	// TODO
	fmt.Fprintf(w, "Items Page")
}
