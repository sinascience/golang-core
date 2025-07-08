package main

import (
	"fmt"
	"time"
)

func sayHelloAsync() {
	time.Sleep(1 * time.Second)
	fmt.Println("Hello from goroutine!")
}

func main() {
	fmt.Println("Hello from main!")
	go sayHelloAsync()
	// What happens if you add a time.Sleep(2 * time.Second) here?
	time.Sleep(2 * time.Second)
}
