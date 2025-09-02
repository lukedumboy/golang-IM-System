Golang-IM-System

A simple Instant Messaging (IM) system implemented in Go, supporting multiple clients over TCP.
This project includes both server and client programs with command-line interaction.

***

✨ Features
- ✅ User online/offline broadcast
- ✅ Public chat room (all users)
- ✅ Private chat between users
- ✅ Rename user
- ✅ Query online users
- ✅ Auto offline when timeout
- ✅ Colored terminal output (optional)

🚀 Usage
1. Start the server
```
go build -o server main.go server.go user.go
./server
```
3. Start the client
```
go build -o client client.go
./client
```
5. Client menu
  - 1 → Public Chat
	- 2 → Private Chat
	- 3 → Rename
	- 4 → Online Users
	- 0 → Exit

📂 Project Structure
```
.
├── server.go      # Main server program
├── user.go        # User object and related methods
├── client.go      # Client program
```

🔧 Technical Highlights

- Built with Go net package (TCP socket programming)
- goroutines for handling multiple connections
- channels for message broadcasting & heartbeat detection
- Command-based message parsing (CMD|RENAME, CMD|TO, etc.)

📌 TODO / Future Plans
- Save chat history (message persistence)
- WebSocket support (for frontend integration)
- Docker deployment