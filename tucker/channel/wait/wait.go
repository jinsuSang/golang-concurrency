package main

import (
	"fmt"
	"sync"
	"time"
)

func square(wg *sync.WaitGroup, ch chan int) {
	// for range 문은 채널에 데이터가 들어오기를 기다리기 때문에
	// 실행되지 않고 모든 고루틴이 멈추게 된다
	for n := range ch {
		fmt.Printf("Square: %d\n", n*n)
		time.Sleep(time.Second)
	}
	wg.Done()
}

func main() {
	var wg sync.WaitGroup
	ch := make(chan int)
	wg.Add(1)
	go square(&wg, ch)
	for i := 0; i < 10; i++ {
		ch <- i * 2
	}
	close(ch) // 채널을 닫음으로써 for range 구문을 종료한다
	wg.Wait()
}
