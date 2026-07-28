// 创建时间：2026/7/14 下午9:05
package main

//手动实例化 http.Server，配置超时参数.不再使用默认全局服务，自定义全部安全超时配置
import (
	"fmt"
	"net/http"
	"time"
)

func index(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("ok"))
}

func main() {
	mux1 := http.NewServeMux()  //创建一张全新、独立的路由映射表，和全局 DefaultServeMux 完全隔离，不会出现路由覆盖冲突
	mux1.HandleFunc("/", index) //在当前这张独立路由表里注册规则：访问根路径 / 执行 index 函数
	//mux1.HandleFunc("/test", 其他函数)
	//mux2.HandleFunc("#a", 其他函数)

	server := &http.Server{
		Addr:           "127.0.0.1:8080", //Addr：服务监听地址与端口，本地只能本机访问 127.0.0.1:8080
		Handler:        mux1,             //指定当前服务使用我们刚创建的独立路由表，彻底抛弃全局路由
		ReadTimeout:    5 * time.Second,  //客户端完整发送一次请求的最大耗时，超时直接断开连接，防御慢速攻击
		WriteTimeout:   10 * time.Second, //服务向客户端返回响应的最长耗时，防止连接卡死不释放
		IdleTimeout:    15 * time.Second, //长连接空闲等待下一次请求的超时时间，空闲超时回收连接释放资源
		MaxHeaderBytes: 1 << 20,          //限制请求头最大 1MB，防止攻击者传递超大 Header 耗尽内存，属于安全限制
	}

	fmt.Println("server start at 8080")
	err := server.ListenAndServe() //启动自定义配置的 HTTP 服务，程序会阻塞在这里，持续监听 8080 端口等待客户端请求
	if err != nil {
		fmt.Println(err)
	}
}
