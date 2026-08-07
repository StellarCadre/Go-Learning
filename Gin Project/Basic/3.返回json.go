// 创建时间：2026/7/30 下午8:27
package main

import "github.com/gin-gonic/gin"

func main() {
	r := gin.Default()
	r.GET("/begin1", func(c *gin.Context) { //匿名函数形式，也可以写成普通函数放在外面
		data1 := map[string]interface{}{ //基本写法1，map[string]interface{}{},可以直接写成gin.H{}
			"name": "tom",
			"age":  18,
			"sex":  "male",
		}
		c.JSON(200, data1)

	})

	r.GET("/begin2", func(c *gin.Context) {
		//或者是将数据包装成一个结构体，然后再返回
		type s struct {
			Name string `json:"name"` //注意这里的写法，必须要写成大写。或者使用json标签和反射
			Age  int    `json:"age"`
			Sex  string `json:"sex"`
		}
		data2 := s{
			Name: "jack",
			Age:  20,
			Sex:  "wemale",
		}
		c.JSON(200, data2)
	})
	r.Run(":9000")
}

/*
本文件演示Gin两种返回JSON数据的方式
1. map[string]interface{} 方式
   等价gin.H（gin.H底层就是该map），快速构造键值对，适合简单临时数据
   缺点：无固定结构约束，字段易写错，大型项目不推荐
2. 自定义结构体方式（企业开发首选）
   ① 结构体字段首字母必须大写，否则外部无法反射读取
   ② json标签`json:"xxx"`可以自定义返回json的字段名，实现大小写转换、别名
   优点：结构规范、类型约束、可读性高，后续配合GORM模型统一使用
核心函数c.JSON(状态码,数据)：自动设置json响应头、序列化数据返回前端
访问地址：
http://127.0.0.1:9000/begin1  map格式json
http://127.0.0.1:9000/begin2  结构体格式json
*/
