package main //声明该文件属于main包，使得go编译器可知这是一个可执行程序，而不是供别人调用的库

import "fmt" //导入go的标准输入输出包，全称是format，提供如格式化输出、读取输入等功能，比如printf

func main() { //定义一个函数，main函数不可有参数   花括号在函数名后，不能换行
	fmt.Printf("Hello World")
}
