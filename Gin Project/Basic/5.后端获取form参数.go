// 创建时间：2026/7/31 下午4:49
package main

//获取form表单
import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	r := gin.Default()
	// 加载两个HTML模板文件，填写自己本地绝对路径
	r.LoadHTMLFiles("C:\\Users\\Aurora\\Desktop\\Go Project\\Gin Project\\Basic\\login5.html",
		"C:\\Users\\Aurora\\Desktop\\Go Project\\Gin Project\\Basic\\index5.html")

	// GET /login：浏览器访问，展示登录表单页面
	r.GET("/login", func(c *gin.Context) {
		// nil 代表不给模板传递任何数据
		c.HTML(http.StatusOK, "login5.html", nil)
	})

	// POST /login：表单提交接口，接收form表单参数
	r.POST("/login", func(c *gin.Context) {
		// c.PostForm("表单name属性值") 获取表单提交参数，即用户在前端输入的账号密码
		username := c.PostForm("username")
		password := c.PostForm("password")
		//带默认值版：
		//username = c.DefaultPostForm("username", "someone")
		//password = c.DefaultPostForm("password", "<PASSWORD>")

		// 将账号密码传给index模板，渲染主页返回浏览器
		c.HTML(http.StatusOK, "index5.html", gin.H{
			"username": username,
			"password": password,
		})
	})

	// 启动服务监听8080端口
	r.Run(":8080")
}

/*
知识点梳理：
1. GET请求：浏览器地址栏直接访问，仅用于展示登录页login.html
2. POST请求：表单点击submit按钮触发，提交账号密码给后端
3. form表单核心要点：
   ① form标签 action="/login" method="post" 定义提交地址和请求方式
   ② input必须写name属性，后端c.PostForm依靠name取值，id仅给label/js使用
   ③ 提交按钮type必须为submit，普通button不会触发表单提交
4. 模板渲染逻辑：
   c.HTML(状态码, 模板文件名, 传给模板的数据)
   模板中 {{ .key }} 读取gin.H中对应的value
5. 访问流程：
   1. 浏览器打开 http://127.0.0.1:8080/login GET请求，加载登录表单
   2. 输入账号密码，点击登录（submit按钮），发起POST /login
   3. 后端接收username、password，渲染index主页并返回
*/
