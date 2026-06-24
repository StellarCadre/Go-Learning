// 创建时间：2026/6/23 下午8:53
package main

import (
	"fmt"
	"io"
	"net"
	"sync"
)

// Server 结构体：定义服务器的核心属性
type Server struct {
	IP   string // 服务器监听的IP地址（如127.0.0.1为本机，0.0.0.0为所有网卡）
	Port int    // 服务器监听的端口号（如8000，需避免占用系统端口）

	OnlineMap map[string]*User //在线用户列表.
	mapLock   sync.RWMutex     //
	Message   chan string      //服务器端广播信息给各用户时使用的channel。

}

// NewServer 构造函数：统一创建、初始化 Server 对象。创建并返回一个Server实例
// 参数：ip（监听IP）、port（监听端口）
// 返回值：指向Server结构体的指针（避免值拷贝，提升性能）
func NewServer(ip string, port int) *Server {
	server := &Server{
		IP:   ip,
		Port: port,

		OnlineMap: make(map[string]*User),
		Message:   make(chan string),
	}
	return server
}

// 将待广播信息发送到服务器的channel的功能
func (s *Server) Broadcast(user *User, message string) { //待发送的对象和待发送的信息
	sendMsg := "[" + user.Address + "]" + user.Name + message //待发送的对象和待发送的信息结合起来
	s.Message <- sendMsg                                      //把待广播信息发到服务器端的channel。借此来实现广播
}

// 服务器端channel给每个在线用户的channel发送的功能
func (s *Server) ListenMessager() {
	for {
		msg := <-s.Message
		s.mapLock.Lock()
		for _, user := range s.OnlineMap {
			user.C <- msg
		}
		s.mapLock.Unlock()
	}
}

// handleConn 连接处理函数：处理单个客户端的连接逻辑
// 接收参数：conn（net.Conn类型，代表客户端与服务器的TCP连接）
// 注意：该方法绑定到Server结构体，可通过s访问服务器的IP/Port等属性
func (s *Server) handleConn(conn net.Conn) {
	//【扩展点】后续可在这里实现：读取客户端发送的数据、向客户端写数据、处理业务逻辑等

	// 打印连接成功的提示，conn.RemoteAddr()可获取客户端的IP+端口
	fmt.Printf("链接建立成功，客户端地址：%s\n", conn.RemoteAddr().String())
	//将上线用户的信息添加到OnlineMap中
	user := NewUser(conn, s)
	//用户业务封装环节，代码替换
	//s.mapLock.Lock()
	//s.OnlineMap[user.Name]=user
	//s.mapLock.Unlock()
	//给所有上线用户广播消息。建议将广播操作封装成一个函数，直接调用，保持代码清晰。
	//s.Broadcast(user,"This user is online")
	//改为：
	user.Online()

	/*
		什么时候要用到go func()：
		   不写 go 的问题（同步阻塞）：代码是从上到下顺序执行，遇到阻塞函数，整个当前函数直接卡在原地，后面代码永远走不到。
		   go func () 作用：把阻塞代码丢到后台并行跑。
		   结合项目分析：
		   场景 1：每个用户独立读消息协程：既要持续收客户端消息，又不能阻塞当前连接逻辑。
		   场景 2：服务全局广播分发协程 go s.ListenMessager()。ListenMessager 内部是 for { msg := <-s.Message }，通道读取会永久阻塞。
		          如果同步调用 s.ListenMessager()，服务直接卡在这，无法执行 Accept 接收新客户端。所以必须 go 丢后台。
		   场景 3：每个用户的消息下发协程 go user.ListenMessage()，ListenMessage 里 msg := <-u.C 阻塞等待广播消息。
		          同步调用会卡住 NewUser 创建流程，必须开协程后台监听自己的私有通道。
		   场景 4：服务接收客户端 go s.handleConn(conn)。handleConn 内部有阻塞的 Read 循环，如果同步执行，Accept 循环会卡死，只能处理一个客户端。
		          开协程才能同时处理成千上万客户端。
	*/

	//接收客户端信息，并把该信息广播给所有在线用户的协程函数
	go func() { // 为当前用户单独启动读取消息协程
		for { // 无限循环，持续等待客户端发数据
			msg := make([]byte, 4096)     // 1. 开辟4096字节缓冲区，用来存放客户端发来的二进制数据
			n, err := user.conn.Read(msg) // 2. 阻塞读取TCP连接数据：没有数据就卡在这等待，不占用CPU
			if n == 0 {                   // 分支1：n==0，代表客户端正常关闭连接（主动关掉CMD/终端）
				//用户业务封装环节，代码替换
				//s.Broadcast(user,"This user is offline")  // 推送下线广播，告知全体用户该用户离线
				//改为：
				user.Offline()

				return
			}
			if err != nil && err != io.EOF { // 分支2：读取出现非EOF异常，网络中断、强制断连、网络报错
				fmt.Printf("read msg failed,err:%v\n", err)
				return
			}
			message := string(msg[:n]) // 3. 截取有效数据：msg缓冲区是4096字节，只用前n个收到的字节
			//用户业务封装环节，代码替换
			//s.Broadcast(user,message)  // 4. 调用全局广播，把这条用户输入的消息推送给所有在线客户端
			//改为：
			user.DoMessage(message)
		}
	}()
	//当前hanleConn阻塞，防止退出
	select {}
}

// Start 启动服务器方法：服务器的核心运行逻辑
func (s *Server) Start() {
	// 1. 监听TCP端口：net.Listen("tcp", 地址字符串)
	// 地址格式：IP:Port（如127.0.0.1:8000），通过fmt.Sprintf拼接
	// 返回值：listener（监听器，用于接收客户端连接）、err（错误信息，非nil表示监听失败）
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.IP, s.Port)) //Listener 只干一件事：蹲在端口等客户上门
	if err != nil {                                                        // 监听失败的错误处理
		fmt.Printf("listen failed,err:%v\n", err)
		return // 监听失败则退出Start方法
	}
	// 2. 延迟关闭监听器：defer确保函数结束时（即使报错）执行listener.Close()
	// 避免端口被长期占用，释放系统资源
	defer listener.Close()

	go s.ListenMessager()

	// 3. 无限循环：持续等待并接收客户端连接（服务器常驻逻辑）
	for {
		// 3.1 接收客户端连接：listener.Accept()会阻塞，直到有客户端连接进来
		// 返回值：conn（客户端连接实例）、err（错误信息，非nil表示接收连接失败）
		conn, err := listener.Accept() //conn = 客户端与服务器之间的专属双向通道
		if err != nil {                // 接收连接失败的错误处理
			fmt.Printf("accept failed,err:%v\n", err)
			continue // 跳过当前循环，继续等待下一个连接
		}
		// 3.2 处理客户端连接：启动协程（go关键字）执行handleConn
		// 为什么用协程？避免单个连接阻塞整个服务器（如客户端不发数据，主线程会卡住）
		// 协程是Go轻量级线程，可同时处理大量连接
		go s.handleConn(conn)
	}
}

/*
服务器端和客户端的通信的流程：
运行服务端代码 → 端口 8000 开始监听，卡住等待；
启动客户端（telnet/nc/ 自己写的 Go 客户端），填写服务端地址端口发起连接；
操作系统完成 TCP 握手，服务端 Accept() 返回 conn，开启协程处理；
两端各自持有 conn，双向收发数据，完成通信；
任意一端关闭连接，通道销毁，通信结束。
*/

/*
服务器端和客户端的通信的前提：
服务器这边需要自己配置监听的 IP 地址和端口，IP 用来指定本机哪个网卡接收连接，
端口用来区分同一台机器上不同的服务程序，配置完成后程序启动就会占用这个 IP 和端口持续等待客户端连接，
另外服务器还需要提前写好客户端连进来之后的数据处理逻辑，保证连接建立后能正常收发数据；
用户也就是客户端这边必须提前知道服务器配置好的 IP 地址（公网ip）或者域名、以及对应的端口号，
只有拿到这两个信息才能主动发起 TCP 连接去和服务器建立通信通道，客户端自身不需要手动配置端口，
操作系统会自动分配随机临时端口，不需要服务器提前知晓，同时客户端不用关心服务器内部怎么处理连接，
只需要按照双方约定的规则发送、接收数据即可。


*/
