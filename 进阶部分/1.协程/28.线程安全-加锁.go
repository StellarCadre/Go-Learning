// 创建时间：2026/6/2 下午8:28

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
var  waittime2 sync.WaitGroup  // 用于等待所有协程函数执行完毕
var lock sync.Mutex  //这是互斥锁的定义
var sum2 int

func add1(){
	lock.Lock()  //加锁操作，保证同一时间只有一个协程函数在执行
	defer lock.Unlock() // 函数退出时自动解锁
	for i := 0; i < 10000; i++ {
		sum2++
	}
	//lock.Unlock()  //解锁操作，解锁后其他协程函数可以继续执行
	waittime2.Done()
}

func sub1(){
	lock.Lock()
	defer lock.Unlock() // 函数退出时自动解锁
	for i := 0; i < 10000; i++ {
		sum2--
	}
	//lock.Unlock()
	waittime2.Done()
}

func main() {

	waittime2.Add(2)
	go add1()
	go sub1()
	waittime2.Wait()
	fmt.Println(sum2)
}

