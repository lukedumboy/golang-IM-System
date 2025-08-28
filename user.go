package main

import (
	"fmt"
	"net"
)

const Green = "\033[32m" // 前景色：绿色
const Red = "\033[31m"   // 前景色：红色
const Reset = "\033[0m"  // 重置：恢复默认颜色

type User struct {
	Name string
	Addr string
	c    chan string
	conn net.Conn
	//将server与user进行关联，可以通过user找到server
	server *Server
}

func (u *User) Online() {
	////新建一名用户
	//user := NewUser(conn)
	//将用户相关信息加入到OnlineMap中
	u.server.mapLock.Lock()
	u.server.OnlineMap[u.Name] = u
	u.server.mapLock.Unlock()
	//fmt.Println("连接建立成功")
	u.server.BroadCast(u, Green+"已上线"+Reset)
}

func (u *User) Offline() {
	//将用户相关信息从OnlineMap中移除
	u.server.mapLock.Lock()
	//u.server.OnlineMap[u.Name] = u
	delete(u.server.OnlineMap, u.Name)
	u.server.mapLock.Unlock()
	//fmt.Println("连接建立成功")
	u.server.BroadCast(u, Red+"已下线"+Reset)
}

func (u *User) HandleMessage(msg string) {
	u.server.BroadCast(u, msg)
}

func NewUser(conn net.Conn, s *Server) *User {
	userAddr := conn.RemoteAddr().String()
	user := &User{
		Name:   userAddr,
		Addr:   userAddr,
		c:      make(chan string),
		conn:   conn,
		server: s,
	}
	go user.ListenMessage()
	return user
}

func (u *User) ListenMessage() {
	for {
		message := <-u.c
		_, err := u.conn.Write([]byte(message))
		if err != nil {
			fmt.Println("Write Error:", err)
			return
		}
	}
}
