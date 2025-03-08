package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
)

type Client struct {
	ServerIp   string
	ServerPort int
	Name       string
	conn       net.Conn
	flag       int
}

func NewClient(serverIp string, serverPort int) *Client {
	// create client
	client := &Client{
		ServerIp:   serverIp,
		ServerPort: serverPort,
		flag:       1000,
	}

	//connect server
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", serverIp, serverPort))
	if err != nil {
		fmt.Println("net,Dial error:", err)
		return nil
	}
	client.conn = conn
	//return
	return client
}

func (client *Client) DealResponse() {
	io.Copy(os.Stdout, client.conn)
}

func (client *Client) PrivateChat() bool {
	var sName string
	var sendMsg string

	client.conn.Write([]byte("who\n"))
	fmt.Println("please input name ")
	fmt.Scanln(&sName)

	fmt.Println("please input your message")

	fmt.Scanln(&sendMsg)

	_, err := client.conn.Write([]byte("to|" + sName + "|" + sendMsg + "\n\n"))
	if err != nil {
		fmt.Println("conn.Write err:", err)
		return false
	}
	return true
}

func (client *Client) PublicChat() bool {
	fmt.Println("please input your message")
	var sendMsg string
	fmt.Scanln(&sendMsg)

	_, err := client.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("conn.Write err:", err)
		return false
	}
	return true
}

func (client *Client) UpdateName() bool {
	fmt.Println("please input new user name")
	fmt.Scanln(&client.Name)

	sendMsg := "rename:" + client.Name + "\n"
	_, err := client.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("conn.Write err:", err)
		return false
	}
	fmt.Println("your name changed! ")
	return true
}

func (client *Client) Run() {
	for client.flag != 0 {
		for client.menu() != true {
		}
		switch client.flag {
		case 1:
			fmt.Println("chose public mode")
			client.PublicChat()
			break
		case 2:
			fmt.Println("chose private mode")
			client.PrivateChat()
			break
		case 3:
			fmt.Println("chose change name mode")
			client.UpdateName()
			break
		}
	}
}

func (client *Client) menu() bool {
	var flag int
	fmt.Println("1. public chat")
	fmt.Println("2. private chat")
	fmt.Println("3. change name")
	fmt.Println("0. quit")

	fmt.Scanln(&flag)

	if flag >= 0 && flag <= 3 {
		client.flag = flag
		return true
	} else {
		fmt.Println("please imput correct number!")
		return false
	}
}

var serverIp string
var serverPort int

func init() {
	flag.StringVar(&serverIp, "ip", "127.0.0.1", "set server ip")
	flag.IntVar(&serverPort, "port", 8888, "set server port")
}

func main() {

	flag.Parse()

	client := NewClient("127.0.0.1", 8888)
	if client == nil {
		fmt.Println(">>>>>>>>>>>> connect failed....")
		return
	}
	fmt.Println(">>>>>>>>>>>> connect success!")

	go client.DealResponse()

	//start
	client.Run()
}
