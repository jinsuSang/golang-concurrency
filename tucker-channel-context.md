# Go 채널과 컨텍스트

- Tucker의 Go 언어 프로그래밍 chapter 25

---

## 채널

- 고루틴끼리 메시지를 전달하는 메시지 큐이다

- 채널 인스턴스 생성

  ```go
  var messages chan string = make(chan string)
  ```

- 채널 데이터 넣기

  ```go
  messages <- "this is a message"
  ```

- 채널 데이터 빼기

  ```go
  var msg string = <- messages
  ```

- 채널 크기

  - 채널 생성시 크기가 0인 채널이 만들어진다
  - 데이터를 보관할 공간이 없기 때문에 빠져나갈 때까지 대기한다
  - 버퍼 채널
    - 내부에 데이터를 보관할 수 있는 메모리 영역을 버퍼라고 부른다 
    - 버퍼가 다 차는 경우 빈자리가 생길 때까지 대기한다
    - 데이터가 제 때 빠지지 않는다면 고루틴이 멈춘다 

- close()

  ```go
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
  ```

- select 문

  - 여러 채널을 동시에 대기할 때 사용한다

  - 하나의 채널에서 데이터를 읽으면 구문을 실행하고 select 문이 종료된다

  - 반복문 처리를 위해 for 반복문을 함께 사용해야 한다 

    ```go
    func square(wg *sync.WaitGroup, ch chan int) {
    	tick := time.Tick(time.Second)
    	terminate := time.After(10 * time.Second)
    
    	for {
    		select {
    		case <-tick:
    			fmt.Println("Tick")
    		case <-terminate:
    			fmt.Println("Terminated")
    			wg.Done()
    			return
    		case n := <-ch:
    			fmt.Printf("Square: %d\n", n*n)
    			time.Sleep(time.Second)
    		}
    	}
    }
    ```

## 컨텍스트

- context 패키지에서 제공하는 기능으로 작업 지시시 작업 가능 시간, 작업 취소 등의 조건을 지시할 수 있는 작업 명세서이다
- 새로운 고루틴으로 작업을 시작할 때 일정 시간 제한을 걸거나 외부에서 작업을 취소할 때 사용된다

### 작업 취소가 가능한 컨텍스트

```go
var wg sync.WaitGroup

func main() {
	wg.Add(1)
	ctx, cancel := context.WithCancel(context.Background())
	go PrintEverySecond(ctx)
	time.Sleep(5 * time.Second)
	cancel() // 5초 후 종료
	wg.Wait()
}

func PrintEverySecond(ctx context.Context) {
	tick := time.Tick(time.Second)
	for {
		select {
		case <-ctx.Done():
			wg.Done()
			return
		case <-tick:
			fmt.Println("Tick")
		}
	}
}
```

### 작업 시간을 설정한 컨텍스트

```go
ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
```

### 특정 값을 설정한 컨텍스트

```go
var wg sync.WaitGroup

func main() {
	wg.Add(1)
	ctx := context.WithValue(context.Background(), "number", 9)
	go square(ctx)
	wg.Wait()
}

func square(ctx context.Context) {
	if v := ctx.Value("number"); v != nil {
		n := v.(int)
		fmt.Println(n * n)
	}
	wg.Done()
}
```

```go
ctx, cancel := context.WithCancel(context.Background())
ctx = context.WithValue(ctx, "number", 9)
ctx = context.WithValue(ctx, "name", "jinsu")
```

