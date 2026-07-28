// 创建时间：2026/7/23 下午3:39
package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
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

	fmt.Println("Day1最终服务启动：127.0.0.1:8080")
	startErr := httpServer.ListenAndServe() //启动自定义配置的 HTTP 服务，程序会阻塞在这里，持续监听 8080 端口等待客户端请求
	if startErr != nil {
		fmt.Printf("启动失败：%v\n", startErr)
	}
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
全阶段通用不变基础（每一代代码都继承）

自定义 &http.Server 超时安全配置
http.NewServeMux 独立路由容器
中间件 Handler 包装执行机制
内存 map 模拟临时数据库（无持久化，重启丢失）
ListenAndServe () 阻塞启动服务，Ctrl+C 关闭程序
*/
