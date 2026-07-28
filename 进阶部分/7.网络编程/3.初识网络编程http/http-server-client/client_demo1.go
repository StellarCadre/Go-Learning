// 创建时间：2026/7/23 下午4:08
package main

import (
	"fmt"
	"io"
	"net/http"
)

/*
http.Client 是 Go 专门用来发网络请求的工具，相当于代码版的 PowerShell iwr、浏览器。
&http.Client{}：创建一个全新的客户端工具
不要每次发请求都新建 &http.Client{}，只全局创建 1 个反复使用，节省网络资源。
举例：
循环连续发 10 次请求，只用阶段 1 全局单个 client：
第 1 次请求：新建 TCP1，用完放进空闲池 → 池内：[TCP1]
第 2 次请求：新建 TCP2，用完放进空闲池 → 池内：[TCP1,TCP2]
第 3 次请求：池子单个主机上限 2，放不下第三条空闲连接
→ 只能把最早的 TCP1 直接销毁，新建 TCP3
第 4 次请求：池子满 2 条 (TCP2,TCP3)，销毁 TCP2，新建 TCP4
第 5 次、第 6 次…… 以此循环
不足：
用的是系统默认 Transport（连接管理器），我们不能自定义池子大小、空闲超时。
虽然全局只用了一个 client，没有重复创建 Client 对象，
但因为单主机空闲连接上限只有 2，并发 / 连续大量请求时，还是会频繁新建、销毁 TCP，复用效果很差。
默认 Transport 空闲连接永久不销毁。
程序长时间跑，大量闲置 TCP 一直占用端口，最终会出现端口耗尽、无法新建连接。
*/
var client1 = &http.Client{}

func main() {
	/*
		http.NewRequest(请求方式, 接口地址, 请求体) 三个参数
		"GET"：请求方式，和我们之前测试接口一致，/ping 是查询接口用 GET
		"http://127.0.0.1:8080/ping"：要访问的服务端地址
		nil：请求体，GET 接口不需要上传数据，填空 nil
		req：我们组装好的请求数据包（里面可以加 Token、参数）
		err：错误信息，如果地址格式写错，err 就不为空
	*/
	req, err := http.NewRequest("GET", "http://127.0.0.1:8080/ping", nil)
	if err != nil {
		fmt.Println("创建请求失败：", err)
		return
	}

	/*
		req.Header：请求的头部信息（服务端 AuthGuard 中间件就是读取这里的 X-Token）
		.Set("键", "值")：添加一组头数据
		作用：服务端校验 Token，不带会直接返回 401 拒绝访问。
	*/
	req.Header.Set("X-Token", "admin123")

	/*
		client.Do(req)：把我们组装好的 req 数据包，通过网络发给 8080 服务端
		resp：服务端返回的完整响应包（状态码、返回文字全部存在这里）
		err：网络错误，比如服务没启动、端口不对、断网
	*/
	resp, err := client1.Do(req)
	if err != nil {
		fmt.Println("请求发送失败：", err)
		return
	}
	/*
		resp.Body：服务端返回的正文数据流（就是接口返回的 JSON 文字）
		网络数据流用完必须关闭，不然会一直占用网络连接，程序卡死
		defer：延迟执行，代表当前函数全部代码执行完，最后自动运行 Close ()
		写在 resp 获取之后，不用手动在每个 return 前写关闭。
	*/
	defer resp.Body.Close()

	/*
		io.ReadAll(resp.Body)：把数据流里所有文字全部读出来，存到 content 变量
		content 类型是字节切片，不能直接打印文字，后面转字符串
	*/
	content, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("读取返回数据失败：", err)
		return
	}

	// 5. 打印状态码 + 返回文本
	fmt.Printf("HTTP状态码：%d\n", resp.StatusCode)
	fmt.Printf("接口返回内容：%s\n", string(content))
}

/*
阶段 1 核心知识点标注（看懂再进阶段 2）
http.NewRequest(请求方式, 地址, 请求体)
第三个参数 nil 代表 GET 无提交数据；POST 传 JSON 时这里填数据流
req.Header.Set() 给请求添加自定义头（对应服务端取 X-Token）
client.Do(req) 真正发起网络请求
defer resp.Body.Close() 硬性规范，不写会大量占用 TCP 连接，程序越跑越卡
io.ReadAll(resp.Body) 读取后端返回的所有字节，转字符串打印
*/

/*
客户端执行顺序：
创建 GET 请求包，填入服务地址
给请求包加上 X-Token=admin123
client 把数据包发给服务端
服务校验 Token、执行 /ping 逻辑，返回 JSON
客户端收到 resp 响应包
读取 resp 里的返回文字
打印状态码和 JSON
函数结束，自动执行 resp.Body.Close () 释放连接
*/
