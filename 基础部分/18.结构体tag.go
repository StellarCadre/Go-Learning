// 创建时间：2026/5/28 下午9:46
package main
//将结果转换为json格式，并传输到前端
import (
	"encoding/json"
	"fmt"
)

type Worker struct {
	Name string `json:"name"`  //后面的json中的内容，表示Name在json中的表示为name
	Age  int    `json:"age"`   //若写成—,则表示不进行json的转换,忽略
}

func main() {
    u1:=Worker{
		Name: "Aurora",
		Age: 18,
	}
	byteData,_:=json.Marshal(u1)
	fmt.Println(string(byteData))

}
