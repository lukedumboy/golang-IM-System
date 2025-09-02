Golang-IM-System

一个使用 Go 语言实现的简单即时通讯（IM）系统，支持通过 TCP 的多客户端通信。

该项目包含服务端和客户端两个程序，并提供命令行交互。

***

✨  功能特性
- ✅ 用户上线/下线广播
- ✅ 公共聊天室（所有用户）
- ✅ 用户间私聊
- ✅ 修改用户名
- ✅ 查询在线用户
- ✅ 超时自动下线
- ✅ 彩色终端输出（可选）

🚀 使用方法
1. 启动服务端
```
go build -o server main.go server.go user.go
./server
```
3. 启动客户端
```
go build -o client client.go
./client
```
5. 客户端菜单
  - 1 → Public Chat
	- 2 → Private Chat
	- 3 → Rename
	- 4 → Online Users
	- 0 → Exit

📂 项目结构
```
.
├── server.go      # Main server program
├── user.go        # User object and related methods
├── client.go      # Client program
```

🔧 技术亮点

- 使用 Go net 包实现（TCP 套接字编程）
- goroutine 并发处理多连接
- channel 用于消息广播与心跳检测
- 基于命令的消息解析（如 CMD|RENAME、CMD|TO 等）

📌 TODO / Future Plans
- 聊天记录保存（消息持久化）
- 支持 WebSocket（方便前端集成）
- Docker 部署