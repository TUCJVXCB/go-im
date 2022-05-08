package main

import (
	"fmt"
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
	user := NewUser(conn)
	msg := "[" + user.addr + "]" + user.name + "加入连接\n"
	s.lock.Lock()
	s.userMap[user.name] = user
	s.lock.Unlock()

	s.broadCast(msg)

	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := conn.Read(buf)
			if n == 0 {
				s.broadCast(user.name + "下线了")
			}
			if err != nil {
				fmt.Println("conn read err", err)
				return
			}
			msg := string(buf)
			s.broadCast(msg)
		}
	}()

	select {}
}

func (s *Server) listen() {
	for {
		msg := <-s.C
		s.broadCast(msg)
	}
}

func (s *Server) broadCast(msg string) {
	s.lock.Lock()
	for _, user := range s.userMap {
		user.C <- msg
	}
	s.lock.Unlock()
}
