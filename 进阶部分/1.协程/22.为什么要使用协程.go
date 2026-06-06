// 创建时间：2026/6/1 下午6:01
package main
/*
 协程一般和带时间延迟的函数有关
 */

import (
	"fmt"
	"time"
)

var wait int


func Shopping(str string){
	fmt.Println(str,"在购物")
	time.Sleep(1*time.Second)
	fmt.Println(str,"购物结束")

	wait--
}

func main() {
	startTime:=time.Now()
	wait=3
	//这种模式属于接力，只有等前一个函数结束之后，才能运行当前函数，容易造成时间线拉长。如以下函数总时间3s。（先来先服务）
    //Shopping("张三")
	//Shopping("李四")
	//Shopping("王五")

	//这种属于Go的协程方式，即同时运行下面的三条语句。（时间片轮转法吗？）
	/*
	但是目前的该方案存在严重的问题，主进程没有上锁，程序在执行完下面的第三条语句后，主进程提前关闭，导致整个程序终止，并没有等全部执行完毕。
    主线程结束，协程函数被迫终止。
	 */
	go Shopping("张三")  //函数前加go，表示该函数是协程函数
	go Shopping("李四")
	go Shopping("王五")

    /*
	为了解决上面的问题，可以在此处进行等待，如time.Sleep(a*time.Second)。
	但是并不知道需要等待多长时间，即a不确定。
	故，可以考虑使用一个变量wait来接收总耗时，然后进行判断,如下：
	但是当前的解决方案，可能遇到这些问题：
	多协程同时修改变量wait，数据错乱。死循环空转 → 疯狂消耗 CPU。代码脆弱，难维护、难扩展。
	综上所述，必须引入go中标准的sync.WaitGroup，见下一个代码文件。
     */
    for {
		if wait==0 {
			break
		}
	}

	fmt.Println("购买业务完成,总耗时：", time.Since(startTime))



}
