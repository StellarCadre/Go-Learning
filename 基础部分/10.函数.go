// 创建时间：2026/5/27 下午6:38
package main

import (
	"fmt"
)

//一段封装了特定功能的代码，用于执行特定任务

// 整数加法器
func add(a int, b int) int { //函数带返回值的函数：后面的int是返回值的类型
	return a + b
}

// 整数减法器
func sub(a int, b int) (int, bool) { //函数带多个返回值的函数 (int,bool)是返回值的类型
	return a - b, true
}

// 整数乘法器
func mul(a int, b int) (miltiple int, statue bool) { //函数带多个返回值的函数。还可以把返回值也定义为变量，由函数体内部的return语句来返回
	miltiple = a * b
	statue = true
	return miltiple, statue //此外，return后面的内容可以不写。因为在返回值定义时已经关联了miltiple和statue
}

func main() {
	var a int
	var b int
	a = 10
	b = 20
	sum := add(a, b)
	fmt.Println(sum)
	difference, flag := sub(a, b)
	printHello()
	fmt.Println(difference, flag)
	product, statue := sub(a, b)
	fmt.Println(product, statue)

	/*
		匿名函数:通常不能在函数中再定义函数，但可以定义一个变量来存储一个函数。
		        特点：没有函数名，直接定义为一个变量的值
	*/
	// 把一个无参、返回 string 的匿名函数，赋值给变量 getName
	var getName = func() string {
		return "Aurora"
	}
	println(getName()) // 调用这个匿名函数

	var setName = func(name string) {
		fmt.Println(name)
	}
	setName("huahua")
	//fmt.Println(setName) //setName本身没有返回值，这样写是错误的，只会输出该函数的地址

}

// 输出功能。函数放到main即可，不论上下。
func printHello() { //函数不带返回值的函数
	println("Hello World!")
}
