package main

import (
	"flag"
	"fmt"
	"net"
)

type Client struct {
	ServerIP   string
	ServerPort int
	Name       string
	conn       net.Conn
	flagCase   int
}

func NewClient(serverIP string, serverPort int, name string) *Client {
	client := &Client{
		ServerIP:   serverIP,
		ServerPort: serverPort,
		Name:       name,
		flagCase:   10,
	}
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", serverIP, serverPort))
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return nil
	}
	client.conn = conn
	return client
}

// menu()用于显示菜单选项
func (client *Client) menu() bool {
	var flagCase int
	fmt.Println("1, Group Chat")
	fmt.Println("2, Private Chat")
	fmt.Println("3, Rename")
	fmt.Println("0, Exit")
	_, err := fmt.Scanf("%d\n", &flagCase)
	if err != nil {
		//fmt.Println("Error reading input:", err)
		return false
	}
	if flagCase >= 0 && flagCase <= 3 {
		client.flagCase = flagCase
		return true
	} else {
		fmt.Println("输入0-3内的数字")
		return false
	}
}

func (client *Client) run() {
	for client.flagCase != 0 {
		for client.menu() != true {
		}
		switch client.flagCase {
		case 1:
			fmt.Println("Group Chat")
			break
		case 2:
			fmt.Println("Private Chat")
			break
		case 3:
			fmt.Println("Rename")
			break
		}
	}
}

// 命令行解析
var serverIP string
var serverPort int
var name string

// init()初始化，用于命令行解析
func init() {
	flag.StringVar(&serverIP, "ip", "127.0.0.1", "Server IP")
	flag.IntVar(&serverPort, "port", 8888, "Server Port")
	flag.StringVar(&name, "name", "Luke", "Name")
}

func main() {
	//命令行解析
	flag.Parse()
	client := NewClient(serverIP, serverPort, name)
	if client == nil {
		fmt.Println("Error connecting to server")
		return
	}
	fmt.Println("Welcome", client.Name)
	//go preventDeadlock()
	client.run()
}
