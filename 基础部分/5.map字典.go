package main

import "fmt"

// 键值对 ，字典不能使用索引，只能使用键来访问值

func changeMap(cityMap map[int]string) {
	cityMap[1] = "张家口"
}

func main() {
	//map的定义和使用
	//var userMap map[int]string //法一：定义map，但是没有初始化，值为nil，不能直接使用，需要先进行初始化
	//userMap := make(map[int]string,10) //法二：make函数来初始化map,10是预分配的内存空间
	//userMap := make(map[int]string) //法三：make函数来初始化map
	//userMap[1001] = "张三"
	//userMap[1002] = "李四"
	//userMap[1003] = "王五"

	var userMap = map[int]string{ //法四：map[key类型]value类型
		1001: "张三",
		1002: "李四",
		1003: "王五",
	}
	fmt.Println(userMap)
	fmt.Println(userMap[1002])
	fmt.Println(userMap[1004])   //不存在的key，返回value类型的零值,如string的空字符串
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

	//遍历map,使用range关键字
	cityMap := map[int]string{
		1: "北京",
		2: "上海",
		3: "广州",
		4: "深圳",
	}
	for key, value := range cityMap {
		fmt.Println(key, value)
	}

	//map的传递：map是引用传递，所以在函数中对map的修改会影响到原map的值。
	changeMap(cityMap)
	fmt.Println(cityMap)

}
