// 创建时间：2026/5/27 下午9:43
package main

import "fmt"

type Class struct{
	Name string
}


type Student struct{  //定义一个结构体
    name string
	age int
	score int
	Class Class   //嵌套Class结构体（匿名嵌套：Student拥有Class的所有属性）  方法是直接把另一个结构体的名字拿过来，   组合
}
//给结构体绑定方法，该方法只属于Student结构体
func (s Student)Study()  { //(s Student)是方法接收者，表示Study这个函数是属于Student类型的。s是接收者变量
	fmt.Println(s.name,"正在学习语文",s.age,"岁",s.score,"分")
}
func (s Student)getclass()  {
	fmt.Println("该学生所在班级是",s.Class.Name)  //或s.Name
}


func main() {
	c1:=Class{   //创建班级的结构体实例
		Name: "1班",
	}

    s1:=Student{  //创建一个学生的结构体实例
		name: "Aurora",   //字段名：普通值
		age: 20,
		score: 100,
		Class: c1,   //结构体字段名：结构体实例
		/*
		c1 是创建好的真实班级（比如 1 班）
		Class 是结构体里的班级位置
		把 c1 这个班级，放进学生的班级位置里！
		 */
	}

	s1.Study()
	s1.getclass()


}
