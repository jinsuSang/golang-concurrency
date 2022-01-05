# Go의 동시성 구성 요소

## 고루틴

다른 코드와 함께 동시에 실행되는 함수이다 

OS스레드가 아니고 런타임에 의해 관리되는 스레드인 그린 스레드도 아니다 

코루틴이라고 불리는 더 높은 수준의 추상화이다.

단순히 동시에 실행되는 서브루틴(함수, 클로저, Go 메서드)으로서 비선점적으로 인터럽트할 수 없다 

고루틴은 잠시 중단 suspend 하거나 재진입 reentry 할 수 있는 여러 지점이 있다 



Go 런타임은 고루틴 실행 시 동작을 관찰해, 고루틴이 멈춰서 대기 block 중일 때 자동으로 일시 중단시키고 대기가 끝나면 다시 시작시킨다 

Go 런타임이 고루틴을 선점 가능하게 해주기는 하지만 고루틴이 멈춘 지점에서만 선점 가능하다 



고루틴을 호스팅하는 Go 메커니즘은 M:N 스케줄러를 구현한 굿으로 M개의 그린 스레드를 N개의 OS 스레드에 매핑한다는 것이다.

사용 가능한 그린 스레드보다 더 많은 고루틴이 있는 경우, 스케줄러는 사용 가능한 스레드들로 고루틴을 분배하고 이 고루틴들이 대기 상태가 되면 다른 고루틴이 실행될 수 있도록 한다 

  Go 는 fork-join 모델이라는 동시성 모델을 따른다 

자식 분기가 다시 부모 분기로 합쳐지는 지점을 합류 지점이라고 한다 

```go
var wg sync.WaitGroup
for _, situation := range []string{"hello", "hi", "good day"} {
    wg.Add(1)
    go func() {
        defer wg.Done()
        fmt.Println(situation)
    }()
}
wg.Wait()
// good day
// good day
// good day
```

Go 런타임은 situation 변수에 대한 참조가 여전히 이루어져 고루틴이 계속 접근할 수 있도록 메모리를 힙으로 옮길 것이라는 사실을 알 수 있다 

고루틴이 실행되기 전에 루프가 종료되어 situation은 마지막 값은 good day에 대한 참조를 저장하고 있는 힙으로 옮겨지고 good day가 세 번 출력 된다 

```go
var wg sync.WaitGroup
for _, situation := range []string{"hello", "hi", "good day"} {
    wg.Add(1)
    go func() {
        defer wg.Done()
        fmt.Println(situation)
    }(situation)
}
wg.Wait()
// hello
// hi
// good day
```

