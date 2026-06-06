// 创建时间：2026/6/4 下午6:35
package main

import (
	"fmt"
	"reflect"
)

type Student struct {
	Name string `json:"name"`
	Age  int `json:"age"`
	IsMan bool `json:"is_man"`
}

func ParseJson(obj any){   //读取结构体中的值，并转换为json字符串的函数
    t:=reflect.TypeOf(obj)  //获取类型,在这里是Student结构体类型
	v:=reflect.ValueOf(obj) //获取值
	//num:=v.NumField() //必须是获取类型或值之后，才能调用 NumField()获取字段数量,结果是3.
    for i:=0;i<v.NumField();i++{  //使用一个for循环，来多次拿字段

        fmt.Println(v.Field(i))  //获取索引为i时的字段的值。字段的值就是自己定义的Tom、20、true等
        fmt.Println()

		typename:=t.Field(i)  //获取索引为i时的字段的名称、标签等。
		fmt.Println(typename.Name)  //获取字段的名称，结果是Name、Age、IsMan
		fmt.Println(typename.Tag)  //获取字段的标签，结果是json: "name"、json: "age"、json: "is_man"
		fmt.Println(typename.Tag.Get("json"))  //获取字段的标签中的json字符串，结果是name、age、is_man

		/*
		有关v.Field(i)和t.Field(i)的区别：
		v.Field(i)是获取字段的值
		t.Field(i)是获取字段的字段名、字段类型、标签（tag）、偏移量
		 */
	}
}


func main() {
   s:=Student{Name: "Tom", Age:  20, IsMan: true,}
   ParseJson(s)
}


/*


三、反射三大核心（必须背）
1. reflect.TypeOf(x)
   → 获取变量的【类型】（是什么类型：int/string/struct...）

2. reflect.ValueOf(x)
   → 获取变量的【值】（变量存储的数据）

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
