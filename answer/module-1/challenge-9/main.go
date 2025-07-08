package main

import (
	"errors"
	"fmt"
)

// TODO: Implement this function.
func GetItemDetails(productID string) (string, error) {
	// If productID is "P123", return "Go Programming Bible" and nil.
	// For any other ID, return an empty string and a new error.

	if productID == "P123" {
		return "Go Programming Bible", nil
	}
	return "", errors.New("product not found")
}

func main() {
	// Test case 1: Success
	name, err := GetItemDetails("P123")
	if err != nil {
		fmt.Printf("Error case 1: %v\n", err)
	} else {
		fmt.Printf("Success case 1: Found %s\n", name) // Expected: Success case 1: Found Go Programming Bible
	}

	// Test case 2: Failure
	name, err = GetItemDetails("P456")
	if err != nil {
		fmt.Printf("Error case 2: %v\n", err) // Expected: Error case 2: product not found
	} else {
		fmt.Printf("Success case 2: Found %s\n", name)
	}
}
