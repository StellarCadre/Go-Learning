// 创建时间：2026/6/2 下午8:15
package main
/*
 线程不安全的意思是：多个协程函数同时对一个变量进行读写操作，会导致数据错乱。
 底层原因：sum1++/sum1-- 并非原子操作，拆解为 CPU 指令包含三步.
 解决方案：可以加互斥锁，也可以使用原子操作。
 */
import (
	"fmt"
	"sync"
)
var  waittime1 sync.WaitGroup  // 用于等待所有协程函数执行完毕
var sum1 int

func add(){
	for i := 0; i < 10000; i++ {
		sum1++
	}
	waittime1.Done()  //协程执行完毕后调用，将等待数减 1
}

func sub(){
	for i := 0; i < 10000; i++ {
		sum1--
	}
	waittime1.Done()
}

func main() {

	waittime1.Add(2)  //标记需要等待 2 个协程完成
    go add()
	go sub()
	waittime1.Wait()  //主协程阻塞，直到等待数变为 0
	fmt.Println(sum1)
}
