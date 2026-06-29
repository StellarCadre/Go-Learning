// 创建时间：2026/6/29 下午8:39
package main

/*
前面你写的 main.go / Server.go / user.go 全部是服务端代码（等待别人来连接、处理聊天、私聊、改名）。
现在这份 client.go 是客户端程序，作用：主动连接你的聊天室服务端，充当聊天用户。
核心作用区分
Server：开门监听 8000 端口，被动等待连接
Client：主动拨号 127.0.0.1:8000，连上服务端，用来发消息、收广播 / 私聊
*/

import (
	"flag"
	"fmt"
	"net"
)

// 结构体存储客户端所有配置与 TCP 连接句柄。
type Client struct {
	ServerIp   string   // 要连接的服务器IP
	ServerPort int      // 服务器端口8000
	Name       string   // 预留：当前客户端昵称（暂时没使用）
	conn       net.Conn // 和服务端的TCP连接对象
}

func NewClient(serverIp string, serverPort int) *Client {
	// 创建客户端对象
	client := &Client{
		ServerIp:   serverIp,
		ServerPort: serverPort,
	}
	//链接服务器，net.Dial：主动发起TCP连接（客户端核心API），net.Dial 三次握手连接服务端
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", serverIp, serverPort)) //组装地址，并返回一个TCP连接对象
	if err != nil {
		fmt.Println("net.Dial err:", err) //连接失败返回 nil
		return nil
	}
	client.conn = conn //成功把连接存入 Client 并返回
	return client
}

var serverIp string
var serverPort int

func init() { //init() 函数在 main() 执行前自动运行，提前注册两个命令行参数,
	flag.StringVar(&serverIp, "ip", "127.0.0.1", "设置服务器IP地址（默认127.0.0.1）")
	flag.IntVar(&serverPort, "port", 8000, "设置服务器端口（默认8000）")
}

func main() {
	//解析命令行参数
	flag.Parse()

	client := NewClient(serverIp, serverPort) //创建客户端。支持默认参数和自定义参数，ip和port
	if client == nil {
		fmt.Println("连接服务器失败")
		return
	}
	fmt.Println("连接服务器成功") //连接失败直接退出；成功打印提示

	//开始执行客户端业务
	select {}
}
