package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
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

func (client *Client) DealResponse() {
	//一旦client.conn有数据就直接copy到Stdout上，永久阻塞监听
	_, err := io.Copy(os.Stdout, client.conn)
	if err != nil {
		fmt.Println("Error reading from server:", err)
	}
	//等价于io.Copy(os.Stdout, client.conn)
	//for {
	//	buff := make([]byte, 1024)
	//	client.conn.Read(buff)
	//	fmt.Println(buff)
	//}

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

func (client *Client) UpdateName() bool {
	fmt.Print("Your Name:")
	// 在读到\n的时候停止读
	reader := bufio.NewReader(os.Stdin)
	line, errRead := reader.ReadString('\n')
	client.Name = strings.TrimSpace(line)
	if errRead != nil {
		return false
	}
	sendMsg := "rename|" + client.Name
	_, err := client.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("conn.Write Error", err)
		return false
	}
	return true
}

func (client *Client) Run() {
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
			client.UpdateName()
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
	//单独开启一个goroutine处理来自服务器的消息
	go client.DealResponse()
	fmt.Println("Welcome", client.Name)
	//go preventDeadlock()
	client.Run()
}
