// 创建时间：2026/6/3 下午5:05
package main

import (
	"encoding/json"
	"fmt"
)

//type Response struct {
//	Code int `json:"code"`
//	Msg string `json:"msg"`
//	Data any `json:"data"`
//	/*
//	Data用了 any 类型，虽然能接收任意数据，但丢失了类型安全：
//	    编译期无法检查赋值给 Data 的类型是否符合预期
//	    后续读取 Data 时需要手动类型断言（data.(Message)），容易出错
//	 */
//}

//使用了泛型后的Response结构体
type Response[T any] struct {
	Code int `json:"code"`
	Msg string `json:"msg"`
	Data T `json:"data"`
}

func main() {

	type Message struct {
		Name string `json:"name"`

	}
	type Information struct {
		Name string `json:"name"`
		Age int `json:"age"`
	}

    //res1:=Response{
	//	Code: 1,
	//	Msg: "success",
	//	Data: Message{
	//		Name: "Tom",
	//	},
	//}
	//res2:=Response{
	//	Code: 1,
	//	Msg: "succ",
	//	Data: Information{
	//		Name: "jack",
	//		Age: 18,
	//	},
	//}
	//
	//byteData,_:=json.Marshal(res1)  //通过 json.Marshal 序列化 Response的内容 为 JSON 字符串并打印
	//fmt.Println(string(byteData))
	//byteData,_=json.Marshal(res2)
	//fmt.Println(string(byteData))

	/* 反序列化操作
	   但是在整个过程中，我们没有指定 Data 的具体类型，是any。
	   所以反序列化时，会根据 JSON 字符串中的内容来推导 Data 的类型，反序列化后却变成了map[name:Tom]
       应当修改Response结构体，用到泛型。
	 */
	//var userResponse Response
	//json.Unmarshal([]byte(`{"code":1,"msg":"success","data":{"name":"Tom"}}`),&userResponse)  //将JSON字符串反序列化，并将结果存储在userResponse中
	//fmt.Println(userResponse)

	//泛型结构体后，执行反序列化操作
	var userResponse Response[Message]
	json.Unmarshal([]byte(`{"code":1,"msg":"success","data":{"name":"Tom"}}`),&userResponse)  //由于结构体使用了泛型，同时上一条语句中也指定了Data的类型为Message，所以这里可以正常反序列化
	fmt.Println(userResponse)
    fmt.Println(userResponse.Data)
	fmt.Println(userResponse.Data.Name)
}







/*
【Go 泛型 学习前瞻】

1. 泛型是什么？
   一句话：让函数/结构体支持“多种类型”，但不用写重复代码。
   以前要给 int、string 各写一套，现在一套通用。

2. 为什么要用泛型？
   • 消除重复代码（比如写一个通用切片工具函数）
   • 保留类型安全（不会像 interface{} 一样丢失类型）
   • 代码更简洁、更通用

3. 泛型用在哪里？
   • 通用工具函数（切片、map、队列、栈、链表）
   • 通用数据结构（通用切片、通用树、通用链表）
   • 不想用 interface{} 丢失类型检查的场景

4. 泛型核心三要素
   • 类型参数 [T any] 或 [T TypeConstraint]
   • 类型约束 any / comparable / 自定义
   • 类型推导（不用手动指定类型，自动推断）

5. 最简单理解
   泛型 = “模板”
   T = 占位符（代表未来传入的任意类型）

6. 新手最容易懵的符号
   [T any] —— 这就是泛型！表示 T 可以是任何类型

7. 工程意义
   让 Go 拥有真正的通用工具库，
   但不丢失类型安全，不降低性能。

一句话总结：
泛型 = 一套代码，支持所有类型，且类型安全！
*/