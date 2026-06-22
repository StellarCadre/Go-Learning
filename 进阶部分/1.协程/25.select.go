// 创建时间：2026/6/2 下午6:41
package main

import (
	"fmt"
	"sync"
	"time"
)

/*
一个协程函数，不只向一个channel发送数据。需要select来选择
*/

var AmountChan = make(chan int)  // 传输金额的通道
var nameChan = make(chan string) // 传输姓名的通道

var doneChan = make(chan struct{}) // 用于标记所有发送协程完成的"信号通道"（无数据，仅用关闭状态）

// 每个send协程会向两个通道分别写入数据
func send(name string, money int, wait1 *sync.WaitGroup) {
	fmt.Println(name, "在发送数据")
	time.Sleep(1 * time.Second)
	fmt.Println(name, "发送数据结束")
	AmountChan <- money // 向金额通道写数据
	nameChan <- name    // 向姓名通道写数据
	wait1.Done()        // 标记当前协程完成
}

func main() {
	var wait1 sync.WaitGroup // 用于等待所有协程函数执行完毕
	start := time.Now()
	wait1.Add(2)
	go send("张三", 2, &wait1)
	go send("李四", 5, &wait1)

	go func() { //使用一个协程函数，专门用来等待所有协程执行完毕，然后关闭信道。
		wait1.Wait()      // 等待所有发送协程完成
		close(AmountChan) // 关闭金额通道
		close(nameChan)   // 关闭姓名通道

		close(doneChan) // 关闭信号通道，用于触发 select 的退出分支
	}()

	var nameList []string
	var AmountList []int

	/*
		写法1：
		不能把两个 for range 都写在 main 主线程里，只能一个放主线程、另一个新开后台协程。
		主线程是串行执行代码，先写 for money := range AmountChan，主线程会直接卡在这个循环里，永远走不到后面 for name := range nameChan，第二个通道完全没机会接收数据，直接卡死。
		功能：开一个后台协程单独接收 AmountChan，主线程单独接收 nameChan，分开两个各自循环读取不同通道，最终能收集到两个通道的数据。
	    可能的问题：
		核心问题 1：两个协程都会阻塞，直到通道关闭。
		后台协程：range 遍历通道规则：通道有数据就读取；通道无数据且未关闭时，永久阻塞；通道关闭，循环自动结束。
		        所以这个后台协程会一直卡在循环里，直到 AmountChan 被 close。
	    主协程：主线程会卡在这一行不动，不会往下执行打印结果代码，必须等到 nameChan 关闭，循环退出，才会执行后面的输出语句。
	    核心问题 2：无法实时同时监听两个通道：
		主线程只盯着 nameChan，看不到 AmountChan 的数据到达。后台协程只盯着 AmountChan，看不到 nameChan 的数据到达。
	    不能在同一个协程内同时监控两个通道哪个先来数据，只能分开两个协程各自看管一个通道。
	    问题 3：拓展性极差：
		如果以后再加第 3、第 4、第 5 个通道，你就要不停新建接收协程，代码到处分散，不好管理。
	*/
	//go func() {  // 协程A：后台收 AmountChan
	//	for money := range AmountChan {
	//		AmountList = append(AmountList, money)
	//	}
	//}()
	//for name:=range nameChan {  // 主协程：收 nameChan
	//	nameList = append(nameList, name)
	//}
	//fmt.Println("业务完成,总耗时：", time.Since(start))
	//fmt.Println("moneyList",AmountList)
	//fmt.Println("nameList",nameList)

	//写法2：
	//下面是使用select来实现两个信道的数据接收。同时，一定要明确for循环的退出条件。
	//只用主线程单个循环，同时监听 AmountChan、nameChan、doneChan 三个通道，不用拆分多个接收协程，解决前面 “一个主线程、一个后台协程分开接收” 的所有缺陷。
	//无限循环，会反复执行 select 多路监听，直到内部执行return跳出循环、结束主线程。
	for {
		//关键特性：select 会同时监视所有 case 通道，哪个通道先有数据就绪，就执行对应分支；两个通道同时来数据则随机选一个执行。
		select {
		case account := <-AmountChan: //监听金额通道，通道有数据到达时触发。执行完立刻回到外层 for，再次进入 select 继续监听所有通道。
			AmountList = append(AmountList, account)
		case name := <-nameChan: //同理。
			nameList = append(nameList, name)
		case <-doneChan: // 来到这里，（doneChan关闭），表示所有发送协程都完成了，可以退出循环了。
			fmt.Println("done")
			fmt.Println("业务完成,总耗时：", time.Since(start))
			fmt.Println("moneyList", AmountList)
			fmt.Println("nameList", nameList)
			return
		}
	}

}
