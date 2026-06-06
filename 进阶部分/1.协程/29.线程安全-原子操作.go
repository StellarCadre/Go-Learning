// 创建时间：2026/6/2 下午8:42


package main
/*
 线程不安全的意思是：多个协程函数同时对一个变量进行读写操作，会导致数据错乱。
 底层原因：sum1++/sum1-- 并非原子操作，拆解为 CPU 指令包含三步.
 解决方案：可以加互斥锁，也可以使用原子操作。
*/
import (
	"fmt"
	"sync"
	"sync/atomic"
)
var  waittime3 sync.WaitGroup  // 用于等待所有协程函数执行完毕


var sum3 int64 // 原子操作要求int64/uint64类型
func add2(){

	for i := 0; i < 10000; i++ {
		atomic.AddInt64(&sum3, 1) // 原子加1
	}
	waittime3.Done()
}

func sub2(){

	for i := 0; i < 10000; i++ {
		atomic.AddInt64(&sum3, -1) // 原子减1
	}

	waittime3.Done()
}

func main() {

	waittime3.Add(2)
	go add2()
	go sub2()
	waittime3.Wait()
	fmt.Println(atomic.LoadInt64(&sum3)) // 原子读取值，输出0
}

/*
线程不安全 vs 线程安全 区别：

1. 线程不安全（不加锁）
   - 多个协程同时对同一个变量 读写
   - sum++ / sum-- 会被CPU指令穿插执行，导致数据错乱
   - 最终结果随机、不正确
   - 代码执行快，但数据不可靠

2. 线程安全（加互斥锁 sync.Mutex）
   - 同一时间只允许一个协程修改变量
   - 加锁 Lock() → 操作变量 → 解锁 Unlock()
   - 避免指令穿插，结果永远正确
   - 代码稍慢，但数据安全可靠
*/


/*
互斥锁 sync.Mutex 与 原子操作 sync/atomic 的区别：

1. 互斥锁（Mutex）
   • 作用：同一时间只允许一个协程进入临界区
   • 适用场景：复杂逻辑、多步骤操作、结构体/多个变量修改
   • 原理：通过锁机制实现并发控制
   • 特点：灵活、通用、使用简单
   • 性能：相对较低（有协程切换开销）

2. 原子操作（atomic）
   • 作用：对单个变量的加减赋值保证CPU指令级原子性
   • 适用场景：简单的数值增减、读取、赋值（int32/int64等）
   • 原理：硬件级原子指令，不可分割
   • 特点：执行极快、无锁、无协程阻塞
   • 性能：非常高（开销远小于锁）

总结：
简单数值操作 → 用 atomic（更快）
复杂逻辑操作 → 用 Mutex（更通用）
*/
