// 创建时间：2026/7/28 下午8:10
package main

import (
	"fmt"
	"net/http"
	"os"
)

//web开发的本质，后端服务器与浏览器、app等的交互，请求与响应

// 往浏览器中传输文字或文本文件
func Say(w http.ResponseWriter, r *http.Request) { //w表示要返回的内容，r表示传入的请求
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	//w.Write([]byte("<h1>Hello World</h1>"))
	/*
		这里要写在前端页面中应该展示的内容，可以加html标签等，花里胡哨的样式等
		w.Write 仅负责向浏览器传输字节数据，浏览器依靠服务端设定的 Content-Type 响应头区分数据类型并对应解析渲染。
		但这样写代码看起来太乱了，可以使用导入文件的方式来展示多个内容
	*/
	file, err := os.ReadFile("C:\\Users\\Aurora\\Desktop\\Go Project\\Gin Project\\Basic\\test.txt")
	if err != nil {
		fmt.Println("文件读取失败")
	}
	w.Write(file)
}

// 往浏览器中传输图片
func img(w http.ResponseWriter, r *http.Request) {
	data, err := os.ReadFile("C:\\Users\\Aurora\\Desktop\\Go Project\\Gin Project\\Basic\\th.jpg")
	if err != nil {
		w.Write([]byte("图片不存在"))
		return
	}
	w.Header().Set("Content-Type", "image/jpeg")
	w.Write(data)
}

//除了文本、文件、图片，还可以传输其他格式的文件，比如音频、视频等，只要浏览器能解析的格式都可以传输

func main() {
	http.HandleFunc("/hello", Say) // 注册路由，当请求路径为****/hello时，执行SayHello函数
	http.HandleFunc("/img", img)
	err := http.ListenAndServe(":8080", nil) // 监听端口8080。和上面一起拼成http://127.0.0.1:8080/hello
	if err != nil {
		fmt.Println("http server start failed")
	}
}

/*
访问方式：
http://127.0.0.1:8080/hello
http://127.0.0.1:8080/img
浏览器访问不存在的地址时，会返回404报错代码
*/
