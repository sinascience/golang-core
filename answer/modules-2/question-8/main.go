package main

import (
	"encoding/json"
	"fmt"
)

type UserProfile struct {
	Username   string
	Email      string
	Attributes map[string]string
}

func main() {
	profile := UserProfile{
		Username: "gopher123",
		Email:    "gopher@example.com",
		Attributes: map[string]string{
			"Country": "USA",
			"Tier":    "Gold",
		},
	}

	// TODO: Marshal the profile struct with an indent for nice formatting.
	// Use json.MarshalIndent()
	// jsonData, err := ...
	jsonData, err := json.MarshalIndent(profile, "","   ")

	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println(string(jsonData))
}