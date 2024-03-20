package main

import (
	"fmt"
	"time"
)

// channel 是 goroutine 之间的通信机制，
// 它可以用来传递数据，也可以用来同步 goroutine 的执行。

func doSome(num int) {
	time.Sleep(time.Millisecond * 100)
	fmt.Printf("doSome %d\n", num)
}

func asyncDoSomeBack(num int, resChan chan string) {
	time.Sleep(time.Millisecond * 100)
	resChan <- fmt.Sprintf("asyncDoSomeBack %d\n", num)
}

func main() { // 主协程

	// 1. 串行执行
	start := time.Now()
	for i := 0; i < 10; i++ {
		doSome(i)
	}
	fmt.Println(">>> sequence use time: ", time.Since(start))
	fmt.Println("=====================================")

	// 2. 并行执行
	start = time.Now()
	for i := 0; i < 10; i++ {
		go doSome(i) // 启用子协程
	}
	// 等待所有协程执行完毕, 否则主协程退出后子协程并不会执行
	time.Sleep(time.Millisecond * 110)
	fmt.Println(">>> parallel use time: ", time.Since(start))
	fmt.Println("=====================================")

	// 3. 并行执行, 并且获取子协程的返回值
	start = time.Now()
	resChan := make(chan string) // unbuffered channel
	for i := 0; i < 10; i++ {
		go asyncDoSomeBack(i, resChan) // 子协程写入channel
	}
	// 这里需要重新开始一个循环, 因为每次启动子协程都会阻塞主协程, 直到子协程写入channel
	// 相当于串行执行, 所以需要在子协程中写入channel, 主协程中读取channel
	for i := 0; i < 10; i++ {
		fmt.Print(<-resChan) // 主协程读取channel
	}
	fmt.Println(">>> async use time: ", time.Since(start))
	fmt.Println(">>> Game Over <<<")
}
