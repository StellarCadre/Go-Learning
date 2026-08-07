package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	r := gin.Default()

	// 加载html模板文件，让Gin可以渲染index.html页面
	r.LoadHTMLFiles("C:\\Users\\Aurora\\Desktop\\Go Project\\Gin Project\\Basic\\index8.html")

	// 首页路由，展示上传表单页面
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	// 文件上传接口，和form标签action="/upload"对应
	r.POST("/upload", func(c *gin.Context) {
		// 读取上传的文件，字符串 "f1" 必须和前端input的name="f1"一致
		file, err := c.FormFile("f1")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"msg": "获取文件失败，请选择文件",
			})
			return
		}

		// 保存文件到服务器本地，file.Filename = 用户上传的原始文件名
		savePath := "./" + file.Filename
		err = c.SaveUploadedFile(file, savePath)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"msg": "文件保存失败",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"msg":      "上传成功",
			"filename": file.Filename,
		})
	})

	r.Run(":8080")
}

/*
====================核心总结====================
1、为什么后端代码必须和前端html内容对应？
① form action="/upload"  ↔ r.POST("/upload")
表单点击提交时，浏览器会自动把数据发送到action填写的地址。
如果前后地址不一致，请求会404，后端接收不到上传请求。

② <input name="f1" ↔ c.FormFile("f1")
一个表单可以同时传递普通文本、多个文件，依靠name区分每一项数据。
浏览器上传时，会给这份文件打上标识f1；后端通过这个标识找到文件。
名字不一致 → Gin无法定位文件，直接报错获取文件失败。

③ enctype="multipart/form-data" 必须写在form中
普通表单默认编码只能传输文本字符串，无法传输二进制文件。
设置这个属性后，浏览器会把文件切割成二进制块传输；
如果缺少这一行，后端无法正常解析上传文件。

④ method="post"
文件体积可能很大，GET请求无法携带大量数据，上传文件强制使用POST。

2、额外重要知识点：
① LoadHTMLFiles：Gin模板功能，用来返回html页面；前后端分离项目一般不需要这个。
② c.SaveUploadedFile：Gin封装好的工具，直接把内存中的文件写入磁盘。
③ 风险点：直接使用 file.Filename 保存存在安全隐患（文件名可能包含非法路径，后续会学习文件名处理）。
④ 当前代码是【单文件上传】，一个表单一次只能上传一个文件；
如果要支持多选文件，需要使用c.FormMultipart()，配合input multiple属性。
⑤ 保存路径 ./ 代表程序运行的根目录，文件会直接生成在项目文件夹内。

3、运行流程梳理：
1.浏览器访问 http://127.0.0.1:8080 触发GET /
2.Gin渲染index.html表单页面展示给用户
3.用户选择文件，点击上传按钮
4.浏览器发起POST请求到 /upload，携带二进制文件
5.后端通过name标识提取文件，保存到本地，返回JSON结果

4.文件上传逻辑：看懂为什么要 “先接收再保存”
用户在前端选中本地文件（文件在自己电脑）
JS / 原生表单把文件二进制流塞进请求，发给 Gin 服务
Gin 收到完整文件数据，临时存到程序内存中
c.FormFile("f1") 读取内存里的文件对象
SaveUploadedFile 把内存中的二进制写入服务器本地磁盘
文件写到服务器本地，代表数据永久存在服务端；再通过 Gin 静态路由绑定 URL 路径，网页就能发送网络请求拉取这份文件并展示。

5. 本代码当前阶段局限：
仅完成文件接收、写入服务器本地磁盘，并向前端返回上传成功信息与文件名。
未配置静态资源路由 r.Static()，浏览器无法通过网络URL请求、加载已保存的文件，
网页不能直接展示图片/下载文件。
若要实现网页访问文件，需要新增静态路由，建立URL地址与服务器本地文件目录的映射关系。
// 新增静态路由：访问 /file/xxx 对应项目根目录 ./xxx
	r.Static("/file", "./")
// 返回可直接在网页打开的完整URL
		fileUrl := "http://127.0.0.1:8080/file/" + file.Filename
		c.JSON(http.StatusOK, gin.H{
			"msg":      "上传成功",
			"filename": file.Filename,
			"url":      fileUrl,
		})
*/
