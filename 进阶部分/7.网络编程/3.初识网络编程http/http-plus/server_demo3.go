// 创建时间：2026/7/14 下午10:22
package main

//自定义 Mux + Go1.22 路径参数
import (
	"fmt"
	"net/http"
	"time"
)

/*
r 只读，是客户端本次请求的完整数据包封装，里面存了整套请求数据，URL、路径、参数、请求头、POST 提交的数据、客户端 IP、请求方法全都有。
常用方法：
r.URL.Path
    拿到完整访问路径，比如访问 /user/100，r.URL.Path 结果就是 /user/100
r.PathValue("id")
    专门提取路由模板 {id} 动态占位内容，只适配 Go1.22 新路由。
*/

func ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("pong"))
}

func getUser(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id") //id 来自对r的解析，不是代码里定义、传递
	/*
		对127.0.0.1:8080/user/100
		路由写的是 {id}，所以括号里填 "id"，就能拿到路径里的内容。
		实际输入了100，所以id拿到的值就是100
	*/
	w.Write([]byte("user id:" + id))

}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/ping", ping)
	mux.HandleFunc("/user/{id}", getUser) //只要地址格式是 /user/任意内容，都会匹配这个路由，进入 getUser 函数,id要自己输入。如127.0.0.1:8080/user/100

	server := &http.Server{
		Addr:           "127.0.0.1:8080",
		Handler:        mux,
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   10 * time.Second,
		IdleTimeout:    15 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	fmt.Println("server start at 8080")
	err := server.ListenAndServe() ////启动自定义配置的 HTTP 服务，程序会阻塞在这里，持续监听 8080 端口等待客户端请求
	if err != nil {
		fmt.Println(err)
	}
}

/*
数据库查询、数据校验、结果返回，全部逻辑都封装在被绑定的函数内；如ping、getUser中
路由只负责分发，不参与任何数据操作
前端只负责输入地址 / 点击按钮发请求，不直接碰数据库，所有数据校验、查询逻辑都在后端
*/
