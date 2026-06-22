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

func Shopping_food(str string, waittime *sync.WaitGroup) {
	fmt.Println(str, "在购物")
	time.Sleep(1 * time.Second)
	fmt.Println(str, "购物结束")

	waittime.Done() //将WaitGroup的计数器值减1
}

func main() {
	var waittime sync.WaitGroup // 声明WaitGroup变量，初始计数器为0

	startTime := time.Now()

	waittime.Add(3) // 这里n=3，因为要启动3个协程(调用3次协程函数)，每个协程最终会调用1次Done()

	go Shopping_food("张三", &waittime) // 启动3个协程，注意传递WaitGroup的指针（避免复制）
	go Shopping_food("李四", &waittime)
	go Shopping_food("王五", &waittime)

	waittime.Wait() //阻塞当前主协程，直到计数器归0（即所有子协程执行完毕），再继续执行后续逻辑。

	fmt.Println("购买业务完成,总耗时：", time.Since(startTime)) // 所有协程执行完成后，才会执行到这里

}

/*
// ======== sync.WaitGroup 核心概念 ========
	// sync.WaitGroup是Go标准库sync包中的同步原语，用于等待一组协程（goroutine）执行完成
	// 核心原理：通过一个内部计数器实现，计数器值可以增加/减少，Wait()会阻塞直到计数器归0
	// 注意事项：
	// 1. WaitGroup实例不能被复制（如传递时必须用指针），否则会导致计数器状态异常
	// 2. 计数器不能为负数，否则会panic
	// 3. 一旦计数器归0，再次调用Add()会重置计数器（不推荐这样使用）

// ======== sync.WaitGroup 核心方法1：Done() ========
	// 作用：将WaitGroup的计数器值减1.必须存在于每个协程函数的函数体中。
	// 注意：必须确保每个Add()对应的协程都调用Done()，否则会导致Wait()永久阻塞或panic
	// 底层：本质是调用Add(-1)，但直接用Done()更语义化

// ======== sync.WaitGroup 核心方法2：Add(n int) ========
	// 作用：将WaitGroup的计数器值增加n
	// 使用时机：必须在启动协程前调用（避免协程已经执行完，Done()先于Add()调用导致计数器负数）

// ======== sync.WaitGroup 核心方法3：Wait() ========
	// 作用：阻塞当前协程（这里是main协程），直到WaitGroup的计数器值归0
	// 使用场景：主线程等待所有子协程执行完成后，再继续执行后续逻辑
	// 底层：会一直阻塞，直到计数器变为0，不会主动退出（所以必须确保所有协程都调用Done()）
    //使用时机：必须在所有子协程启动后调用（避免主协程提前退出，导致子协程没有执行完毕）

// 补充说明：
	// 1. WaitGroup适用于“一对多”的同步场景（如主线程等待多个子协程）
          WaitGroup 的定位：单纯同步等待工具，只管 “等不等”，不管 “数据互通”。
	// 2. 如果需要在协程间传递数据，比如解决如下问题，需要用到channel：
          主线程定义变量，子协程可读，但并发修改会有竞争问题。
          子协程内部生成的结果，主线程没法直接获取；
          多个子协程之间也无法互相交换计算结果。
	// 3. 全局WaitGroup也可以使用，但推荐按需声明（减少耦合）

协程间传递数据 = 在不同 goroutine 之间安全地发送 / 接收数据，实现数据互通。
核心载体就是 channel（管道）：
一个协程往 channel 塞数据（发送）；
另一个协程从 channel 取数据（接收）；
底层自带同步机制，天然解决并发读写安全问题。
*/
