// ... (reuse code from Challenge 3) ...
package main

import "fmt"

// TODO: Define Notifier interface and ConsoleNotifier struct
type Notifier interface {
	Send(message string)
}
type ConsoleNotifier struct {
	
}
func (c ConsoleNotifier) Send(message string) {
	fmt.Println("[CONSOLE NOTIFIER]:", message)
}
// TODO: Define EmailNotifier struct and implement the Send method
type EmailNotifier struct{

}
func (e EmailNotifier) Send(message string) {
	fmt.Println("[EMAIL NOTIFIER]:", message)
}
// TODO: Define OrderService struct
// It should hold a Notifier.
type OrderService struct {
	notifier Notifier
}

// NewOrderService is the constructor that "injects" the dependency.
func NewOrderService(notifier Notifier) *OrderService {
	// ...
	return &OrderService{notifier: notifier}
}

// PlaceOrder should use the injected notifier.
func (s *OrderService) PlaceOrder() {
	s.notifier.Send("Your order has been placed!")
}

// func main() {
// 	// Create the concrete dependency
// 	consoleNotifier := ConsoleNotifier{}
// 	// Inject it into the service
// 	orderService := NewOrderService(consoleNotifier)
// 	orderService.PlaceOrder() // Expected: [CONSOLE NOTIFIER]: Your order has been placed!
// }


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
	emailNotifier := EmailNotifier{}
	orderService2 := NewOrderService(emailNotifier)
	fmt.Println("Using Email Notifier:")
	// orderService2.PlaceOrder()
	orderService2.PlaceOrder()

}