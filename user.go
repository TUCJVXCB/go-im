package main

import "net"

type User struct {
	name string
	addr string
	C    chan string
	conn net.Conn
}

func NewUser(conn net.Conn) *User {
	u := &User{
		name: conn.RemoteAddr().String(),
		addr: conn.RemoteAddr().String(),
		C:    make(chan string),
		conn: conn,
	}

	go u.receiveMsg()

	return u
}

func (user *User) receiveMsg() {
	for {
		msg := <-user.C
		user.conn.Write([]byte(msg))
	}
}
