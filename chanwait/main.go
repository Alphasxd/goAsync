package main

import (
	"fmt"
	"sync"
	"time"
)

// waitGroup 用于等待一组 goroutine 完成
// 用 wg 重构 chandemo/main.go, 因为 sleep 在实际应用中是不可控的
func main() {
	start := time.Now()

	wg := &sync.WaitGroup{} // 一般使用指针传递，避免 wg 拷贝

	for i := 0; i < 10; i++ {
		wg.Add(1) // 每启动一个 goroutine，就把 wg 的计数器加 1
		go func(num int, wg *sync.WaitGroup) {
			time.Sleep(time.Millisecond * 100)
			fmt.Printf("goroutine %d done\n", num)
			wg.Done() // 每个 goroutine 完成后，就把 wg 的计数器减 1
		}(i, wg)
	}

	wg.Wait() // 等待 wg 的计数器减为 0
	fmt.Printf(">>> all goroutine done, time cost: %v\n", time.Since(start))
}
