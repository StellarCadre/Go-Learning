// 创建时间：2026/5/27 下午9:43
package main

import (
	"fmt"
)

type Class struct {
	Name string
}

type Student struct { //定义一个结构体
	name  string
	age   int
	score int
	Class Class //嵌套Class结构体（匿名嵌套，也叫做继承，Student拥有Class的所有属性）  方法是直接把另一个结构体的名字拿过来，   组合
}

// 结构体方法使用结构体：
// 给结构体绑定方法，该方法只属于Student结构体
func (s Student) Study() { //(s Student)是方法接收者，表示Study这个函数是属于Student类型的。s是接收者变量。结构体方法之读取内容
	fmt.Println(s.name, "正在学习语文", s.age, "岁", s.score, "分")
}
func (s Student) getclass() {
	fmt.Println("该学生所在班级是", s.Class.Name) //或s.Name
}

func (s Student) setage1(age int) { //结构体方法之修改内容，这里看似可以修改结构体的age，但实际上是修改了结构体的副本s，不会改变原始值。
	s.age = s.age - 1
} //使用如下形式才可以修改：
func (s *Student) setage2(age int) {
	s.age = s.age - 1
}

// 继承
type GoodSStudent struct {
	Student     // 继承Student结构体的所有属性和方法
	id      int //新增score字段
}

func (g GoodSStudent) Study() { //重写了Student结构体的Study方法
	fmt.Println("正在学习语文的是", g.name)
}
func (g GoodSStudent) getage() {
	fmt.Println("他的年龄是", g.age)
} //子类自己新增的方法

// 非结构体方法使用结构体：值传递和指针传递
func change1(stu Student) { //以值传递的形式接收结构体仍然是值传递，不会改变原始值
	stu.name = "linhuahua"
	stu.age = 22
	stu.score = 95
	stu.Class.Name = "2班"
}

func change2(stu *Student) { //以指针传递的形式接收结构体，会改变原始值
	stu.name = "linhuahua"
	stu.age = 25
	stu.score = 95
	stu.Class.Name = "2班"
}

func main() {
	c1 := Class{ //创建班级的结构体实例
		Name: "1班",
	}

	s1 := Student{ //创建一个学生的结构体实例
		name:  "Aurora", //字段名：普通值
		age:   20,
		score: 100,
		Class: c1, //结构体字段名：结构体实例
		/*
			c1 是创建好的真实班级（比如 1 班）
			Class 是结构体里的班级位置
			把 c1 这个班级，放进学生的班级位置里！
		*/
	}
	fmt.Println("正常打印测试")
	s1.Study()
	s1.getclass()

	fmt.Println("结构体方法修改结构体之打印测试1")
	s1.setage1(11)
	s1.Study()
	fmt.Println("结构体方法修改结构体之打印测试2")
	s1.setage2(30)
	s1.Study()
	fmt.Println("非结构体方法修改结构体之打印测试1")
	change1(s1)
	s1.Study()
	fmt.Println("非结构体方法修改结构体之打印测试2")
	change2(&s1)
	s1.Study()

	fmt.Println("父类子类测试")
	g1 := GoodSStudent{ //创建一个好学生（学生的子类）的结构体实例，同时他能拿到Student的所有属性和方法
		Student: Student{name: "liuna", age: 44, score: 144, Class: c1},
		id:      0001,
	}
	g1.getage()   //调用自己的方法
	g1.Study()    //调用重写父类后的方法
	g1.getclass() //调用继承至父类的方法

}
