// 创建时间：2026/7/31 下午7:08
package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

/*
参数绑定核心概念：
自动将前端请求携带的所有参数，统一映射赋值到自定义结构体中，替代分散的c.Query/c.PostForm/c.Param等写法
*/

type UserInfo struct {
	Username string `form:"username"  json:"username"`
	Password string `form:"password"  json:"password"`
}

func main() {
	r := gin.Default()
	r.GET("/user", func(c *gin.Context) { // 接口1：GET请求，接收URL Query参数
		/*
			常规方案：
			若请求中参数很多，每次将请求中携带的参数解析之后，都需要初始化结构体，稍微繁琐
		*/
		//username:= c.Query("username") //从链接中，取出用户输入的username==xxx和password=<KEY>的值
		//password:= c.Query("password")
		//user:=UserInfo{
		//	Username: username, //将拿到的值赋值给结构体中的字段
		//	Password: password,
		//}
		//fmt.Println(user)
		/*
			使用绑定方案：
			c.ShouldBind(&u) ：将请求中携带的参数解析之后，有username和password的部分，赋值给结构体变量u
			这里涉及对结构体的修改，所以需要传入结构体变量的地址。
			由于ShouldBind并不知道如何将请求中携带的参数，赋值给结构体的相应部分，所以需要在结构体中，为字段添加标签。即反射。

		*/
		var u UserInfo   //定义一个结构体变量
		c.ShouldBind(&u) //将传入的内容使用shouldbind绑定到结构体变量中。

		fmt.Println(u)
		c.JSON(200, gin.H{
			"message": "success",
		})
	})

	r.POST("/form", func(c *gin.Context) { //接口2：POST表单提交（form-urlencoded）,比如前端的登录、注册窗口
		var u UserInfo
		c.ShouldBind(&u)
		fmt.Println(u)
		c.JSON(200, gin.H{
			"message": "success",
		})
	})

	r.POST("/json", func(c *gin.Context) { // 接口3：POST JSON请求体（前后端分离主流格式）
		var u UserInfo
		c.ShouldBind(&u)
		fmt.Println(u)
		c.JSON(200, gin.H{
			"message": "success",
		})
	})

	r.Run(":8080")
}

/*
痛点：每种参数格式必须调用单独方法，代码分散、重复，参数多时代码冗长。
解决方案：Gin 参数绑定 ShouldBind

一、结构体标签说明
1. form:"username"：匹配GET查询参数、POST表单中的字段名username
2. json:"username"：匹配前端POST提交JSON体内的key：username
3. 原理：Gin通过反射读取标签，自动对应请求数据与结构体字段
4. 强制要求：结构体字段首字母必须大写，否则反射无法读取赋值

二、c.ShouldBind 自动适配规则
根据请求头Content-Type自动选择解析方式：
1. GET请求：自动使用Query绑定规则，读取?后的查询参数
2. POST请求
   - Content-Type=application/json：解析JSON请求体，匹配json标签
   - Content-Type=x-www-form-urlencoded：解析表单数据，匹配form标签

三、测试访问示例
1. Query参数GET访问
http://127.0.0.1:8080/user?username=xiaowang&password=123456
2. 表单POST访问
Postman选择x-www-form-urlencoded，表单key填username、password，请求地址 /form
3. JSON POST访问
Postman raw->JSON，填写{"username":"test","password":"666666"}，请求地址 /json

四、核心优势
1. 一套结构体兼容三种主流传参格式，统一接收参数
2. 无需逐个调用c.Query/c.PostForm，减少重复代码
3. 可扩展binding标签实现参数自动校验（必填、长度、数字限制等）
*/
