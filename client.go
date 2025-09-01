package main

import (
	"bufio"
	"flag"
	"fmt"
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
	//等价于io.Copy(os.Stdout, client.conn)
	//for {
	//	buff := make([]byte, 1024)
	//	client.conn.Read(buff)
	//	fmt.Println(buff)
	//}
	reader := bufio.NewReader(client.conn)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Connection closed by server.")
			//client.conn.Close()
			return
		}
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "CMD|") {
			// Handle command messages here, for now just print
			if line[:11] == "CMD|OFFLINE" {
				client.conn.Close()
				os.Exit(0)
			}
		} else {
			fmt.Println(line)
		}
	}
}

// menu()用于显示菜单选项
func (client *Client) menu() bool {
	var flagCase int
	fmt.Println("1, Public Chat")
	fmt.Println("2, Private Chat")
	fmt.Println("3, Rename")
	fmt.Println("4, Online Users")
	fmt.Println("0, Exit")
	_, err := fmt.Scanf("%d\n", &flagCase)
	if err != nil {
		//fmt.Println("Error reading input:", err)
		return false
	}
	if flagCase >= 0 && flagCase <= 4 {
		client.flagCase = flagCase
		return true
	} else {
		fmt.Println("输入0-4内的数字")
		return false
	}
}

func (client *Client) PublicChat() bool {
	var chatMsg string
	fmt.Println("Public Chat (exit to quit)")
	_, err := fmt.Scanln(&chatMsg)
	if err != nil {
	}
	for chatMsg != "exit" {
		if len(chatMsg) != 0 {
			sendMsg := chatMsg
			_, err := client.conn.Write([]byte(sendMsg))
			if err != nil {
				return false
			}
		}
		//chatMsg = ""
		fmt.Println("Public Chat (exit to quit)")
		_, err := fmt.Scanln(&chatMsg)
		if err != nil {
		}
	}
	return true
}

func (client *Client) SelectUsers() {
	sendMsg := "CMD|WHO"
	_, err := client.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("conn Write error:", err)
	}
}

func (client *Client) PrivateChat() bool {
	var chatUserName string
	var chatMsg string
	client.SelectUsers()
	fmt.Println("to whom you want to chat (exit to quit)")
	fmt.Scanln(&chatUserName)

	for chatUserName != "exit" {
		fmt.Println("input your message:(exit to quit)")
		fmt.Print("->" + chatUserName + ":")
		fmt.Scanln(&chatMsg)
		for chatMsg != "exit" {
			if len(chatMsg) != 0 {
				sendMsg := "CMD|TO|" + chatUserName + "|" + chatMsg
				_, err := client.conn.Write([]byte(sendMsg))
				if err != nil {
					return false
				}
			}
			chatMsg = ""
			fmt.Print("->" + chatUserName + ":")
			fmt.Scanln(&chatMsg)
		}
		client.SelectUsers()
		fmt.Println("to whom you want to chat (exit to quit)")
		fmt.Scanln(&chatUserName)
	}
	return true
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
	sendMsg := "CMD|RENAME|" + client.Name
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
			client.PublicChat()
			break
		case 2:
			client.PrivateChat()
			break
		case 3:
			client.UpdateName()
			break
		case 4:
			client.SelectUsers()
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
