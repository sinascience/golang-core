package main

import "fmt"

type Product struct {
	Name  string
	Price int32
	Stock int16
}

// TODO: Implement this function.
func CalculateInventoryValue(inventory []Product) int64 {
	var totalValue int64 = 0
	// Loop through the inventory.
	// For each product, multiply its Price by its Stock.
	// Remember to convert to int64.
	// Add the result to totalValue.
	for _, value := range inventory {
		price64 := int64(value.Price)
		stock64 := int64(value.Stock)
		totalValue += price64 * stock64
	}
	return totalValue
}

func main() {
	inventory := []Product{
		{Name: "Laptop", Price: 15000000, Stock: 10},
		{Name: "Mouse", Price: 250000, Stock: 50},
		{Name: "Keyboard", Price: 750000, Stock: 30},
	}

	totalValue := CalculateInventoryValue(inventory)
	fmt.Printf("Total Inventory Value: %d\n", totalValue) // Expected: Total Inventory Value: 185000000
}