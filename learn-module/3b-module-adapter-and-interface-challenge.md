# Module 3: Interfaces & Architecture - Challenges

These challenges focus on the architectural concepts from Module 3: interfaces, dependency injection, and handling different types of HTTP requests.

### Submission Guidelines

Follow the same structure as previous modules: `answers/module-3/question-X/main.go`.

-----

### Challenge 1: Basic Interface

**Task**: Define an interface `Speaker` with one method: `Speak() string`. Create two different structs, `Human` and `Dog`, that both satisfy the `Speaker` interface. The `Human` should say "Hello", and the `Dog` should say "Woof". Write a function that accepts a `Speaker` and prints the result of its `Speak` method.

**Code Starter**:

```go
package main

import "fmt"

// TODO: Define the Speaker interface

// TODO: Define Human and Dog structs and implement the Speak method for each

// This function can accept any type that satisfies the Speaker interface.
func makeItSpeak(s Speaker) {
	fmt.Println(s.Speak())
}

func main() {
	human := Human{}
	dog := Dog{}

	makeItSpeak(human) // Expected: Hello
	makeItSpeak(dog)   // Expected: Woof
}
```

-----

### Challenge 2: Interface for Calculations

**Task**: Define an interface `Shape` with one method: `Area() float64`. Create two structs, `Circle` (with a `Radius` field) and `Rectangle` (with `Width` and `Height` fields), that both implement the `Shape` interface. Print the area of one of each.

**Code Starter**:

```go
package main

import (
	"fmt"
	"math"
)

// TODO: Define the Shape interface

// TODO: Define Circle and Rectangle structs and implement the Area() method

func main() {
	circle := Circle{Radius: 5}
	rectangle := Rectangle{Width: 10, Height: 5}

	fmt.Printf("Circle Area: %f\n", circle.Area())     // Expected: ~78.54
	fmt.Printf("Rectangle Area: %f\n", rectangle.Area()) // Expected: 50.0
}
```

-----

### Challenge 3: Simple Dependency Injection

**Task**: Create a `Notifier` interface with a `Send(message string)` method. Create a `ConsoleNotifier` struct that implements this interface by printing the message to the console. Then, create an `OrderService` that takes a `Notifier` in its constructor. The `OrderService` should have a method `PlaceOrder` that calls the notifier.

**Code Starter**:

```go
package main

import "fmt"

// TODO: Define Notifier interface and ConsoleNotifier struct

// TODO: Define OrderService struct
// It should hold a Notifier.
type OrderService struct {

}

// NewOrderService is the constructor that "injects" the dependency.
func NewOrderService(notifier Notifier) *OrderService {
	// ...
}

// PlaceOrder should use the injected notifier.
func (s *OrderService) PlaceOrder() {
	// s.notifier.Send("Your order has been placed!")
}

func main() {
	// Create the concrete dependency
	consoleNotifier := ConsoleNotifier{}
	// Inject it into the service
	orderService := NewOrderService(consoleNotifier)

	orderService.PlaceOrder() // Expected: [CONSOLE NOTIFIER]: Your order has been placed!
}
```

-----

### Challenge 4: Swapping Dependencies

**Task**: Using the code from Challenge 3, create a *second* notifier called `EmailNotifier` that also implements the `Notifier` interface. This one should print `[EMAIL]: Sending email with message: ...`. In your `main` function, show how you can create a second `OrderService` with this new notifier without changing the `OrderService` code at all.

**Code Starter**:

```go
// ... (reuse code from Challenge 3) ...

// TODO: Define EmailNotifier struct and implement the Send method

func main() {
	// --- Scenario 1 ---
	consoleNotifier := ConsoleNotifier{}
	orderService1 := NewOrderService(consoleNotifier)
	fmt.Println("Using Console Notifier:")
	orderService1.PlaceOrder()

	fmt.Println("\n--- Scenario 2 ---")
	// --- Scenario 2 ---
	// TODO: Create an EmailNotifier and inject it into a new OrderService
	// emailNotifier := ...
	// orderService2 := ...
	fmt.Println("Using Email Notifier:")
	// orderService2.PlaceOrder()
}
```