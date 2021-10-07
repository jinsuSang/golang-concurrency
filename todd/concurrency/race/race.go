package main

import (
	"fmt"
	"runtime"
	"sync"
)

func main() {
	fmt.Println("CPU NUM: ", runtime.NumCPU())
	fmt.Println("Goroutines: ", runtime.NumGoroutine())

	counter := 0
	const gs = 1000

	var wg sync.WaitGroup
	wg.Add(gs * 2)

	for i := 0; i < gs; i++ {
		go func() {
			defer wg.Done()
			v := counter
			// time.Sleep(time.Second)
			runtime.Gosched()
			v++
			counter = v
		}()
		go func() {
			defer wg.Done()
			v := counter
			// time.Sleep(time.Second)
			runtime.Gosched()
			v--
			counter = v
		}()
		fmt.Println("counter: ", counter)
	}
	wg.Wait()
	fmt.Println("Goroutines: ", runtime.NumGoroutine())
	fmt.Println(counter)
}
