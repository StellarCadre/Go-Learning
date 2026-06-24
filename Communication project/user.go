// 创建时间：2026/6/23 下午9:53
package main

import (
	"net"
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

func (User *User) Online() {
	User.server.mapLock.Lock()
	User.server.OnlineMap[User.Name] = User
	User.server.mapLock.Unlock()
	User.server.Broadcast(User, "This user is online")
}

func (User *User) Offline() {
	User.server.mapLock.Lock()
	delete(User.server.OnlineMap, User.Name)
	User.server.mapLock.Unlock()
	User.server.Broadcast(User, "This user is offline")
}

// 客户端业务处理逻辑
func (User *User) DoMessage(msg string) {
	User.server.Broadcast(User, msg)
}

// 监听当前每个用户各自的User channel，一旦服务器端广播了消息到channel，就发送给用户
func (u *User) ListenMessage() {
	for { //持续监听
		msg := <-u.C                     //把得到的信息放到msg变量中保存
		u.conn.Write([]byte(msg + "\n")) //发送给用户
	}
}
