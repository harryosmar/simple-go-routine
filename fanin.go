package main

import (
	"fmt"
	"sync"
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

	// FANOUT
	c1 := sq3(in, doneChan)
	c2 := sq3(in, doneChan)
	c3 := sq3(in, doneChan)
	c4 := sq3(in, doneChan)
	c5 := sq3(in, doneChan)


	// FANIN : Consume the merged output from c1-c5
	for n := range merge(c1, c2, c3, c4, c5) {
		fmt.Println(n)
	}

	select {
	case done := <-doneChan:
		elapsed := time.Since(start)
		fmt.Println("DONE", done, elapsed)
	}
}

func merge(cs ...<-chan int) <-chan int {
	var wg sync.WaitGroup
	out := make(chan int)

	// Start an output goroutine for each input channel in cs.  output
	// copies values from c to out until c is closed, then calls wg.Done.
	output := func(c <-chan int) {
		for n := range c {
			out <- n
		}
		wg.Done()
	}
	wg.Add(len(cs))
	for _, c := range cs {
		go output(c)
	}

	// Start a goroutine to close out once all the output goroutines are
	// done.  This must start after the wg.Add call.
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

func sq3(in <-chan int, doneChan chan<- bool) <-chan int {
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