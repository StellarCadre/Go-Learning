// 创建时间：2026/6/3 下午4:45
package main

import "fmt"

//func plus(a int, b int) int {
//	return a + b
//}  由上面这个升级为下面这个：

// 类型约束：int 或 uint
func add[T int| uint ](a T, b T) T {  //这里的T是类型占位符（类型参数），可以是int或uint。[T int|uint] 是类型参数列表
    return a + b
}

//更多类型版本
func test1[T int|float64,K string|int|float32] (a T,b K){ //这里的T是int或float64，K是string或int或float32
	fmt.Println(a)
	fmt.Println(b)
}

//给类型包装成整体，不必都写出来
type TntType interface {
	int|int8|int16|int32|int64
	uint|uint8|uint16|uint32|uint64
}
func test2[T TntType](a T,b T){
	fmt.Println(a)
	fmt.Println(b)
}




func main() {

	result1:=add(1, 2)  // 自动推导 T 为 int
	fmt.Println(result1)

    var u1,u2=uint(1),uint(3)
	result2:=add(u1, u2)  // 自动推导 T 为 uint
	fmt.Println(result2)

}



/*
【Go 泛型 学习前瞻】

1. 泛型是什么？
   一句话：让函数/结构体支持“多种类型”，但不用写重复代码。
   以前要给 int、string 各写一套，现在一套通用。

2. 为什么要用泛型？
   • 消除重复代码（比如写一个通用切片工具函数）
   • 保留类型安全（不会像 interface{} 一样丢失类型）
   • 代码更简洁、更通用

3. 泛型用在哪里？
   • 通用工具函数（切片、map、队列、栈、链表）
   • 通用数据结构（通用切片、通用树、通用链表）
   • 不想用 interface{} 丢失类型检查的场景

4. 泛型核心三要素
   • 类型参数 [T any] 或 [T TypeConstraint]
   • 类型约束 any / comparable / 自定义
   • 类型推导（不用手动指定类型，自动推断）

5. 最简单理解
   泛型 = “模板”
   T = 占位符（代表未来传入的任意类型）

6. 新手最容易懵的符号
   [T any] —— 这就是泛型！表示 T 可以是任何类型

7. 工程意义
   让 Go 拥有真正的通用工具库，
   但不丢失类型安全，不降低性能。

一句话总结：
泛型 = 一套代码，支持所有类型，且类型安全！
*/

