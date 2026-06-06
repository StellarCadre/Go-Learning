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


var AmountChan = make(chan int)  //信道1
var nameChan = make(chan string) //信道2

var doneChan = make(chan struct{})  // 该信道不是用来发送数据的，而是用来关闭下面的for循环

func  send(name string,money int,wait1 *sync.WaitGroup)  {
	 fmt.Println(name,"在发送数据")
	 time.Sleep(1*time.Second)
	 fmt.Println(name,"发送数据结束")
	 AmountChan <- money
	 nameChan <- name
	 wait1.Done()
}


func main() {
	var wait1 sync.WaitGroup   // 用于等待所有协程函数执行完毕
	start := time.Now()
	wait1.Add(2)
	go send("张三",2,&wait1)
	go send("李四",5,&wait1)

	go func() { //使用一个协程函数，专门用来等待所有协程执行完毕，然后关闭信道。
		wait1.Wait()
		close(AmountChan)
		close(nameChan)

		close(doneChan)
	}()

	var nameList []string
	var AmountList []int

	/*
	下面这个部分，能够实现两个信道的数据接收。
	但是，
	主协程会死死卡在 nameChan 上，直到通道关闭，否则永远不退出
	后台协程会死死卡在 AmountChan 上，直到通道关闭，否则永远不退出
	无法同时监听两个通道，只能一个前台、一个后台
	只能先后处理通道，不能同时处理。
	 */
	//go func() {  //// 协程A：后台收 AmountChan
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

	//下面是使用select来实现两个信道的数据接收。同时，一定要明确for循环的退出条件。
    for{
		select {
		case account:=<-AmountChan:
			AmountList = append(AmountList, account)
		case name:=<-nameChan:
			nameList = append(nameList, name)
		case <-doneChan:  // 当doneChan关闭时，for循环会跳出
			fmt.Println("done")
			fmt.Println("业务完成,总耗时：", time.Since(start))
			fmt.Println("moneyList",AmountList)
			fmt.Println("nameList",nameList)
			return
		}
	}




}
