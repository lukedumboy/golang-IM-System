package main

import (
	"fmt"
	"io"
	"net"
	"sync"
	"time"
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

// BroadcastLoop 监听广播消息，并转发给所有在线用户
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
	theMessage := "[" + user.Addr + "]" + user.Name + ":" + msg + "\n"
	s.Message <- theMessage
}

// Handler 单独将已经建立的连接拿出来，作为一条单独的线程维护
func (s *Server) Handler(conn net.Conn) {
	//新建一名用户
	user := NewUser(conn, s)
	//将用户相关信息加入到OnlineMap中
	user.Online()

	//监听用户是否活跃的channel
	isAlive := make(chan bool)

	go func() {
		buff := make([]byte, 4096)
		for {
			//从conn中读取数据到buff中
			n, err := conn.Read(buff)
			if n == 0 {
				user.Offline()
				return
			}
			if err != nil && err != io.EOF {
				fmt.Println("Read error:", err)
				return
			}
			msg := string(buff[:n])
			//用户处理消息
			user.HandleMessage(msg)
			//向channel isAlive发送信息，代表这是一个活跃用户
			isAlive <- true
		}
	}()

	//阻塞当前Handler
	//select {}
	for {
		select {
		case <-isAlive:
			//当前用户为活跃用户，重置定时器
		case <-time.After(time.Second * 30):
			//超出时间，下线处理
			//强制关闭当前user
			//???第一个观察点，就算SendMsg注释之后仍然报错，意味着与SendMsg无关
			user.SendMsg("你已因超出时间下线")
			s.mapLock.Lock()
			delete(s.OnlineMap, user.Name)
			s.mapLock.Unlock()
			close(user.c)

			_ = conn.Close()
			return
		}
	}
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
