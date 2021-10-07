package main

import (
	"fmt"
	"runtime"
	"sync"
)

var wg sync.WaitGroup

func main() {
	fmt.Println("OS\t\t", runtime.GOOS)
	fmt.Println("ARCH\t\t", runtime.GOARCH)
	fmt.Println("CPUs\t\t", runtime.NumCPU())
	fmt.Println("Goroutines\t\t", runtime.NumGoroutine())

	wg.Add(2)
	go func01()
	go func02()


	fmt.Println("CPUs\t\t", runtime.NumCPU())
	fmt.Println("Goroutines\t\t", runtime.NumGoroutine())
	wg.Wait()
}

func func01() {
	defer wg.Done()
	for i := 0; i < 100; i++ {
		fmt.Println("A: ", i)
	}
}

func func02() {
	defer wg.Done()
	for i := 0; i < 100; i++ {
		fmt.Println("B: ", i)
	}
}
