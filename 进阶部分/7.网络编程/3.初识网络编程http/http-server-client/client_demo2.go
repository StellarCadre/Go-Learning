// 创建时间：2026/7/23 下午5:05
package main

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

// 自定义TCP连接池 Transport（连接管理器）
var customTransport = &http.Transport{
	MaxIdleConns:    20,               //全局空闲电话线最多存放 20 根，超过 20 根的空闲连接直接关掉，避免内存越占越多。
	MaxConnsPerHost: 10,               //针对同一个地址 127.0.0.1:8080，最多同时建立 10 条通话线路。 防止程序疯狂并发请求，把对方服务打崩限流。
	IdleConnTimeout: 30 * time.Second, // 空闲连接30秒未使用自动销毁释放
}

// 全局唯一Client，绑定自定义连接池
var client2 = &http.Client{
	Transport: customTransport, //http.Client 默认自带一个系统 Transport，参数不可修改。 我们自己写好连接池规则后，手动赋值给 Client，让这个客户端全程使用我们定制的连接管理规则。
}

func main() {
	// 1. 构造GET请求
	req, err := http.NewRequest("GET", "http://127.0.0.1:8080/ping", nil)
	if err != nil {
		fmt.Println("创建请求失败：", err)
		return
	}

	// 2. 携带鉴权Token头
	req.Header.Set("X-Token", "admin123")

	// 3. 发送请求
	resp, err := client2.Do(req)
	if err != nil {
		fmt.Println("请求发送失败：", err)
		return
	}
	defer resp.Body.Close()

	// 4. 读取返回内容
	content, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("读取返回数据失败：", err)
		return
	}

	// 打印结果
	fmt.Printf("HTTP状态码：%d\n", resp.StatusCode)
	fmt.Printf("接口返回内容：%s\n", string(content))
}
