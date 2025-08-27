package main

import (
	"fmt"
	"io"
	"net"
	"sync"
)

// Server 定义服务器的结构体，具有什么信息，是模版
type Server struct {
	IP   string
	Port int

	OnlineMap map[string]*User
	mapLock   sync.RWMutex
	Message   chan string
}

// NewServer 创建新服务器，将不同的信息填入这个模版
func NewServer(ip string, port int) *Server {
	return &Server{
		IP:        ip,
		Port:      port,
		OnlineMap: make(map[string]*User),
		Message:   make(chan string),
	}
}

// BroadcastLoop 监听广播消息，并转发给所有
func (s *Server) BroadcastLoop() {
	for {
		msg := <-s.Message
		s.mapLock.Lock()
		for _, user := range s.OnlineMap {
			user.c <- msg
		}
		s.mapLock.Unlock()
	}
}

// BroadCast  广播消息
func (s *Server) BroadCast(user *User, msg string) {
	theMessage := "[" + user.Addr + "]" + user.Name + ":" + msg
	s.Message <- theMessage
}

// Handler 单独将已经建立的连接拿出来，作为一条单独的线程维护
func (s *Server) Handler(conn net.Conn) {
	//新建一名用户
	user := NewUser(conn)
	//将用户相关信息加入到OnlineMap中
	s.mapLock.Lock()
	s.OnlineMap[user.Name] = user
	s.mapLock.Unlock()
	//fmt.Println("连接建立成功")
	s.BroadCast(user, "已上线\n")

	go func() {
		buff := make([]byte, 4096)
		for {
			n, err := conn.Read(buff)
			if n == 0 {
				s.BroadCast(user, "已下线\n")
				return
			}
			if err != nil && err != io.EOF {
				fmt.Println("Read error:", err)
				return
			}
			// ???去除末尾的\n，首先为什么要去除，再就是为什么读到的信息里面会有\n
			msg := string(buff[:n])
			//s.Message <- msg
			s.BroadCast(user, msg)
		}
	}()

	//阻塞当前Handler
	select {}
}

// Start 服务器的启动操作，先规定使用什么协议，再确定特定的端口监听进入的连接
func (s *Server) Start() {
	//socket listen操作
	Listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.IP, s.Port))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func() {
		_ = Listener.Close()
	}() //()为调用函数，类似于 server.Start()

	go s.BroadcastLoop()

	for {
		conn, err := Listener.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		go s.Handler(conn)
	}
}
