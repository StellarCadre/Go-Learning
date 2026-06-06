package main

import "fmt"

func main() {
	//switch写法一
	var age int
	println("请输入一个年龄：")
	fmt.Scan(&age)
	switch {
	case age < 0: //case后面可以跟表达式，表达式为true，则执行case分支
		println("输入错误")
		fallthrough //默认情况下，switch语言只会执行第一个满足的case分支。若要继续执行后面的分支，需要使用fallthrough关键字
	case age < 18:
		println("未成年")
		fallthrough
	case age < 24:
		println("青年")
		fallthrough
	case age < 60:
		println("中年")
		fallthrough
	default:
		println("老年")
	}

	//switch写法二
	var week int
	println("请输入星期(数字)：")
	fmt.Scan(&week)
	switch week {
	case 1: //直接判断输入是什么，不需要表达式，一致就行case分支
		println("星期一")
	case 2:
		println("星期二")
	case 3:
		println("星期三")
	case 4:
		println("星期四")
	case 5:
		println("星期五")
	case 6:
		println("星期六")
	case 7:
		println("星期日")
	default:
		println("输入错误")

	}

}
