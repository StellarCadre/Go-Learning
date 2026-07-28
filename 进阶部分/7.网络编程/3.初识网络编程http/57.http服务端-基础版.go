// 创建时间：2026/6/6 下午8:51
package main

import (
	"fmt"
	"net/http"
)

func Index(writer http.ResponseWriter, request *http.Request) {
	//request 只读，承载客户端全部请求信息（路径、GET 参数、POST 表单、请求头、Body 等）
	//writer 写入响应（返回文本 / JSON/HTML、设置状态码、响应头）
	fmt.Println(request.URL.Path)       //获取客户端当前访问的具体路径（比如他访问的是 / 还是 /login）。
	writer.Write([]byte("Hello World")) //发送数据。把 "Hello World" 这句话发送给正在访问的客户端。
}

func main() {
	http.HandleFunc("/", Index) //规定规则。告诉程序：当有人访问网址根目录 / 时，就去执行 Index 这个函数里的代码。
	/*
		任何能处理 HTTP 请求的东西，都必须是 Handler。不管是普通函数、自定义结构体，最终都要变成 Handler 才能交给服务器处理请求。
		http.Handler（接口）:规定只要一个类型有ServeHTTP(w http.ResponseWriter, r *http.Request)方法，它就是 Handler。
		func Index本身只是普通函数，不自带ServeHTTP 方法，本身不满足 Handler 接口规则，不能直接丢给路由。
		Go 提供了适配器 HandlerFunc：把普通函数包一层，自动补上 ServeHTTP 方法，让普通函数变成合法 Handler。
		等同于mux.Handle("/", http.HandlerFunc(index))
		请求到来时，服务器会自动调用func Index方法，执行业务代码。
	*/
	/*
			DefaultServeMux（全局路由） vs NewServeMux（独立路由）
			路由是什么？
			  路由就是一个“地址 - 处理函数对照表”，当有人访问时，根据规则，执行对应的代码。如访问 / → 执行 index 函数；访问 /ping → 执行 ping 函数。
			  服务器收到请求，靠路由查表，找到对应代码执行。
			全局路由 http.DefaultServeMux：所有规则都注册到全局路由上，全局路由会自动匹配请求，并执行对应的代码。http.HandleFunc("/", index)
			  问题：项目大、分多个文件时，所有人都往这一张表里注册路由，很容易出现两个地方注册同一个路径，发生覆盖冲突。

			独立路由mux := http.NewServeMux()：每个项目有自己的独立路由，项目内部使用独立路由，不影响其他项目。
		      新建一张全新、独立的路由对照表，和全局 DefaultServeMux 完全分开，互不干扰：mux := http.NewServeMux() mux.HandleFunc("/", index)
			  好处：可以创建多张路由表，做路由分组，不会互相污染，正式项目都这么用。
	*/
	fmt.Println("web server listening at http://127.0.0.1:8080")
	http.ListenAndServe("127.0.0.1:8080", nil) //启动服务器并卡住程序。让程序在 8080 端口一直运行着，等待别人来访问。
	/*
		第二个参数需要传一张路由表（本质是 Handler）。
		传 nil：服务器自动拿全局路由 DefaultServeMux 使用:http.ListenAndServe("127.0.0.1:8080", nil) 等价逻辑：用全局那张共享路由表处理所有请求。
	    传自定义 mux:完全抛弃全局路由，只用自己新建的这张表。
		           server := &http.Server{
		           Handler: mux, // 传入自己新建的独立路由表
		           }
		           server.ListenAndServe()
	*/
	//内部自动循环 Accept 连接，每一个请求自动开 goroutine，天然并发，不用手动开协程
}

/*
 控制台打印：  web server listening at http://127.0.0.1:8080
*/

/*
=======================================================================
【Go 网络编程 - HTTP 协议核心学习笔记】
=======================================================================
一、认知升级：HTTP 与 TCP 的根本区别
1. 层次不同：TCP 是传输层（底层水管），HTTP 是应用层（水管里流送的包裹）。
2. 边界问题：TCP 是流式的（会粘包）；HTTP 是有边界的（通过请求头里的 Content-Length 明确告诉对方有多长，彻底解决粘包）。
3. 交互模式：TCP 可以全双工（同时互相发）；HTTP 是严格的【请求-响应】模型（Request-Response），客户端发一次，服务端回一次，然后连接往往就结束了。

二、为什么说 Go 的 HTTP 库是“降维打击”？
回忆一下我们在纯 TCP 编程时遇到的最大痛点：
- 痛点 1：单线程卡死。
  👉 HTTP 库的魔法：`net/http` 在底层【自动为每一个客户端请求开启了一个 goroutine】！你根本不需要手动写 `go func()`，它天生就是高并发的！
- 痛点 2：拆解数据极其麻烦（手写 buffer、判断 EOF）。
  👉 HTTP 库的魔法：直接给你封装成了现成的 Request（请求对象）和 Response（响应对象），直接拿属性就行。

三、HTTP 服务端核心要素（必练）
1. 核心包：`net/http`
2. 注册路由：http.HandleFunc("路径", 处理函数)
   - 作用：告诉服务器，当用户访问某个路径（比如 "/hello"）时，该执行哪段代码。
3. 启动服务：http.ListenAndServe("地址:端口", nil)
   - 作用：相当于 TCP 的 Listen + 死循环 Accept，并且自动处理并发。
4. 处理函数签名（极其重要，必须背下）：
   func(w http.ResponseWriter, r *http.Request)
   - r (Request)：从这里读取客户端发来的所有信息（方法、URL、参数、请求体）。
   - w (ResponseWriter)：往这里写入你想返回给客户端的数据（HTML、JSON、纯文本）。

四、HTTP 客户端核心要素（必练）
1. 发送 GET 请求：resp, err := http.Get("http://127.0.0.1:8080/hello")
2. 发送 POST 请求：http.Post("URL", "数据格式", 字节流)
3. 资源释放（新手必踩坑）：
   - 收到响应后，无论如何必须关闭响应体：`defer resp.Body.Close()`，否则会严重内存泄漏！
4. 读取响应数据：
   - 使用 `io.ReadAll(resp.Body)` 一次性把服务端返回的数据全读出来。

五、你的下一步练手计划：
1. 写一个 HTTP 服务端，注册一个 "/ping" 路由，当浏览器访问它时，网页上显示 "pong"。
2. 写一个 HTTP 客户端，用代码去请求你刚刚写的服务端，并把返回的 "pong" 打印在控制台。
3. 尝试从 `r.URL.Query()` 中获取浏览器的请求参数（比如 /ping?name=zhangsan）。
=======================================================================
*/
