package main

import (
	"fmt"
	"time"
)

func main() {
	start := time.Now()
	doneChan := make(chan bool)
	in := make(chan int, 5)
	in <- 1
	in <- 2
	in <- 3
	in <- 4
	in <- 5
	close(in)

	var result []<-chan int
	for i := 0; i < len(in); i++ {
		result = append(result, sq2(in, doneChan))
	}

	for _, output := range result {
		fmt.Println(<-output)
	}

	select {
	case done := <-doneChan:
		elapsed := time.Since(start)
		fmt.Println("DONE", done, elapsed)
	}
}

func sq2(in <-chan int, doneChan chan<- bool) <-chan int {
	out := make(chan int)
	go func() {
		for n := range in {
			time.Sleep(1 * time.Second)
			out <- n * n
		}
		close(out)
		doneChan <- true // mark as done, so it could trigger `select case done`
	}()
	return out
}
