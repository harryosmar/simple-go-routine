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

	c1 := sq(in, doneChan)
	c2 := sq(in, doneChan)
	c3 := sq(in, doneChan)
	c4 := sq(in, doneChan)
	c5 := sq(in, doneChan)
	fmt.Println(<-c1, <-c2, <-c3, <-c4, <-c5)

	select {
	case done := <-doneChan:
		elapsed := time.Since(start)
		fmt.Println("DONE", done, elapsed)
	}
}

func sq(in <-chan int, doneChan chan<- bool) <-chan int {
	out := make(chan int, len(in))
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