package main

import "fmt"

// 键值对 ，字典不能使用索引，只能使用键来访问值
func main() {
	var userMap = map[int]string{ //map[key类型]value类型
		1001: "张三",
		1002: "李四",
		1003: "王五",
	}
	fmt.Println(userMap)
	fmt.Println(userMap[1002])
	fmt.Println(userMap[1004])   //不存在的key，返回value类型的零值,string的空字符串
	value1 := userMap[1004]      //使用value1变量来接收1004键对应的值
	value1, ok1 := userMap[1004] //使用value1变量来接收1004键对应的值，ok1变量来接收是否存在的布尔值
	fmt.Println(value1, ok1)     //结果：空字符串 和 false
	value2 := userMap[1003]      //使用value2变量来接收1003键对应的值
	value2, ok2 := userMap[1003] //使用value2变量来接收1003键对应的值，ok2变量来接收是否存在的布尔值
	fmt.Println(value2, ok2)     //结果：王五 和 true

	//对应位置的值进行修改
	userMap[1001] = "张三丰"
	fmt.Println(userMap)

	//删除键值对
	delete(userMap, 1002)
	fmt.Println(userMap)
}
