// 创建时间：2026/8/3 下午7:11
package main

import "github.com/gin-gonic/gin"

// 路由组思想
func main() {
	r := gin.Default()
	/*
		举例：
		用户模块下，有多个路径，常规写法：
		r.GET("/user/login", func(c *gin.Context) {})
		r.GET("/user/register", func(c *gin.Context) {})
		r.GET("/user/logout", func(c *gin.Context) {})
		这样，当模块增多时，路由会越来越多，不利于管理和维护。
	*/
	usergroup := r.Group("/user") //  /user是前缀
	{
		usergroup.GET("/login", func(c *gin.Context) {})
		usergroup.GET("/register", func(c *gin.Context) {})
		usergroup.GET("/logout", func(c *gin.Context) {})
	}
	r.Run(":8080")

}
