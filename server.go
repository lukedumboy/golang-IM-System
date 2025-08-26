package main

import (
	"fmt"
	"net"
)

// Server 定义服务器的结构体，具有什么信息，是模版
type Server struct {
	IP   string
	Port int
}

// NewServer 创建新服务器，将不同的信息填入这个模版
func NewServer(ip string, port int) *Server {
	return &Server{IP: ip, Port: port}
}

// Handler 单独将已经建立的连接拿出来，作为一条单独的线程维护
func (s *Server) Handler(conn net.Conn) {
	fmt.Println("连接建立成功")
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

	for {
		conn, err := Listener.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		go s.Handler(conn)
	}
}
