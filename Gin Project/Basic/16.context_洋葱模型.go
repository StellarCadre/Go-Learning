// 创建时间：2026/8/19 下午6:14
package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func mid1(c *gin.Context) {
	fmt.Println("mid1 >>> 【前置】执行")
	c.Next() // 去执行后续逻辑，等待后续全部跑完，才回到下面
	fmt.Println("mid1 <<< 【后置】执行")
}

func mid2(c *gin.Context) {
	fmt.Println("mid2 >>> 【前置】执行")
	c.Next()
	fmt.Println("mid2 <<< 【后置】执行")
}

func main() {
	r := gin.Default()

	// 注册顺序 mid1，mid2，然后业务handler
	r.GET("/demo", mid1, mid2, func(c *gin.Context) {
		fmt.Println("--------Handler业务逻辑执行--------")
		c.JSON(200, gin.H{"msg": "ok"})
	})

	fmt.Println("服务启动 :8080")
	r.Run(":8080")
}

/*
启动服务后，访问：http://127.0.0.1:8080/demo
*/
