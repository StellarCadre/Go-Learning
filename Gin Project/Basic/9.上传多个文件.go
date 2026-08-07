// 创建时间：2026/8/2 下午7:49
package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func main() {
	r := gin.Default()
	// 处理multipart forms提交文件时默认的内存限制是32 MiB
	// 可以通过下面的方式修改
	// router.MaxMultipartMemory = 8 << 20 // 8 MiB

	// 加载html模板文件
	r.LoadHTMLFiles("C:\\Users\\Aurora\\Desktop\\Go Project\\Gin Project\\Basic\\index9.html")

	// 访问首页，展示多文件上传表单页面
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index9.html", nil)
	})

	// 多文件上传接口
	r.POST("/upload", func(c *gin.Context) {
		/*
			MultipartForm：解析当前请求中 enctype="multipart/form-data" 格式的完整表单数据，返回一个 MultipartForm 结构体对象，包含两类数据：
			    form.File：所有上传文件的数组（多文件上传核心）
			    form.Value：表单里普通文本输入框的数据（username、password 这类）
		*/
		form, err := c.MultipartForm()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"msg": "表单解析失败！"})
			return
		}

		/* form.File["file"]：获取前端name="file"下所有文件，文件切片数组
		form.File是一个 map，key：前端 <input name="xxx"> 里的 name 值value：同名输入框上传的所有文件切片（数组）
		form.File["file"]取 map 中 key 为 "file" 的文件数组
		*/

		files := form.File["file"]

		// 循环遍历所有上传文件，逐个保存
		for index, file := range files {
			log.Println("上传文件名称：", file.Filename)
			// index后缀，避免同名文件覆盖
			dst := fmt.Sprintf("C:/tmp/%s_%d", file.Filename, index)
			// 保存文件到目标路径
			err = c.SaveUploadedFile(file, dst)
			if err != nil {
				log.Printf("文件 %s 保存失败：%v\n", file.Filename, err)
			}
		}

		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("%d files uploaded!", len(files)),
		})
	})

	r.Run() // 默认监听 :8080
}

/*
====================核心知识点总结====================
1、单文件 VS 多文件上传API区别
① c.FormFile("key")
适用：一次仅上传一个文件，返回单个*multipart.FileHeader
② c.MultipartForm() → form.File["key"]
适用：多选多文件上传，返回 []*multipart.FileHeader 文件切片，配合for range循环处理

2、前端html对应修改（实现多选文件）
<input type="file" name="file" multiple>
重点：增加 multiple 属性，允许一次性选择多个文件；name="file" 和后端form.File["file"]保持一致

3、MaxMultipartMemory原理
文件上传时，小于阈值的数据放内存；超大文件，超出部分落地临时文件，避免大文件耗尽服务器内存。

4、防止覆盖文件方案
dst路径拼接index编号：同一个名称的多个文件上传时，不会互相覆盖。

5、注意事项
① 需要提前手动创建 C:/tmp 文件夹！文件夹不存在会导致保存失败。
② 前后端分离场景下，Vue前端构造FormData时，多次append同名key即可实现多文件上传。
③ 同样需要请求头 enctype="multipart/form-data"，和单文件上传要求一致。

*/

/*
1. c.MultipartForm()
解析multipart/form-data类型的完整表单，返回MultipartForm结构体
结构体包含两部分：
  form.File：存储所有上传文件，map结构
  form.Value：存储表单普通文本参数（输入框文字）
err接收解析失败的错误，用于参数校验

2. form.File["file"]
form.File是map，键对应前端input标签的name属性
form.File["file"] 取出name="file"上传的全部文件，返回文件切片数组
单文件只用取第一个元素，多文件需要循环遍历切片批量保存

3. 配套前端要求
<input type="file" name="file" multiple>
name值必须和这里括号内的字符串完全匹配，否则读取不到文件
*/
