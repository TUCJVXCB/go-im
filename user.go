package main

import (
	"net"
	"strings"
)

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

func (u *User) Online() {
	u.server.lock.Lock()
	u.server.userMap[u.name] = u
	u.server.lock.Unlock()

	u.server.broadCast(u, "已上线")
}

func (u *User) offline() {
	u.server.lock.Lock()
	delete(u.server.userMap, u.name)
	u.server.lock.Unlock()

	u.server.broadCast(u, "已退出群聊")
}

func (u *User) sendMessage(msg string) {
	u.conn.Write([]byte(msg + "\n"))
}

func (u *User) DoMessage(msg string) {
	if "who" == msg {
		for _, user := range u.server.userMap {
			onlineMsg := "[" + user.addr + "]" + user.name + ":在线...\n"
			u.sendMessage(onlineMsg)
		}
	} else if strings.HasPrefix(msg, "rename_") {
		_, newName, _ := strings.Cut(msg, "rename_")
		if u.server.userMap[newName] != nil {
			u.sendMessage("用户名已重复~\n")
			return
		}
		delete(u.server.userMap, u.addr)
		u.name = newName
		u.server.userMap[newName] = u
		u.sendMessage("您已成功修改用户名:" + newName)
	} else {
		u.server.broadCast(u, msg)
	}
}

func (u *User) receiveMsg() {
	for {
		msg := <-u.C
		u.conn.Write([]byte(msg + "\n"))
	}
}
