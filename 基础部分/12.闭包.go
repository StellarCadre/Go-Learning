package main
/*
保存私有变量，不销毁、不污染全局
制造 “带记忆” 的函数（计数器、累加器）
批量生成功能相似的函数（工厂函数）
在延迟 / 异步里安全携带外部变量
 */
import "fmt"

// 🔥 这是一个 工厂函数：返回一个【闭包函数】
// 闭包 = 内部函数 + 它捕获的外部变量 num

//外层函数：只是用来创造环境的（放变量 + 放内部函数）
func addFactory() func() int {   //这里的func（）和int是addFactory的两个返回值类型
	// 👇 这个变量会被闭包【永久持有】，不会被销毁
	num := 0

	// 内部函数 + 它能访问到的外部变量称为闭包
	return func() int {
		num++       // 闭包内部可以修改外部的 num
		return num  // 闭包内部可以访问外部的 num
	}
}

func main() {
	// 1. 创建一个闭包（得到一个函数）
	// add 现在是一个闭包，它带着自己的 num
	add := addFactory()  //只在这里调用了 1 次addFactory()，创建 num = 0，返回闭包函数，退出 addFactory，但 num 没有被销毁！被闭包带走了！

	// 2. 多次调用，你会发现 num 会一直累加！
	fmt.Println(add()) // 1   后面这些是调用闭包，不会再进 addFactory，所以 num 不会重新初始化
	fmt.Println(add()) // 2
	fmt.Println(add()) // 3
	fmt.Println(add()) // 4
}