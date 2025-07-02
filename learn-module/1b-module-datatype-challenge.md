# Module 1: Core Go Concepts - Challenges (with Starters)

Welcome to the challenges for Module 1\! The purpose of these exercises is to solidify your understanding of Go's fundamental concepts like data types, structs, functions, and control flow.

Each challenge provides a starter code block to guide you.

### Submission Guidelines

1.  Create a new root folder named `answers`.
2.  Inside `answers`, create a folder for this module: `module-1`.
3.  For each question, create a subfolder and a `main.go` file (e.g., `answers/module-1/question-1/main.go`).
4.  Copy the starter code into your `main.go` and complete the required logic.

-----

### Challenge 1: Data Types

**Task**: Declare variables for `age`, `population`, and `worldPopulation` using the most memory-efficient integer types. Assign them values and print them.

**Code Starter**:

```go
package main

import "fmt"

func main() {
	// TODO: Declare variables with the correct types.
	// var age ...
	// var population ...
	// var worldPopulation ...

	// Assign values
	age = 35
	population = 850000
	worldPopulation = 8100000000

	// The fmt.Printf function's %T verb prints the type of the variable.
	fmt.Printf("Age: %d (Type: %T)\n", age)                 // Expected: Age: 35 (Type: int8)
	fmt.Printf("Population: %d (Type: %T)\n", population)       // Expected: Population: 850000 (Type: int32)
	fmt.Printf("World Population: %d (Type: %T)\n", worldPopulation) // Expected: World Population: 8100000000 (Type: int64)
}
```

-----

### Challenge 2: Structs

**Task**: Define a `Book` struct. Create an instance of this struct, fill it with data, and print it.

**Code Starter**:

```go
package main

import "fmt"

// TODO: Define the Book struct here.
// It should have fields for Title (string), Author (string),
// Pages (int16), and ISBN (string).
type Book struct {

}

func main() {
	// TODO: Create an instance of the Book struct.
	myBook := Book{
		// Fill in the details for your favorite book
	}

	// The %+v verb prints the struct with field names.
	fmt.Printf("%+v\n", myBook)
	// Example Expected Output: {Title:The Go Programming Language Author:Alan A. A. Donovan & Brian W. Kernighan Pages:380 ISBN:978-0134190440}
}
```

-----

### Challenge 3: String Formatting

**Task**: Create a function that accepts `eventName` and `year` and returns a formatted welcome string.

**Code Starter**:

```go
package main

import "fmt"

// TODO: Implement this function.
func formatWelcomeMessage(eventName string, year int) string {
	// Use fmt.Sprintf to create the welcome message.
	return ""
}

func main() {
	message := formatWelcomeMessage("Go Developer Day", 2025)
	fmt.Println(message) // Expected: Welcome to the Go Developer Day 2025!
}
```

-----

### Challenge 4: Enum with `iota`

**Task**: Create an enum for `Months` using `iota`. Write a function that converts a `Month` to its string name.

**Code Starter**:

```go
package main

import "fmt"

type Month uint8

// TODO: Use iota to define the months from January to December.
const (
	_ Month = iota // Ignore 0
	January
	// ... continue for all 12 months
)

// TODO: Implement this function to convert a Month to its string representation.
func GetMonthName(m Month) string {
	// Hint: A switch statement or a slice of strings would work well here.
	return ""
}

func main() {
	fmt.Println(GetMonthName(January)) // Expected: January
}
```

-----

### Challenge 5: Calculations & Type Conversion

**Task**: Write a function `CalculateTotal` that accepts `quantity` (`int8`) and `price` (`int32`) and returns the total (`int64`).

**Code Starter**:

```go
package main

import "fmt"

// TODO: Implement this function.
func CalculateTotal(quantity int8, price int32) int64 {
	// Remember to convert both quantity and price to int64 before multiplying.
	return 0
}

func main() {
	var quantity int8 = 10
	var price int32 = 25000
	total := CalculateTotal(quantity, price)
	fmt.Printf("Total: %d (Type: %T)\n", total) // Expected: Total: 250000 (Type: int64)
}
```

-----

### Challenge 6: Pointers

**Task**: Write a function `MarkAsPaid` that takes a pointer to a boolean and changes its value to `true`.

**Code Starter**:

```go
package main

import "fmt"

// TODO: Implement this function.
func MarkAsPaid(isPaid *bool) {
	// Change the value that the pointer is pointing to.
}

func main() {
	isUnpaid := false
	fmt.Printf("Status before: %v\n", isUnpaid) // Expected: Status before: false

	// Pass the memory address of isUnpaid to the function.
	MarkAsPaid(&isUnpaid)

	fmt.Printf("Status after: %v\n", isUnpaid) // Expected: Status after: true
}
```

-----

### Challenge 7: Slices and Structs

**Task**: Using the `Book` struct from Challenge \#2, create a slice of books and loop through it to print each title.

**Code Starter**:

```go
package main

import "fmt"

type Book struct {
	Title  string
	Author string
	Pages  int16
	ISBN   string
}

func main() {
	// TODO: Create a slice of Book structs and add at least 3 books.
	library := []Book{
		// Add book 1
		// Add book 2
		// Add book 3
	}

	fmt.Println("Library Titles:")
	// TODO: Write a for loop to iterate over the library
	// and print the Title of each book.
}
```

-----

### Challenge 8: Control Flow (`if/else`)

**Task**: Write a function `GetCategoryType` that returns a string describing the `ProductCategory`.

**Code Starter**:

```go
package main

import "fmt"

type ProductCategory uint8
const (
	_ ProductCategory = iota
	Goods
	Service
	Subscription
)

// TODO: Implement this function.
func GetCategoryType(category ProductCategory) string {
	// Use an if/else or switch statement.
	// If the category is Goods, return "Physical Item".
	// Otherwise, return "Digital Item or Service".
	return ""
}

func main() {
	fmt.Printf("Goods: %s\n", GetCategoryType(Goods))         // Expected: Goods: Physical Item
	fmt.Printf("Service: %s\n", GetCategoryType(Service))     // Expected: Service: Digital Item or Service
}
```

-----

### Challenge 9: Functions with Multiple Returns

**Task**: Create a function `GetItemDetails` that simulates fetching data, returning a value and a potential error.

**Code Starter**:

```go
package main

import (
	"errors"
	"fmt"
)

// TODO: Implement this function.
func GetItemDetails(productID string) (string, error) {
	// If productID is "P123", return "Go Programming Bible" and nil.
	// For any other ID, return an empty string and a new error.
	return "", nil
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
```

-----

### Challenge 10: Putting It All Together

**Task**: Create a `Product` struct and a slice of products. Write a function to calculate the total inventory value.

**Code Starter**:

```go
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
```