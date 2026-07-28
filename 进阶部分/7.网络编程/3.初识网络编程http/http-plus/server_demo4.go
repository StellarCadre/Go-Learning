// 创建时间：2026/7/16 下午8:02
package main

import (
	"fmt"
	"net/http"
	"strconv"
	"time"
)

// 模拟数据库用户数据
var userData = map[int]string{
	1001: "张三",
	1002: "李四",
}

// 业务处理函数1
func ping1(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("pong 健康检测成功"))
}

// 业务处理函数2：拿到输入的id，查询数据库userData，返回结果
func getUser1(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	uid, err := strconv.Atoi(idStr)
	if err != nil {
		w.Write([]byte("错误：ID必须是数字"))
		return
	}
	name, ok := userData[uid]
	if !ok {
		w.Write([]byte(fmt.Sprintf("id=%d 的用户不存在", uid)))
		return
	}
	w.Write([]byte(fmt.Sprintf("查询成功：id=%d，用户名=%s", uid, name)))
}

// 中间件A：访问日志
func LogMiddleware(next http.Handler) http.Handler {
	// 返回实现Handler接口的匿名函数
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		// 请求进入时打印日志
		fmt.Printf("[请求开始] 方法:%s 路径:%s 客户端IP:%s\n", r.Method, r.URL.Path, r.RemoteAddr)
		//这一段在路由匹配、执行业务函数之前运行。 可以做：打印日志、校验 token 鉴权、跨域 header、限流、黑名单拦截。

		// 执行后续路由/业务函数
		next.ServeHTTP(w, r)
		//next 就是我们传入的 mux（路由） 执行 next.ServeHTTP(w,r) 才会去匹配 /ping、/user/{id}，调用你的 ping /getUser 函数。

		// 业务代码执行完，才会走到这里，统一记录耗时、状态码。
		cost := time.Since(start)
		fmt.Printf("[请求结束] 耗时:%s 路径:%s\n\n", cost, r.URL.Path)
	})
}

func main() {
	//下面的代码从上到下执行，只是提前把路由和中间件组装好，只是「预制零件」，不会处理任何浏览器请求。
	mux := http.NewServeMux()
	// 注册路由
	mux.HandleFunc("/ping1", ping1)
	mux.HandleFunc("/user/{id}", getUser1)

	// 用中间件包裹路由mux，所有请求都会先走日志中间件，包括这些业务处理函数
	wrapHandler := LogMiddleware(mux)

	server := &http.Server{
		Addr:           "127.0.0.1:8080",
		Handler:        wrapHandler, // 传入包装后的Handler，不再直接传mux
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   10 * time.Second,
		IdleTimeout:    15 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	fmt.Println("服务启动：127.0.0.1:8080")
	err := server.ListenAndServe()
	if err != nil {
		fmt.Printf("服务启动失败：%v\n", err)
	}
}

/*
用中间件，每个业务函数都写日志:
如果不包装中间件，只能把日志代码复制粘贴到 ping、getUser 每一个业务函数里，代码极度重复、难维护。
新增 10 个接口，就要复制 10 遍日志代码；以后想改日志格式，10 个函数全部手动改，极易漏改。
中间件包装（只写一次，全局生效）:
日志逻辑只写在 LogMiddleware 里，只写一遍。
把整个 mux 丢进去包装，所有路由自动走日志逻辑，业务函数干干净净，只写核心业务。
*/

/*
浏览器请求 → wrapHandler（LogMiddleware）打印前置日志 → next.ServeHTTP () → mux 匹配路由执行业务函数 → 回到中间件打印耗时 → 响应返回浏览器
*/

/*
测试访问方式（浏览器 /curl）
健康接口：http://127.0.0.1:8080/ping1
用户查询：http://127.0.0.1:8080/user/1001
*/
