package main

import "fmt"

func main() {

	//中断式，卫语句（最推荐的写法）
	fmt.Print("请输入一个成绩：")
	var score int
	fmt.Scan(&score)
	if score >= 90 {
		println("优秀")
		return
	} else if score >= 80 {
		println("良好")
		return
	} else if score >= 60 {
		println("及格")
		return
	} else {
		println("不及格")
		return
	}

	//嵌套式，if语句可以嵌套在if语句中，但是一般不推荐嵌套，容易造成代码阅读性差
	//多条件式,&&是与，||是或，!是非。短路运算：&&左边为false，右边不执行，||左边为true，右边不执行
	fmt.Print("请输入一个年龄：")
	var age int
	fmt.Scan(&age)
	if age > 0 && age < 18 {
		println("未成年")
	} else if age >= 18 && age < 30 {
		println("青年")
	} else if age >= 30 && age < 60 {
		println("中年")
	} else {
		println("老年")
	}
}
