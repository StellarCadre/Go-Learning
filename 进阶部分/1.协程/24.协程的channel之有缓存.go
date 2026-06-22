// 创建时间：2026/6/22 下午6:39
// 创建时间：2026/6/1
package main

import (
	"fmt"
	"sync"
	"time"
)

/*
有缓冲Channel特性：
1. 声明时指定缓冲区长度，允许发送方先写入N个数据而无需接收方立即接收
2. 缓冲区未满时，发送操作不会阻塞；缓冲区满时，发送方会阻塞等待接收
3. 仍需在所有发送操作完成后关闭Channel，避免接收方读取时永久阻塞
*/

// 声明带缓冲的Channel，缓冲区长度为3（匹配协程数量）
var msgChan = make(chan string, 3)

// sendMsg 向Channel发送消息（模拟业务逻辑）
// str：发送方标识；msg：发送的消息内容；wg：等待组，用于等待所有协程完成
func sendMsg(str string, msg string, wg *sync.WaitGroup) {
	defer wg.Done() // 协程结束时将计数器减1（defer确保必执行）

	fmt.Printf("[%s] 开始发送消息\n", str)
	time.Sleep(1 * time.Second) // 模拟业务耗时

	// 发送数据到有缓冲Channel：缓冲区未满时不会阻塞
	msgChan <- fmt.Sprintf("[%s]：%s", str, msg) //负责拼接字符串，生成一段文本，发送存入 channel 缓冲区，交给接收协程读取。
	fmt.Printf("[%s] 消息发送完成\n", str)
}

func main() {
	var wg sync.WaitGroup
	start := time.Now()

	// 启动3个协程发送消息
	wg.Add(3)
	go sendMsg("协程1", "Hello Channel", &wg)
	go sendMsg("协程2", "缓冲Channel测试", &wg)
	go sendMsg("协程3", "Go并发编程", &wg)

	// 匿名协程：等待所有发送协程完成后关闭Channel
	// 核心：必须在所有发送操作完成后关闭，否则接收方遍历会阻塞
	go func() {
		wg.Wait()      // 跑在单独后台协程。专门等发送任务全部结束，负责关闭管道。关闭本身和主线程无关
		close(msgChan) // 关闭Channel，告知接收方无更多数据。关闭后，无法再向Channel写入数据，但可以继续读取Channel已有数据。
		fmt.Println("✅ 所有消息发送完成，关闭Channel")
	}()

	//下面的这两种情况，在使用带缓冲的channel时，都会正常执行。但情况2，若使用无缓冲的channel时，会出现死锁。
	// 情况1：主线程立刻进入for range循环持续接收通道数据，收发逻辑并发运行
	//var msgList []string
	//// 遍历有缓冲Channel：Channel关闭后，遍历会自动退出
	//for msg := range msgChan {
	//	fmt.Printf("📥 接收消息：%s\n", msg)
	//	msgList = append(msgList, msg)
	//}

	//情况2：主线程先执行wg.Wait()阻塞等待所有发送协程全部执行完毕，等待结束后再读取通道内所有值
	// 带缓冲通道：数据可先存入缓冲区，发送协程无需等待接收方就能完成发送、执行Done，不会死锁
	// 无缓冲通道：主线程卡在wg.Wait()阶段没有接收逻辑，发送协程执行发送时找不到接收方，全部阻塞在发送语句，无法执行Done，wg永久等待，最终触发死锁
	wg.Wait() //阻塞主线程，必须等待所有协程执行完毕
	var msgList []string
	// 遍历有缓冲Channel：Channel关闭后，遍历会自动退出
	for msg := range msgChan {
		fmt.Printf("📥 接收消息：%s\n", msg)
		msgList = append(msgList, msg)
	}

	// 输出最终结果
	fmt.Println("\n===== 执行结果 =====")
	fmt.Printf("总耗时：%v\n", time.Since(start))
	fmt.Printf("接收的消息列表：%v\n", msgList)
}
