// 创建时间：2026/7/23 下午7:45
package main

// ========== 阶段3新增：创建3秒超时上下文 ==========
import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

// 自定义TCP连接池 Transport
var customTransport1 = &http.Transport{
	MaxIdleConns:    20,
	MaxConnsPerHost: 10,
	IdleConnTimeout: 30 * time.Second,
}

// 全局唯一Client，绑定自定义连接池
var client = &http.Client{
	Transport: customTransport1,
}

func main() {
	// ========== 阶段3新增：创建3秒超时上下文 ==========
	/*
		context.Background()
		作用：生成一个空白根上下文，无超时、无限制，是 main 函数里固定标准写法。
		类比：一张空白白纸，用来承载超时规则。
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		context.WithTimeout(父上下文, 超时时间)：给空白白纸加上「3 秒自动过期」规则
		返回两个变量：
		ctx：携带 3 秒超时规则的上下文对象，后面要传给请求
		cancel：手动取消函数，用来释放上下文资源
		含义：这次网络请求最多等待 3 秒，3 秒没收到服务返回，直接强制断开请求，报错退出。
	*/
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	// 函数结束前释放ctx资源，必须写
	//cancel() 作用：手动销毁 ctx，释放底层 goroutine、定时器资源
	//defer：函数全部执行完，自动执行 cancel ()
	//规范要求：只要用了WithTimeout/WithCancel，必须写defer cancel()，否则会内存泄漏。
	defer cancel()

	// 构造带ctx的请求，替换原来的NewRequest
	req, err := http.NewRequestWithContext(ctx, "GET", "http://127.0.0.1:8080/ping", nil)
	//把刚才带 3 秒超时的ctx传入请求，这条请求才会受 3 秒限制。
	//如果还用旧的NewRequest，写再多 ctx 也不会生效。
	if err != nil {
		fmt.Println("创建请求失败：", err)
		return
	}

	// 携带鉴权Token头
	req.Header.Set("X-Token", "admin123")

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("请求发送失败：", err)
		return
	}
	defer resp.Body.Close()

	// 读取返回内容
	content, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("读取返回数据失败：", err)
		return
	}

	// 打印结果
	fmt.Printf("HTTP状态码：%d\n", resp.StatusCode)
	fmt.Printf("接口返回内容：%s\n", string(content))
}

/*
为什么需要这个东西？不写会有什么问题？
场景：服务端接口卡顿、卡死，10 秒都不返回数据。
阶段 1、2 代码：客户端会无限阻塞等待，goroutine 一直挂着，请求多了程序卡死；
阶段 3 ctx 超时：到 3 秒自动切断网络连接，抛出超时错误，协程正常释放，不会堆积。
举个生活例子
你打电话给别人（发请求）：
无 ctx：对方一直不接，你永远举着电话等；
ctx 3 秒超时：等待 3 秒没人接，系统自动挂断电话，不用一直耗着。
*/

/*
完整执行逻辑演示
创建 3 秒超时 ctx，注册延迟释放 cancel
把 ctx 绑定到 GET 请求
client.Do 发送请求
服务 3 秒内返回：正常读取数据，函数结束自动 cancel 释放 ctx
服务超过 3 秒无响应：底层自动断开 TCP，Do 直接返回超时 error，后续代码不执行
*/
