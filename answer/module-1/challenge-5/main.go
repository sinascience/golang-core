package main

import "fmt"

// TODO: Implement this function.
func CalculateTotal(quantity int8, price int32) int64 {
	// Remember to convert both quantity and price to int64 before multiplying.

	total := int64(price) * int64(quantity)

	return total
}

func main() {
	var quantity int8 = 10
	var price int32 = 25000
	total := CalculateTotal(quantity, price)
	fmt.Printf("Total: %d (Type: %T)\n", total, total) // Expected: Total: 250000 (Type: int64)
}
