package main

import (
	"fmt"
	"io"
	"net"
	"sync"
)

type Server struct {
	IP   string
	Port int

	userMap map[string]*User
	lock    sync.RWMutex
	C       chan string
}

func NewServer(IP string, port int) *Server {
	server := &Server{
		IP:      IP,
		Port:    port,
		userMap: make(map[string]*User),
		C:       make(chan string),
	}
	return server
}

func (s *Server) Start() {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.IP, s.Port))

	if err != nil {
		fmt.Println("net Listen err:", err)
		return
	}

	defer listener.Close()

	go s.listen()

	for {
		accept, err := listener.Accept()
		if err != nil {
			fmt.Println("net connect err:", err)
			continue
		}
		go s.handle(accept)
	}
}

func (s *Server) handle(conn net.Conn) {
	user := NewUser(conn, s)
	user.Online()

	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := conn.Read(buf)
			if n == 0 {
				user.offline()
				return
			}
			if err != nil && err != io.EOF {
				fmt.Println("conn read err", err)
				return
			}
			msg := string(buf[:n-1])
			user.DoMessage(msg)
		}
	}()

	select {}
}

func (s *Server) listen() {
	for {
		msg := <-s.C
		s.lock.Lock()
		for _, user := range s.userMap {
			user.C <- msg
		}
		s.lock.Unlock()
	}
}

func (s *Server) broadCast(user *User, msg string) {
	sendMsg := "[" + user.addr + "]" + user.name + ":" + msg
	s.C <- sendMsg
}
