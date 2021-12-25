package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

func main() {
	var wg sync.WaitGroup

	var num int64 = 0
	gs := 100
	wg.Add(gs)

	for i := 0; i < gs; i++ {
		go func() {
			atomic.AddInt64(&num, 1)
			fmt.Println(num)
			wg.Done()
		}()
	}
	wg.Wait()
	fmt.Println(num)
}
