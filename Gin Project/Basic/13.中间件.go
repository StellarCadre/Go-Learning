// 创建时间：2026/8/3 下午7:23
package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"time"
)

/*
用户自定义的中间件，常用来实现登录认证、权限校验、日志记录、数据分页、耗时统计等功能
可以定义多个中间件，按顺序执行，也可以在某个中间件中终止后续中间件的执行
定义的中间件函数类型要是gin.HandlerFunc，即 func(*gin.Context)
请求流转示意图：
                                 ────────▶ /index
                                 ────────▶ /home
客户端请求 ────────▶ 【中间件】  ────────▶ /shop
                                 ────────▶ /video
                                 ────────▶ /user

客户端 ◀──────── 【中间件】◀──────── 业务处理返回

流程说明：
1. 所有接口请求先经过中间件，再到达对应的业务路由
2. 中间件可以在请求到达接口前做校验；也可以在接口处理完成后，对返回结果加工
3. 如果中间件校验失败（未登录、token无效），可以直接拦截请求，不再向后传递，不执行接口逻辑
中间件核心特点：
① 前置处理：请求抵达业务接口之前执行（登录校验、参数打印、跨域处理）
② 后置处理：业务接口执行完毕后，还能再次执行（统一封装返回数据）
③ 可拦截：调用c.Abort() 直接终止请求，不再向后执行后续中间件与业务函数
*/

func Hello(c *gin.Context) { //这里的func是HandlerFunc类型，不过其不是中间件函数，是业务处理函数，只不过不是匿名函数形式。
	name, _ := c.Get("name")
	c.JSON(200, gin.H{
		"message": "hello world",
		"name":    name,
	})
}

func CountTime(c *gin.Context) { //这是中间件
	fmt.Println("这是计时中间件")
	start := time.Now()
	c.Next() //调用后续的中间件或处理函数,直到所有函数都执行完毕，然后执行下面的代码
	//c.abort()//终止后续处理函数,直接跳到下面的代码
	cost := time.Since(start)
	fmt.Println("cost:", cost)
	fmt.Println("计时中间件退出")
}
func M2(c *gin.Context) {
	fmt.Println("这是其他中间件")
	c.Set("name", "aurora") //这样是给c设置值，在后续的中间件或处理函数中可以获取到这个值，和c.Get("name")搭配
	c.Next()
	fmt.Println("M2中间件退出")
}

func main() {
	r := gin.Default()
	//r.Use(CountTime,M2) 为所有路由（全局）注册中间件
	r.GET("/index", CountTime, M2, Hello) // 为单个路由注册中间件
	/*为路由组注册中间件
	usergroup:=r.Group("/user",CountTime,M2)
		{
			usergroup.GET("/login", func(c *gin.Context) {})
			usergroup.GET("/register", func(c *gin.Context) {})
			usergroup.GET("/logout", func(c *gin.Context) {})
		}
	*/
	r.Run(":9000")
}

/*
http://127.0.0.1:9000/index
流程：
CountTime前置打印 → c.Next()
        ↓
M2前置打印 → c.Next()
        ↓
Hello业务函数执行
        ↓
M2后置打印（M2中间件退出）
        ↓
CountTime后置打印（cost、计时中间件退出）
*/

/*
r:=gin.Default()默认使用了Logger和Recovery中间件，
Logger中间件会打印请求日志，包括请求方法、路径、状态码、响应时间、请求头、请求体等信息
Recovery中间件会捕获panic异常，打印堆栈信息，避免服务崩溃。
若要关闭默认中间件，可以使用gin.New()创建不带默认中间件的路由实例。
*/

/*
知识点：Gin中间件/处理器内开启goroutine的上下文规则
1. *gin.Context 属于请求专属的可变上下文，它是复用对象，并不是每次请求新建。
   一次HTTP请求处理完毕之后，gin会回收、重置该Context，供给下一条请求继续使用。

2. 风险：如果你直接在新开启的goroutine当中使用原始指针 c *gin.Context
   主协程已经结束本次请求、回收清空了上下文；子goroutine再读取c，读到的是下一个请求的数据，引发数据错乱、并发竞态、panic崩溃。

3. c.Copy() 的作用
   创建一份全新、独立、只读的上下文副本，拷贝当前请求全部参数、header、请求数据。
   新goroutine操作副本，不受原请求生命周期影响，原context回收重置不会干扰子协程。

4. 使用规范
   - 主协程：正常使用原始c
   - go func()开启的异步协程：先执行 copyCtx := c.Copy()，异步代码全部使用 copyCtx
*/
