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
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
)

// 结构体存储客户端所有配置与 TCP 连接句柄。
type Client struct {
	ServerIp   string   // 要连接的服务器IP
	ServerPort int      // 服务器端口8000
	Name       string   // 预留：当前客户端昵称（暂时没使用）
	conn       net.Conn // Conn 是 Go 标准库的TCP 连接接口，代表你客户端和服务端之间建立好的一条双向网络通道。类比：一根双向水管，一边是你的客户端，另一边是聊天室服务端。

	flag int //当前用户选择的模式
}

func NewClient(serverIp string, serverPort int) *Client {
	// 创建客户端对象
	client := &Client{
		ServerIp:   serverIp,
		ServerPort: serverPort,

		flag: 999,
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

// 菜单显示功能
func (client *Client) Menu() bool {
	fmt.Println("1.公聊模式")
	fmt.Println("2.私聊模式")
	fmt.Println("3.更新用户名")
	fmt.Println("0.退出")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	var flag int
	fmt.Sscan(scanner.Text(), &flag)
	if flag >= 1 && flag <= 3 {
		client.flag = flag
		return true
	} else if flag == 0 {
		client.flag = 0
		return true
	} else {
		fmt.Println("请输入正确的选项")
		return false
	}
}

var serverIp string
var serverPort int

func init() { //init() 函数在 main() 执行前自动运行，提前注册两个命令行参数,ip和port
	flag.StringVar(&serverIp, "ip", "127.0.0.1", "设置服务器IP地址（默认127.0.0.1）")
	flag.IntVar(&serverPort, "port", 8000, "设置服务器端口（默认8000）")
}

// 处理server端发送的消息,打印
func (client *Client) DealResponseMessage() {
	buf := make([]byte, 4096)
	for {
		n, err := client.conn.Read(buf) //client.conn.Read(buf) 阻塞等待：服务器没发消息时，代码卡在这不动； 一旦服务端下发数据，就把数据读到 buf 缓冲区里。
		if err != nil {
			fmt.Println("\n服务器连接断开，程序退出")
			os.Exit(0)
		}
		// 只打印本次读到的新数据，不打印缓冲区脏数据
		fmt.Print(string(buf[:n]))
	}
}

// client改名操作,并发送给server
func (client *Client) Rename() bool { //   直接输入新昵称
	fmt.Println("【只需要输入纯昵称，不要加/rename前缀】")
	fmt.Println("请输入新的用户名：")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	input := scanner.Text()

	// 输入exit直接退出改名，回到菜单
	if input == "exit" {
		fmt.Println("放弃改名，返回主菜单")
		return true
	}
	client.Name = input
	if client.Name == "" {
		fmt.Println("用户名不能为空！")
		return false
	}
	sendSeg := fmt.Sprintf("/rename %s\n", client.Name) //要发给服务器的指定文本内容。
	_, err := client.conn.Write([]byte(sendSeg))        //[]byte(sendSeg)：把字符串转为字节切片 网络传输只能传字节，不能直接发字符串，必须转 []byte.client.conn.Write(字节切片)
	//conn.Write作用：把字节数据通过 TCP 连接发送给服务端程序
	if err != nil {
		fmt.Println("rename err:", err)
		return false
	}
	fmt.Println("改名指令已发送，等待服务器反馈")
	return true
}

// client选择公聊操作,并发送给server
func (client *Client) PublicChat() bool {
	fmt.Println("=====公聊模式=====")
	fmt.Println("输入消息发送，输入 exit 退出公聊回到菜单")
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("请输入广播消息：")
		scanner.Scan()
		ChatMsg := scanner.Text()
		// 退出公聊模式
		if ChatMsg == "exit" {
			fmt.Println("退出公聊模式")
			return true
		}
		if ChatMsg == "" {
			fmt.Println("发送消息不能为空！")
			continue
		}
		// 追加换行，防止服务端消息粘包
		sendSeg := fmt.Sprintf("%s\n", ChatMsg)
		_, err := client.conn.Write([]byte(sendSeg))
		if err != nil {
			fmt.Println("公聊消息发送失败 err:", err)
			return false
		}
		fmt.Println("在线广播内容发送成功")
	}
}

// client私聊操作,并发送给server
func (client *Client) PrivateChat() bool {
	fmt.Println("=====私聊模式=====")
	fmt.Println("输入 exit 随时退出私聊返回菜单")
	scanner := bufio.NewScanner(os.Stdin)

	// 第一步：先输入一次私聊目标用户
	fmt.Print("请输入对方昵称：")
	scanner.Scan()
	targetName := scanner.Text()
	if targetName == "exit" {
		fmt.Println("退出私聊模式")
		return true
	}
	if targetName == "" {
		fmt.Println("昵称不能为空，退出私聊")
		return false
	}
	// 第二步：固定目标，循环发送多条私聊消息
	for {
		scanner.Scan()
		msg := scanner.Text()

		if msg == "exit" {
			fmt.Println("退出私聊模式")
			return true
		}
		if msg == "" {
			fmt.Println("消息不能为空，请重新输入！")
			continue
		}
		// 拼接标准私聊指令 /to 目标昵称 消息
		sendSeg := fmt.Sprintf("/to %s %s\n", targetName, msg)
		_, err := client.conn.Write([]byte(sendSeg))
		if err != nil {
			fmt.Println("私聊发送失败，连接已断开")
			return false
		}
		fmt.Printf("私聊【%s】发送成功\n", targetName)
	}
}

// client端，具体的业务逻辑
func (client *Client) Run() {
	for client.flag != 0 {
		for client.Menu() != true {
		}
		switch client.flag {
		case 1:
			client.PublicChat()
		case 2:
			client.PrivateChat()
		case 3:
			client.Rename()
		}
	}
	fmt.Println("程序正常退出")
	client.conn.Close() // 关闭TCP连接释放资源
}

func main() {
	//解析命令行参数
	flag.Parse()

	client := NewClient(serverIp, serverPort) //创建客户端。支持默认参数和自定义命令行参数，ip和port
	if client == nil {
		fmt.Println("连接服务器失败")
		return
	}
	fmt.Println("连接服务器成功") //连接失败直接退出；成功打印提示

	go client.DealResponseMessage() //单独开一个goroutine去处理服务器端发送的消息

	//开始执行客户端业务
	client.Run()
	//一直循环展示菜单、等你输入数字 1/2/3/0、输入昵称 / 聊天文字、调用 conn.Write 发数据给服务器。
	//这段代码运行时，整个程序会停在 fmt.Scanf / scanner.Scan() 等待你输入，此时不会去读服务器发过来的任何文字。
}

/*
【控制台输入避坑警示，务必牢记】
1. fmt.Scanf 与 fmt.Scanln 严禁混用！底层存在输入缓冲区残留换行符 \n 问题
   - Scanf("%d", &num)：只读取匹配数字，不会吃掉末尾回车换行，\n 留在缓冲区
   - Scanln()：读取整行，若缓存已有 \n，会直接读到空字符串，跳过等待用户输入
   故障现象：菜单选数字后进入输入环节，直接跳过输入、返回菜单疯狂提示输入错误选项（本次bug根源）

2. bufio.Scanner 是交互式控制台程序最优方案，解决上述缓冲区污染问题
   优势：
   ① 按完整一行读取输入，自动丢弃末尾换行符，无缓存残留
   ② 支持带空格、中文的长文本（昵称、聊天消息等）
   ③ 统一一套输入API，无需混用多个fmt输入函数，逻辑统一
   ④ 可拿到原始输入字符串，自由转换数字/字符串，自定义格式校验

3. 使用场景区分
   ✅ 优先用Scanner：多轮菜单交互、输入带空格内容、连续多次输入（聊天室客户端）
   ⚠️ 仅临时极简单步程序可用Scanf：一次性读取数字后程序直接退出，无后续输入
*/

/*
以 /rename 小明 改名功能为例：分步完整执行流程
启动服务端
运行 main.go → 创建 Server → 调用 server.Start() 开启端口监听，阻塞等待客户端连接。
启动客户端连接服务端
NewClient 使用 net.Dial 建立 TCP 连接，得到 client.conn。
执行 go client.DealResponseMessage() 开启后台协程，持续 conn.Read 等待服务端消息。
菜单选择 3，进入 Rename ()
Scanner 读取你输入的昵称 小明，拼接字符串：/rename 小明\n。
执行 client.conn.Write([]byte(sendSeg))，通过 TCP 网络把字节流发给服务端。
服务端 server.go 接收数据
Start() 里的 Accept 早已捕获你的客户端连接，在 handleConn 循环执行 conn.Read(buf)，读到客户端发送的 /rename 小明。
转字符串后调用 user.DoMessage(msg) 交给业务层。
user.go 处理改名逻辑
DoMessage 判断前缀 /rename，分割取出新昵称，更新 u.Name。
调用 u.conn.Write() 将 改名成功，你的新昵称：小明 发回客户端。
客户端后台协程打印结果
DealResponseMessage 中 conn.Read 读到服务端返回的提示文字，直接打印到控制台，你看到改名成功提示。
一轮结束，回到菜单，可继续操作。
*/
