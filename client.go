package main

import (
	"flag"
	"fmt"
	"net"
	"time"
)

type Client struct {
	ServerIP   string
	ServerPort int
	Name       string
	conn       net.Conn
}

func NewClient(serverIP string, serverPort int, name string) *Client {
	client := &Client{
		ServerIP:   serverIP,
		ServerPort: serverPort,
		Name:       name,
	}
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", serverIP, serverPort))
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return nil
	}
	client.conn = conn
	return client
}

var serverIP string
var serverPort int
var name string

func init() {
	flag.StringVar(&serverIP, "ip", "127.0.0.1", "Server IP")
	flag.IntVar(&serverPort, "port", 8888, "Server Port")
	flag.StringVar(&name, "name", "Luke", "Name")
}

func preventDeadlock() {
	for {
		time.Sleep(1 * time.Second)
	}
}

func main() {
	//命令行解析
	flag.Parse()
	client := NewClient(serverIP, serverPort, name)
	if client == nil {
		fmt.Println("Error connecting to server")
		return
	}
	fmt.Println("Welcome ", client.Name)
	go preventDeadlock()
	select {}
}
