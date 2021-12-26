package main

import "fmt"

func main() {
	c := make(chan int)
	// send
	go send(c)

	// receive
	receive(c)

}

// send
func send(c chan<- int) {
	c <- 39
}

// receive
func receive(c <-chan int) {
	fmt.Println(<-c)
}
