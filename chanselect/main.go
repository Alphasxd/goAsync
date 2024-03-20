package main

import (
	"fmt"
	"time"
)

// MsgServer 消息服务器
// 业务逻辑 channel: msgChan
// 关闭信号 channel: quitChan
// 业务逻辑和关闭信号分开，是更好的设计
type MsgServer struct {
	msgChan chan string
	// quitChan 用于接收关闭信号，sruct{} 是空结构体，不占用内存
	// 专门用于传递关闭信号，不需要传递任何数据，通常用于通知其他 goroutine 停止工作
	// 这种方式的好处是不需要使用锁，因为 struct{}{} 是不可变的，不需要加锁
	quitChan chan struct{}
}

// NewMsgServer 创建消息服务器
func NewMsgServer() *MsgServer {
	return &MsgServer{
		make(chan string, 100),
		make(chan struct{}),
	}
}

// sendMsg 发送消息，将格式化后的消息发送到消息通道
func sendMsg(ms *MsgServer, msg string, num int) {
	for i := 0; i < num; i++ {
		ms.msgChan <- fmt.Sprintf("message %d: %s", i, msg)
	}
}

// handleMsg 处理消息，从消息通道中读取消息并打印
func (ms *MsgServer) handleMsg(msg string) {
	fmt.Println(msg)
}

// work 工作，从消息通道中读取消息并处理，直到 quitChan 收到关闭信号
func (ms *MsgServer) work() {
	// 标签，用于 break 跳出循环
msgloop:
	for {
		// select 语句，用于同时处理多个通道，实现多路复用
		select {
		case msg := <-ms.msgChan:
			ms.handleMsg(msg)
		case <-ms.quitChan:
			// break loop tag，而不是跳出 select
			break msgloop
		}
	}
	fmt.Println(">>> Game Over <<<")
}

// stop 停止，向 quitChan 发送关闭信号
func (ms *MsgServer) stop() {
	// close(ms.quitChan)
	// struct{}{} 是空结构体的实例
	ms.quitChan <- struct{}{}
}

func main() {
	ms := NewMsgServer()
	go func() {
		// 1 秒后停止消息服务器，调用 stop 方法
		time.Sleep(time.Second)
		ms.stop()
	}()
	sendMsg(ms, "hello", 10)
	ms.work()
}
