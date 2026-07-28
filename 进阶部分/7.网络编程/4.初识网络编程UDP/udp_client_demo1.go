// 创建时间：2026/7/26 下午9:12
// udp_client_demo1.go UDP客户端（DialUDP绑定单一服务端写法）
package main

import (
	"fmt"
	"net"
	"time"
)

/*
=======================================================================
UDP客户端 vs TCP客户端 核心差异前置说明（对照TCP客户端注释）
=======================================================================
 1. 连接本质区别
    TCP客户端 net.Dial("tcp", addr)：执行三次握手，建立真实长连接Conn
    UDP客户端 net.DialUDP()：无握手！仅本地记录目标服务地址，逻辑上绑定，底层仍是无连接报文
 2. 读写API差异
    TCP：conn.Write() / conn.Read()，连接绑定对方地址，无需额外传地址
    UDP Dial模式：clientConn.Write()/Read()，封装后写法和TCP一致；
    原生UDP模式：必须 WriteToUDP(数据, 目标地址)，每次携带地址
 3. 下线感知差异
    TCP：服务端关闭，Read直接返回 io.EOF，能感知对方下线
    UDP：无断开信号，丢包/服务下线Read不会报错，无法自动感知离线
 4. 并发读写缺陷（同TCP单线程客户端短板）
    当前代码单线程串行：先发消息→阻塞等回包，无法同时接收服务端主动推送消息
    优化方案：开goroutine单独循环Read，实现读写分离全双工

=======================================================================
*/
func main() {
	// 1. ResolveUDPAddr：将字符串地址转为UDP地址结构体
	// 对标TCP：net.ResolveTCPAddr，仅协议参数改为udp
	//解析TCP地址：
	serverAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:9090")
	if err != nil {
		fmt.Println("解析服务端地址失败：", err)
		return
	}

	// 2. DialUDP 创建UDP客户端对象,是客户端核心函数
	// 参数说明：
	// 第1参数：协议udp
	// 第2参数：本地地址nil → 系统自动随机分配本机端口
	// 第3参数：远端服务地址
	// 对比TCP net.Dial：TCP会握手建连接，UDP仅记录目标地址，无网络交互
	clientConn, err := net.DialUDP("udp", nil, serverAddr)
	if err != nil {
		fmt.Println("创建UDP客户端连接失败：", err)
		return
	}
	// defer延迟关闭，程序退出释放端口资源，和TCP客户端defer conn.Close逻辑一致
	defer func() {
		_ = clientConn.Close()
		fmt.Println("\nUDP客户端通道已关闭")
	}()

	fmt.Println("UDP客户端启动，目标服务：127.0.0.1:9090")
	// 接收缓冲区：存放服务端返回报文，单次最大读取1024字节，和TCP buf用法完全相同
	buf := make([]byte, 1024)

	// 循环发送10条测试消息
	for i := 1; i <= 10; i++ {
		sendMsg := fmt.Sprintf("第%d条UDP测试消息", i)
		// 3. Write发送数据，DialUDP绑定后不用传目标地址
		// 写法和TCP conn.Write([]byte)完全一致，简化编码
		_, err = clientConn.Write([]byte(sendMsg)) //发送给服务端
		if err != nil {
			fmt.Println("发送消息失败：", err)
			continue
		}
		fmt.Printf("客户端已发送：%s\n", sendMsg)

		// 4. Read阻塞读取服务端回显报文
		// 阻塞特性同TCP Read：无数据时程序卡在这一行等待
		// 区别：TCP对方关闭会返回io.EOF，UDP永远不会返回EOF
		n, err := clientConn.Read(buf) //读取服务端返回的数据
		if err != nil {
			fmt.Println("接收回显失败：", err)
			continue
		}
		// buf[:n]截取有效字节，避免缓冲区空零值乱码，TCP/UDP通用写法
		fmt.Printf("收到服务端回包：%s\n\n", string(buf[:n]))

		// 间隔1秒发送下一条
		time.Sleep(1 * time.Second)
	}

	fmt.Println("全部消息发送完毕，客户端退出")
}

/*
客户端、服务端都要用 ResolveUDPAddr，但用途完全不一样
`ResolveUDPAddr` 只是一个地址格式转换工具：
接收字符串 `"ip:port"`，输出程序能识别的 `*net.UDPAddr` 结构体。
- 服务端用它生成「自己要监听的地址」；
- 客户端用它生成「要连接的服务端地址」；
两者只是使用场景不同，函数本身没有区别。
1. 服务端调用 ResolveUDPAddr
作用：
把字符串地址转成 `*net.UDPAddr`，**用来绑定本机端口 9090**，对外提供服务，监听所有客户端报文。
相当于：告诉操作系统，我要占本机 `127.0.0.1:9090` 这个端口收消息。
2. UDP 客户端 DialUDP 里的 ResolveUDPAddr
作用：
把服务端的 IP 端口字符串转成 `*net.UDPAddr`，**标记要通信的远端目标地址**。
第二个参数 `nil` 代表客户端本机端口由系统随机分配，不需要手动绑定固定端口。
*/

/*
TCP和UDP有关服务端和客户端的区别
1. 服务端：
- TCP：net.ResolveTCPAddr解析地址+net.ListenTCP监听
- UDP：net.ResolveUDPAddr解析地址+net.ListenUDP监听
2. 客户端：
- TCP：net.Dial()客户端核心函数，与服务端建立 TCP 连接
- UDP：net.ResolveUDPAddr解析地址+net.DialUDP客户端核心函数，创建UDP客户端对象
对于为何TCP客户端无net.ResolveUDPAddr，因为已经封装到net.Dial()函数中了
*/

/*
ListenUDP / DialUDP 才是 UDP 真正核心（创建操作系统 Socket）
1.服务端：net.ListenUDP ("udp", listenAddr)
入参：解析好的本机 UDP 地址 listenAddr
底层行为：
1. 向 OS 申请一个 UDP socket；
2. 将 socket 绑定到本机指定 IP + 端口（9090）；
3. 返回 `*UDPConn`，持有这个 socket，用来持续接收全网客户端报文。
有了这个对象，才能调用 ReadFromUDP 收消息、WriteToUDP 回复消息。

2.客户端：net.DialUDP ("udp", nil, serverAddr)
参数解释：
- 第二个参数 nil：本地地址，nil 代表让系统随机分配本机临时端口；
- 第三个参数 serverAddr：解析好的远端服务地址。
底层行为：
1. 申请 UDP socket；
2. 本地随机端口绑定；
3. **仅在本地内存记录远端地址，没有网络握手报文**（UDP 无连接）；
4. 返回 `*UDPConn`，后续调用 Write () 时自动带上预先存好的远端地址发送。
*/

/*
客户端+服务端完整流程：
UDP 服务端流程
1. ResolveUDPAddr：字符串地址转结构体（纯转换）
2. ListenUDP：创建 UDP 套接字 + 绑定端口（核心，占用端口资源）
3. 循环 ReadFromUDP / WriteToUDP：收发报文

UDP 客户端 Dial 模式流程
1. ResolveUDPAddr：解析服务端地址为结构体（纯转换）
2. DialUDP：创建客户端套接字、缓存远端地址（核心）
3. 循环 Write / Read：和服务端通信
*/
