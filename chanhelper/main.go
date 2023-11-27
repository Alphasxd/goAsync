package main

import "fmt"

func main() {
	msgChan := make(chan int, 10) // buffered channel
	for i := 0; i < 10; i++ {
		msgChan <- i
	}

	// close 函数，用于关闭 channel，仅适用于发送者和双向channel
	// 由 sender 关闭，而不是 receiver，因为 receiver 可能不知道什么时候关闭
	// 关闭channel后，无法再向其发送任何数据，否则会引发panic，但是可以继续从channel接收数据
	// 当关闭的channel中最后一个元素被接收后， x, ok := <-ch 中的 ok 会被设置为 false，x 会是对应类型的零值
	close(msgChan) // 这里在主线程中关闭，主线程作为 sender

	// range 函数，用于迭代不断操作的 channel，直到 channel 被关闭
	for item := range msgChan {
		fmt.Println(item)
	}

	// 也可以使用 for 循环来迭代 channel，但是需要手动判断 channel 是否关闭
	//for {
	//	item, ok := <-msgChan
	//	if !ok {
	//		break
	//	}
	//	fmt.Println(item)
	//}

	fmt.Println(">>> Game Over <<<")
}