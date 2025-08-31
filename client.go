package main

import (
	"fmt"
	"net"
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

func main() {
	client := NewClient("127.0.0.1", 8888, "Luke")
	if client == nil {
		fmt.Println("Error connecting to server")
		return
	}
	fmt.Println("Welcome ", client.Name)
	select {}
}
