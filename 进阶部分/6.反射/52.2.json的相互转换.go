// 创建时间：2026/6/21 下午7:57
package main

import (
	"encoding/json"
	"fmt"
)

type Movie struct {
	Title  string   `json:"title"`
	Year   int      `json:"year"`
	Price  int      `json:"price"`
	Actors []string `json:"actors"`
}

func main() {
	movie1 := Movie{Title: "喜剧之王", Year: 2000, Price: 10, Actors: []string{"周星驰", "梁朝伟"}}
	//编码：将结构体转换为json字符串
	jsonStr, err := json.Marshal(movie1) //json.Marshal在使用时会拿到字段对应的json标签
	if err != nil {
		fmt.Println("json序列化失败:", err)
		return
	}
	fmt.Println(jsonStr)         //json.Marshal返回的 [] byte 类型（字节切片） 在终端直接打印时，显示的是每个字符对应的 ASCII/UTF-8 编码数值，而非字符串本身。
	fmt.Println(string(jsonStr)) //将 [] byte 类型转换为 string 类型，然后打印。

	//解码：将json字符串转换为结构体
	movie2 := Movie{}
	json.Unmarshal(jsonStr, &movie2) //将jsonStr转换为结构体并赋值给movie2
	fmt.Println(movie2)

}
