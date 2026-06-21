// 创建时间：2026/6/4 下午9:30
package main

import (
	"fmt"
	"reflect"
)

type Teacher struct {
	Name  string `json:"name"`
	Age   int    `json:"age"`
	IsMan bool   `json:"is_man"`
}

// 给结构体绑定方法，该方法只属于Teacher结构体
func (t Teacher) Call() { // 方法名 Call 是大写，属于导出方法
	fmt.Println("你觉得我会被调用吗？1")
}

func (t Teacher) Add() {
	fmt.Println("你觉得我会被调用吗？2")
}

func (t Teacher) add() { //由于方法名是小写，属于非导出方法，无法被反射获取
	fmt.Println("你觉得我会被调用吗？3")
}

//func (t Teacher)Print(str string)  {  //若该方法也要使用，下面的代码需要升级:需要在反射调用方法前，先判断方法的参数数量，再根据参数数量传递对应参数。
//	fmt.Println("打印输出为：",str)
//}

func CallMethord(obj any) {
	v := reflect.ValueOf(obj).Elem()     //解引用后获取字段的值
	t := reflect.TypeOf(obj).Elem()      //解引用后获取字段的字段名、字段类型、标签（tag）、偏移量
	for i := 0; i < v.NumMethod(); i++ { //遍历结构体字段  v.NumMethod()表示获取v的方法数量
		m := t.Method(i) //这样获取的方法，只看、只读方法的名字、信息（不能调用！）因为是使用t得到类型
		fmt.Println(m.Name)

		//调用结构体中的方法（无参传nil/空切片，有参需包装为[]reflect.Value）
		method := v.Method(i) //这样获取的方法，能执行、能调用真正的函数,因为是使用v得到值
		method.Call([]reflect.Value{})
		//method.Call([]reflect.Value{reflect.ValueOf("你好")})
	}
}

func main() {
	teacher := Teacher{Name: "MaWeika", Age: 500, IsMan: false}
	CallMethord(&teacher)
}

/*
反射调用结构体方法 核心知识点汇总

1. 核心API：reflect.Value.Call()
   专门用于通过反射执行 函数 / 结构体方法
   调用格式：methodValue.Call(参数切片)
   返回值：[]reflect.Value（接收方法返回的所有结果）

2. 获取结构体方法
   - v.MethodByName("方法名")：根据字符串获取方法（最常用）
   - v.Method(索引)：按顺序获取方法
   注意：只有大写导出的方法才能被反射获取

3. 方法调用规则
   - 结构体方法如果是 值接收者：传值、传指针都能调用
   - 结构体方法如果是 指针接收者：必须传指针才能调用
   - 调用前必须判断方法是否存在：IsValid()

4. 参数传递格式
   参数必须包装为：[]reflect.Value
   无参：nil 或 []reflect.Value{}
   有参：[]reflect.Value{reflect.ValueOf(参数1), reflect.ValueOf(参数2)}

5. 返回值接收
   方法返回值会以 []reflect.Value 返回
   使用 .Interface() 转回普通类型使用

6. 重要判断
   - IsValid()：方法是否存在
   - IsNil()：是否为nil方法
   可避免反射调用panic
*/
