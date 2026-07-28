// 创建时间：2026/7/14 下午10:01
package main

import (
	"fmt"
	"net/http"
	"time"
)

/*
理解 Handler 接口，区分 HandlerFunc 和原生 Handler
自定义路由测试：
浏览器访问 127.0.0.1:8080/hello，输出：implement Handler interface
浏览器访问 127.0.0.1:8080/func，输出：HandlerFunc adapter
*/

type HelloHandler struct{}

func (h HelloHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("implement Handler interface"))
}

func testFunc(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("HandlerFunc adapter"))
}

func main() {
	mux := http.NewServeMux()
	mux.Handle("/hello", HelloHandler{})
	/*
		自定义结构体 HelloHandler 实现了 ServeHTTP，原生满足 http.Handler，可直接传入 Handle
	*/
	mux.HandleFunc("/func", testFunc)
	/*
		普通函数没有 ServeHTTP 方法，标准库自动通过 HandlerFunc 适配器包装，转换成合法 Handler
	*/

	server := &http.Server{
		Addr:           "127.0.0.1:8080",
		Handler:        mux,
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   10 * time.Second,
		IdleTimeout:    15 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	fmt.Println("server start at 8080")
	err := server.ListenAndServe()
	if err != nil {
		fmt.Println(err)
	}
}
