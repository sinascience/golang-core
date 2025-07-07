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
	// Dikarenakan pada saat kita menambahkan time.Sleep(2 * time.Second) maka kita akan menunggu selama 2 detik untuk menjalankan program. 
	// Jika tidak ada time.Sleep(2 * time.Second) maka program akan langsung berhenti saat main selesai mengerjakan fungsi lain yang berjalan secara parallel
	// Yang mana jika kita menunggu sekitar 2 detik maka kita akan melihat output Hello from goroutine yang memerlukan waktu tunggu 1 detik sebelum dieksekusi
}