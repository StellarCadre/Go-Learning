// 创建时间：2026/6/4 下午5:32
package main

import (
	"fmt"
	"reflect"
)

func GetType(obj any) { //这里obj使用any，any是空接口，表示可以接受任何类型。等同于i interface{}
	t := reflect.TypeOf(obj) //获取obj的类型
	switch t.Kind() {        //获取底层类型
	case reflect.Struct: //若t.Kind()是结构体
		fmt.Println("这是一个结构体")
	case reflect.Int:
		fmt.Println("这是一个int") //若t.Kind()是int
	case reflect.String:
		fmt.Println("这是一个string")
	}
}

func GetValue(obj any) {
	v := reflect.ValueOf(obj) //获取到值,然后根据值的类型进行判断并输出
	switch v.Kind() {
	case reflect.Int:
		fmt.Println("这是一个int", v.Int()) //或其是int型，使用Int()方法获取值
	case reflect.String:
		fmt.Println("这是一个string", v.String()) //或其是string型，使用String()方法获取值
	case reflect.Struct:
		fmt.Println("这是一个结构体,这里的结构体方法放后面讲")
	}

}

func main() {
	type Person struct {
		Name string
		Age  int
	}
	p := Person{Name: "Tom", Age: 20}

	GetType(1)
	GetType("hello")
	GetType(p)

	GetValue(1)
	GetValue("hello")
	GetValue(p)

}

/*
=============================================
Go 反射（reflect）核心知识点（初学者必看）
=============================================

一、反射是什么？
1. 反射是 Go 语言提供的一种机制，允许程序在【运行时】
   动态获取变量的【类型信息】和【值信息】。
2. 可以动态操作变量：修改变量值、获取结构体字段、调用方法等。
3. 不需要提前知道变量类型，程序运行时才“看”变量是什么。

一句话总结：
反射 = 程序运行时，动态查看/修改变量的类型与值。

====================================================

二、为什么要用反射？（使用场景）
1. 编写【通用工具函数】（不用管传入的变量是什么类型）
2. JSON 序列化/反序列化（底层全是反射）
3. GORM 数据库映射（全靠反射）
4. Gin、Viper 等框架参数绑定、配置映射
5. 动态处理结构体、切片、map 等复杂类型

一句话总结：
反射 = 写通用代码、框架代码的核心技术。

====================================================

三、反射三大核心（必须背）
1. reflect.TypeOf(x)
   → 获取变量x的【类型】（是什么类型：int/string/struct...）

2. reflect.ValueOf(x)
   → 获取变量x的【值】（变量存储的数据）

3. 反射对象.Interface()
   → 将反射值转回普通变量

====================================================

四、反射最常用方法（必会）
1. 获取类型
   typ := reflect.TypeOf(x)

2. 获取值
   val := reflect.ValueOf(x)

3. 获取底层类型（int/struct/ptr/slice）
   typ.Kind()

4. 获取结构体字段数量
   typ.NumField()

5. 获取结构体第 i 个字段
   typ.Field(i)

6. 获取值（根据类型调用）
   val.Int()
   val.String()
   val.Bool()

7. 修改变量值（必须传指针 + Elem()）
   val.Elem().SetInt(100)
   val.Elem().SetString("hello")

8. 调用结构体方法
   val.Method(0).Call(nil)

====================================================

五、反射三大铁律（90% 错误都来自这里）
1. 要修改变量，必须传入【指针】
2. 修改值必须使用 Elem() 获取指针指向的对象
3. 结构体只有【大写导出字段】才能被反射访问

====================================================
*/
