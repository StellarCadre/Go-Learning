// 创建时间：2026/6/23 下午8:53
package main

// 项目入口
func main() {
	server := NewServer("127.0.0.1", 8000) // 1. 创建Server实例：指定监听的IP（127.0.0.1表示本机）和端口（8000）
	server.Start()                         // 2. 启动服务器：开始监听端口、接收客户端连接
}
