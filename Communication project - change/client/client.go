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
	"context" // 新增：上下文包，用来控制操作超时、取消
	"errors"  // 新增：错误判断专用包
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"time" // 新增：KeepAlive、超时时间单位
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
	// 初始化客户端结构体，存储服务地址、连接句柄、菜单标识
	client := &Client{
		ServerIp:   serverIp,
		ServerPort: serverPort,
		flag:       999,
	}
	//旧 链接服务器，net.Dial：主动发起TCP连接（客户端核心API），net.Dial 三次握手连接服务端
	//conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", serverIp, serverPort)) //组装地址，并返回一个TCP连接对象
	//if err != nil {
	//	fmt.Println("net.Dial err:", err) //连接失败返回 nil
	//	return nil
	//}
	//client.conn = conn //成功把连接存入 Client 并返回
	//return client

	//新
	// 拼接完整服务地址 格式：IP:端口
	addr := fmt.Sprintf("%s:%d", serverIp, serverPort)
	// 1. 构造Dialer，开启TCP底层保活.
	/*
		  ==================== 知识点：net.Dialer 详解 ====================
			net.Dialer 是拨号配置结构体，控制握手行为、内核保活.原生 net.Dial 内部就是用默认Dialer
			自定义Dialer可以手动控制两大关键能力：
			1. Timeout：整个三次握手拨号的最大等待时长，超时直接失败，不会无限阻塞
			2. KeepAlive：操作系统内核TCP保活探测间隔
			   连接成功后持续生效，内核每隔指定时间自动发探测包，检测对方是否掉线（拔网线/休眠）
			   探测无响应，内核直接标记连接失效，应用层快速收到报错
	*/
	dialer := &net.Dialer{
		Timeout:   5 * time.Second,  // 拨号超时：5秒连不上直接判定失败
		KeepAlive: 30 * time.Second, // 内核每30秒发送一次TCP保活探测包
	}
	// 2. 创建上下文，用于控制拨号超时取消
	/*
		  ==================== 知识点：context.WithTimeout 详解 ====================
			context 上下文：统一管理一段代码生命周期、超时、手动取消
			context.Background()：根上下文，无过期、不可取消，程序入口标准根节点
			WithTimeout(parent, 时长)：基于父上下文创建带自动过期的子上下文
			返回两个值：
				ctx：上下文对象，传给DialContext用于控制拨号超时
				cancel：取消函数，调用即可手动终止当前操作
			defer cancel() 延迟执行：函数无论正常/异常退出，都释放ctx资源，避免泄漏
	*/
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second) //创建一个有 5 秒有效期的控制器，后续其他函数（DialContext）可以接收这个控制器，一旦超时就停止执行操作。
	/*
			context.WithTimeout (父 ctx, 超时时间)做了什么?
			  输入两个参数：parent：父上下文（这里传根 context.Background()）。timeout：超时时长 5 * time.Second
			  返回两个值：ctx：全新子上下文（我们命名 ctx，用来传递超时信号）。cancel：取消函数，类型 func()，调用它就可以手动终止这个 ctx。
			内在逻辑：
		      调用 WithTimeout 后，内部会开启一个隐形定时器：计时 5 秒，时间一到，这个 ctx 自动标记为「已过期 / 取消」。
		      任何使用这个 ctx 的代码，都能感知到超时。两种方式让 ctx 失效：等待 5 秒自动超时（被动）。手动执行 cancel() 函数，立刻失效（主动）
	*/
	defer cancel()

	// 3. DialContext带上下文建立连接
	/*
					  ==================== 知识点：DialContext 详解 ====================
						dialer.DialContext(ctx, 协议, 地址)
						和旧版 net.Dial 核心区别：绑定上下文ctx
						下面场景会立刻终止拨号：
				        1.握手顺利完成，函数返回可用 conn；执行 defer cancel()，销毁 ctx 内部定时器，无资源泄漏。
				        2.目标 IP 无响应，5 秒倒计时结束.ctx 自动触发取消信号，DialContext 立刻中断握手；返回超时错误，进入 PrintNetErr 打印超时提示，函数 return；defer cancel() 执行，清理定时器。
		                3.代码手动提前调用 cancel ()
	*/
	conn, err := dialer.DialContext(ctx, "tcp", addr) //受上下文倒计时管控的 TCP 拨号函数
	if err != nil {
		// 拨号出错，调用统一错误处理函数，区分错误类型友好打印
		PrintNetErr("拨号连接服务端", err)
		return nil
	}
	// 连接创建成功，TCP通道存入客户端结构体，供收发消息使用
	client.conn = conn
	return client
}

// PrintNetErr 区分各类网络错误，标准化输出
/*
==================== 知识点：PrintNetErr 统一错误处理详解 ====================
作用：不再单纯打印原始错误字符串，自动分类网络故障，方便调试与业务区分
用到两个核心错误API：
1. errors.As(err, &目标结构体)：
   作用：提取错误err内部的结构体类型，用于细分网络专属错误,是否是*net.OpError（网络操作错误）
   net.OpError 是标准库专门封装的网络错误结构体，所有 dial/read/write 失败都会包装成它；
   "dial"：建立连接失败（端口关闭、IP 不通）
   "read"：读取客户端消息失败（断网）
   "write"：向客户端发消息失败
   使用前提：提前定义 var opErr *net.OpError 指针变量，传地址给 errors.As
2. errors.Is(err, 标准错误常量)：判断错误链里是否存在指定固定错误值，只匹配常量，不能拿结构
   context.DeadlineExceeded：ctx 倒计时到期超时
   io.EOF：TCP 对端正常关闭连接（用户手动退出，不是故障）
分支逻辑：
- 拨号类错误：端口未开、地址错误
- 上下文超时：操作时间超限
- EOF：正常下线，不是故障
- 其余：未知原始错误直接打印
*/
func PrintNetErr(oper string, err error) { // 参数 oper：当前执行的操作名称（如拨号、读取消息），用于日志区分场景.// 参数 err：网络操作返回的原始错误对象
	// 无错误直接返回，无需处理
	if err == nil {
		return
	}
	var opErr *net.OpError      // 定义网络操作错误变量，接收断言结果
	if errors.As(err, &opErr) { // 判断当前错误是否是标准网络操作错误
		if opErr.Op == "dial" { // opErr.Op 标记当前网络操作类型：dial/ read / write
			fmt.Printf("[网络错误]%s：目标端口未开放/地址不可达\n", oper)
			return
		}
	}
	if errors.Is(err, context.DeadlineExceeded) { // 判断错误链里是否包含「上下文超时」错误
		fmt.Printf("[网络超时]%s操作超时\n", oper)
		return
	}
	if errors.Is(err, io.EOF) { // 判断错误是否为EOF：对方正常关闭连接（主动exit退出）
		fmt.Printf("[正常断开]%s：对方主动关闭连接\n", oper)
		return
	}
	fmt.Printf("[未知网络错误]%s：%v\n", oper, err) // 无法分类的原始错误，直接打印完整信息
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
	//旧buf := make([]byte, 4096)
	//新：创建bufio.Reader封装TCP连接，实现流式按行读取,
	reader := bufio.NewReader(client.conn)
	for {
		//旧n, err := client.conn.Read(buf) //client.conn.Read(buf) 阻塞等待：服务器没发消息时，代码卡在这不动； 一旦服务端下发数据，就把数据读到 buf 缓冲区里。
		//新 按行读取（直到\n结束），自动处理粘包，返回完整的一行消息
		msg, err := reader.ReadString('\n')
		if err != nil {
			PrintNetErr("接收服务端消息", err)
			os.Exit(0)
		}

		//打印完整的一行消息，无脏数据、无粘包
		fmt.Print(msg)
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
===== Context 精简完整知识点 =====
1. 四种创建方式
- Background()：顶层根上下文，无超时、不可取消，程序入口使用
- WithTimeout：带自动倒计时，拨号、查询、单次短请求使用
- WithCancel：无定时，仅手动调用cancel关闭，长连接协程专用
- WithDeadline：指定固定时间点过期，定时任务场景

2. 强制编码规范
调用以上三个派生函数得到cancel，必须defer cancel()，防止定时器、goroutine泄漏

3. Context 内置三个方法
Done()：返回只读channel，取消/超时后通道关闭
Err()：获取终止原因：DeadlineExceeded 超时 / Canceled 手动取消
Deadline()：返回截止时间，无超时第二个返回false

4. 层级连锁特性
父上下文取消，所有子上下文自动一并取消，适合批量关闭连接

5. 标准库统一设计
网络DialContext、HTTP NewRequestWithContext、SQL QueryContext均支持ctx，用于管控阻塞操作超时

6. 使用场景区分
✅ 适合：TCP读写、数据库、HTTP、阻塞循环、外部进程调用
❌ 不适合：简单内存运算、无阻塞瞬时逻辑
*/
