package main

import "fmt"

// TODO: Define the Speaker interface
type Speaker interface {
	Speak() string
}

// TODO: Define Human and Dog structs and implement the Speak method for each
type Human struct {
	speak string
}
type Dog struct {
	speak string
}

func (h Human) Speak() string {
	return h.speak
}

func (d Dog) Speak() string {
	return d.speak
}

// This function can accept any type that satisfies the Speaker interface.
func makeItSpeak(s Speaker) {
	fmt.Println(s.Speak())
}

func main() {
	human := Human{speak: "Hello"}
	dog := Dog{speak: "Woof"}

	makeItSpeak(human) // Expected: Hello
	makeItSpeak(dog)   // Expected: Woof
}
