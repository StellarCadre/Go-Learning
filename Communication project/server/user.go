// 创建时间：2026/6/23 下午9:53
package main

/*
user.go 是单个客户端的业务封装层。


*/
import (
	"net"
	"strings"
)

type User struct {
	Name    string
	Address string
	C       chan string //用于服务器端给当前用户发送信息的中转站。每人单独一个。每个channel要包装成一个协程goroutine。
	conn    net.Conn    //表示客户端连接实例，conn = 客户端与服务器之间的专属双向通道。

	server *Server //当前用户所属的服务器
}

// NewUser 创建User实例
func NewUser(conn net.Conn, s *Server) *User {
	userAddr := conn.RemoteAddr().String() //获取客户端地址

	user := &User{
		Name:    userAddr, //地址作为用户名
		Address: userAddr, //地址作为地址
		C:       make(chan string),
		conn:    conn,

		server: s, //当前用户所属的服务器
	}
	//每创建一个用户，就会自动监听对应的channel的goroutine。各自用户不干扰，因为各用户分属不同协程。
	go user.ListenMessage()
	return user
}

func (User *User) Online() { //完成用户上线全套业务，对外只暴露一行调用，隐藏锁、map 操作、广播
	User.server.mapLock.Lock()
	User.server.OnlineMap[User.Name] = User
	User.server.mapLock.Unlock()
	User.server.Broadcast(User, "zaixian")
}

func (User *User) Offline() { //客户端断开时清理资源、下线通知
	User.server.mapLock.Lock()
	delete(User.server.OnlineMap, User.Name)
	User.server.mapLock.Unlock()
	User.server.Broadcast(User, "lixian")
}

// 客户端业务处理逻辑，处理客户端发给服务器的上行数据
func (u *User) DoMessage(msg string) {
	// 1. 在线用户查询：统一用 /list 标准指令
	trimMsg := strings.TrimSpace(msg) // 剔除回车、空格、换行，纯净匹配指令
	if trimMsg == "/list" {
		u.server.mapLock.RLock()       // 加读锁，并发安全读取在线用户map，只读场景用读锁提升并发性能
		var onlineList strings.Builder // 使用strings.Builder字符串拼接，相比+拼接效率更高，适合多行文本组装
		onlineList.WriteString("=====当前在线用户=====\r\n")
		for _, user := range u.server.OnlineMap { // 遍历服务端全部在线用户，拼接每个用户地址+昵称在线状态
			onlineList.WriteString("[" + user.Address + "] " + user.Name + " 在线\r\n")
		}
		u.server.mapLock.RUnlock()

		u.conn.Write([]byte(onlineList.String())) // 仅将在线列表单发至发起查询的当前客户端，不全局广播
		return
	}
	// 2. 改名指令：/rename 新名字，空格分割，修复越界、锁、判空
	if strings.HasPrefix(trimMsg, "/rename ") {
		newName := strings.TrimSpace(strings.TrimPrefix(trimMsg, "/rename ")) // 截取/rename前缀后的内容，并再次去空格，提取纯净新用户名
		if newName == "" {                                                    // 校验：新用户名不能为空
			u.conn.Write([]byte("改名失败：用户名不能为空\r\n"))
			return
		}
		// 读锁先判断名字占用
		u.server.mapLock.RLock()
		_, exist := u.server.OnlineMap[newName]
		u.server.mapLock.RUnlock()
		if exist {
			u.conn.Write([]byte("改名失败：用户名已被占用\r\n"))
			return
		}
		// 写锁修改在线map键值
		u.server.mapLock.Lock()            // 加写锁，独占修改在线用户map（增删map必须写锁，保证并发安全）
		delete(u.server.OnlineMap, u.Name) // 删掉map中原用户名对应的旧键，OnlineMap键为用户名，必须同步删除旧索引
		u.Name = newName                   // 更新当前User结构体内部的昵称字段
		u.server.OnlineMap[u.Name] = u     // 以新用户名为键，重新将当前用户存入在线map，完成键名更替
		u.server.mapLock.Unlock()

		u.conn.Write([]byte("改名成功，新昵称：" + newName + "\r\n")) // 单发成功提示给改名操作者本人
		u.server.Broadcast(u, "用户修改昵称为："+newName)            // 全局广播改名事件，告知所有在线用户该用户昵称变更
		return
	}

	// 3.私聊指令：/to 用户名 消息内容
	//执行主体：*User 类型对象 u，代表当前发消息的客户端用户
	//触发前提：客户端输入内容经过 strings.TrimSpace(msg) 去除首尾空格、换行，存入 trimMsg
	if strings.HasPrefix(trimMsg, "/to ") { // 匹配私聊指令：必须以 /to 空格开头，区分普通聊天
		//按空格分割字符串，最多拆成3段
		parts := strings.SplitN(trimMsg, " ", 3)
		if len(parts) < 3 {
			u.conn.Write([]byte("格式错误！正确格式：/to 用户名 私聊内容\r\n"))
			return
		}
		targetName := parts[1] // 目标私聊对象的用户名
		priContent := parts[2] // 私聊正文内容

		// 读锁保护读取在线用户
		u.server.mapLock.RLock()                            //mapLock.RLock()：加读共享锁，多客户端可同时读在线表，不阻塞其他查询，性能优于写锁
		targetUser, exist := u.server.OnlineMap[targetName] //u.server：当前用户绑定的服务端实例，OnlineMap 是服务端全局在线用户字典：map[用户名]*User
		/*
			用 targetName 作为键查找，得到两个返回值:
			   targetUser：目标用户的 User 结构体指针（包含对方 TCP 连接、专属消息通道 C）
			   exist：布尔值，true = 用户在线，false = 用户离线
		*/
		u.server.mapLock.RUnlock()
		if !exist {
			u.conn.Write([]byte("私聊失败：该用户不在线\r\n"))
			return
		}

		// 组装私聊消息，带发送人标记
		priMsg := "[私聊]" + u.Name + "对你说：" + priContent //u.Name：发送者自身昵称，用来让接收方知道是谁发来的私聊
		select {
		case targetUser.C <- priMsg: //拼接完整提示文案，存入 priMsg，这条消息只会投递到目标用户的专属通道
		default:
		}
		// 给自己回显发送记录
		u.conn.Write([]byte("你私聊【" + targetName + "】：" + priContent + "\r\n"))
		return
	}

	// 4. 普通聊天消息全局广播
	u.server.Broadcast(u, trimMsg)
}

// 处理服务器推给客户端的下行广播数据。监听当前每个用户各自的User channel，一旦服务器端广播了消息到channel，就发送给用户。
func (u *User) ListenMessage() {
	for { //持续监听
		msg := <-u.C                     //把得到的信息放到msg变量中保存
		u.conn.Write([]byte(msg + "\n")) //发送给用户
	}
}
