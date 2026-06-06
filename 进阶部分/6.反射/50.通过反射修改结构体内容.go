// 创建时间：2026/6/4 下午7:39

package main

import (
	"fmt"
	"reflect"
	"strings"
)

type Person struct {
	Name1 string `big:"-"`
	Name2 string
}

func Modify(obj any){   //修改结构体中的值，并转换为json字符串的函数
	t := reflect.TypeOf(obj)  //获取类型,在这里是Person结构体类型
	v := reflect.ValueOf(obj) //获取值
	if t.Kind() == reflect.Ptr {  // 如果是指针，解引用
		t = t.Elem()
		v = v.Elem()
	}
	if t.Kind() != reflect.Struct {
		return
	}
	// 判断是否为结构体
	if t.Kind() != reflect.Struct {
		return
	}
	for i:=0;i<v.NumField();i++{  //使用一个for循环，来多次拿字段
		//获取big标签
		big := t.Field(i).Tag.Get("big")
		if big == "-" {
			continue
		}
		values := v.Field(i)  //获取索引为i时的字段的值。字段的值就是自己定义的Jack,Jerry等
		values.SetString(strings.ToUpper(values.String()))  // 将字段值改为大写

		// 打印时字段名变大写。注意，无法修改结构体中，字段的名称，因为结构体的字段名称是只读的。只能在打印时，临时将字段名称转换为大写。
		fieldName := t.Field(i).Name  //获取索引为i时的字段的名称
		upperName := strings.ToUpper(fieldName)
		fmt.Println(upperName)
	}
}


func main() {
	s:=Person{Name1: "Jack", Name2: "Jerry",}
	fmt.Println("修改前：", s) // 打印修改前的结构体
	Modify(&s)
	fmt.Println("修改后：", s) // 打印修改后的结构体
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

