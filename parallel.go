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
	wg := sync.WaitGroup{}
	
	wg.Add(1)
	go sum(1, c, &wg)
	wg.Add(1)
	go sum(2, c, &wg)
	wg.Add(1)
	go sum(3, c, &wg)
	wg.Add(1)
	go sum(4, c, &wg)
	
	go func() {
		wg.Wait()
		close(c)
	}()
	
	for res := range c {
		fmt.Println(res)
	}
	
	elapsed := time.Since(start)

	fmt.Println("DONE", elapsed)
}
