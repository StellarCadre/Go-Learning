// 创建时间：2026/7/30 上午10:39
package main

import (
	"github.com/gin-gonic/gin"
)

func say(c *gin.Context) { //c表示上下文，包含请求和响应的所有内容
	c.JSON(200, gin.H{ //返回json格式的数据,gin.H是map[string]any的简写
		"message": "success",
	})
	/*扩展
	// 返回纯文本
	c.String(200, "hello text")
	// 返回html页面
	c.HTML(200, "login5.html", gin.H{"title":"首页"})
	// 返回文件（下载/图片）
	c.File("test.jpg")
	// 无数据，仅返回状态码
	c.Status(204)
	*/
}

func main() {
	r := gin.Default()   //创建一个默认的gin路由引擎
	r.GET("/hello", say) //当用户使用get请求访问/hello时，执行say函数

	/*
		RESTful风格测试
	*/
	r.GET("/book", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "get success",
		})
	})
	r.POST("/book", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "post success",
		})
	})
	r.PUT("/book", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "put success",
		})
	})
	r.DELETE("/book", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "delete success",
		})
	})

	r.Run(":9000") //监听端口,不写默认8080

}

/*
启动方式：
http://127.0.0.1:9000/hello
*/

/*
=====================================
原生net/http(1.http and web use.go) VS Gin框架(2.gin and web use.go) 对比总结
1. 底层关系
Gin底层依旧封装net/http，只是简化大量重复模板代码，二者都基于HTTP协议做Web服务。

2. 路由注册差异
原生http：统一HandleFunc，不区分请求方法GET/POST，需要手动判断r.Method，同一路径无法区分多种操作；
Gin：r.GET/r.POST/r.PUT/r.DELETE分开注册路由，天然支持RESTful规范，同一路径绑定不同请求方法代表增删改查。

3. 处理函数入参不同
原生处理函数：w http.ResponseWriter, r *http.Request，手动操作响应头、字节流输出内容；
Gin处理函数：c *gin.Context，封装请求+响应全部能力，统一获取参数、返回数据。

4. 返回数据写法
原生：手动设置Header、调用w.Write([]byte)写字节，返回JSON还要手动json序列化；
Gin：内置c.JSON(状态码, gin.H{})一行直接返回标准JSON，无需手动处理编码、请求头。

5. 默认配套能力
原生http：无日志、无崩溃恢复，程序panic直接退出，404页面简陋；
gin.Default()自带日志中间件、崩溃恢复中间件，自动打印接口访问日志，捕获panic不宕机。

6. 启动监听
原生：http.ListenAndServe(":8080", nil)；
Gin：r.Run(":9000")，底层封装ListenAndServe，写法更简洁。

7. 适用场景
原生http：适合理解HTTP底层原理、简单静态文件/图片传输演示；
Gin：企业前后端分离API开发，天然适配RESTful接口，开发效率更高。
=====================================
*/
