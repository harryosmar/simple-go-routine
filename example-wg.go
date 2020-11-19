package main

import (
	"fmt"
	"sync"
	"time"
)

func sum(num int, c chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()
	time.Sleep(1 * time.Second)
	c <- num + 100
}

func main() {
	start := time.Now()
	c := make(chan int)
	wgCond := sync.WaitGroup{}
	
	wgCond.Add(1)
	go sum(1, c, &wgCond)
	wgCond.Add(1)
	go sum(2, c, &wgCond)
	wgCond.Add(1)
	go sum(3, c, &wgCond)
	wgCond.Add(1)
	go sum(4, c, &wgCond)
	
	go func() {
		wgCond.Wait()
		close(c)
	}()
	
	for res := range c {
		fmt.Println(res)
	}
	
	elapsed := time.Since(start)

	fmt.Println("DONE", elapsed)
}
