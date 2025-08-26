package main

import "net"

type User struct {
	Name string
	Addr string
	c    chan string
	conn net.Conn
}

func NewUser(conn net.Conn) *User {
	userAddr := conn.RemoteAddr().String()
	user := &User{
		Name: userAddr,
		Addr: userAddr,
		c:    make(chan string),
		conn: conn,
	}
	go user.ListenMessage()
	return user
}

func (u *User) ListenMessage() {
	for {
		message := <-u.c
		u.conn.Write([]byte(message))
	}
}
