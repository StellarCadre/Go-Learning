// 创建时间：2026/6/2 下午8:05
package main

/*
 协程超时的意思是：给协程设一个 “最大等待时间”，超过时间没做完，就不等了！
 */

import (
	"fmt"
	"time"
)

var done = make(chan struct{})

func event()  {
    fmt.Println("事件执行")
	time.Sleep(2*time.Second)
	fmt.Println("事件执行结束")
}

func main() {
    go event()

	select {
	case <-done:
		fmt.Println("事件执行完成")
	case <-time.After(3*time.Second):
		fmt.Println("事件执行超时")
	}
}
