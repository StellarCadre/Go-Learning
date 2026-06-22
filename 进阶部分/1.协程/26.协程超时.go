// 创建时间：2026/6/2 下午8:05
package main

/*
 协程超时的意思是：给协程设一个 “最大等待时间”，超过时间没做完，就不等了！
*/

import (
	"fmt"
	"time"
)

var done = make(chan struct{}) //声明一个struct{}类型的无缓冲通道done，用于标记协程任务是否完成

func event() {
	fmt.Println("事件执行")
	time.Sleep(2 * time.Second)
	fmt.Println("事件执行结束")
	done <- struct{}{} // 任务完成后，向done通道发送空信号

}

func main() {
	go event()

	select {
	case <-done: // 监听done通道：若收到信号，说明任务完成。
		fmt.Println("事件执行完成")
	case <-time.After(3 * time.Second): //当event的执行时间超过等待时间，则会执行超时处理，终止通道。
		fmt.Println("事件执行超时")
	}
}
