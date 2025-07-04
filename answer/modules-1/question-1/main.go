package main

import "fmt"

func main() {
	// TODO: Declare variables with the correct types.
	// var age ...
	// var population ...
	// var worldPopulation ...

	// Assign values
	var age int8 = 35
	var population int32 = 850000
	var worldPopulation int64 = 8100000000

	// The fmt.Printf function's %T verb prints the type of the variable.
	fmt.Printf("Age: %d (Type: %T)\n", age, age)                 // Expected: Age: 35 (Type: int8)
	fmt.Printf("Population: %d (Type: %T)\n", population, population)       // Expected: Population: 850000 (Type: int32)
	fmt.Printf("World Population: %d (Type: %T)\n", worldPopulation, worldPopulation) // Expected: World Population: 8100000000 (Type: int64)
}