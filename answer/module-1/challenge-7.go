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
		{Title: "The Go Programming Language", Author: "Alan A. A. Donovan", Pages: 380, ISBN: "978-0134190440"},
		{Title: "Clean Code", Author: "Robert C. Martin", Pages: 464, ISBN: "978-0132350884"},
		{Title: "Introduction to Algorithms", Author: "Thomas H. Cormen", Pages: 1312, ISBN: "978-0262033848"},
	}

	fmt.Println("Library Titles:")
	// TODO: Write a for loop to iterate over the library
	// and print the Title of each book.
	for _, book := range library {
		fmt.Println(book.Title)
	}
}
