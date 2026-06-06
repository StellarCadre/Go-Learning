// 创建时间：2026/5/28 下午4:32
package main
//在结构体中使用指针，可以修改结构体中的值
import "fmt"

type Class1 struct{
	Name string
}


type Student1 struct { //定义一个结构体
	name   string
	age    int
	score  int
	Class1  Class1 //嵌套Class结构体（匿名嵌套：Student拥有Class的所有属性）  方法是直接把另一个结构体的名字拿过来，

}

//给结构体绑定方法，该方法只属于Student结构体
func (s Student1)Study1()  { //(s Student)是方法接收者，表示Study这个函数是属于Student类型的。s是接收者变量
	fmt.Println(s.name,"正在学习语文",s.age,"岁",s.score,"分")
}
func (s Student1)Getclass1()  {
	fmt.Println("该学生所在班级是",s.Class1.Name)  //或s.Name
}


func (s Student1) SetnameValue(str string){ //这样属于值传递，只会修改s1的副本，不会修改s1本身，所以无法实现名字修改。
	s.name=str
}
/*
(s Student1)是真正要修改的对象本身
(str string)是方法的参数
在Go，前一个括号中放的是接收者变量，后一个括号中放的是参数。故要想修改结构体中的值，应当传递结构体（前者）的地址
举例：
1：func (s *Student1) change(str *string)  等价于c中的：void change(Student &s, string &str)
2：func (s *Student1) change(str string)   等价于c中的：void change(Student &s, string str)
3：func (s Student1) change(str *string)   等价于c中的：void change(Student s, string &str)
要想实现修改，必须选1
第一个括号：决定 能不能改结构体
第二个括号：决定 能不能改传进来的字符串
 */
func (s Student1)SetnamePoint1(str *string){  //错误用法
	s.name=*str
}
func (s *Student1)GetnamePoint2(str string){  //正确用法
	s.name=str
}


func main() {
	c1:=Class1{   //创建班级的结构体实例
		Name: "6班",
	}

	s1:=Student1{  //创建一个学生的结构体实例
		name: "zhangyue",   //字段名：普通值
		age: 18,
		score: 87,
		Class1: c1,   //结构体字段名：结构体实例

	}

	s1.Study1()
	s1.Getclass1()
	var str string
	fmt.Scan(&str)
	s1.SetnameValue(str)
	s1.Study1()
    s1.SetnamePoint1(&str)
	s1.Study1()
	s1.GetnamePoint2(str)
	s1.Study1()

}