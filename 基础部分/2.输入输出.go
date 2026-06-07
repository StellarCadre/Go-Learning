package main

import "fmt"

func main() {
	//是否换行
	fmt.Printf("我是不换行")       //输出后不换行
	fmt.Println("我是换行", "没坐") //输出后换行

	//格式化输出,n是换行标志
	fmt.Printf("%s\n", "前面s是字符串标志")  //字符串标志，不带双引号
	fmt.Printf("%q\n", "abc")        //字符串标志,带双引号
	fmt.Printf("%d\n", 520)          //整数标志
	fmt.Printf("%f\n", 3.1415926)    //小数标志，保留所有小数
	fmt.Printf("%0.3f\n", 3.1415926) //小数标志，保留3位小数
	fmt.Printf("%t\n", true)         //布尔标志
	fmt.Printf("%c\n", 'a')          //字符标志，不带单引号

	fmt.Println("______________________________________________")
	fmt.Printf("%T %T\n", "string", 1.3) //类型标志,能够打印出变量类型
	s := 123
	fmt.Printf("%T\n", s)

	fmt.Printf("%v\n", 13)     //%v (Value) 按照变量的“默认格式”输出其值13
	fmt.Printf("%v\n", "")     //空字符串在控制台不显示
	fmt.Printf("%#v\n", "")    //# 标志会强制打印出符合 Go 语言语法的字面量表示。
	fmt.Printf("%+v\n", "fsf") //显示字段名。%+v 主要是为结构体 (Struct) 设计的
	//如果你定义了一个结构体 user := User{Name: "Alice", Age: 20}，使用 %v 只会输出 {Alice 20}，而使用 %+v 则会输出 {Name:Alice Age:20}，这在分析复杂的数据封装时极其有用。

	fmt.Printf("%U\n", 'a') //Unicode编码标志
	fmt.Printf("%b\n", 123) //二进制标志，打印出123对应的二进制数
	fmt.Printf("%o\n", 123) //八进制标志，打印出123对应的八进制数
	fmt.Printf("%x\n", 123) //十六进制标志，打印出123对应的十六进制数
	fmt.Printf("%X\n", 123) //十六进制标志，打印出123对应的十六进制数，大写

	var a = 123
	fmt.Printf("%p\n", &a) //指针标志，&a是取a的地址

	//将格式化输出赋值给变量
	a1 := fmt.Sprintf("%s", "明天")
	a2 := fmt.Sprintf("%s", "放假")
	fmt.Println(a1 + a2)

	//输入
	fmt.Println("请输入一个学号：")
	var a3 int                  //该变量必须在fmt.Scan之前声明，否则会报错，用于接收控制台的输入
	fmt.Scan(&a3)               //scan拿到刚定义的变量的地址，然后把控制台的输入赋值给该变量
	fmt.Println("你输入的整数是：", a3) //并输出

	fmt.Println("请输入学生姓名:")
	var name string
	fmt.Scan(&name)
	fmt.Println("你输入的学生姓名是：", name)

	//输入并加错误判断
	fmt.Println("请输入学生年龄:")
	var age int
	n, err := fmt.Scan(&age) //n是输入的数量，err是错误信息
	fmt.Println(n, err, age) //最后的输出结果分为：1：正确版：实际的数量 <nil> 实际的值 2：错误版：0 expected integer 错误时的默认值0
}
