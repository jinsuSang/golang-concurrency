package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

var wg sync.WaitGroup

func diningProblem(name string, first, second *sync.Mutex, firstName, secondName string)  {
	for i := 0; i < 1000; i++ {
		fmt.Println(name, "저녁 식사 준비")
		first.Lock()
		fmt.Println(name, firstName, "획득")
		second.Lock()
		fmt.Println(name, secondName, "획득")

		fmt.Println(name, "저녁 식사")
		time.Sleep(time.Duration(rand.Intn(10)) * time.Millisecond)

		second.Unlock()
		first.Unlock()
	}
	wg.Done()
}

func main() {
	rand.Seed(time.Now().UnixNano())

	wg.Add(2)
	fork := &sync.Mutex{}
	spoon := &sync.Mutex{}

	go diningProblem("A", fork, spoon, "fork", "spoon")
	go diningProblem("B", spoon, fork, "spoon", "fork")


	wg.Wait()
}
