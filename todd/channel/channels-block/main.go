package main

import "fmt"

func main() {
	c := make(chan int)
	c <- 50
	fmt.Println(<-c)
}
