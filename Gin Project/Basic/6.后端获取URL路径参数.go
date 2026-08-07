// 创建时间：2026/7/31 下午5:56
package main

import "github.com/gin-gonic/gin"

func main() {
	r := gin.Default()

	// 路由1：双路径参数 /:user/:id
	// :user、:id 为动态路径占位符，参数直接嵌入URL路径中
	r.GET("/:user/:id", func(c *gin.Context) {
		// c.Param("占位名") 读取路径参数，返回字符串
		name := c.Param("user")
		id := c.Param("id")
		// 组装JSON返回给前端/Postman
		c.JSON(200, gin.H{
			"name": name,
			"id":   id,
		})
	})

	// 路由2：多段路径参数 /blog/:year/:month/:day
	// 适合按日期、层级资源查询的场景（日志、文章归档）
	r.GET("/blog/:year/:month/:day", func(c *gin.Context) {
		year := c.Param("year")
		month := c.Param("month")
		day := c.Param("day")
		c.JSON(200, gin.H{
			"year":  year,
			"month": month,
			"day":   day,
		})
	})

	// 启动服务监听8000端口
	r.Run(":8000")
}

/*
## 知识点：Path路径参数（URL路径参数）
1. 定义语法
   路由中使用 :xxx 标记动态占位，/分割多个参数段，无?、&符号
   例：/:user/:id 、/blog/:year/:month/:day

2. 获取方式
   c.Param("参数占位名称")，统一返回string类型，数字需手动转int

3. 测试访问地址
   ① 用户路由：http://127.0.0.1:8000/tom/123456
      name=tom  id=123456
   ② 博客日期路由：http://127.0.0.1:8000/blog/2026/7/31
      year=2026 month=7 day=31

4. 匹配规则
   路径段数量必须完全对应，少一段/多一段都会404，匹配严格
   例：访问 /tom 只会匹配失败，不会匹配 /:user/:id

5. 和之前Query参数区分
   - Path路径参数：写在URL路径里，用于标识唯一资源（用户ID、日期、文章ID）
   - Query查询参数：写在?后，用于筛选、分页等附加条件

回顾前面学习的各类参数获取方式汇总：
【1】Path路径参数
路由写法：r.GET("/user/:username/:id", ...)
获取函数：c.Param("参数名")
特点：参数嵌入URL路径内部，路径分段严格匹配，段数不一致直接404
测试地址示例：http://127.0.0.1:8080/user/zhangsan/1001

【2】Query查询参数
路由写法：r.GET("/user", ...)
获取函数：c.Query("参数名") / c.DefaultQuery()
特点：参数放在URL ? 之后，用&分隔，属于可选筛选条件
测试地址示例：http://127.0.0.1:8080/user?name=zhangsan&age=20

【3】Header请求头参数
获取函数：c.GetHeader("键名")
特点：放在请求头部，一般用于传递Token、设备信息等

【4】Form表单参数（POST表单提交）
获取函数：c.PostForm("参数名")
特点：前端form表单提交，Content-Type为application/x-www-form-urlencoded

----------------------------------------------------------------------
痛点总结：
不同类型的参数，需要调用完全不同的方法获取：
路径参数 → c.Param()
Query参数 → c.Query()
表单参数 → c.PostForm()
请求头 → c.GetHeader()
后续如果接收JSON请求体，又要额外写解析代码。
方式分散、代码重复，无法统一处理。

解决方案：Gin 参数绑定（ShouldBind）见下一节
利用结构体 + 标签tag，自动识别请求数据类型（JSON / Form表单 / Query），
自动把请求参数映射到结构体变量中，一套结构体可以兼容多种传参格式，统一接收参数，大幅简化代码。
同时支持参数校验（binding:"required"），快速判断参数是否传递。
*/
