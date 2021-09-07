package main

import (
	"fmt"
	"sync"
)

var wg sync.WaitGroup

func SumAtoB(a, b int) {
	sum := 0
	for i := a; i < b; i++ {
		sum += i
	}
	fmt.Printf("%d\n", sum)
	wg.Done()
}

func main() {
	wg.Add(20)
	for i := 0; i < 20; i++ {
		go SumAtoB(1, 1000000)
	}
	wg.Wait()
}
