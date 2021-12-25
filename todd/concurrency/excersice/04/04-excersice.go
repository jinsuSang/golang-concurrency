package main

import (
	"fmt"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	var m sync.Mutex

	num := 0
	gs := 100
	wg.Add(gs)

	for i := 0; i < gs; i++ {
		go func() {
			m.Lock()
			v := num
			v++
			num = v
			fmt.Println(num)
			m.Unlock()
			wg.Done()
		}()
	}
	wg.Wait()
	fmt.Println(num)
}
