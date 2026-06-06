// 创建时间：2026/6/3 下午4:04
package main

import "fmt"
/*
用来捕获那些不明显的中断错误，如数组越界,空指针。使用defer。
recover 只能救 panic，而且必须满足 2 个条件:必须写在 defer 里面,必须在 发生 panic 的同一个 goroutine。
recover的触发机理：

*/

func read(){


    defer func(){  // 1. defer 函数先被注册（此时未执行）
		err:=recover()  // 4. 这里执行 recover:发现当前 goroutine 正处于 panic 状态,捕获这个 panic,把 panic 停止 → 程序不再崩溃,返回 panic 的错误信息.
		fmt.Println(err)
	}()

	var List = []int{1,2}  // 2. 正常执行代码
	fmt.Println(List[2]) //3.数组越界触发 panic，当前 goroutine 进入 panic 状态。panic 一触发，立刻发生 3 件事：当前函数立刻停止执行、当前函数 panic 之后的所有代码，永远不执行、立即开始执行 当前函数里所有已经注册的 defer 函数。
	fmt.Println("这句话永远不会打印")  	// 5. panic 之后的代码 永远不会执行
}


func main() {
    read()
	//下面是其他正常的语句。但是都不会被执行，因为read()已经终止了。那么就需要一种恢复机制，来继续执行其他语句。
	fmt.Println("main")
}

/*

【recover 工程使用核心规则】
1. 触发 panic 后：当前函数立即停止执行 → 执行 defer → recover 捕获错误
2. 恢复成功：当前函数退出，回到上层调用处，程序继续执行
3. 恢复失败：程序直接崩溃中断

【重点：工程里绝对不用每个函数都写 defer+recover！】
正确用法：
   只在【顶层入口函数】写 1 次 defer+recover
   下面所有子函数、底层函数发生的 panic 都会被统一捕获

适用入口：main、HTTP接口、RPC接口、goroutine、定时任务
下层函数：只管正常写业务，不需要写任何 recover

【完整流程总结】
panic 触发 → 停止当前函数 → 执行 defer → recover 捕获错误
→ 清除崩溃状态 → 退出当前函数 → 回到上层继续执行
*/

