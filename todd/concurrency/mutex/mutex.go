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

	var mu sync.Mutex

	for i := 0; i < gs; i++ {
		go func() {
			defer wg.Done()
			mu.Lock()
			v := counter
			// time.Sleep(time.Second)
			runtime.Gosched()
			v++
			counter = v
			mu.Unlock()
		}()
		go func() {
			defer wg.Done()
			mu.Lock()
			v := counter
			// time.Sleep(time.Second)
			runtime.Gosched()
			v--
			counter = v
			mu.Unlock()
		}()
		fmt.Println("counter: ", counter)
	}
	wg.Wait()
	fmt.Println("Goroutines: ", runtime.NumGoroutine())
	fmt.Println(counter)
}
