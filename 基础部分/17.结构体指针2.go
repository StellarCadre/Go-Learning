// 创建时间：2026/5/28 下午8:13
package main

import "fmt"

type UserInfo struct{
	name string
	age int
}

func (u *UserInfo)Setname(str string)  {
    u.name=str
}
func (u *UserInfo)Setage(age int)  {
	u.age=age
}

func main() {
    u1:=UserInfo{
		name: "Aurora",
		age: 18,
	}
	fmt.Println("他的名字是",u1.name,"他的年龄是",u1.age)
	u1.Setname("linhuahua")
	u1.Setage(21)
	fmt.Println("他的名字是",u1.name,"他的年龄是",u1.age)

}

