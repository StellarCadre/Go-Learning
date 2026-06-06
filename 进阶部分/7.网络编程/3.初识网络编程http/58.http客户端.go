// 创建时间：2026/6/6 下午9:49
package main

import (
	"fmt"
	"io"
	"net/http"
)

/*
这个客户端代码，本质上就是一个“没有图形界面的极简浏览器,和用户将链接复制到浏览器并打开差不多”。
 */



func main() {
     response,err:=http.Get("http://127.0.0.1:8080")  //瞬间通过网络向 8080 端口发射了一个 HTTP 请求。你在浏览器的地址栏里输入网址，并按下回车键。它主动去寻找服务端，并带回了服务端的全部回复（存放在 response 变量里）。
	 if err!=nil{
		 fmt.Println("请求失败")
		 return
	 }
	 byteData,_:=io.ReadAll(response.Body)  //客户端收到响应，代码走到 io.ReadAll，把 "Hello World" 拆解出来，最后通过 fmt.Println 打印在你的黑框框（控制台）里。
	 fmt.Println(string(byteData))
}
