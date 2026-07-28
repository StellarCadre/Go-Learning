// 创建时间：2026/7/22 下午7:21
package main

/*
阶段 4、5 的业务函数是简易文本返回，适合入门看懂流程；
阶段 6 要做真实前后端对接，统一用 JSON 格式返回，这是前后端项目标准，所以代码全部重构
*/
import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

// 模拟数据库存储会员
var memberStore = map[int]string{
	1001: "张三",
	1002: "李四",
}

// 自增主键ID
var autoId = 1003

// 接收新增会员的JSON结构体
type AddMemberForm struct {
	UserName string `json:"name"`
}

// 健康检测接口 GET /ping
func ping3(res http.ResponseWriter, reqInfo *http.Request) { //reqInfo是客户端发来的，res是将要发送回去的
	/*
		旧写法：resp.Write([]byte("service is running"))
	    新写法：
	*/

	/*
		作用：给返回给客户端的数据包，提前加一个头部说明。
		浏览器收到这个头就知道：后面的正文是 JSON、中文不会乱码。
		这一步只是设置标记，还没真正把数据发出去。
		res 代表要发给前端的响应对象
		 res.Header() 获取响应头（响应的附加信息）
		 .Set(key, value) 设置一条头信息
		 Content-Type 固定标识：代表响应内容的类型
		 application/json 标准标识，代表内容是 JSON
		 charset=utf-8 解决中文乱码，不加的话张三李四可能变成问号乱码
	*/
	res.Header().Set("Content-Type", "application/json;charset=utf-8")

	/*
	    json.NewEncoder(res)
		创建一个 JSON 转换器，指定输出目的地是 res（响应流）。
		意思：转换完的文本，直接写到要发给客户端的数据包里。
		.Encode(你的map数据)
		把大括号里的 map 键值对，自动翻译成标准 JSON 字符串：
	*/
	_ = json.NewEncoder(res).Encode(map[string]any{
		"code":    200,
		"message": "服务正常运行",
	})
}

// 查询会员详情 GET /member/{id}
func getUser3(res http.ResponseWriter, reqInfo *http.Request) { //resp：响应输出对象 req：客户端请求信息
	idStr := reqInfo.PathValue("id") //获取路径 ID、转数字
	mid, parseErr := strconv.Atoi(idStr)
	if parseErr != nil {
		//旧：resp.Write([]byte("参数错误：ID必须为数字"))
		//新：
		res.Header().Set("Content-Type", "application/json;charset=utf-8") //作用：给返回给客户端的数据包，提前加一个头部说明。
		/*
			res.WriteHeader(400) 设置HTTP状态码400，代表「客户端参数错误」
			旧代码没有状态码，前端无法区分错误类型
		*/
		res.WriteHeader(400)
		_ = json.NewEncoder(res).Encode(map[string]any{
			"code":    400,
			"message": "参数错误，ID必须为数字",
		})
		return
	}
	name, exist := memberStore[mid]
	if !exist {
		//解释同上
		res.Header().Set("Content-Type", "application/json;charset=utf-8")
		res.WriteHeader(404)
		_ = json.NewEncoder(res).Encode(map[string]any{
			"code":    404,
			"message": fmt.Sprintf("不存在编号%d的会员", mid),
		})
		return
	}

	res.Header().Set("Content-Type", "application/json;charset=utf-8")
	_ = json.NewEncoder(res).Encode(map[string]any{
		"code": 200,
		"data": map[string]any{
			"memberId": mid,
			"userName": name,
		},
	})
}

// 新增会员接口 POST /member/create
// addMember 是POST 接口：前端主动提交一段 JSON 数据给后端，后端读取、校验、保存，再返回结果。
func addMember(res http.ResponseWriter, reqInfo *http.Request) { //res：用来返回 JSON 给前端.reqInfo：客户端全部请求信息，这里多了一个关键：reqInfo.Body 前端提交的 JSON 数据包
	/*
		reqInfo.Method 获取请求方式：GET/POST/PUT/DELETE
		http.MethodPost 是 Go 内置常量，等价于字符串 "POST"
		405 是标准 HTTP 状态码：请求方法不允许
		逻辑：如果用户用浏览器直接 GET 访问 /member/create，直接返回错误，终止函数 return
	*/
	if reqInfo.Method != http.MethodPost {
		res.Header().Set("Content-Type", "application/json;charset=utf-8")
		res.WriteHeader(405)
		_ = json.NewEncoder(res).Encode(map[string]any{
			"code":    405,
			"message": "该接口仅支持POST请求",
		})
		return
	}

	// 解析JSON请求体
	var form AddMemberForm //创建一个空结构体，用来存放前端传过来的数据
	/*
		json.NewDecoder(reqInfo.Body)
		reqInfo.Body：前端 POST 上传的 JSON 数据流（输入流）
		NewDecoder：创建解码器，专门读取「前端发给后端的 JSON」
		（区分：NewEncoder 是后端发给前端；NewDecoder 前端传给后端）
		.Decode(&form)
		&form：取结构体地址，把解析出来的数据填入 form 变量
		返回值 decodeErr：解析失败会报错（比如前端传的不是合法 JSON）
	*/
	decodeErr := json.NewDecoder(reqInfo.Body).Decode(&form)
	/*
		defer reqInfo.Body.Close()
		defer：函数执行结束前一定会执行
		Body 是网络流，用完必须关闭，不关闭会造成连接泄漏、服务卡顿
		规范：只要读取 req.Body，必须写这一行
	*/
	defer reqInfo.Body.Close()
	//触发场景：前端传了乱的文本、残缺 JSON、非 JSON 格式数据 统一返回 400 参数错误，模板和之前所有错误返回完全一致。
	if decodeErr != nil {
		res.Header().Set("Content-Type", "application/json;charset=utf-8")
		res.WriteHeader(400)
		_ = json.NewEncoder(res).Encode(map[string]any{
			"code":    400,
			"message": "JSON数据格式错误",
		})
		return
	}

	/*
		解析 JSON 成功后，取出 form.UserName（前端传的 name）
		如果是空字符串，直接返回错误，不执行新增逻辑。
	*/
	if form.UserName == "" {
		res.Header().Set("Content-Type", "application/json;charset=utf-8")
		res.WriteHeader(400)
		_ = json.NewEncoder(res).Encode(map[string]any{
			"code":    400,
			"message": "用户名不能为空",
		})
		return
	}

	/*写入模拟存储
	autoId 全局自增变量，初始 1003
	拿当前 autoId 作为新会员 ID，把前端传的用户名存入 map
	autoId++ 下次新建自动 + 1，保证 ID 不重复
	*/
	newMid := autoId
	memberStore[newMid] = form.UserName
	autoId++
	//201 标准 HTTP 状态码：资源创建成功（新增数据专用） 复用你已经学会的 JSON 返回模板，data 带回刚创建的 ID 和姓名
	res.Header().Set("Content-Type", "application/json;charset=utf-8")
	res.WriteHeader(201)
	_ = json.NewEncoder(res).Encode(map[string]any{
		"code":    201,
		"message": "会员创建成功",
		"data": map[string]any{
			"memberId": newMid,
			"userName": form.UserName,
		},
	})
}

// 中间件1：请求日志记录（代码同demo5中一样）
func RequestLogger(innerLogic http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, reqInfo *http.Request) {
		startTs := time.Now()
		fmt.Printf("[请求日志] 请求方式:%s 访问路径:%s 客户端地址:%s\n",
			reqInfo.Method, reqInfo.URL.Path, reqInfo.RemoteAddr)

		innerLogic.ServeHTTP(res, reqInfo)

		useTime := time.Since(startTs)
		fmt.Printf("[请求日志] 请求完成 耗时:%s 路径:%s\n\n", useTime, reqInfo.URL.Path)
	})
}

// 中间件2：Token权限校验（代码同demo5中一样）
func AuthGuard(innerLogic http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, reqInfo *http.Request) {
		tokenContent := reqInfo.Header.Get("X-Token")
		if tokenContent != "admin123" {
			res.Header().Set("Content-Type", "application/json;charset=utf-8")
			res.WriteHeader(http.StatusUnauthorized)
			_ = json.NewEncoder(res).Encode(map[string]any{
				"code":    401,
				"message": "权限不足，Token无效",
			})
			return
		}
		fmt.Println("[权限校验] Token校验通过")
		innerLogic.ServeHTTP(res, reqInfo)
	})
}

func main() {
	routePool := http.NewServeMux()
	routePool.HandleFunc("/ping", ping3)
	routePool.HandleFunc("/member/{id}", getUser3)
	routePool.HandleFunc("/member/create", addMember)

	// 中间件嵌套：外层鉴权，内层日志
	serverCore := AuthGuard(RequestLogger(routePool))

	httpServer := &http.Server{
		Addr:           "127.0.0.1:8080",
		Handler:        serverCore,
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   10 * time.Second,
		IdleTimeout:    15 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	fmt.Println("服务启动地址：127.0.0.1:8080")
	startErr := httpServer.ListenAndServe()
	if startErr != nil {
		fmt.Printf("服务启动失败：%v\n", startErr)
	}
}

/*
统一用 JSON 格式返回的通用固定模板（三段式）
// 第一段：必写！声明返回数据是JSON，解决中文乱码
res.Header().Set("Content-Type", "application/json;charset=utf-8")

// 第二段：可选！设置HTTP状态码，正常200可以省略，错误必须写
res.WriteHeader(数字)

// 第三段：必写！构造map，自动转JSON返回给前端
_ = json.NewEncoder(res).Encode(map[string]any{
    "code": 业务状态码,
    "message": "提示文字",
    // 成功场景额外加 data 存放业务数据，错误不需要data
})
*/

/*
完整执行流程（一次 POST 请求）
请求经过鉴权中间件 Token 校验、日志中间件，进入 addMember
判断请求方式不是 POST → 返回 405 错误
是 POST，则读取前端上传的 JSON 数据流 reqInfo.Body
用解码器把 JSON 映射到结构体 form
解析失败 → 返回 JSON 格式错误
解析成功，判断用户名为空 → 返回参数错误
校验全部通过，分配自增 ID，存入 map 模拟数据库
返回 201 状态码 + JSON，告知前端创建成功，返回新会员信息
*/

/*
对比阶段 5 缺失的核心新能力
阶段 5 只有 GET，只能从 URL 拿数据；阶段 6 支持 POST 接收前端提交 JSON
新增 json.NewDecoder 解析请求入参（反向 JSON 操作）
学习 defer 关闭请求流，网络编程必备规范
REST 规范：POST 新增数据返回 201 状态码，区分普通查询 200
结构体 + json 标签，规范化接收前端参数，代替手动字符串拆分
*/

/*
启动方式：
$postBody = '{"name":"王五"}'
(iwr -Uri "http://127.0.0.1:8080/member/create" -Method Post -Headers @{"X-Token"="admin123";"Content-Type"="application/json"} -Body $postBody).Content
第一行：$postBody = '{"name":"王五"}'
$postBody：PowerShell 里定义变量，$是变量开头标识，等价 Go 里 var postBody = "xxx"
'{"name":"王五"}'：单引号包裹的字符串，就是我们要传给后端接口的 JSON 数据
作用：把要提交的 JSON 先存到变量里，后面命令直接引用，不用重复写一大段 JSON
第二行：iwr = Invoke-WebRequest，Windows 自带发请求工具
-Uri "地址"：指定要访问的接口地址
-Method Post：声明这次请求是 POST 方式（新增数据专用）
-Headers @{键="值";键="值"}：请求头，两个必须项
X-Token: admin123：鉴权中间件校验用，没有会被拦截
Content-Type: application/json：告诉后端，我提交的 Body 是 JSON 格式，后端才能用json.NewDecoder正常解析
-Body $postBody：把上面变量存的 JSON，作为请求体传给接口，对应代码里 reqInfo.Body
.Content：只打印接口返回的 JSON 文本，隐藏一堆无关请求详情
*/

/*
4种请求方式的差异：
GET：查询数据，不修改服务器数据（查会员、健康检测）
不会携带请求体 Body，数据放 URL 里，浏览器直接输入地址默认就是 GET
POST：新增数据，修改服务器数据（创建会员、注册、提交表单）
必须携带 JSON / 表单 Body，浏览器地址栏无法直接发起 POST
PUT：修改已有数据
DELETE：删除数据
举例：
ping3 接口 /ping
作用：单纯检测服务能不能连通，只查询状态，用 GET，不需要新增 / 修改数据，不用限制 POST
getUser3 接口 /member/{id}
作用：根据 ID 查询已有会员，只读数据，标准 GET 接口，任何 GET 请求都合法，无需判断 Method
addMember 接口 /member/create
作用：新增会员，会修改内存里的memberStore存储，规范要求只能用 POST
如果有人用 GET 访问这个新增接口，逻辑会错乱，所以必须加判断拦截，返回 405 方法不允许
*/
