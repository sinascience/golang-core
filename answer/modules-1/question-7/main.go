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
		{
			Title:  "Book 1",
			Author: "People 1",
			Pages:  100,
			ISBN:   "999-99999999",
		},
		// Add book 2
		{
			Title:  "Book 2",
			Author: "People 2",
			Pages:  200,
			ISBN:   "999-99999999",
		},
		// Add book 3
		{
			Title:  "Book 3",
			Author: "People 3",
			Pages:  300,
			ISBN:   "999-99999999",
		},
	}

	fmt.Println("Library Titles:")
	// TODO: Write a for loop to iterate over the library
	for _, books := range library{
		fmt.Println(books.Title)
	}
	// and print the Title of each book.
}