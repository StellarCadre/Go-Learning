// 创建时间：2026/7/27 下午8:36
// udp_heartbeat_client.go UDP心跳客户端
package main

import (
	"fmt"
	"net"
	"time"
)

/*
=====================================================================
知识点前置说明 & TCP对比
 1. UDP客户端天然缺陷：无连接，无法感知服务端是否宕机
    TCP：服务端关闭，Read直接返回io.EOF
    UDP：服务下线/丢包，Read持续阻塞卡死，必须设置【读超时 SetReadDeadline】
 2. 心跳逻辑：定时向外发送ping，等待服务端ack应答
 3. 超时意义：长时间收不到ack，判定服务不可达

=====================================================================
*/
func main() {
	// 1.地址转换工具：字符串地址转为UDP结构体
	serverAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:9090")
	if err != nil {
		fmt.Println("解析服务端地址失败：", err)
		return
	}

	// 2.DialUDP 创建UDP客户端核心套接字
	clientConn, err := net.DialUDP("udp", nil, serverAddr)
	if err != nil {
		fmt.Println("创建UDP客户端失败：", err)
		return
	}
	defer func() {
		_ = clientConn.Close()
		fmt.Println("\n客户端通道关闭")
	}()
	fmt.Println("UDP心跳客户端启动，目标：127.0.0.1:9090")

	buf := make([]byte, 1024)

	// --------------------------【新增改动1】周期性定时器发送心跳 --------------------------
	// 每2秒发送一次心跳包
	ticker := time.NewTicker(2 * time.Second) //创建周期定时器，每隔 2 秒自动向 ticker.C 通道推送时间信号
	defer ticker.Stop()

	for range ticker.C { //每收到一次定时信号，完整执行一轮for,给服务端发送一个ping
		// 发送心跳 ping
		_, err := clientConn.Write([]byte("ping"))
		if err != nil {
			fmt.Println("心跳发送失败：", err)
			continue
		}
		fmt.Println("✅ 发送心跳：ping")

		// --------------------------【新增改动2】设置读取超时 --------------------------
		// 设置1秒读取时限：超过1秒没有ack返回，直接报错退出Read
		/*
			UDP 无连接,如果服务端关闭、网线断开、报文丢包：Read() 会永久阻塞卡死，整个客户端停滞不动，没有任何报错。
			所以，SetReadDeadline给连接clientConn设置全局读写截止时间，但作用域是下一次 IO 操作：当前时间一秒后无应答返回超时错误，不再阻塞。
		*/
		err = clientConn.SetReadDeadline(time.Now().Add(1 * time.Second))
		if err != nil {
			fmt.Println("设置读超时失败：", err)
			continue
		}

		n, err := clientConn.Read(buf)
		if err != nil {
			// 超时错误代表没有收到ack，大概率丢包/服务离线
			fmt.Println("❌ 等待ack超时，未收到服务应答\n")
			continue
		}
		fmt.Printf("✅ 收到服务端应答：%s\n\n", string(buf[:n]))
	}
}

/*
为什么每次循环都要重新 SetReadDeadline，保证只生效一次 Read:
Go 中SetReadDeadline(t)传入的是绝对时间戳，不是相对时长：
假设当前时间 10:00:00.000，Add(1s) → 截止时间 10:00:01.000；
本轮 Read 在10:00:00.200收到 ack，Read 执行完毕；
这个截止时间10:00:01.000永久失效，不会自动重置；
等到 2 秒后下一轮循环（10:00:02.000），如果不重新设置 deadline，旧截止时间早已过期，下一次 Read 会立刻直接报超时，完全无法等待数据。
所以代码把SetReadDeadline写在每轮 Read 之前，强制刷新全新的绝对截止时间，保证每一次 Read 都拥有独立 1 秒等待窗口。
*/

/*
当前for range ticker.C中存在的不足：
每一轮 ticker 循环，只会执行一次 Write 发 ping，并且只 Read 读取缓冲区里第一条报文。
缓冲区剩下的数据包会留在内核，等到下一轮循环，即2秒后才有机会读。无论服务器发的多频繁。
例：
假设客户端发完 ping 后，服务端瞬间连续发 3 条 ack：
3 条 ack 全部进入操作系统内核 UDP 接收缓冲区排队；
当前 Read 只会读取第 1 条 ack，打印应答；
本轮循环结束，回到 for range ticker.C 阻塞等待下一个 2 秒；
剩下 2 条 ack 存在缓冲区，不会被处理；
2 秒后进入下一轮循环，发新 ping，再 Read 取出第 2 条旧 ack，逻辑错乱。
当前代码只适合一问一答标准心跳场景：客户端发 1 个 ping，服务端回 1 个 ack。
可以考虑在 Read 外层套无限内层循环，读到超时为止：
// 内层循环：持续读，直到缓冲区空、触发超时
    for {
        n, err := clientConn.Read(buf)
        if err != nil {
            break // 超时，无更多数据，退出内层
        }
        fmt.Printf("收到应答：%s\n", string(buf[:n]))
    }
*/
