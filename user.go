package main

import (
	"fmt"
	"net"
	"strings"
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

// Online 提示用户上线，并加入到OnlineMap
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

// Offline 提示用户已下线，并从OnlineMap中移除
func (u *User) Offline() {
	//将用户相关信息从OnlineMap中移除
	u.server.mapLock.Lock()
	//u.server.OnlineMap[u.Name] = u
	delete(u.server.OnlineMap, u.Name)
	u.server.mapLock.Unlock()
	//fmt.Println("连接建立成功")
	u.server.BroadCast(u, Red+"已下线"+Reset)
}

// SendMsg 给当前user对应的客户端发送消息
func (u *User) SendMsg(msg string) {
	_, err := u.conn.Write([]byte(msg + "\n"))
	if err != nil {
		fmt.Println("Write Error", err)
	}
}

func (u *User) HandleMessage(msg string) {
	if msg == "who" {
		u.server.mapLock.Lock()
		for _, user := range u.server.OnlineMap {
			onlineMsg := "[" + user.Addr + "]" + user.Name + ":" + Green + "Online" + Reset
			u.SendMsg(onlineMsg)
		}
		u.server.mapLock.Unlock()
	} else if len(msg) > 7 && msg[:7] == "rename|" {
		newName := strings.Split(msg, "|")[1]
		_, exist := u.server.OnlineMap[newName]
		if exist {
			u.SendMsg("该用户名已存在")
		} else {
			u.server.mapLock.Lock()
			//更改原有名字说对应的索引，并不影响储存结构和位置，只是更改了索引关系
			delete(u.server.OnlineMap, u.Name)
			u.server.OnlineMap[newName] = u
			u.server.mapLock.Unlock()
			u.Name = newName
			u.SendMsg("新用户名:" + newName)
		}
	} else {
		u.server.BroadCast(u, msg)
	}
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
