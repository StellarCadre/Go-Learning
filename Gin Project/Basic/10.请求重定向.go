// 创建时间：2026/8/2 下午8:09
package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	r := gin.Default()

	// ========== 方式一：HTTP 重定向（客户端重定向） ==========
	// 当浏览器访问 /index 时，会收到 301 状态码，浏览器自动跳转到百度
	r.GET("/index", func(c *gin.Context) {
		// 注释掉的代码：原本会返回 JSON 数据
		//c.JSON(200, gin.H{
		//	"message": "ok",
		//})

		// Redirect：HTTP 重定向
		// 参数1：http.StatusMovedPermanently = 301（永久重定向）
		// 参数2：目标 URL（注意：原代码中 "http://www.baidu.com " 末尾有多余空格，应该去掉）
		// 效果：浏览器地址栏会变成 www.baidu.com，是客户端发起的新请求
		c.Redirect(http.StatusMovedPermanently, "http://www.baidu.com")
	})

	// ========== 方式二：服务端内部转发（路由重定向） ==========
	// 当访问 /a 时，服务端内部将请求转发给 /b 处理，浏览器地址栏不变
	r.GET("/a", func(c *gin.Context) {
		// 修改当前请求的 URL 路径为 /b
		c.Request.URL.Path = "/b"
		// 让 Gin 重新处理这个修改后的请求（相当于内部转发）
		// 注意：这里不会重新匹配路由中间件，只是将请求交给对应的 handler
		r.HandleContext(c)
	})

	// /b 是实际处理请求的路由，返回 JSON 数据
	r.GET("/b", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"test": "跳转成功", // 访问 /a 时，浏览器会显示这个 JSON
		})
	})

	// 启动服务器，监听 8080 端口
	r.Run(":8080")
}

/*
启动方式：
http://127.0.0.1:8080/index   → 301 跳转到百度
http://127.0.0.1:8080/a      → 内部转发到 /b，返回 {"test":"跳转成功"}
*/
