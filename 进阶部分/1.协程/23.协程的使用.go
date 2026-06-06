// 创建时间：2026/6/1 下午7:00
package main
/*
 协程一般和带时间延迟的函数有关
*/

import (
	"fmt"
	"sync"
	"time"
)

func Shopping_food(str string,waittime *sync.WaitGroup){
	fmt.Println(str,"在购物")
	time.Sleep(1*time.Second)
	fmt.Println(str,"购物结束")

	waittime.Done()
}

func main() {
	var  waittime sync.WaitGroup

	startTime:=time.Now()

	waittime.Add(3)  // 下面调用了几次协程函数，这里就加几次

	go Shopping_food("张三",&waittime)  //函数前加go，表示该函数是协程函数
	go Shopping_food("李四",&waittime)
	go Shopping_food("王五",&waittime)

	waittime.Wait()

	fmt.Println("购买业务完成,总耗时：", time.Since(startTime))

//此外，对于简单程序，也可将var  waittime sync.WaitGroup 写在外面，作为全局变量

//我若想在外面拿到协程函数内容处理后的结果，可以用channel，见24.24.协程的channel.go

}
