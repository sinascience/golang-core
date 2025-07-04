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
			mu.Lock()
			counter++
			mu.Unlock()
		}()
	}

	wg.Wait()
	// If you run this without the mutex, the count will be less than 1000.
	// With the mutex, it should be exactly 1000.
	fmt.Printf("Final counter: %d\n", counter)
}
