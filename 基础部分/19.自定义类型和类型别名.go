// 创建时间：2026/5/28 下午9:59
package main

import "fmt"
/*
自定义数据类型，指使用type关键字定义的新类型，分为基本类型的别名、结构体和函数等组合而成的自定义类型。能够帮助我们更好地抽象和封装数据。
结构体就是自定义类型

自定义类型：可以绑定方法，打印出的类型是自己。类型别名不用转换。
类型别名：不能绑定方法，打印出的类型是原始类型
 */

type Code1 int  //自定义类型
type Code2 = int  //类型别名


func(c Code1)Getcode1(){
	fmt.Println("自定义类型：可以绑定方法")
}
//func(c Code2)Getcode2(){
//	fmt.Println("类型别名：不能绑定方法，所以错误")
//}

func main() {
    var c1 Code1 =1
	fmt.Printf("%v,%T\n",c1,c1)  //类型是自己Code1
	var c2 Code2 =1
	fmt.Printf("%v,%T\n",c2,c2)  //类型是原始类型int
}
