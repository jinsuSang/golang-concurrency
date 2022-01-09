# Go  동시성 패턴

### for-select 루프 

```go
for { // Either loop infinitely or range over something
    select {
    // Do some work with channels
    }
}
```

#### 채널에서 반복 변수 보내기

순회할 수 있는 것을 채널 값으로 변환하려고 하는 경우 

```go
for _, s := range []string{"a", "b", "c"} {
    select {
    case <-done:
        return
    case stringStream <- s:
    }
}
```

####  멈추기를 기다리면서 무한 대기

```go 
for {
    select {
    case <-done:
        return
    default:
    }

    // Do non-preemptable work
}

for {
    select {
    case <-done:
        return
    default:
        // Do non-preemptable work
    }
}
```



### 고루틴 누수 방지

#### 고루틴 종료 경로

- 작업 완료
- 복구할 수 없는 에러로 인해 더 이상 작업을 계속할 수 없을 때 
- 작업을 중단하라는 요청을 받았을 때 



#### 고루틴 누수 예시

```go 
doWork := func(strings <-chan string) <-chan interface{} {
    completed := make(chan interface{})
    go func() {
        defer fmt.Println("doWork exited.")
        defer close(completed)
        for s := range strings {
            // Do something interesting
            fmt.Println(s)
        }
    }()
    return completed
}

doWork(nil)
// Perhaps more work is done here
fmt.Println("Done.")
```

- strings 채널은 실제로 어떠한 문자열도 쓰지 않는다
- doWork 를 포함하는 고루틴은 이 프로세스가 지속되는 동안 메모리에 남아 있다 
- dowork 과 메인 고루틴 안에서 고루틴을 조인하면 데드락 상태가 발생한다 



부모 고루틴이 자식 고루틴에게 취소 신호를 보낼 수 있도록 신호를 설정해야 한다 

```go
doWork := func(
  done <-chan interface{},
  strings <-chan string,
) <-chan interface{} { 
    terminated := make(chan interface{})
    go func() {
        defer fmt.Println("doWork exited.")
        defer close(terminated)
        for {
            select {
            case s := <-strings:
                // Do something interesting
                fmt.Println(s)
            case <-done: 
                return
            }
        }
    }()
    return terminated
}

done := make(chan interface{})
terminated := doWork(done, nil)

go func() { 
    // Cancel the operation after 1 second.
    time.Sleep(1 * time.Second)
    fmt.Println("Canceling doWork goroutine...")
    close(done)
}()

<-terminated 
fmt.Println("Done.")
```



#### 채널에 값을 쓰려는 시도를 차단하는 고루틴의 경우 

```go
newRandStream := func() <-chan int {
    randStream := make(chan int)
    go func() {
        defer fmt.Println("newRandStream closure exited.") 1
        defer close(randStream)
        for {
            randStream <- rand.Int()
        }
    }()

    return randStream
}

randStream := newRandStream()
fmt.Println("3 random ints:")
for i := 1; i <= 3; i++ {
    fmt.Printf("%d: %d\n", i, <-randStream)
}
```

- 해결책

```go
newRandStream := func(done <-chan interface{}) <-chan int {
    randStream := make(chan int)
    go func() {
        defer fmt.Println("newRandStream closure exited.")
        defer close(randStream)
        for {
            select {
            case randStream <- rand.Int():
            case <-done:
                return
            }
        }
    }()

    return randStream
}

done := make(chan interface{})
randStream := newRandStream(done)
fmt.Println("3 random ints:")
for i := 1; i <= 3; i++ {
    fmt.Printf("%d: %d\n", i, <-randStream)
}
close(done)

// Simulate ongoing work
time.Sleep(1 * time.Second)
```



### or 채널

```go
var or func(channels ...<-chan interface{}) <-chan interface{}
or = func(channels ...<-chan interface{}) <-chan interface{} { 1
    switch len(channels) {
    case 0: 2
        return nil
    case 1: 3
        return channels[0]
    }

    orDone := make(chan interface{})
    go func() { 4
        defer close(orDone)

        switch len(channels) {
        case 2: 5
            select {
            case <-channels[0]:
            case <-channels[1]:
            }
        default: 6
            select {
            case <-channels[0]:
            case <-channels[1]:
            case <-channels[2]:
            case <-or(append(channels[3:], orDone)...): 6
            }
        }
    }()
    return orDone
}   
```

```go
sig := func(after time.Duration) <-chan interface{}{ 
    c := make(chan interface{})
    go func() {
        defer close(c)
        time.Sleep(after)
    }()
    return c
}

start := time.Now() 
<-or(
    sig(2*time.Hour),
    sig(5*time.Minute),
    sig(1*time.Second),
    sig(1*time.Hour),
    sig(1*time.Minute),
)
fmt.Printf("done after %v", time.Since(start)) 
```



### 에러 처리

```go
checkStatus := func(
  done <-chan interface{},
  urls ...string,
) <-chan *http.Response {
    responses := make(chan *http.Response)
    go func() {
        defer close(responses)
        for _, url := range urls {
            resp, err := http.Get(url)
            if err != nil {
                fmt.Println(err) 1
                continue
            }
            select {
            case <-done:
                return
            case responses <- resp:
            }
        }
    }()
    return responses
}

done := make(chan interface{})
defer close(done)

urls := []string{"https://www.google.com", "https://badhost"}
for response := range checkStatus(done, urls...) {
    fmt.Printf("Response: %v\n", response.Status)
}
```

- 해결책

```go
type Result struct { 
    Error error
    Response *http.Response
}
checkStatus := func(done <-chan interface{}, urls ...string) <-chan Result { 
    results := make(chan Result)
    go func() {
        defer close(results)

        for _, url := range urls {
            var result Result
            resp, err := http.Get(url)
            result = Result{Error: err, Response: resp} 
            select {
            case <-done:
                return
            case results <- result: 
            }
        }
    }()
    return results
}

done := make(chan interface{})
defer close(done)

urls := []string{"https://www.google.com", "https://badhost"}
for result := range checkStatus(done, urls...) {
    if result.Error != nil { 5
        fmt.Printf("error: %v", result.Error)
        continue
    }
    fmt.Printf("Response: %v\n", result.Response.Status)
}
```

### 파이프라인

시스템에서 추상화를 구성하는 데 사용하는 도구

스트림이나 데이터에 대한 일괄 처리 작업들을 처리해야 할 때 사용한다 

``` go
multiply := func(value, multiplier int) int {
    return value * multiplier
}

add := func(value, additive int) int {
    return value + additive
}

ints := []int{1, 2, 3, 4}
for _, v := range ints {
    fmt.Println(multiply(add(multiply(v, 2), 1), 2))
}
```

#### 파이프라인 구축의 모범 사례

```go
generator := func(done <-chan interface{}, integers ...int) <-chan int {
    intStream := make(chan int)
    go func() {
        defer close(intStream)
        for _, i := range integers {
            select {
            case <-done:
                return
            case intStream <- i:
            }
        }
    }()
    return intStream
}

multiply := func(
  done <-chan interface{},
  intStream <-chan int,
  multiplier int,
) <-chan int {
    multipliedStream := make(chan int)
    go func() {
        defer close(multipliedStream)
        for i := range intStream {
            select {
            case <-done:
                return
            case multipliedStream <- i*multiplier:
            }
        }
    }()
    return multipliedStream
}

add := func(
  done <-chan interface{},
  intStream <-chan int,
  additive int,
) <-chan int {
    addedStream := make(chan int)
    go func() {
        defer close(addedStream)
        for i := range intStream {
            select {
            case <-done:
                return
            case addedStream <- i+additive:
            }
        }
    }()
    return addedStream
}

done := make(chan interface{})
defer close(done)

intStream := generator(done, 1, 2, 3, 4)
pipeline :  

for v := range pipeline {
    fmt.Println(v)
}
```



### 팬 아웃, 팬 인

- 팬 아웃

  파이프라인 입력을 처리하기 위해 여러 고루틴들을 시작하는 프로세스 

- 팬 인

  여러 결과를 하나의 채널로 결합하는 프로세스 

#### 

- 단계가 이전에 계산한 값에 의존하지 않는다 
- 단계를 실행하는데 시간이 오래 걸린다 

순서 독립성은 동시에 실행되는 해당 단계의 복사본이 어떤 순서로 실행되는지 어떤 순서로 리턴할지를 보장하지 않기 때문에 중요하다 

```go
fanIn := func(
    done <-chan interface{},
    channels ...<-chan interface{},
) <-chan interface{} { 1
    var wg sync.WaitGroup 2
    multiplexedStream := make(chan interface{})

    multiplex := func(c <-chan interface{}) { 3
        defer wg.Done()
        for i := range c {
            select {
            case <-done:
                return
            case multiplexedStream <- i:
            }
        }
    }

    // Select from all the channels
    wg.Add(len(channels)) 4
    for _, c := range channels {
        go multiplex(c)
    }

    // Wait for all the reads to complete
    go func() { 5
        wg.Wait()
        close(multiplexedStream)
    }()

    return multiplexedStream
}
```



### or-done 채널

시스템에서 서로 다른 부분의 채널들로 작업하게 되는 경우 

파이프라인과 달리 작업중인 코드가 done 채널을 통해 취소될 때 채널이 어떻게 동작할지 단언할 수 없다 

고루틴이 취소됐다는 것이 읽어오는 채널 역시 취소됐음을 의미하는지는 알 수 없다 

```
loop:
for {
    select {
    case <-done:
        break loop
    case maybeVal, ok := <-myChan:
        if ok == false {
            return // or maybe break from for
        }
        // Do something with val
    }
}
```

```go
orDone := func(done, c <-chan interface{}) <-chan interface{} {
    valStream := make(chan interface{})
    go func() {
        defer close(valStream)
        for {
            select {
            case <-done:
                return
            case v, ok := <-c:
                if ok == false {
                    return
                }
                select {
                case valStream <- v:
                case <-done:
                }
            }
        }
    }()
    return valStream
}
```



### tee 채널

채널에서 들어오는 값을 분리해 코드베이스의 별개의 두 영역으로 보내고자 할 수도 있다 

채널에서 사용자 명령 스트림을 가져와서 이를 실행해줄 누군가에게 이 명령을 보내고 또 나중에 감사를 위해 명령을 기록할 누군가에게 보낼 수 있다 

```go
tee := func(
    done <-chan interface{},
    in <-chan interface{},
) (_, _ <-chan interface{}) { <-chan interface{}) {
    out1 := make(chan interface{})
    out2 := make(chan interface{})
    go func() {
        defer close(out1)
        defer close(out2)
        for val := range orDone(done, in) {
            var out1, out2 = out1, out2 1
            for i := 0; i < 2; i++ { 2
                select {
                case <-done:
                case out1<-val:
                    out1 = nil 3
                case out2<-val:
                    out2 = nil 3
                }
            }
        }
    }()
    return out1, out2
}
```



### bridge 채널

연속된 채널로부터 값을 사용하고 싶을 때 

```go
<-chan <-chan interface{}
```

```go
bridge := func(
    done <-chan interface{},
    chanStream <-chan <-chan interface{},
) <-chan interface{} {
    valStream := make(chan interface{}) 1
    go func() {
        defer close(valStream)
        for { 2
            var stream <-chan interface{}
            select {
            case maybeStream, ok := <-chanStream:
                if ok == false {
                    return
                }
                stream = maybeStream
            case <-done:
                return
            }
            for val := range orDone(done, stream) { 3
                select {
                case valStream <- val:
                case <-done:
                }
            }
        }
    }()
    return valStream
}
```

한 요소가 쓰여진 연속된 10개 채널을 생성하고 이 채널들을 bridge 함수로 전달하는 예제

```go
genVals := func() <-chan <-chan interface{} {
    chanStream := make(chan (<-chan interface{}))
    go func() {
        defer close(chanStream)
        for i := 0; i < 10; i++ {
            stream := make(chan interface{}, 1)
            stream <- i
            close(stream)
            chanStream <- stream
        }
    }()
    return chanStream
}

for v := range bridge(nil, genVals()) {
    fmt.Printf("%v ", v)
}
```



### 대기열 사용

파이프라인이 아직 준비되지 않았더라도 파이프라인에 대한 작업을 받아들이는 것

대기열을 성급하게 도입하면 데드락이나 라이브락 같은 동기화 문제가 드러나지 않을 수 있다 

대기열 사용은 프로그램의 총 실행 속도를 거의 높여주지 않는다 

다른 방식으로 동작하도록 허용할 뿐이다 

``` go
done := make(chan interface{})
defer close(done)

zeros := take(done, 3, repeat(done, 0))
short := sleep(done, 1*time.Second, zeros)
long := sleep(done, 4*time.Second, short)
pipeline := long
```

1. 끊임없이 0 스트림을 생성하는 반복 단계
2. 3개의 아이템을 받으면 그 이전 단계를 취소하는 단계
3. 1초간 슬립하는 짧은 단계
4. 4초간 슬립하는 긴 단계







## Context

```go
type Context interface {

    // Deadline returns the time when work done on behalf of this
    // context should be canceled. Deadline returns ok==false when no
    // deadline is set. Successive calls to Deadline return the same
    // results.
    Deadline() (deadline time.Time, ok bool)

    // Done returns a channel that's closed when work done on behalf
    // of this context should be canceled. Done may return nil if this
    // context can never be canceled. Successive calls to Done return
    // the same value.
    Done() <-chan struct{}

    // Err returns a non-nil error value after Done is closed. Err
    // returns Canceled if the context was canceled or
    // DeadlineExceeded if the context's deadline passed. No other
    // values for Err are defined.  After Done is closed, successive
    // calls to Err return the same value.
    Err() error

    // Value returns the value associated with this context for key,
    // or nil if no value is associated with key. Successive calls to
    // Value with the same key returns the same result.
    Value(key interface{}) interface{}
}
```

- 호출 그래프상의 분기를 취소하기 위한 API 제공
- 호출 그래프를 따라 요청 범위 데이터를 전송하기 위한 데이터 저장소의 제공



```go
func WithCancel(parent Context) (ctx Context, cancel CancelFunc)
func WithDeadline(parent Context, deadline time.Time) (Context, CancelFunc)
func WithTimeout(parent Context, timeout time.Duration) (Context, CancelFunc)
```

```go
func main() {
    var wg sync.WaitGroup
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    wg.Add(1)
    go func() {
        defer wg.Done()

        if err := printGreeting(ctx); err != nil {
            fmt.Printf("cannot print greeting: %v\n", err)
            cancel()
        }
    }()

    wg.Add(1)
    go func() {
        defer wg.Done()
        if err := printFarewell(ctx); err != nil {
            fmt.Printf("cannot print farewell: %v\n", err)
        }
    }()

    wg.Wait()
}

func printGreeting(ctx context.Context) error {
    greeting, err := genGreeting(ctx)
    if err != nil {
        return err
    }
    fmt.Printf("%s world!\n", greeting)
    return nil
}

func printFarewell(ctx context.Context) error {
    farewell, err := genFarewell(ctx)
    if err != nil {
        return err
    }
    fmt.Printf("%s world!\n", farewell)
    return nil
}

func genGreeting(ctx context.Context) (string, error) {
    ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
    defer cancel()

    switch locale, err := locale(ctx); {
    case err != nil:
        return "", err
    case locale == "EN/US":
        return "hello", nil
    }
    return "", fmt.Errorf("unsupported locale")
}

func genFarewell(ctx context.Context) (string, error) {
    switch locale, err := locale(ctx); {
    case err != nil:
        return "", err
    case locale == "EN/US":
        return "goodbye", nil
    }
    return "", fmt.Errorf("unsupported locale")
}

func locale(ctx context.Context) (string, error) {
    if deadline, ok := ctx.Deadline(); ok { 1
        if deadline.Sub(time.Now().Add(1*time.Minute)) <= 0 {
            return "", context.DeadlineExceeded
        }
    }

    select {
    case <-ctx.Done():
        return "", ctx.Err()
    case <-time.After(1 * time.Minute):
    }
    return "EN/US", nil
}
```



#### Context 데이터 저장 및 조회

```go
func main() {
    ProcessRequest("jane", "abc123")
}

func ProcessRequest(userID, authToken string) {
    ctx := context.WithValue(context.Background(), "userID", userID)
    ctx = context.WithValue(ctx, "authToken", authToken)
    HandleResponse(ctx)
}

func HandleResponse(ctx context.Context) {
    fmt.Printf(
        "handling response for %v (%v)",
        ctx.Value("userID"),
        ctx.Value("authToken"),
    )
}
```

