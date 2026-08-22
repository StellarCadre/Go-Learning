// 创建时间：2026/8/19 下午5:16
package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

// simulateAuthMiddleware 模拟鉴权中间件
func simulateAuthMiddleware(c *gin.Context) {
	fmt.Println("【中间件】开始执行，模拟解析token，得到用户id = 5")
	// 将用户ID存入本次请求上下文，key:"userId"，值uint(5)
	// Set：往当前请求的随身储物袋存入一条数据
	c.Set("userId", uint(5))
	// c.Next()：放行，执行排在后面的Handler函数
	// 等后面Handler全部执行完毕之后，程序会回到这一行之后继续执行
	c.Next()
	fmt.Println("【中间件】handler执行完毕，回到中间件后置逻辑")
}

func main() {
	r := gin.Default()
	// 路由：访问/test，先走simulateAuthMiddleware，再走后面匿名handler
	r.GET("/test", simulateAuthMiddleware, func(c *gin.Context) {
		// Handler中读取context里面保存的数据
		val, exists := c.Get("userId")
		// exists 判断储物袋里面有没有"userId"这条记录
		if !exists {
			c.JSON(200, gin.H{"msg": "没有找到userId"})
			return
		}
		// val目前类型是interface{}，不知道真实类型，需要类型断言转为uint
		userId, ok := val.(uint)
		if !ok {
			c.JSON(200, gin.H{"msg": "类型转换失败"})
			return
		}
		fmt.Printf("【Handler】读取到userId = %d\n", userId)
		c.JSON(200, gin.H{
			"code":   0,
			"msg":    "读取上下文成功",
			"userId": userId,
		})
	})
	fmt.Println("服务启动 :8080")
	r.Run(":8080")
}

/*
启动服务后，访问：http://127.0.0.1:8080/test
*/

/*
代码介绍：
c (*gin.Context)
内部自带一张临时的请求‑私有 map（键值表）
这个存储空间，只属于当前这一次 HTTP 请求。
用户 A 发起请求,如存的 userId，用户 B 完全看不到；请求结束，里面所有数据自动清空。

Set(键名, 要存进去的数据)
可以存任意类型：uint、string、结构体都可以

Get(键名)
返回两个值：
val：存进去的数据
exists：bool，true代表找到了这个key；false代表不存在

为什么必须类型断言？
Set 的第二个参数类型是 interface{}（空接口）
空接口可以接收任意数据类型：uint、int、string、bool、结构体全都能往里丢。
但是！当你用 c.Get()拿回来的时候，Go 并不知道你当初存进去的到底是什么类型！
所以返回的val，永远是 interface{} 类型，不能直接使用，需要进行断言分析。
*/

/*
有关数据流通分析：
前端发送 HTTP 数据包（header、body、query 参数），前端没有 c 对象。
Gin 接收到这次请求，在服务端内存新建一份独立的 *gin.Context（c）。
Gin 把前端传来的请求信息解析好，放进这个 c 里面。
所以我们可以：
c.GetHeader() 读取请求头
c.Query() 获取 url 参数
c.ShouldBindJSON(&req) 把 body 映射到 DTO
c.Set() 往 c 里面存数据
c.Get读出刚才Set写进去的草稿数据
eg：JWT 场景简单复述
前端只传 token 字符串（放在 Header）。
Gin 新建 c，把 header 装进去。
中间件从 c 读取 header 拿到 token，后端解析 token 算出userId。
c.Set("userId", userId)：后端把算出来的 id 写到 c 的草稿区。
handler 中c.Get()拿到 userId，做类型断言，只把 uint 数字传给 service，绝不传 c。
*/

/*
c.Set/c.Get 使用时机与真实业务场景
前提回顾：
请求过来，Gin 底层建好c，把前端 http 数据 (header/query/body) 装入 c。
然后把c依次交给：中间件 A → 中间件 B → Handler。
读取前端传入的数据：用 c.GetHeader()、c.Query()、c.ShouldBindJSON()，不需要 Set。
c.Set()：存后端计算生成的数据（不是前端直接给的），放在 c 内部Keys草稿区，供本次请求后续函数读取；前端完全看不到。
c.Get()：读取前面某个函数Set存入的草稿数据。
使用时机：同一个 HTTP 请求内部，不同函数之间要共享后端计算出来的数据，就用 Set‑Get。
仅限于 middleware、handler；service/repository 严禁使用，拿到值之后把普通类型传下去。

场景 1：鉴权中间件（最常用，后面 JWT 就要写这个）
时机：中间件 Set，Handler 取
中间件从c.GetHeader("Authorization")读取前端传过来的 token 字符串（这是前端的数据）
后端解析 token，计算得到userId，这个 userId 不是前端直接提交，是后端解密算出来的
c.Set("userId", userId)存入草稿
走到 Handler，c.Get("userId")取出，类型断言得到 uint
Handler 把单纯数字userId传给 service，不传递 c
*/
