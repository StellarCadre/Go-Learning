// 创建时间：2026/8/17 下午6:51
package main

import (
	"learn_layer_demo/config"
	"learn_layer_demo/handler"
	"learn_layer_demo/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化数据库
	err := config.InitDB()
	if err != nil {
		panic("数据库连接失败：" + err.Error())
	}

	r := gin.Default()
	// 用户相关路由
	userGroup := r.Group("/user")
	{
		userGroup.GET("/:id", handler.GetUserById)
		userGroup.POST("", handler.CreateUser)
		userGroup.PUT("", handler.UpdateUser)
		userGroup.DELETE("/:id", handler.DeleteUser)
	}

	apiGroup := r.Group("/api")
	//登录接口，无鉴权,生成token并返回给前端
	apiGroup.POST("/login", handler.Login) //http://127.0.0.1:8080/api/login

	//需要鉴权的路由组
	authGroup := apiGroup.Group("/")
	authGroup.Use(middleware.JWTAuthMiddleware) //.Use()的含义：给当前整个路由组，挂载中间件。
	{
		// 全部接口都自动经过JWT校验
		authGroup.GET("/profile", handler.GetProfile) //http://127.0.0.1:8080/api/profile 代表获取当前登录用户资料的页面接口。
		//authGroup.GET("/order/list", handler.GetOrderList)  //我的订单
		//authGroup.POST("/article", handler.CreateArticle)   //发布文章
		//authGroup.GET("/cart", handler.GetCart)              //购物车
	}

	r.Run(":8080")
}

/*
===================== 接口测试说明 =====================
1. 查询单个用户（GET 请求）
   地址：http://127.0.0.1:8080/user/1
   方式：浏览器直接访问，或 Postman 选择 GET 方式发送

2. 新增用户（POST 请求）
   地址：http://127.0.0.1:8080/user
   注意：浏览器默认只能发 GET，必须使用 Postman/Apifox 等接口工具
   Postman 操作步骤：
   ① 请求方式下拉选择 POST
   ② 切换到 Body 标签页
   ③ 左侧选择 raw，右侧下拉选择 JSON
   ④ 编辑区填入请求体：
      {
          "name": "张三",
          "age": 20
      }
   ⑤ 点击右上角 Send 发送

3. 更新用户（PUT 请求）
   地址：http://127.0.0.1:8080/user
   操作：同 POST，Body 选 raw-JSON，请求体需携带 ID 字段
正确方式：
{
    "ID":1,
    "name":"张三_修改后",
    "age":28
}
错误方式：会被validator校验拦截
{
    "ID":0,
    "name":"张三",
    "age":20
}

4. 删除用户（DELETE 请求）
   地址：http://127.0.0.1:8080/user/5
   操作：请求方式选 DELETE，直接发送即可
=====================================================
*/

/*
新内容：鉴权
====================================================================================
路由说明：
apiGroup := r.Group("/api")
为所有接口添加统一的URL前缀 /api。

apiGroup.POST("/login", handler.Login)
登录接口，属于公开接口，不需要携带JWT令牌鉴权。

authGroup := apiGroup.Group("/")
创建隶属于/api下的子路由组；
authGroup.Use(middleware.JWTAuthMiddleware)
给当前authGroup分组下所有接口绑定JWT鉴权中间件，访问分组内接口必须携带合法Token；
authGroup.GET("/profile", handler.GetProfile)
获取用户信息接口，属于受保护接口，必须完成登录拿到令牌之后才可访问。

-------------------------------------------------------------------------------------
后端启动方式：
在项目根目录执行 go run main.go，服务启动后监听地址：http://127.0.0.1:8080

-------------------------------------------------------------------------------------
Postman完整调用流程：
第一步：调用登录接口 http://127.0.0.1:8080/api/login
1. 请求方法选择 POST
2. 填写请求URL：http://127.0.0.1:8080/api/login
3. Authorization标签页设置为 No Auth，登录接口不需要鉴权头
4. 切换到 Body -> raw -> JSON，提交账号密码JSON数据，一定要是数据库中存在的，例如：
{
    "username":"test",
    "password":"123456"
}
5. 点击Send发送请求
6. 账号密码校验成功，后端返回JSON响应，响应内包含token字段，示例输出：
{
    "code": 200,
    "msg": "登录成功",
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxMSwiZXhwIjoxNzg3MzIwNTc4LCJpYXQiOjE3ODczMTMzNzh9.-FmQ8fiE6iwlNHZnG2xxRAJpPQJbhCDa3kxMH2-ar40"
}
7. 将返回的token字符串完整复制保存下来。

-------------------------------------------------------------------------------------
第二步：携带Token访问受保护接口 http://127.0.0.1:8080/api/profile
1. 请求方法选择 GET
2. 填写请求URL：http://127.0.0.1:8080/api/profile
3. 切换到 Authorization标签页，Auth Type选择 Bearer Token
4. 在Token输入框中粘贴刚刚登录获取到的完整令牌字符串，不要手动输入Bearer，Postman会自动生成Authorization请求头
5. GET接口无需填写Body内容，直接点击Send发送
6. Token校验通过，后端返回当前登录用户的个人信息，示例输出：
{
    "code": 200,
    "msg": "查询成功",
    "data": {
        "user_id": 11,
        "username": "test",
        "nickname": "测试用户"
    }
}
-------------------------------------------------------------------------------------
现实业务场景模拟含义：
整个流程模拟网页/APP用户登录系统。用户提交账号密码登录成功后，后端生成一条
属于该用户唯一的JWT身份令牌(token)，并将token以JSON响应的形式返回给前端。
前端接收到token后，会将其存储在浏览器localStorage、sessionStorage或者移动端
本地缓存中，作为该用户后续访问受保护资源的身份凭证。

当用户点击访问个人中心这类需要登录权限的页面时，前端会在本次HTTP请求的
Authorization请求头中自动带上该token（格式：Bearer token值），随同请求一并发送至后端。
请求到达Gin后端后，JWTAuthMiddleware鉴权中间件从gin.Context中读取到该请求头，
解析并校验token的签名有效性、是否过期；校验通过后从中提取出用户编号user_id，
并存入gin.Context向下传递。随后业务处理器取出该user_id作为查询条件去数据库检索
对应用户的数据，最终将用户信息返回前端渲染展示。
如果请求没有携带token、token篡改无效或者令牌已经过期，鉴权中间件将直接拦截请求，
不再执行业务查询逻辑，向前端返回401未授权，拒绝用户访问该受保护接口。
====================================================================================
*/
