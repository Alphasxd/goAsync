package main

import (
	"fmt"
	"sync"
)

type Counter struct {
	sync.Mutex
	num int // int64
}

func (c *Counter) Inc() {
	c.Lock()
	defer c.Unlock()
	c.num++
	// atomic.AddInt64(&c.num, 1)
}

func main() {
	c := &Counter{}
	wg := &sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			c.Inc()
			wg.Done()
		}()
	}
	wg.Wait()
	fmt.Println(c.num)
}
