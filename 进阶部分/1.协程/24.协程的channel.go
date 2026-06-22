// 创建时间：2026/6/1 下午7:32
package main

import (
	"fmt"
	"sync"
	"time"
)

/*
使用channel（信道），能够解决：
主线程定义变量，子协程可读，但并发修改会有竞争问题。
协程内部生成的结果，主线程没法直接获取；
多个子协程之间也无法互相交换计算结果。

由不能传数据的协程调整为能传数据的协程，需要设置的内容：
1.首先是信道的声明和初始化。
2.然后是sender和receiver的写法。
3.最后最关键的是一定要在所有子协程执行完毕后，关闭信道，避免阻塞。可以将关闭信道的逻辑和原本协程等待写在一个匿名函数中，放在主进程中。
*/

// 声明并初始化一个信道moneyChan
var moneyChan = make(chan int) //长度为0。为0表示无缓冲信道，只有sender和receiver都准备好时，才能传递数据。

func pay(str string, paymoney int, wait *sync.WaitGroup) {
	fmt.Println(str, "购物开始")
	time.Sleep(1 * time.Second)
	fmt.Println(str, "购物结束")

	moneyChan <- paymoney //sender：将数据paymoney写入信道moneyChan。若仅有sender而没有receiver，则会一直阻塞。

	wait.Done() //将WaitGroup的计数器值减1
}

func main() {
	var wait sync.WaitGroup
	start := time.Now()
	wait.Add(3)
	go pay("张三", 2, &wait)
	go pay("李四", 5, &wait)
	go pay("王五", 3, &wait)

	//这里的整个逻辑是：我使用一个协程，匿名的，专门用来等待所有协程执行完毕，然后关闭信道。
	//我去后台等着，等所有协程跑完，我再关闭 channel！
	//我不耽误你主协程继续往下跑！
	go func() { //两个功能：1.等待所有协程执行完毕；2.关闭信道
		wait.Wait()      //阻塞当前主协程，直到计数器归0
		close(moneyChan) //等所有协程跑完，我再关闭 channel
	}()

	var moneyList []int

	for money := range moneyChan { //遍历多个moneyChan。若信道关闭，则for循环会自动跳出  //receiver：从信道读取数据。若仅有receiver而没有sender，则会一直阻塞。
		fmt.Println("收到的钱：", money)
		moneyList = append(moneyList, money)
	} //等同于下面这些
	//for{
	//	money,ok:=<- moneyChan  //receiver：从信道读取数据。若仅有receiver而没有sender，则会一直阻塞。
	//	fmt.Println("收到的钱：",money,"是否正常：",ok)
	//	if !ok{  //不正常表示信道关闭了
	//		break
	//	}
	//}
	/*
	    这里可能阻塞的原因：
		3 个协程执行完 pay 函数后，会向 moneyChan 发送 3 次数据。
		main 函数的 for 循环能接收这 3 次数据，但接收完 3 次数据后，for 循环会继续执行 money,ok:=<- moneyChan。
		此时，无缓冲 Channel 没有任何发送方（3 个协程已执行完）
		Channel 未被关闭，因此接收操作会阻塞.
		要解决该问题，需要出现ok为false的情况，然后关闭信道。
		需要再开辟一个协程，专门用于关闭信道。等待所有pay协程执行完后，关闭Channel，见上方第36行。
	*/

	fmt.Println("购买业务完成,总耗时：", time.Since(start))
	fmt.Println("moneyList", moneyList)
}

/*
主协程：启动匿名协程 → 它去后台等待
主协程：马上执行 for开始收数据
3 个 pay 1.协程：并发执行，陆续发送数据
主协程：一边收数据，一边等待
匿名协程：等到所有 pay 结束 → 关闭 channel
主协程：发现 channel 关闭 → 退出循环
*/
