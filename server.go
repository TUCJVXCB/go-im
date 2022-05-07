package main

import (
	"fmt"
	"net"
)

type Server struct {
	IP   string
	Port int
}

func NewServer(IP string, port int) *Server {
	server := &Server{
		IP:   IP,
		Port: port,
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
	fmt.Println("新建连接..")
}
