// 创建时间：2026/7/22 上午9:54
/*
核心知识点
中间件可以多层嵌套包装，执行顺序：外层先执行前置逻辑，内层后执行；后置逻辑顺序相反
演示两层中间件：鉴权（校验请求头 Token）→ 日志 → mux 路由业务
不用修改任何原有业务代码，新增功能只新增中间件即可
*/
package main

import (
	"fmt"
	"net/http"
	"strconv"
	"time"
)

// 模拟数据集
var memberData = map[int]string{
	1001: "张三",
	1002: "李四",
}

// 业务接口：连通检测
func ping2(resp http.ResponseWriter, req *http.Request) {
	resp.Write([]byte("service is running"))
}

// 业务接口：查询成员信息
func getUser2(resp http.ResponseWriter, req *http.Request) {
	idText := req.PathValue("id")
	memberId, err := strconv.Atoi(idText)
	if err != nil {
		resp.Write([]byte("参数错误：ID必须为数字"))
		return
	}
	name, exists := memberData[memberId]
	if !exists {
		resp.Write([]byte(fmt.Sprintf("未找到编号 %d 的成员", memberId)))
		return
	}
	resp.Write([]byte(fmt.Sprintf("查询成功 编号:%d 姓名:%s", memberId, name)))
}

// 参数 innerHandler = 你传入的内层处理器，可以是 mux，也可以是另一个中间件。
// 返回值 http.Handler：包装完、带日志逻辑的新处理器
// 中间件A：访问日志
func AccessRecorder(innerHandler http.Handler) http.Handler {
	// 返回一个匿名HandlerFunc（实现ServeHTTP）
	return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		// 1. 前置逻辑：请求进来先执行（拦截、校验、打印日志）
		startTime := time.Now()
		fmt.Printf("[访问日志] 方式:%s 地址:%s 客户端:%s\n",
			req.Method, req.URL.Path, req.RemoteAddr) //打印请求信息：请求方法 GET/POST、访问路径、客户端 IP

		// 2. 放行：innerHandler.ServeHTTP(resp, req)
		innerHandler.ServeHTTP(resp, req) //执行内层的中间件/路由；不写这行=直接截断请求

		cost := time.Since(startTime)
		fmt.Printf("[访问日志] 请求完成 耗时:%s 路径:%s\n\n", cost, req.URL.Path)
	})
	//进入日志中间件 → 打印请求开始日志 → innerHandler.ServeHTTP() 放行去 mux 执行业务 → 业务执行完毕回到日志中间件 → 打印耗时日志。
}

// 入参 innerHandler 这里不再是 mux，而是日志中间件 AccessRecorder。 意思：鉴权放行后，下一层走日志逻辑。
// 中间件B：身份校验
func TokenChecker(innerHandler http.Handler) http.Handler {
	// 返回一个匿名HandlerFunc（实现ServeHTTP）
	return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		// 1. 前置逻辑：请求进来先执行（拦截、校验、打印日志）
		tokenVal := req.Header.Get("X-Token")
		if tokenVal != "admin123" {
			resp.WriteHeader(http.StatusUnauthorized)
			resp.Write([]byte("权限拒绝，Token无效"))
			return
			/*
				拦截逻辑：
				从请求头取出键 X-Token 的值；
				判断 token 不等于规定的 admin123：
				设置响应状态码 401（未授权）；
				给前端返回错误文字；
				return 直接终止函数，不会执行 innerHandler.ServeHTTP
			*/
		}
		fmt.Println("[权限校验] Token合法")
		// 2. 放行：innerHandler.ServeHTTP(resp, req)
		//    执行内层的中间件/路由；不写这行=直接截断请求
		innerHandler.ServeHTTP(resp, req)
	})
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/ping", ping2)
	mux.HandleFunc("/member/{id}", getUser2)

	// 嵌套组装中间件：外层TokenChecker，内层AccessRecorder
	allHandler := TokenChecker(AccessRecorder(mux))

	svr := &http.Server{
		Addr:           "127.0.0.1:8080",
		Handler:        allHandler,
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   10 * time.Second,
		IdleTimeout:    15 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	fmt.Println("服务启动地址：127.0.0.1:8080")
	err := svr.ListenAndServe()
	if err != nil {
		fmt.Printf("服务启动异常：%v\n", err)
	}
}

/*
1. 合法请求（携带正确 Token，走完鉴权 + 日志 + 业务全套中间件）
只打印接口返回内容（推荐）
(iwr -Uri "http://127.0.0.1:8080/member/1001" -Headers @{"X-Token"="admin123"}).Content
完整请求信息（包含状态码、响应头、内容）
iwr -Uri "http://127.0.0.1:8080/member/1001" -Headers @{"X-Token"="admin123"}
健康检测接口 /ping 合法请求
(iwr -Uri "http://127.0.0.1:8080/ping" -Headers @{"X-Token"="admin123"}).Content

2. 非法请求 1：不带任何 Token（直接拦截，不进入日志中间件）
执行会抛出 401 报错，控制台无任何日志打印，返回内容：权限拒绝，Token无效
(iwr -Uri "http://127.0.0.1:8080/member/1001").Content
非法请求 2：携带错误 Token 值（拦截）
(iwr -Uri "http://127.0.0.1:8080/member/1001" -Headers @{"X-Token"="wrongtoken"}).Content
*/
