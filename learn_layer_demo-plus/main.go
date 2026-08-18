// 创建时间：2026/8/17 下午6:51
package main

import (
	"learn_layer_demo/config"
	"learn_layer_demo/handler"

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
