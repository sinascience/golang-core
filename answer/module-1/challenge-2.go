package main

import "fmt"

// TODO: Define the Book struct here.
// It should have fields for Title (string), Author (string),
// Pages (int16), and ISBN (string).
type Book struct {
	Title, Author string
	Pages         int16
	ISBN          string
}

func main() {
	// TODO: Create an instance of the Book struct.
	myBook := Book{
		Title:  "The Go Programming Language",
		Author: "Alan A. A. Donovan & Brian W. Kernighan",
		Pages:  380,
		ISBN:   "978-0134190440",
	}

	// The %+v verb prints the struct with field names.
	fmt.Printf("%+v\n", myBook)
	// Example Expected Output: {Title:The Go Programming Language Author:Alan A. A. Donovan & Brian W. Kernighan Pages:380 ISBN:978-0134190440}
}
