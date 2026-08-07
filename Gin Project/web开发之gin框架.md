web介绍：
web是基于http协议进行交互的应用网络
通过浏览器或APP访问的各种资源

gin项目完整流程：
1.创建go.mod
go.mod
2.下载 Gin 依赖包
go get github.com/gin-gonic/gin
3.代码中导入
"github.com/gin-gonic/gin"


原风格：
r.GET("/book",...)
r.GET("/create book",...)
r.GET("/update book",...)
r.GET("/delete book",...)
RESTful开发风格，搭配Postman工具作为客户端的测试工具：
r.GET("/book",...) 读取、查询
r.POST("/book",...) 创建
r.PUT("/book",...) 更新
r.FELETE("/book",...) 删除