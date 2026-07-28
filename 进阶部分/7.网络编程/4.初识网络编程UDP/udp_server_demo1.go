// 创建时间：2026/7/26 下午5:20
// udp_server_demo1.go 阶段A：最简UDP回显服务
package main

import (
	"fmt"
	"net"
)

func main() {
	// 1. 定义监听地址
	/*
		对标 TCP：net.ResolveTCPAddr("tcp", "127.0.0.1:8080")
		解析UDP地址：协议类型是udp，地址是127.0.0.1:9090（本地回环+9090端口）
		作用完全一样：把字符串 IP 端口转为系统可识别的地址结构体；
	*/
	listenAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:9090")
	if err != nil {
		fmt.Println("解析地址失败:", err)
		return
	}

	// 2.创建UDP监听器
	/*
		TCP：ListenTCP 生成 Listener，循环 Accept 拿到独立 Conn；
		UDP：ListenUDP 直接返回唯一udpConn，没有单独监听器、没有 Accept；
		defer udpConn.Close()：程序退出关闭端口释放资源，和 TCP conn.Close() 资源释放逻辑一致，防止端口占用泄漏。
	*/
	udpConn, err := net.ListenUDP("udp", listenAddr)
	if err != nil {
		fmt.Println("启动UDP服务失败:", err)
		return
	}
	defer udpConn.Close()
	fmt.Println("UDP服务启动，监听 127.0.0.1:9090")

	// 缓冲区：存放接收的数据，最多一次接收1024字节
	/*
		和 TCP 代码 var buf = make([]byte, 1024) 作用一模一样：开辟内存缓冲区，临时存放网络读取的二进制数据，单次最大接收 1024 字节。
	*/
	buf := make([]byte, 1024)

	//ReadFromUDP 阻塞读取报文
	/*
		和 TCP 服务外层for{Accept()}逻辑一致：服务永久运行，持续等待客户端消息。
		TCP 循环核心是接收新连接；UDP 循环核心是接收报文，无需处理连接接入。
		TCP：n, err := conn.Read(buf)，仅返回读取字节数 + 错误，客户端地址绑定在 Conn 内部，不用额外获取；
		UDP ReadFromUDP 返回 3 个值：
		n：本次报文有效字节长度，和 TCP 一致；
		remoteAddr：发送这条消息的客户端完整地址（UDP 无连接，必须手动保存）；
		err：网络读取错误；
		关键点：UDP 不会返回io.EOF，客户端下线、断网不会触发任何错误，代码不会自动跳出循环，这是和 TCP 最大区别（TCP 客户端关闭会返回 EOF）。
	*/
	for {
		// 3. 阻塞等待接收客户端报文
		// n：读到的字节长度；remoteAddr：发送方客户端地址
		n, remoteAddr, err := udpConn.ReadFromUDP(buf)
		if err != nil {
			fmt.Println("接收数据异常:", err)
			continue
		}

		// 截取有效数据，转换字符串
		recvData := string(buf[:n])
		fmt.Printf("收到【%v】消息：%s\n", remoteAddr, recvData)

		// 4. WriteToUDP：向对应客户端原路返回数据
		_, err = udpConn.WriteToUDP(buf[:n], remoteAddr)
		if err != nil {
			fmt.Println("回复消息失败:", err)
		}
	}
}

/*
流程：
解析 UDP 地址，绑定 9090 端口创建 udpConn；
开启无限循环阻塞等待客户端发送报文；
收到消息，同时拿到发送方客户端地址；
打印消息内容；
使用刚才拿到的客户端地址，原路把消息发回；
循环持续等待下一条报文，全程只用一个 udpConn 接待所有客户端。
*/

/*
TCP 与 UDP 完整流程对比：监听器、循环逻辑、为什么 UDP 没有 Accept：
一、底层核心本质先行区分
1.TCP：面向连接
通信前必须建立专属通道（三次握手），通信结束销毁通道（四次挥手）。每个客户端 = 独立专属连接对象。
2.UDP：无连接数据报
没有 “建立通道” 这个步骤，客户端直接发包；服务端只开一个端口，所有客户端消息全部发到这一个端口，每条报文自带发送方地址。

二、监听器对象对比（TCP Listen vs UDPConn）
1. TCP 存在两套对象：listen + Conn，
listen, err := net.ListenTCP("tcp", addr)
conn, err := listen.Accept()
listen：独立监听器，只负责监听端口、接收三次握手、创建新连接；   大门 (listen)  全程只创建一次
conn：单个客户端专属通信通道，绑定唯一客户端 IP + 端口；       客户座位 (conn) for循环迎接多个客户  每来一个客户端，就生成一个全新 conn

2. UDP 只有一套对象：UDPConn
addr, _ := net.ResolveUDPAddr("udp", "127.0.0.1:9090")
udpConn, _ := net.ListenUDP("udp", addr) // 唯一对象，兼顾监听+收发  早餐摊，欢迎所有人来，不需要for来连接每一个客户
UDP 不需要监听每个客户，ListenUDP 直接返回 UDPConn：
这一个对象同时承担两件事：
监听 9090 端口，接收所有客户端报文；
收发所有客户端的数据。
不存在 “为每个客户端新建连接” 的概念，全程只复用这一个 UDPConn。

三、为什么 UDP 完全不需要 Accept ()？
Accept 的作用（TCP 专属）
listener.Accept() 的底层逻辑：阻塞等待客户端发起 TCP 三次握手；握手完成，创建专属 TCP 连接 conn；返回 conn，给程序单独读写这个客户端。
Accept 存在的唯一意义：建立专属连接。
UDP 不需要 Accept 的 3 个根本原因
无握手流程
UDP 没有三次握手，客户端不需要提前登记建立通道，直接发送报文，不存在 “连接建立” 这一步，自然不需要 Accept 来等待握手。
无专属连接对象
TCP 一个客户端对应一个 conn；UDP 所有客户端共用同一个 UDPConn，不需要为新客户端生成新对象。
每条报文自带发送方地址
TCP 的 conn 内部已经绑定客户端地址，读写时不用额外携带；
UDP 每次调用 ReadFromUDP 都会同步返回 remoteAddr，直接拿到客户端地址，不需要提前通过 Accept 绑定。
一句话：Accept 是用来 “建立连接” 的，UDP 根本没有连接，所以不需要这个函数。

四、两者无限 for 循环内部工作对比
1. TCP 服务端双层循环（单线程阻塞版）
外层 for：接待新客户，生成专属连接；
内层 for：和当前已接入的单个客户端持续聊天；
2. UDP 服务端单层循环（仅一层 for）
只有一层循环，只做一件事：不断接收任意客户端的报文。
没有 “新客户端接入” 的步骤，任何客户端随时可以发包进来。
*/

/*
当前缺陷：
ReadFromUDP 阻塞时，下一条报文排队，并发高时处理缓慢；
优化：收到报文后开 goroutine 异步处理，不阻塞下一次 Read。
*/

/*
启动方式：使用Nmap的ncat.exe
ncat -u 127.0.0.1 9090
*/
