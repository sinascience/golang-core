package main

import "fmt"

// TODO: Implement this function.
func CalculateTotal(quantity int8, price int32) int64 {
	// Remember to convert both quantity and price to int64 before multiplying.
	quantity64 := int64(quantity)
	price64:= int64(price)
	return quantity64 * price64
}

func main() {
	var quantity int8 = 10
	var price int32 = 25000
	total := CalculateTotal(quantity, price)
	fmt.Printf("Total: %d (Type: %T)\n", total, total) // Expected: Total: 250000 (Type: int64)
}