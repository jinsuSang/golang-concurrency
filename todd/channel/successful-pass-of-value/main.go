package main

import "fmt"

func main() {
	c := make(chan int)
	go func() { c <- 50 }()
	fmt.Println(<-c)
}
