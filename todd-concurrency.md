# Go 동시성 기초 및 활용

- Todd 의 section20 Concurrency

---

- [Concurrency is not parallelism](https://go.dev/blog/waza-talk)
- [Concurrency vs Parallelism](https://www.youtube.com/watch?v=Y1pgpn2gOSg)

- 동시성 - 한 번에 많은 일을 처리하는 것
- 병렬성 - 많은 일을 한 번에 처리하는 것 

## runtime

### runtime 패키지

```go
fmt.Println("OS\t\t", runtime.GOOS)
fmt.Println("ARCH\t\t", runtime.GOARCH)
fmt.Println("CPUs\t\t", runtime.NumCPU())
fmt.Println("Goroutines\t\t", runtime.NumGoroutine())
```

### WaitGroup

```go
var wg sync.WaitGroup
```

- wg.Add() - wait group 추가하기
- wg.Wait() - wait group 이 모두 Done 할 때까지 대기
- wg.Done() - 수행이 끝남을 알림 