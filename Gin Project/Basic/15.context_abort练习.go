// 创建时间：2026/8/19 下午6:07
package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

// 模拟鉴权失败的中间件
func authFailMiddleware(c *gin.Context) {
	fmt.Println("【中间件】检测未登录，拦截请求")
	// 终止整条请求链，后面的handler不会运行
	c.Abort()
	// ⚠️ Abort 不会自动返回响应！必须手动返回JSON
	c.JSON(401, gin.H{"code": 401, "msg": "未登录，禁止访问"})
}

func main() {
	r := gin.Default()
	r.GET("/secret", authFailMiddleware, func(c *gin.Context) {
		// 这一段代码永远不会执行，因为前面执行了Abort
		fmt.Println("Handler被执行了！")
		c.JSON(200, gin.H{"msg": "成功访问私密资源"})
	})

	r.Run(":8080")
}

/*
启动服务后，访问：http://127.0.0.1:8080/secret
*/
