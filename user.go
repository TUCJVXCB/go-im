package main

import "net"

type User struct {
	name   string
	addr   string
	C      chan string
	conn   net.Conn
	server *Server
}

func NewUser(conn net.Conn, server *Server) *User {
	userAddr := conn.RemoteAddr().String()
	u := &User{
		name:   userAddr,
		addr:   userAddr,
		C:      make(chan string),
		conn:   conn,
		server: server,
	}

	go u.receiveMsg()

	return u
}

func (user *User) Online() {
	user.server.lock.Lock()
	user.server.userMap[user.name] = user
	user.server.lock.Unlock()

	user.server.broadCast(user, "已上线")
}

func (user *User) offline() {
	user.server.lock.Lock()
	delete(user.server.userMap, user.name)
	user.server.lock.Unlock()

	user.server.broadCast(user, "已退出群聊")
}

func (user *User) DoMessage(msg string) {
	user.server.broadCast(user, msg)
}

func (user *User) receiveMsg() {
	for {
		msg := <-user.C
		user.conn.Write([]byte(msg + "\n"))
	}
}
