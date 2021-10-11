# Go 동시성 기초 및 활용

- Todd 의 section20 Concurrency

---

- [Concurrency is not parallelism](https://go.dev/blog/waza-talk)
- [Concurrency vs Parallelism](https://www.youtube.com/watch?v=Y1pgpn2gOSg)

- 동시성 - 한 번에 많은 일을 처리하는 것
- 병렬성 - 많은 일을 한 번에 처리하는 것 

## runtime 패키지

```go
fmt.Println("OS\t\t", runtime.GOOS)
fmt.Println("ARCH\t\t", runtime.GOARCH)
fmt.Println("CPUs\t\t", runtime.NumCPU())
fmt.Println("Goroutines\t\t", runtime.NumGoroutine())
```

## WaitGroup

```go
var wg sync.WaitGroup
```

- wg.Add() - wait group 추가하기
- wg.Wait() - wait group 이 모두 Done 할 때까지 대기
- wg.Done() - 수행이 끝남을 알림 

## Go Statements

- https://golang.org/ref/spec#Go_statements

## Race Condition

```bash
$ go run -race race.go
```

- 경쟁 조건 확인

```go
runtime.Gosched()
```

- 다른 고루틴이 CPU를 사요할 수 있도록 양보한다

## [Mutex](https://pkg.go.dev/sync#Mutex)

- 뮤텍스는 상호 제외 잠금입니다. 뮤텍스의 0 값은 잠금 해제된 뮤텍스 입니다. 처음 사용한 후에는 뮤텍스를 복사해서는 안됩니다.

## Atomic

- atomic 패키지는 동기화 알고리즘 구현에 유용한 저수준 atomic  메모리 초기 요소를 제공합니다

- 사용에 유의해야 하며 특수한 저수준 응용 프로그램을 제외하고 동기화는 채널이나 동기화 패키지 기능을 수행하는 것이 좋습니다. 통신을 통해 메모리를 공유하고 메모리를 공유하여 통신하지 않습니다

  

