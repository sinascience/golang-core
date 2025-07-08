package main

import "fmt"

// TODO: Define Notifier interface and ConsoleNotifier struct
type Notifier interface {
	Send(message string)
}
type ConsoleNotifier struct{}
type EmailNotifier struct{}

func (c ConsoleNotifier) Send(message string) {
	fmt.Println("[CONSOLE NOTIFIER]", message)
}
func (e EmailNotifier) Send(message string) {
	fmt.Println("[EMAIL]: Sending email with message:", message)
}

// TODO: Define OrderService struct
// It should hold a Notifier.
type OrderService struct {
	notifier Notifier
}

// NewOrderService is the constructor that "injects" the dependency.
func NewOrderService(notifier Notifier) *OrderService {
	return &OrderService{notifier: notifier}
}

// PlaceOrder should use the injected notifier.
func (s *OrderService) PlaceOrder() {
	s.notifier.Send("Your order has been placed!")
}

func main() {
	// --- Scenario 1 ---
	consoleNotifier := ConsoleNotifier{}
	orderService1 := NewOrderService(consoleNotifier)
	fmt.Println("Using Console Notifier:")
	orderService1.PlaceOrder()

	fmt.Println("\n--- Scenario 2 ---")
	// --- Scenario 2 ---
	// TODO: Create an EmailNotifier and inject it into a new OrderService
	emailNotifier := EmailNotifier{}
	orderService2 := NewOrderService(emailNotifier)
	fmt.Println("Using Email Notifier:")
	orderService2.PlaceOrder()
}
