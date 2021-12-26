package main

import "fmt"

func main() {
	even := make(chan int)
	odd := make(chan int)
	quit := make(chan bool)

	go send(even, odd, quit)

	receive(even, odd, quit)
}

func receive(even, odd <-chan int, quit <-chan bool) {
	for {
		select {
		case v := <-even:
			fmt.Println("from the even channel: ", v)
		case v := <-odd:
			fmt.Println("from the odd channel: ", v)
		case v, ok := <-quit:
			if !ok {
				fmt.Println("from comma ok: ", v, ok)
				return
			}
			fmt.Println("from comma ok: ", v)
		}
	}
}

func send(even, odd chan<- int, quit chan<- bool) {
	for i := 1; i <= 100; i++ {
		if i%2 == 0 {
			even <- i
		} else {
			odd <- i
		}
	}
	close(quit)
}
