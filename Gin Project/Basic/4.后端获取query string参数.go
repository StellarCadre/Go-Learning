// 创建时间：2026/7/31 下午4:01
package main

import "github.com/gin-gonic/gin"

// 如：https://www.bing.com/search?query=QBLH
func main() {
	r := gin.Default()
	// 演示获取URL查询参数（Query String，地址?后面拼接的参数）
	r.GET("/web", func(c *gin.Context) {
		// 1. c.Query("参数名")：直接获取query参数，无传参时返回空字符串
		name := c.Query("query")

		// 可选写法1：c.DefaultQuery(key,默认值)，参数不存在/为空时自动使用默认值
		// name := c.DefaultQuery("query", "someone")

		// 可选写法2：c.GetQuery(key)，第二个返回值标记参数是否存在，手动自定义逻辑
		// name, ok := c.GetQuery("query")
		// if !ok {
		// 	name = "someone"
		// }

		// 获取第二个查询参数，多个参数用 & 分隔
		age := c.Query("age")

		// 将获取到的参数组装JSON返回
		c.JSON(200, gin.H{
			"name": name,
			"age":  age,
		})
	})
	// 启动服务，监听8000端口
	r.Run(":8000")
}

/*
# 知识点：Query String 查询参数
1. 参数格式：接口地址?参数名1=值1&参数名2=值2
   示例1（单参数）：http://127.0.0.1:8000/web?query=tom
   示例2（多参数）：http://127.0.0.1:8000/web?query=tom&age=18
2. 三种获取Query参数的方法
   (1) c.Query("key")
      参数不存在/为空，返回空字符串，无法区分“没传”和“传了空”
   (2) c.DefaultQuery("key", "默认值")
      参数缺失时直接填充预设默认值，开发最简写法
   (3) c.GetQuery("key") (val string, exists bool)
      exists=true：前端传递了该参数；exists=false：前端未传递该参数，适合精细化判断
3. 调试工具：浏览器、Postman均可访问测试，无需前端页面
*/
