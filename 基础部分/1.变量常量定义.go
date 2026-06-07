package main

import "fmt"

/*
重点：
go语言强制要求所有定义的变量必须被使用，否则是编译错误，而不是警告，只对局部变量适用。
*/

var time1 = "2026.5.15" //全局变量声明
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
	var id1 int //声明但未初始化，默认值为0
	id1 = 2
	var name1 string = "liu" //声明并初始化一个值
	var name2 = "xiao"       //省略类型，go会自动推导类型
	var a, b = 4, 5

	fmt.Println(id1) //PrintLn是打印函数，会自动换行
	fmt.Println(name1)
	fmt.Println(name2)
	fmt.Println(a + b)

	fmt.Printf("type of id1: %T\n", id1)     //使用%T能够打印变量类型
	fmt.Printf("type of name1: %T\n", name1) //Printf是格式化输出函数，可以指定输出格式，这里是%T，表示输出变量类型
	fmt.Printf("type of name2: %T\n", name2)
	fmt.Printf("type of a: %T\n", a)

	//变量定义方法二：短变量声明法， 变量名:= 植  ,能够自动识别类似直接赋值，只能在函数内使用，但是go最常用最推荐的写法
	id2 := 1
	id3 := 3
	fmt.Println(id2 + id3)

	//常量的定义方法： const 常量名 类型 =植   被定义后就不能被修改了，可以不使用
	const pai float32 = 3.1415

	//const还可用于定义枚举
	const (
		Sunday = iota //Sundat的值为0，iota是go语言中，const的常量计数器
		Monday        //Monday的值为1，且iota的值会自动加1，为1
		Tuesday
		Wednesday
		Thursday
		Friday
		Saturday
	)

	//命名规范： 一个包中，只有首字母大写的变量、函数、方法、属性才能被其他包访问

}
