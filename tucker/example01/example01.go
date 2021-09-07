package main

import (
	"fmt"
	"time"
)

func printKorean() {
	koreans := []rune{'가', '나', '다', '라', '마', '바', '사'}
	for _, korean := range koreans {
		time.Sleep(300 * time.Millisecond)
		fmt.Printf("%c ", korean)
	}
}

func printNumbers() {
	for i := 0; i < 20; i++ {
		time.Sleep(100 * time.Millisecond)
		fmt.Printf("%d ", i)
	}
}

func main() {
	go printKorean()
	go printNumbers()
	time.Sleep(3 * time.Second)
}
