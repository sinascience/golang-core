package main

import "fmt"

// TODO: Implement this function.
func formatWelcomeMessage(eventName string, year int) string {
	// Use fmt.Sprintf to create the welcome message.
	yearString := fmt.Sprintf("%d", year)
	return "Welcome to the " + eventName + " " + yearString
}

func main() {
	message := formatWelcomeMessage("Go Developer Day", 2025)
	fmt.Println(message) // Expected: Welcome to the Go Developer Day 2025!
}
