# Simple Go routine example 

```go
package main

import (
	"fmt"
	"time"
)

func main() {
	doneChan := make(chan bool)
	params := make(chan int, 5)
	params <- 1
	params <- 2
	params <- 3
	params <- 4
	params <- 5
	// close channel, there is no new value send to this channel.
	// NOTE : Even when channel is closed, go routine can still receive msg, but block for sent msg
	// When using `range` to get msg from channel, we need to close channel to indicate that there is no value.
	// If not closed, then `range` will continue to listen to the open channel
	close(params)

	pipeline(params, doneChan)

	select {
	case done := <-doneChan:
		fmt.Println("DONE", done)
	}
}

func pipeline(params <-chan int, doneChan chan bool) {
	go func() {
		for _ = range step3(step2(step1(params))) {
			//fmt.Println(result)
		}

		doneChan <- true // mark as done, so it could trigger `select case done`
		close(doneChan)  // close channel, all process has been completed
	}()
}

func step1(params <-chan int) <-chan string {
	ch := make(chan string, len(params))

	go func() {
		for param := range params {
			time.Sleep(1 * time.Second)
			result := fmt.Sprintf("step1(%v)", param)
			fmt.Println(result)
			ch <- result
		}

		close(ch) // close channel
	}()

	return ch
}

func step2(params <-chan string) <-chan string {
	ch := make(chan string, len(params))

	go func() {
		for param := range params {
			time.Sleep(1 * time.Second)
			result := fmt.Sprintf("step2(%v)", param)
			fmt.Println(result)
			ch <- result
		}

		close(ch) // close channel
	}()

	return ch
}

func step3(params <-chan string) <-chan string {
	ch := make(chan string, len(params))

	go func() {
		for param := range params {
			time.Sleep(1 * time.Second)
			result := fmt.Sprintf("step3(%v)", param)
			fmt.Println(result)
			ch <- result
		}

		close(ch) // close channel
	}()

	return ch
}
```

output
```
step1(1)
step2(step1(1))
step1(2)
step2(step1(2))
step1(3)
step3(step2(step1(1)))
step1(4)
step2(step1(3))
step3(step2(step1(2)))
step2(step1(4))
step3(step2(step1(3)))
step1(5)
step3(step2(step1(4)))
step2(step1(5))
step3(step2(step1(5)))
DONE true

Process finished with exit code 0
```
