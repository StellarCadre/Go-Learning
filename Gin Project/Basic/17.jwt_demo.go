// 创建时间：2026/8/19 下午6:39
package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"time"
)

/*
HS256 属于对称签名算法：同一把密钥，干两件事：签名、校验签名。
jwtSecret 就是这一把 “钥匙 / 印章”，是一段只有咱们后端服务器知道的保密字符串。
生成 token 的时候（登录 GenerateToken）：
类比盖章：拿你的专属印章（jwtSecret），在文件 (A.B) 末尾盖上章 C，产出完整 token A.B.C，交给前端。
👉印章只在服务器手里，前端拿不到这个密钥！前端只拿到盖好章的 token 字符串。
解析 token 的时候（鉴权中间件 ParseToken）：
前端把 A.B.C 传回来。后端拿到 A、B、C 三段：
后端拿出同一把密钥 jwtSecret；
拿 A+B，用 HS256 重新计算一遍，算出一个全新的签名 C_new；
拿算出来的C_new 和前端给过来的 C 做对比：
如果 C_new == C → ✅印章对上，内容没有被篡改；再检查是否过期。
如果不相等 → ❌文件被人修改过！直接判定 token 非法，拒绝访问。
eg：
学校（后端）手里独有一枚印章【jwtSecret】。
学生登录成功，学校打印一张纸条（Header+Payload，写着学生 id、有效期），盖上学校私有的印章，把纸条交给学生（前端拿到 token）。
学生下次来办事，带上这张纸条（http 请求头带上 token）。
学校收到纸条：
①看纸条上写的学生 id、有效期（Payload 可以被任何人读）；
②拿自己手里原装印章，重新盖一次，对比纸条上面的章印。
如果章对不上：说明纸条被人涂改伪造，直接拒绝办理业务。

如果 jwtSecret 泄露了会发生什么？
坏人拿到你的密钥！
坏人就可以自己伪造 token：随便写任意 userId，用泄露的密钥盖上合法印章，生成完全合法的 token，冒充任意用户登录你的系统。
所以密钥是敏感机密，绝对不能写死在代码提交 git，demo 写死只是为了练习。正式项目放到 yaml 配置文件。

签名不保护 “看不见内容”，签名保护内容不能被篡改。
允许你看纸条写了什么；
但是不允许你偷偷修改纸条上的 userId 还蒙混过关。
*/
var jwtSecret = []byte("my-secret-key-123456")

/*
MyClaims是JWT 的载荷 (Payload) 结构体，就是 token 里面存放业务数据的地方。

	 jwt.RegisteredClaims是jwt/v5 库内置的结构体，等价于把jwt.RegisteredClaims里面所有字段，直接合并到MyClaims里面。
	 jwt.RegisteredClaims内置常用字段（标准 JWT 声明）：
	    ExpiresAt *NumericDate    //token 过期时间，到时间之后 token 直接失效，拒绝解析。
		IssuedAt  *NumericDate    //什么时候签发的这个 token。
		NotBefore *NumericDate    // 生效时间 nbf
		Issuer    string          // 签发者
		Subject   string          // 主题
	    等等
	 故MyClaims现在拥有两部分字段：
	    自己写的业务字段：UserId uint（我们要存的用户 ID）
	    匿名嵌入过来：ExpiresAt、IssuedAt等 jwt 标准时间字段
*/
type MyClaims struct {
	UserId uint `json:"user_id"`
	jwt.RegisteredClaims
}

// GenerateToken 传入userId，生成JWT token字符串返回给前端用户
func GenerateToken(userId uint) (string, error) {

	// ①组装claims载荷：业务数据UserId + jwt标准时间配置
	claims := MyClaims{
		UserId: userId, // 业务数据：把登录用户id放进去，后面解析token就能拿到这个id

		// 给嵌入进来的RegisteredClaims赋值过期、签发时间
		RegisteredClaims: jwt.RegisteredClaims{
			// jwt.NewNumericDate：把Go的time.Time转成jwt库需要的时间格式
			// time.Now().Add(2*time.Hour) → 当前时间往后推2小时，就是过期时刻
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(2 * time.Hour)), // 2小时后token过期
			IssuedAt:  jwt.NewNumericDate(time.Now()),                    // 签发时间：就是现在
		},
	}
	/*
		JWT 的 token 由三部分组成，用.点分隔
		Header（头部）：记录使用什么签名算法，这里就是 HS256
		Payload（载荷）：就是我们的MyClaims，包含user_id、过期时间 exp、签发时间 iat；Base64 可解码，不加密
		Signature（签名）：header+payload经过密钥jwtSecret加密算出来的签名。

		服务端校验的时候：拿到前端传过来 token，用同一个密钥重新计算签名，和 token 里面的 signature 对比。
		如果 token 被人篡改过（比如手动改 payload 里面的 user_id），重新算出来的签名对不上 → 判断 token 非法。
		密钥jwtSecret作用：
		生成 token 的时候：用它算出签名，打在 token 尾部
		解析 token 的时候：用它校验签名是否合法
		所以这个密钥不能泄露！demo 写死；真实项目放到配置文件 yaml，不能硬编码写代码里，更不能提交到 git。
	*/

	/*
		②在Go 内存中创建一个 *jwt.Token 对象token，仅仅是内存对象！还没有生成最终字符串，没有做签名！指定签名算法, 载荷对象claims:
		    jwt.SigningMethodHS256：代表后面要用 HS256 对称算法做签名。（同一个密钥用来签名、也用来校验）
		    claims 就是上面组装好的全部载荷数据
		  执行完毕后，token的Header（头部）和Payload（载荷）就组装完毕，Signature（签名）还没有
	*/
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	/*
		③使用密钥 jwtSecret 对token做签名，输出最终的token字符串：
		   SignedString(jwtSecret)的作用：
		     1.
		     把内存中的 Header 结构体做 Base64URL 编码 → 得到第一串字符串 A
		     把内存中的 Payload (claims) 结构体序列化为 JSON，再做 Base64URL 编码 → 得到第二串字符串 B
		     2.
		     拿 A.B（用点拼接 A 和 B），再传入密钥jwtSecret，执行 HS256 算法，计算出签名 Signature 字节流，
		       再做 Base64URL 编码得到字符串 C
		     3.
		     最后把三部分用.拼接：A.B.C，这一长串就是最终返回的 token 字符串！
	*/
	return token.SignedString(jwtSecret)
}

// ParseToken 解析前端传过来的token字符串A.B.C。做解析 + 签名校验 + 过期校验。
// 参数 tokenString：前端HTTP请求头携带的一长串 "A.B.C" 的jwt字符串
// 返回值：解析通过返回载荷里面的userId；失败返回 error。
func ParseToken(tokenString string) (uint, error) {
	/*
		jwt.ParseWithClaims：jwt库提供的解析函数:
		   参数1：tokenString 前端传过来完整的token字符串
		   参数2：&MyClaims{} 传入我们自定义载荷结构体的指针，库会把解析出来的数据填充进这个对象
		   参数3：是一个回调函数，【用来返回校验签名所使用的密钥】
			     jwt库内部会调用这个回调；我们在这里把服务端保密的 jwtSecret 返回给库，用来校验签名Signature
		   返回值：
		        token：是jwt库的Token结构体
		           有.Claims、.Header、.Valid等字段
		           其中，解析成功后的参数2MyClaims会被放到.Claims字段中
		              .Header要存放比如算法 HS256
		        err：解析过程中发生的错误：格式错误、签名不对、已经过期都会产生err
	*/
	token, err := jwt.ParseWithClaims(tokenString, &MyClaims{},
		/*
			func(parsedHeaderToken *jwt.Token)是上层函数内部的回调函数：
			   入参：parsedHeaderToken *jwt.Token，是接收前端的 token 字符串，先解码Header 部分，得到算法（HS256），构建的一个半成品的*jwt.Token对象
			   然后在内部可以使用这个parsedHeaderToken，这里没用到
			   返回值(interface{}, error)：
			       第一个返回值：返回校验签名需要的密钥，我们返回jwtSecret
			       第二个返回值：获取密钥过程如果出错，可以返回 error；demo 没有错误返回nil
		*/
		func(parsedHeaderToken *jwt.Token) (interface{}, error) {
			return jwtSecret, nil //属于回调函数func(token *jwt.Token)
		})
	/*
		jwt.ParseWithClaims函数完整流程：
		   切分 token 字符串，解码 Header，生成半成品*jwt.Token
		   ✅调用我们写的回调函数，把半成品 token 传给回调的入参
		   执行回调内部代码 return jwtSecret,nil → 密钥交给 jwt 库
		   jwt 拿到密钥，校验签名、校验过期。校验通过后，再填充 MyClaims 数据，最后组装完整的*jwt.Token
	*/

	// 如果err不为nil：解析直接失败
	// 可能原因：token格式不对、签名校验失败（被篡改）、token已经过期
	if err != nil {
		return 0, err // 返回0，同时把错误向外抛出，上层鉴权中间件接收这个err做拦截
	}

	/*
		token.Claims：是interface{}类型，库解析完的数据存放在这里，需要做类型断言转回我们的 *MyClaims
		ok：断言是否成功；token.Valid：jwt库内置标记，true代表签名合法 && 没有过期
	*/
	if claims, ok := token.Claims.(*MyClaims); ok && token.Valid {
		// 断言成功，校验全部通过，从载荷取出业务字段UserId返回
		return claims.UserId, nil
	}
	// 走到这里：类型断言失败，或者token.Valid=false，判定token无效
	return 0, fmt.Errorf("token无效")
}

// JWTAuthMiddleware JWT鉴权中间件
// c：Gin上下文，里面包含本次HTTP请求全部信息：请求头、请求体、url参数等
func JWTAuthMiddleware(c *gin.Context) {
	// 1. c.GetHeader("Authorization")：从HTTP请求头中读取名叫 Authorization 的值
	// 前端约定：把token放在这个请求头里，格式：Authorization: Bearer 此处放token字符串
	// 举个例子 authHeader 拿到的完整字符串："Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9......"
	authHeader := c.GetHeader("Authorization")
	// 判断：如果前端完全没有携带这个请求头，authHeader就是空字符串
	if authHeader == "" {
		c.Abort() // 终止后续handler执行，不再走到/profile的业务代码，后续的函数也不会执行
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "请求头缺少Authorization"})
		return // 结束当前中间件函数
	}

	var tokenStr string // 用来存放剥离完 "Bearer " 前缀之后，纯净的A.B.C的token字符串
	// 2. fmt.Sscanf：字符串格式化读取，提取 "Bearer " 后面那一大串token
	// 模板 "Bearer %s"：
	// 要求输入字符串必须以 Bearer 空格开头，把后面的内容扫描取出，写入&tokenStr
	// 示例：输入"Bearer abc123" → tokenStr = "abc123"
	// 返回值：第一个是成功解析的字段数量；第二个err代表格式是否匹配
	_, err := fmt.Sscanf(authHeader, "Bearer %s", &tokenStr)
	// err != nil：字符串不匹配模板，比如少写Bearer、少空格
	// tokenStr == ""：解析出来token是空
	if err != nil || tokenStr == "" {
		c.Abort() // 截断请求链，不执行业务handler
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "Authorization格式错误，应为Bearer token"})
		return
	}

	// 3.调用我们自己写好的ParseToken函数，传入纯净的tokenStr，做解析、签名校验、过期校验
	// 返回 userId：解析成功拿到用户id；err不为nil代表token有问题（篡改、过期、格式错误）
	userId, err := ParseToken(tokenStr)
	if err != nil {
		c.Abort()
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "token非法或已过期"})
		return
	}
	// 4.✅全部校验通过！把解析出来的userId存入c的草稿区Keys，供后面/profile的handler读取使用
	c.Set("userId", userId)
	// 5.c.Next()：放行！执行后面挂载的handler，也就是/profile的业务处理函数
	c.Next()
}

func main() {
	r := gin.Default()
	//模拟登录时给用户分配专属的token
	r.POST("/login", func(c *gin.Context) {
		// 1.定义本地匿名结构体req，用来接收前端POST提交过来的JSON数据
		// json标签：前端json的key和结构体字段做映射
		// 前端传 {"username":"admin","password":"123456"}
		// → req.Username = "admin"，req.Password = "123456"
		var req struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}

		// c.ShouldBindJSON(&req)
		// ①c是Gin底层收到请求后新建出来的对象，里面已经装好了前端POST的body（JSON字节流）
		// ②ShouldBindJSON读取c里面的body，按照json标签映射，把值填充到req结构体
		// ③&req 传的是地址，函数内部直接修改req里面的字段
		// ④err != nil代表：JSON格式错误、缺少字段、类型不匹配
		if err := c.ShouldBindJSON(&req); err != nil {
			// 参数错误，直接返回json给前端
			c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "参数错误"})
			return // return终止当前handler后续代码执行
		}
		// 模拟业务校验：对比账号密码，真实项目这里要查询数据库
		if req.Username == "admin" && req.Password == "123456" {
			// 账号密码正确！
			//传入用户ID=1001，得到经过签名等处理后的token字符串，之后将其返回给前端用户
			token, _ := GenerateToken(1001)
			// c.JSON：构造响应，把数据以JSON格式返回给浏览器（前端）
			// 把生成好的token返回给前端保存（前端之后访问受保护接口就带上这个token）
			c.JSON(http.StatusOK, gin.H{
				"code":  0,
				"msg":   "登录成功",
				"token": token,
			})
			return
		}
		// 走到这里：账号密码不对
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "账号密码错误"})
	})

	//模拟登录后，用户执行内部操作时，需要进行token验证（鉴权中间件）：
	//登录 →生成 token（payload 存 userId）→前端携带 token 请求接口 →中间件解析 token 拿到 userId →handler 用 userId 查询数据库得到用户数据。
	r.GET("/profile", JWTAuthMiddleware, func(c *gin.Context) {
		val, exists := c.Get("userId")
		if !exists {
			c.JSON(500, gin.H{"msg": "获取用户id失败"})
			return
		}
		userId, ok := val.(uint)
		if !ok {
			c.JSON(500, gin.H{"msg": "类型断言失败"})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"code":    0,
			"msg":     "获取个人信息成功",
			"user_id": userId,
		})
	})
	fmt.Println("服务启动 :8080")
	r.Run(":8080")
}

/*
启动:http://127.0.0.1:8080/login
{"username":"admin","password":"123456"}
*/

/*
完整闭环再梳理一遍，带上密钥角色
用户登录，账号密码正确。
组装 claims：userId、过期时间。
NewWithClaims 创建内存*jwt.Token对象，填充 Header、claims。
调用 SignedString(jwtSecret)：使用服务端保密密钥做 HS256 签名，生成 A.B.C token 字符串。返回给前端。
👉前端拿到 A.B.C；前端不知道 jwtSecret。
前端后续请求，在 Header 带上这个 token。
鉴权中间件拿到 token 字符串，调用ParseToken(tokenStr)。
ParseToken 内部，传入同一个jwtSecret，重新计算签名和 token 携带的签名比对。
签名匹配、没有过期：解析取出 userId，c.Set("userId",userId)；放行c.Next()。
签名不匹配 / 过期：c.Abort()拦截，返回 401。
*/
