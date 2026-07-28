// 创建时间：2026/7/23 下午3:39
package main

/*
GET /ping 简单健康接口（适合入门测试客户端）
GET /user/{id} 路径参数查询
GET /user/list 分页查询（带 URL 参数）
POST /user/add 提交 JSON，用来练习 POST 请求
*/
import (
	"context" // 【阶段4优雅停机新增】
	"encoding/json"
	"fmt"
	"net/http"
	"os"        // 【阶段4优雅停机新增】
	"os/signal" // 【阶段4优雅停机新增】
	"strconv"
	"syscall" // 【阶段4优雅停机新增】
	"time"
)

// 模拟内存存储
var userStorage = map[int]string{
	1001: "张三",
	1002: "李四",
}
var autoIncrementID = 1003

// 接收新增用户JSON结构体
type CreateUserForm struct {
	UserName string `json:"username"`
}

// ===================== 业务接口 =====================
// 1.GET /ping 健康检测
func healthCheck(res http.ResponseWriter, reqInfo *http.Request) {
	res.Header().Set("Content-Type", "application/json;charset=utf-8")
	_ = json.NewEncoder(res).Encode(map[string]any{
		"code":    200,
		"message": "服务运行正常",
	})
}

// 2.GET /user/list 分页接口，读取URL查询参数 ?page=1&size=2
func queryUserPage(res http.ResponseWriter, reqInfo *http.Request) {
	// 从URL获取查询参数
	pageStr := reqInfo.URL.Query().Get("page")
	sizeStr := reqInfo.URL.Query().Get("size")
	// 默认值：page=1 size=5
	page := 1
	size := 5
	if pageStr != "" {
		p, err := strconv.Atoi(pageStr)
		if err == nil && p > 0 {
			page = p
		}
	}
	if sizeStr != "" {
		s, err := strconv.Atoi(sizeStr)
		if err == nil && s > 0 && s <= 50 {
			size = s
		}
	}
	// 将map数据转为切片，方便分页
	var userList []map[string]any
	for uid, name := range userStorage {
		userList = append(userList, map[string]any{
			"userId":   uid,
			"userName": name,
		})
	}
	total := len(userList)
	startIndex := (page - 1) * size
	endIndex := startIndex + size
	if startIndex > total {
		startIndex = total
	}
	if endIndex > total {
		endIndex = total
	}
	pageData := userList[startIndex:endIndex]
	res.Header().Set("Content-Type", "application/json;charset=utf-8")
	_ = json.NewEncoder(res).Encode(map[string]any{
		"code":    200,
		"message": "分页查询成功",
		"data": map[string]any{
			"list":  pageData,
			"page":  page,
			"size":  size,
			"total": total,
		},
	})
}

// 3.POST /user/add 接收JSON新增用户
func addUser(res http.ResponseWriter, reqInfo *http.Request) {
	if reqInfo.Method != http.MethodPost {
		res.Header().Set("Content-Type", "application/json;charset=utf-8")
		res.WriteHeader(405)
		_ = json.NewEncoder(res).Encode(map[string]any{
			"code":    405,
			"message": "接口仅支持POST请求",
		})
		return
	}
	var form CreateUserForm
	decodeErr := json.NewDecoder(reqInfo.Body).Decode(&form)
	defer reqInfo.Body.Close()
	if decodeErr != nil {
		res.Header().Set("Content-Type", "application/json;charset=utf-8")
		res.WriteHeader(400)
		_ = json.NewEncoder(res).Encode(map[string]any{
			"code":    400,
			"message": "JSON格式错误",
		})
		return
	}
	if form.UserName == "" {
		res.Header().Set("Content-Type", "application/json;charset=utf-8")
		res.WriteHeader(400)
		_ = json.NewEncoder(res).Encode(map[string]any{
			"code":    400,
			"message": "用户名不能为空",
		})
		return
	}
	newID := autoIncrementID
	userStorage[newID] = form.UserName
	autoIncrementID++
	res.Header().Set("Content-Type", "application/json;charset=utf-8")
	res.WriteHeader(201)
	_ = json.NewEncoder(res).Encode(map[string]any{
		"code":    201,
		"message": "用户创建成功",
		"data": map[string]any{
			"userId":   newID,
			"userName": form.UserName,
		},
	})
}

// 4.GET /user/{id} 路径参数查询单个用户
func querySingleUser(res http.ResponseWriter, reqInfo *http.Request) {
	idStr := reqInfo.PathValue("id")
	uid, parseErr := strconv.Atoi(idStr)
	if parseErr != nil {
		res.Header().Set("Content-Type", "application/json;charset=utf-8")
		res.WriteHeader(400)
		_ = json.NewEncoder(res).Encode(map[string]any{
			"code":    400,
			"message": "id必须为数字",
		})
		return
	}
	name, exist := userStorage[uid]
	if !exist {
		res.Header().Set("Content-Type", "application/json;charset=utf-8")
		res.WriteHeader(404)
		_ = json.NewEncoder(res).Encode(map[string]any{
			"code":    404,
			"message": fmt.Sprintf("不存在id=%d的用户", uid),
		})
		return
	}
	res.Header().Set("Content-Type", "application/json;charset=utf-8")
	_ = json.NewEncoder(res).Encode(map[string]any{
		"code":    200,
		"message": "查询成功",
		"data": map[string]any{
			"userId":   uid,
			"userName": name,
		},
	})
}

// ===================== 中间件 =====================
// 请求日志中间件
func RequestLogger1(innerLogic http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, reqInfo *http.Request) {
		startTs := time.Now()
		fmt.Printf("[访问日志] 方式:%s 路径:%s 客户端:%s\n",
			reqInfo.Method, reqInfo.URL.Path, reqInfo.RemoteAddr)
		innerLogic.ServeHTTP(res, reqInfo)
		useTime := time.Since(startTs)
		fmt.Printf("[访问日志] 请求结束 耗时:%s\n\n", useTime)
	})
}

// Token鉴权中间件
func AuthGuard1(innerLogic http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, reqInfo *http.Request) {
		tokenContent := reqInfo.Header.Get("X-Token")
		if tokenContent != "admin123" {
			res.Header().Set("Content-Type", "application/json;charset=utf-8")
			res.WriteHeader(http.StatusUnauthorized)
			_ = json.NewEncoder(res).Encode(map[string]any{
				"code":    401,
				"message": "Token无效，拒绝访问",
			})
			return
		}
		fmt.Println("[权限校验] Token合法")
		innerLogic.ServeHTTP(res, reqInfo)
	})
}
func main() {
	routePool := http.NewServeMux()
	// 注册全部路由
	routePool.HandleFunc("/ping", healthCheck)
	routePool.HandleFunc("/user/list", queryUserPage)
	routePool.HandleFunc("/user/add", addUser)
	routePool.HandleFunc("/user/{id}", querySingleUser)
	// 链式中间件：鉴权外层，日志内层
	coreHandler := AuthGuard1(RequestLogger1(routePool))
	httpServer := &http.Server{
		Addr:           "127.0.0.1:8080",
		Handler:        coreHandler,
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   10 * time.Second,
		IdleTimeout:    15 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	/*
		原版问题：
		ListenAndServe() 是阻塞函数：执行到这里，主线程直接卡住，后面所有代码永远不会运行；
		程序只能傻傻监听端口，没有任何机会监听 Ctrl+C 关闭信号；
		按 Ctrl+C 会直接暴力杀死进程，正在处理的请求直接中断。
		故将原本的
		startErr := httpServer.ListenAndServe()
		if startErr != nil {
			fmt.Printf("启动失败：%v\n", startErr)
		}
		替换为如下代码：
	*/

	// =====================【阶段4 全部新增优雅停机代码 开始】=====================
	// 1. 把阻塞的 ListenAndServe 丢到子协程里运行，主线程不会被卡住，主线程可以继续往下执行监听关闭信号的代码。
	go func() { //把
		fmt.Println("Day1最终服务启动：127.0.0.1:8080")
		startErr := httpServer.ListenAndServe()
		// 主动调用 server.Shutdown() 优雅关闭服务时，会返回 ErrServerClosed，不用打印报错
		if startErr != nil && startErr != http.ErrServerClosed {
			fmt.Printf("服务启动异常：%v\n", startErr)
			os.Exit(1)
		}
	}()

	// 2. 创建信号管道，监听系统终止信号
	sigChan := make(chan os.Signal, 1)
	// 注册需要捕获的信号：Ctrl+C(SIGINT)、容器/服务器终止信号(SIGTERM)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	fmt.Println("服务运行中，按下 Ctrl+C 可触发优雅停机")

	// 阻塞主线程，等待关闭信号
	/*
		从通道里读取信号。
		道里没有信号时：主线程卡在这一行，原地等待，不会结束程序；
		按下 Ctrl+C / 容器下发终止信号：通道收到信号，这行代码执行完毕，程序继续往下走停机逻辑。
	*/
	<-sigChan
	fmt.Println("\n已收到关闭信号，开始执行优雅停机流程...")

	// 3. 创建停机上下文，设置10秒兜底超时，防止请求卡死无法关闭
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	// 执行优雅关闭逻辑
	// 1. 拒绝新连接 2. 等待现有请求处理完成 3. 超时强制关闭
	/*
		传入带超时的 context，如果等待时间超过 ctx 设定时长（我们设 10 秒），会强制断开所有空闲 / 卡住的连接
	*/
	shutdownErr := httpServer.Shutdown(shutdownCtx)
	if shutdownErr != nil {
		fmt.Printf("优雅停机失败：%v\n", shutdownErr)
	} else {
		fmt.Println("服务优雅关闭完成，所有业务请求处理完毕")
	}
	// =====================【阶段4 全部新增优雅停机代码 结束】=====================
}

/*
全套代码学习递进总览（demo1 → demo2 → demo3 → demo4 → demo5 → demo6 → demo7）
阶段 1 server_demo1：最基础自定义 HTTP 服务
新增知识点：
抛弃默认全局服务，手动实例化 &http.Server 自定义服务
完整安全超时配置：ReadTimeout/WriteTimeout/IdleTimeout、限制请求头 MaxHeaderBytes
http.NewServeMux 创建独立隔离路由表，避免全局路由冲突
基础 HandleFunc 注册简单根路由，w.Write 返回纯文本
核心结构：自定义服务 + 独立 mux + 单业务函数
阶段 2 server_demo2：理解 Handler 接口本质
新增知识点：
区分两种 Handler 实现方式
结构体自定义类型，手动实现 ServeHTTP 原生 Handler
普通函数通过 HandlerFunc 适配器自动转为 Handler（日常开发主流）
mux.Handle () 接收原生 Handler、mux.HandleFunc () 接收普通业务函数
核心作用：搞懂中间件底层原理（所有中间件本质都是包装 Handler）
阶段 3 server_demo3：Go1.22 动态路径参数
新增知识点：
动态路由模板 /user/{id} 写法
r.PathValue ("占位名") 提取 URL 路径上的动态值
区分 r.URL.Path（完整原始路径）与 PathValue（提取占位参数）
核心作用：实现按 ID 查询单条资源的接口基础
阶段 4 server_demo4：单层中间件 + 内存模拟数据库
新增知识点：
map 全局变量模拟简易内存数据库（程序关闭数据丢失）
strconv.Atoi 字符串转数字，参数合法性校验、不存在判断
单层中间件开发：日志中间件 LogMiddleware
中间件执行规则：
next.ServeHTTP () 前：前置逻辑（打印请求信息）
next.ServeHTTP () 后：后置逻辑（统计耗时）
全局统一日志，不用每个接口重复写打印代码
执行链路：请求 → 日志中间件 → mux 路由 → 业务函数
阶段 5 server_demo5：多层嵌套中间件（鉴权 + 日志）
新增知识点：
中间件嵌套组装：外层鉴权 TokenChecker，内层日志 AccessRecorder
allHandler := TokenChecker (AccessRecorder (mux))
请求头读取 r.Header.Get ("X-Token")，Token 权限校验拦截
http.StatusUnauthorized (401) 未授权状态码，校验失败直接 return 截断请求
多层中间件执行顺序规则：
前置：由外到内；后置：由内到外
合法请求链路：鉴权前置打印合法 → 日志前置打印请求 → 业务执行 → 日志后置耗时
非法请求链路：鉴权校验失败直接返回，不进入日志、路由、业务
阶段 6 server_demo6：标准化 JSON 前后端接口（核心工程化改造）
新增知识点：
引入 encoding/json 包，两套 JSON 工具
json.NewEncoder (res)：后端输出 JSON 给前端
json.NewDecoder (reqInfo.Body)：解析前端 POST 提交的 JSON
统一全局 JSON 返回模板三段式：
① res.Header () 设置 Content-Type 解决中文乱码
② 可选 res.WriteHeader () 设置对应 HTTP 状态码 (400/404/405/201)
③ json.NewEncoder 输出带 code/message/data 标准 JSON 结构
POST 新增接口规范：限制 reqInfo.Method == http.MethodPost，405 拦截非法请求方式
结构体 + json 标签接收前端 JSON 参数，defer reqInfo.Body.Close () 释放网络流
全局自增 ID autoId，实现新增数据存入内存 map
改动说明：所有旧纯文本返回全部重构为标准 JSON 格式，错误、成功统一返回结构
阶段 7 server_demo7（最终综合整合版）
唯一全新知识点：URL 查询参数 + 分页逻辑
reqInfo.URL.Query ().Get ("key") 读取 URL?page=1&size=5 这类查询参数
参数默认值兜底、数值范围限制（防止传入非法分页数字）
map 转切片、下标截取实现分页列表返回
整合全部过往知识点：
双层中间件鉴权 + 日志完整保留
三种请求参数全覆盖：
① 路径参数 /user/{id} 单条查询
② Query 查询参数 /user/list?page=1 分页列表
③ POST Body JSON /user/add 新增数据
REST 接口规范落地：GET 只读、POST 新增，状态码规范区分各类错误
统一变量命名、统一 JSON 返回模板，一套代码覆盖本阶段全部基础 HTTP 接口知识点

====================【本次阶段4升级新增知识点】====================
1. 新增依赖包：context / os / os/signal / syscall，用于信号监听、停机超时上下文
2. goroutine 异步启动服务，解除 ListenAndServe 阻塞，主线程负责监听关闭信号
3. signal.Notify 捕获系统关闭信号 SIGINT(Ctrl+C)、SIGTERM(容器终止)
4. httpServer.Shutdown() 优雅关闭核心方法：停止接收新请求、等待存量请求完成
5. context.WithTimeout 设置停机兜底超时，避免程序永久卡死无法退出
6. 区分正常关闭错误 ErrServerClosed，屏蔽无关报错日志
全阶段通用不变基础（每一代代码都继承）
自定义 &http.Server 超时安全配置
http.NewServeMux 独立路由容器
中间件 Handler 包装执行机制
内存 map 模拟临时数据库（无持久化，重启丢失）
ListenAndServe () 阻塞启动服务，未改造前 Ctrl+C 暴力杀死请求
*/

/*
原版流程
main 主线程
创建路由、server
打印启动文字
执行 ListenAndServe → 主线程卡死，等待请求
没有后续代码，Ctrl+C 直接杀进程
改造后流程
main 主线程
创建路由、server
子协程启动 HTTP 服务，主线程走到信号监听代码；
创建信号通道，监听 Ctrl+C、容器终止信号；
主线程阻塞在 <-sigChan 等待；
用户按下 Ctrl+C → 信号进入通道，主线程唤醒；
创建 10 秒超时上下文；
调用 Shutdown：停止接收新请求，等待存量请求；
10 秒内全部请求跑完：打印关闭成功；
超过 10 秒还有请求卡住：打印停机失败；
main 函数执行完毕，程序正常退出。
*/
