// 创建时间：2026/7/27 下午6:18
// udp_heartbeat_server.go UDP心跳服务端
package main

import (
	"fmt"
	"net"
	"sync"
	"time"
)

/*
=======================================================================
知识点前置说明 & TCP对比
1. UDP天然缺陷：无连接，客户端下线/断网不会发送任何关闭信号，服务端无法主动感知
   TCP：客户端正常断开 Read() 返回io.EOF，可以直接删除连接；
   UDP：只能依靠【应用层心跳】判断设备在线状态。
2. 并发风险：
   主线程协程：接收报文 → 修改map
   清理协程：定时遍历、删除map
   Go内置map不支持多goroutine同时读写，不加锁直接程序崩溃
   解决方案：sync.Mutex 互斥锁，同一时间只允许一个协程操作map
=======================================================================
*/

// 新：在线设备表：key=客户端地址，value=最后一次收到心跳的时间
/*
为什么需要这两个？结合 TCP 对比:
TCP 天然区分客户端：每个客户端对应独立 conn，断开直接丢弃连接；
UDP 共用一个 udpConn，没有连接标识，必须手动用 map 保存「哪个 IP 端口在线」；
并发冲突问题：
主线程：收到消息 → 往 map 写入更新时间（写 map）
后台清理协程：每秒遍历、删除离线客户端（读 + 删 map）
Go 原生 map不允许多 goroutine 同时读写，不加锁直接程序崩溃 fatal error: concurrent map writes
sync.Mutex 互斥锁：同一时间只允许一个协程操作 map。
*/
var onlineDevice = make(map[*net.UDPAddr]time.Time)
var lock sync.Mutex // 保护map操作的互斥锁

// 后台清理协程：定时清理离线设备函数
func cleanOfflineDevice() {
	ticker := time.NewTicker(1 * time.Second)
	// 每1秒扫描一次.创建周期性定时器打点器，每隔 1 秒自动向 ticker.C 通道发送当前时间，实现定时循环逻辑。
	defer ticker.Stop()

	/*
		ticker.C 是一个只读 channel，每到 1 秒周期，通道会塞入当前时间；
		for range 会持续阻塞等待通道数据，收到一次就完整执行一轮大括号内的清理逻辑，无限循环。
		流程：
		程序启动，卡在 for range ticker.C 等待 1 秒；1 秒到，通道收到时间，进入循环体执行清理判断；
		清理代码跑完，再次阻塞等待下一个 1 秒；往复循环，实现每秒自动巡检在线设备 map。
	*/
	for range ticker.C {
		lock.Lock() // 加锁，独占map
		now := time.Now()
		// 遍历所有在线设备
		for addr, lastPingTime := range onlineDevice {
			// 当前时间 - 最后心跳 >5秒 → 判断离线
			if now.Sub(lastPingTime) > 5*time.Second {
				//Go 内置函数，删除 map 中指定 key，清理离线客户端记录，减少内存占用，避免 map 无限膨胀。
				delete(onlineDevice, addr)
				fmt.Printf("设备离线：%v，已清理\n", addr)
			}
		}
		lock.Unlock() // 解锁，释放map给其他协程使用
		/*
			为什么必须加锁？
			Go 原生 map非并发安全，当前程序存在两个 goroutine 同时操作 onlineDevice 共享 map：
			主线程 goroutine：收到客户端消息 → onlineDevice[remoteAddr] = time.Now()（写 map）
			本清理 goroutine：遍历 map + delete 删除 key（读 + 写 map）
			如果不加锁，两个协程同时读写 map，程序直接崩溃报错
		*/
	}
}

func main() {
	// 1. 地址转换工具：字符串地址转为UDP结构体（仅格式转换，无网络操作）
	listenAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:9090")
	if err != nil {
		fmt.Println("解析地址失败:", err)
		return
	}

	// 2. 核心函数：创建UDP套接字、绑定9090端口
	udpConn, err := net.ListenUDP("udp", listenAddr)
	if err != nil {
		fmt.Println("启动UDP服务失败:", err)
		return
	}
	defer udpConn.Close()
	fmt.Println("UDP心跳服务启动，监听 127.0.0.1:9090")

	/*
		新：启动独立协程，用于后台持续清理离线设备
		作用：
		主线程专心收客户端报文，单独开一条 goroutine 后台定时扫描在线列表，两个逻辑并行运行，互不阻塞。
	*/
	go cleanOfflineDevice()

	buf := make([]byte, 1024)
	for {
		// 阻塞等待客户端报文
		n, remoteAddr, err := udpConn.ReadFromUDP(buf)
		if err != nil {
			fmt.Println("接收数据异常:", err)
			continue
		}
		recvMsg := string(buf[:n])
		fmt.Printf("收到设备【%v】消息：%s\n", remoteAddr, recvMsg)

		/*
			新：加锁：更新客户端最后心跳时间
			逻辑：
			只要收到任意客户端发来的报文（ping 心跳），就刷新它在 map 里的最后在线时间，定时器扫描时不会判定离线。
		*/
		lock.Lock()
		onlineDevice[remoteAddr] = time.Now() //  新增、更新客户端最后心跳时间
		lock.Unlock()

		// 返回ack应答给客户端
		ackMsg := []byte("ack")
		_, err = udpConn.WriteToUDP(ackMsg, remoteAddr)
		if err != nil {
			fmt.Printf("向 %v 返回ack失败：%v\n", remoteAddr, err)
		}
	}
}

/*
【基础 UDP 回显服务流程】
解析地址、创建 udpConn 绑定端口
单一层 for 循环阻塞等待客户端报文
收到消息 → 打印 → 原样消息返回客户端
全程不记录任何客户端状态，客户端退出无感知
【UDP 心跳服务完整流程】
全局创建 map 存在线设备、互斥锁保护并发
解析地址、创建 udpConn 绑定端口（和基础版一致）
go cleanOfflineDevice() 启动后台定时清理协程
主线程循环：
ReadFromUDP 接收客户端消息
打印消息
加锁，更新该客户端最后心跳时间
回复固定 ack 应答
后台每秒执行一次扫描：
遍历 map，5 秒无心跳的客户端直接删除，打印离线日志
*/

/*
====================================================================================
Ticker 完整知识点注释：原理、作用、对比time.Sleep、资源释放、time.Tick坑点
1. Ticker底层原理
time.NewTicker(d) 创建周期性打点器，底层会注册系统全局定时器堆，内部维护一个缓冲容量1的只读通道ticker.C；
后台runtime定时协程会每隔固定间隔d，自动向通道写入当前时间time.Time；
for range ticker.C 会阻塞等待通道信号，收到一次信号就完整执行一轮循环内清理逻辑，永久循环触发周期任务。

2. Ticker核心作用
专门用于长期后台周期性任务（本场景：每秒扫描在线设备、清理离线UDP客户端）；
一次创建，持续复用，适合服务常驻后台巡检、定时上报、定时清理等场景。

3. 为什么不用 time.Sleep 循环（不推荐写法）
错误示例：
for {
	// 离线清理逻辑
	time.Sleep(1 * time.Second)
}
缺陷：误差叠加漂移问题
Sleep模式的周期 = 业务执行耗时 + 固定1秒休眠；
如果清理遍历map耗时200ms，本轮总时长1.2s，下一轮又叠加耗时，长期运行定时间隔越来越长，精度持续偏移；
Ticker是基准时间驱动，严格按照固定1秒间隔触发，不受循环内代码执行耗时影响，定时精度稳定，无漂移误差。

4. defer ticker.Stop() 知识点
ticker.Stop() 作用：关闭底层系统定时器，停止向通道发送定时信号，释放runtime定时器资源；
若不调用Stop，ticker会持续占用系统调度资源，函数退出后定时器后台仍运行，造成内存/调度资源泄漏；
defer保证：无论函数正常执行完毕、异常退出、goroutine终止，一定会执行Stop，杜绝资源遗漏释放。

5. time.Tick() 简写严重坑点
新手简写写法 for range time.Tick(1*time.Second)，底层等价NewTicker，但无返回ticker句柄，无法手动Stop；
长期后台服务禁止使用，函数退出后定时器无法回收，永久资源泄漏；
生产标准规范：统一使用NewTicker + defer Stop 成对写法。
====================================================================================
*/
