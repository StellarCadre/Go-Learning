// 创建时间：2026/8/3 下午4:56
package main

import "github.com/gin-gonic/gin"

func main() {
	r := gin.Default()
	//访问/index的Get请求会走该方法
	r.GET("/index", func(c *gin.Context) { //获取资源
		c.JSON(200, gin.H{
			"message": "get",
		})
	})

	r.POST("/index", func(c *gin.Context) { //新增资源
		c.JSON(200, gin.H{
			"message": "post",
		})
	})

	r.DELETE("/index", func(c *gin.Context) { //删除资源
		c.JSON(200, gin.H{
			"message": "delete",
		})
	})

	r.PUT("/index", func(c *gin.Context) { //修改、更新资源
		c.JSON(200, gin.H{
			"message": "put",
		})
	})

	// r.Any：匹配该路径下所有请求方式（GET/POST/PUT/DELETE等）
	r.Any("/index", func(c *gin.Context) {})

	// NoRoute：统一捕获所有没有注册的路由请求，自定义404返回内容
	// 访问不存在地址时触发，替代Gin默认404页面
	r.NoRoute(func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "not found o",
		})
	})

	r.Run(":8080")
}

/*
==================== 路由基础知识点总结 ====================
1. 常用HTTP请求方法路由（RESTful规范对应）
   r.GET(path, handler) ：查询/获取数据
   r.POST(path, handler)：新增数据
   r.PUT(path, handler) ：全量更新数据
   r.DELETE(path, handler)：删除数据
   同一路径可以绑定多种请求方法，分别处理不同业务。

2. r.Any(path, handler)
   匹配该路径下任意请求方式，不区分GET/POST/PUT/DELETE，
   适合统一接收所有类型请求的场景，日常接口开发很少使用。

3. r.NoRoute() 自定义404处理
   当浏览器访问项目中未定义的路由地址时，会进入NoRoute注册的函数；
   可以自定义返回JSON、跳转页面、返回统一错误格式，优化用户体验。

4. 匹配规则
   路由严格区分请求路径与请求方法；
   若路径存在对应GET/POST等路由，优先匹配对应方法；
   无任何匹配时，进入NoRoute逻辑。

5. 拓展说明
   当前是基础单路由写法，后续路由组Group可以批量统一管理一类接口前缀，
   适合项目接口分层（如/api/user、/api/blog）。
*/
