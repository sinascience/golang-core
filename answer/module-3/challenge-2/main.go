package main

import (
	"fmt"
	"math"
)

// TODO: Define the Shape interface
type Shape interface {
	Area() float32
}

// TODO: Define Circle and Rectangle structs and implement the Area() method
type Circle struct {
	Radius float32
}
type Rectangle struct {
	Width  float32
	Height float32
}

func (c Circle) Area() float32 {
	return math.Pi * (c.Radius * c.Radius)
}

func (r Rectangle) Area() float32 {
	return r.Height * r.Width
}

func main() {
	circle := Circle{Radius: 5}
	rectangle := Rectangle{Width: 10, Height: 5}

	fmt.Printf("Circle Area: %f\n", circle.Area())       // Expected: ~78.54
	fmt.Printf("Rectangle Area: %f\n", rectangle.Area()) // Expected: 50.0
}
