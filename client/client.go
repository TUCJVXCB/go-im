package main

import (
	"fmt"
	"io"
	"net"
	"os"
)

type Client struct {
	ServerIP   string
	ServerPort int
	name       string
	conn       net.Conn
	flag       int
}

func NewClient(serverIP string, serverPort int) *Client {
	client := &Client{
		ServerIP:   serverIP,
		ServerPort: serverPort,
		flag:       999,
	}
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", serverIP, serverPort))
	if err != nil {
		fmt.Println("net.Dial error:", err)
		return nil
	}
	client.conn = conn
	return client
}

func (client *Client) menu() bool {
	var flag int
	fmt.Println("1.公聊模式")
	fmt.Println("2.私聊模式")
	fmt.Println("3.更新用户名")
	fmt.Println("0.退出")

	fmt.Scanln(&flag)

	if flag >= 0 && flag <= 3 {
		client.flag = flag
		return true
	} else {
		fmt.Println(">>>>请输入合法范围内的数字<<<<")
		return false
	}
}

var ip, port string

//func init() {
//	flag.StringVar(&ip, "ip", "127.0.0.1", "设置服务器Ip地址")
//}

func (client *Client) Run() {
	for client.flag != 0 {
		for client.menu() != true {
		}

		switch client.flag {
		case 1:
			fmt.Println("公聊模式选择...")
			client.handlePublicChat()
			break
		case 2:
			fmt.Println("私聊模式选择...")
			break
		case 3:
			fmt.Println("更新用户名选择...")
			client.handleRename()
			break
		}
	}
}

func (client *Client) handlePublicChat() {
	fmt.Println("请输入您要发送的消息:")
	var msg string

	for {
		fmt.Scanln(&msg)
		if msg == "exit" {
			break
		} else {
			client.conn.Write([]byte(msg + "\n"))
		}
	}
}

func (client *Client) handleRename() {
	fmt.Println("请输入您要修改的名称:")
	var name string

	for {
		_, err := fmt.Scanln(&name)
		if err != nil {
			fmt.Println("read from stdin err", err)
			return
		}
		if name == "" {
			fmt.Println("输入为空，请重新输入")
		} else {
			client.conn.Write([]byte("rename_" + name + "\n"))
			break
		}
	}
}

func (client *Client) DealWithResponse() {
	io.Copy(os.Stdout, client.conn)
}

func main() {
	client := NewClient("127.0.0.1", 8080)
	if client == nil {
		fmt.Println("连接服务器失败...")
		return
	}
	fmt.Println("连接服务器成功...")

	go client.DealWithResponse()

	client.Run()

}
