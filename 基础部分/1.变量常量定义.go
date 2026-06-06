package main

import "fmt"

/*
重点：
go语言强制要求所有定义的变量必须被使用，否则是编译错误，而不是警告，只对局部变量适用。
*/

var time1 = "2026.5.15" //此为全局变量
var (
	s1 string = "2026"
	s2 string = "5.15"
)

func main() {

	/* go定义变量和其他语言不同,而且没有符号
	int a;
	int b;
	*/

	//变量定义方法一：var 变量名 类型   以及变体
	var id1 int
	id1 = 2
	var name1 string = "liu"
	var name2 = "xiao"
	var a, b = 4, 5

	fmt.Println(id1)
	fmt.Println(name1)
	fmt.Println(name2)
	fmt.Println(a + b)

	//变量定义方法二：短变量声明法， 变量名:= 植  ,能够自动识别类似直接赋值，只能在函数内使用，但是go最常用最推荐的写法
	id2 := 1
	id3 := 3
	fmt.Println(id2 + id3)

	//常量的定义方法： const 常量名 类型 =植   被定义后就不能被修改了
	const pai float32 = 3.1415

	//命名规范： 一个包中，只有首字母大写的变量、函数、方法、属性才能被其他包访问

}
