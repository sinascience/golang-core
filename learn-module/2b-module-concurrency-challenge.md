# Module 2: Concurrency & Data Handling - Challenges

These challenges will help you master the core concepts of this module: goroutines, waitgroups, channels for communication, and handling JSON data.

### Submission Guidelines

Follow the same structure as Module 1: create a folder for each question inside `answers/module-2/`.

-----

### Challenge 1: Simple Goroutine

**Task**: Write a `main` function that prints "Hello from main\!". Then, start a goroutine that waits for 1 second and then prints "Hello from goroutine\!". Notice what happens when you run the program. Does the second message always appear? Why or why not?

**Code Starter**:

```go
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
}
```

-----

### Challenge 2: Using `sync.WaitGroup`

**Task**: Fix the program from Challenge 1 using a `sync.WaitGroup`. The `main` function must not exit until the goroutine has finished its work.

**Code Starter**:

```go
package main

import (
	"fmt"
	"sync"
	"time"
)

func sayHelloWithWaitGroup(wg *sync.WaitGroup) {
	defer wg.Done() // This will be called when the function returns.
	time.Sleep(1 * time.Second)
	fmt.Println("Hello from goroutine!")
}

func main() {
	var wg sync.WaitGroup // Create a WaitGroup.

	fmt.Println("Hello from main!")

	wg.Add(1) // Increment the WaitGroup counter.
	go sayHelloWithWaitGroup(&wg)

	wg.Wait() // Block until the WaitGroup counter is zero.
	fmt.Println("Main function finished.")
}
```

-----

### Challenge 3: Basic Channels

**Task**: Write a program where a goroutine generates a message and sends it to the `main` function over a channel. The `main` function should wait to receive the message and then print it.

**Code Starter**:

```go
package main

import (
	"fmt"
	"time"
)

func produceMessage(ch chan string) {
	time.Sleep(2 * time.Second)
	ch <- "Message from the producer!" // Send message into the channel.
}

func main() {
	messageChannel := make(chan string) // Create a channel of type string.

	go produceMessage(messageChannel)

	fmt.Println("Waiting for a message...")
	// Block and wait to receive a message from the channel.
	receivedMessage := <-messageChannel
	fmt.Printf("Received: %s\n", receivedMessage)
}
```

-----

### Challenge 4: Understanding `defer`

**Task**: Write a function that simulates opening a resource (like a file or database connection) and uses `defer` to guarantee it gets closed.

**Code Starter**:

```go
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
```

-----

### Challenge 5: Fixing a Race Condition with a Mutex

**Task**: The following code has a race condition. Multiple goroutines try to increment a shared counter at the same time, leading to an incorrect final value. Your task is to fix it using a `sync.Mutex`.

**Code Starter**:

```go
package main

import (
	"fmt"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	var mu sync.Mutex // The Mutex
	counter := 0

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			// TODO: Lock the mutex before changing the counter.
			// TODO: Unlock the mutex after changing the counter.
			counter++
		}()
	}

	wg.Wait()
	// If you run this without the mutex, the count will be less than 1000.
	// With the mutex, it should be exactly 1000.
	fmt.Printf("Final counter: %d\n", counter)
}
```

-----

### Challenge 6: JSON Marshaling

**Task**: Take a Go map and use the `encoding/json` package to "marshal" it into a JSON byte slice, then print it as a string.

**Code Starter**:

```go
package main

import (
	"encoding/json"
	"fmt"
)

func main() {
	categoryCount := map[string]int{
		"Goods":   150,
		"Service": 80,
	}

	// TODO: Marshal the map into a JSON byte slice.
	// jsonData, err := ...

	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println(string(jsonData)) // Expected: {"Goods":150,"Service":80}
}
```

-----

### Challenge 7: JSON Unmarshaling

**Task**: Take a JSON string and "unmarshal" it into a Go `map[string]int`.

**Code Starter**:

```go
package main

import (
	"encoding/json"
	"fmt"
)

func main() {
	jsonString := `{"Laptops": 25, "Keyboards": 100, "Mice": 150}`
	var productStock map[string]int

	// TODO: Unmarshal the jsonString into the productStock map.
	// err := ...

	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Printf("Stock for Keyboards: %d\n", productStock["Keyboards"]) // Expected: Stock for Keyboards: 100
}
```

-----

### Challenge 8: Complex Struct to JSON

**Task**: Define a nested struct for a user profile and marshal it to a nicely formatted (indented) JSON string.

**Code Starter**:

```go
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

	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println(string(jsonData))
}
```

-----

### Challenge 9: Goroutine with Channels for Results

**Task**: Write a function that performs a "slow calculation" in a goroutine and returns the result to the caller via a channel.

**Code Starter**:

```go
package main

import (
	"fmt"
	"time"
)

// This function starts a goroutine and immediately returns the channel.
func slowCalculation() chan int {
	resultChannel := make(chan int)
	go func() {
		time.Sleep(2 * time.Second)
		resultChannel <- 100 * 5 // The slow work
	}()
	return resultChannel
}

func main() {
	fmt.Println("Starting slow calculation...")
	resultCh := slowCalculation()

	// Wait for the result from the channel.
	result := <-resultCh

	fmt.Printf("Calculation finished. Result: %d\n", result)
}
```

-----

### Challenge 10: Putting It All Together (Mini Report Simulator)

**Task**: Simulate the module's main task. Create a `Report` struct. Write a function that listens for new sales on a channel and updates the report safely using a mutex.

**Code Starter**:

```go
package main

import (
	"fmt"
	"sync"
	"time"
)

type Sale struct {
	Product string
	Amount  int
}

type Report struct {
	TotalSales int
	mu         sync.Mutex // To protect TotalSales
}

// updateReport listens on the channel and updates the report.
func (r *Report) updateReport(salesChannel chan Sale, wg *sync.WaitGroup) {
	defer wg.Done()
	for sale := range salesChannel {
		r.mu.Lock()
		fmt.Printf("Processing sale of %s for $%d\n", sale.Product, sale.Amount)
		r.TotalSales += sale.Amount
		r.mu.Unlock()
		time.Sleep(50 * time.Millisecond) // Simulate work
	}
}

func main() {
	var wg sync.WaitGroup
	report := Report{}
	salesChannel := make(chan Sale, 5) // A buffered channel

	wg.Add(1)
	go report.updateReport(salesChannel, &wg)

	// Simulate incoming sales from multiple sources.
	salesChannel <- Sale{Product: "Laptop", Amount: 1200}
	salesChannel <- Sale{Product: "Mouse", Amount: 25}
	salesChannel <- Sale{Product: "Keyboard", Amount: 75}
	salesChannel <- Sale{Product: "Monitor", Amount: 300}

	close(salesChannel) // Close the channel to signal no more sales are coming.
	wg.Wait()           // Wait for the report processor to finish.

	fmt.Printf("Final Report - Total Sales: $%d\n", report.TotalSales) // Expected: Final Report - Total Sales: $1600
}
```