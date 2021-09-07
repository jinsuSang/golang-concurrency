package main

import (
	"fmt"
	"sync"
	"time"
)

var mutex sync.Mutex
var wg sync.WaitGroup

type Account struct {
	Balance int
}

func DepositAndWithdraw(account *Account) {
	mutex.Lock()
	defer mutex.Unlock()
	if account.Balance < 0 {
		panic(fmt.Sprintf("Balance should not be negative value: %d", account.Balance))
	}
	account.Balance += 1000
	fmt.Println(account.Balance)
	time.Sleep(time.Millisecond)
	account.Balance -= 1000
	fmt.Println(account.Balance)
}

func main() {

	account := &Account{0}
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func() {
			for {
				DepositAndWithdraw(account)
			}
			wg.Done()
		}()
	}
	wg.Wait()
}
