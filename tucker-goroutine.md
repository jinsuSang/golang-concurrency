# Go 동시성 기초 및 활용

- Tucker의 Go 언어 프로그래밍 chapter 24

---

## 스레드

- 고루틴은 경량 스레드로 함수나 명령을 동시에 실행할 때 사용한다.
- `main()` 함수도 고루틴이다
- 프로세스는 한 개 이상의 스레드를 가지며 싱글 스레드, 멀티 스레드 프로세스로 나누어 진다
- 스레드가 CPU 코어를 빠르게 교대로 점유하면서 모든 스레드가 동작하는 것처럼 보인다

## 컨텍스트 스위치

- CPU 코어가 여러 스레드를 전환하면서 수행시 발생하는 비용을 컨텍스트 스위치 비용이라고 한다
- 전환을 위해 저장되는 스레드 정보는 스레드 컨텍스트라고 한다 
- 많은 스레드를 전개하면 성능이 저하된다. 하지만 golang은 CPU 코어마다 OS 스레드 하나만 할당해 사용해 컨텍스트 스위치 비용이 발생하지 않는다

## 고루틴 사용

```go
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

// 0 1 가 2 3 4 나 5 6 7 다 8 9 10 라 11 12 13 마 14 15 바 16 17 18 사 19
```

- 메인 함수가 종료되면 고루틴들은 모두 종료되니 주의해야 한다

### WaitGroup

```go
var wg sync.WaitGroup

func SumAtoB(a, b int){
   sum := 0
   for i := a; i < b; i++ {
      sum += i
   }
   fmt.Printf("%d\n", sum)
   wg.Done()
}

func main() {
   wg.Add(20)
   for i := 0; i < 20; i++ {
      go SumAtoB(1, 1000000)
   }
   wg.Wait() // 작업 완료 대기
}
```

## Goroutine 동작 방법

- 고루틴은 명령을 수행하는 단일 흐름으로 OS 스레드를 이용하는 경량 스레드이다
- CPU 코어가 두 개인 경우 

#### 고루틴 하나일 때

- `main()` 루틴만 존재하면 OS 스레드를 하나 만들어 첫 번째 코어와 연결한다. 
- OS 스레드에서 고루틴을 실행하게 된다

#### 고루틴 두 개일 때

- 두번째 OS 스레드를 생성하여 두 번째 고루틴을 실행한다

#### 고루틴 세 개 이상일 때

- CPU 코어가 두 개이므로 고루틴이 대기한다
- 고루틴이 모두 완료되면 고루틴이 교체된다

### 시스템 콜 호출 

- 대기 상태인 고루틴에 CPU 코어와 스레드를 할당하면 CPU 자원 낭비가 발생한다 
- 그래서 Golang 에서는 대기 상태로 보내고 실행을 기다리는 다른 루틴에 CPU 코어와 OS 스레드를 할당하여 실행된다
- 코어와 스레드는 변경되지 않고 고루틴만 옮겨 다니므로 컨텍스트 스위칭 비용이 발생하지 않는다

## 동시성 프로그래밍 주의점과 해결방법

- 동일 메모리 자원에 여러 고루틴이 접근할 때 문제가 발생한다

### 뮤텍스

- mutual exclusion, 상호 배제

- Lock 을 획득한 고루틴만 실행을 한다

  문제점

  1. 동시성 프로그래밍을 통항 성능 향상이 되지 않는다.
  2. 데드락이 발생할 가능성이 있다

  ```go
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
  ```

  #### 데드락 발생 상황

  - A 가 포크를 가지고 B가 스푼을 가지는 순간 A와 B는 다른 자원을 가질 수 없다

  ```go
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
  ```

  - 멀티코어 환경에서 여러 고루틴을 사용하여 성능을 향상시킬 수 있다
  - 같은 메모리 자원에 여러 고루틴이 접근하면 문제가 발생한다
  - 뮤텍스를 활용하는 방법이 있으나 성능 향상이 안되고 데드락이 발생할 수 있다

## 자원 관리 기법

- 영역 분할

- 역할 분할

  ### 영역 분할

  - jobList 슬라이스를 활용하여 고루틴이 서로 접근하지 못하도록 설정하였다 

  ```go
  type Job interface {
     Do()
  }
  
  type SquareJob struct {
     index int
  }
  
  func (j *SquareJob) Do() {
     fmt.Println(j.index, "작업 시작")
     time.Sleep(1 * time.Second)
     fmt.Println(j.index, j.index*j.index)
  }
  
  func main() {
     var jobList [10]Job
  
     for i := 0; i < 10; i++ {
        jobList[i] = &SquareJob{i + 1}
     }
  
     var wg sync.WaitGroup
     wg.Add(10)
  
     for i := 0; i < 10; i++ {
        job := jobList[i]
        go func() {
           job.Do()
           wg.Done()
        }()
     }
  
     wg.Wait()
  }
  ```

