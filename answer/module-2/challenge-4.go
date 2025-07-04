package main

import "fmt"

func processResource() {
	fmt.Println("1. Opening resource...")
	// Use defer to ensure the resource is closed.
	// This line will be executed right before the function returns.
	defer fmt.Println("4. Closing resource (deferred).")

	fmt.Println("2. Processing resource...")
	fmt.Println("3. Finished processing.")
}

func main() {
	processResource()
}
