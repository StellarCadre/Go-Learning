// 创建时间：2026/5/29 下午6:26
package main

import "fmt"

/*
接口是一组仅包含方法名、参数、返回值的方法的集合，未实现具体方法
可以解决相似功能（仅传入对象不同 ）的方法冗余的问题，实现多态。如下：
*/
type Cat struct {
	name string
}
type Dog struct {
	name string
}

func (c Cat) Sing() {
	fmt.Println(c.name, "在唱歌")
}
func (d Dog) Sing() {
	fmt.Println(d.name, "在唱歌")
}

func sing(c Cat) {
	c.Sing()
}

func main() {
	c1 := Cat{name: "小花"}
	sing(c1)
	d1 := Dog{name: "小狗"}
	sing(d1)
	//因为之前定义的sing函数，需要的参数是Cat类型，所以这里传入Dog会报错.那么，为了解决该问题，一般来说可以再定义一个sing函数，传入Dog类型，但是这样就造成了方法冗余，所以接口就应运而生。
}
