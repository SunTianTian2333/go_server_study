package main

import (
	"fmt"
	"net"
	"sync"
)

type Server struct {
	Ip   string
	Port int

	// online user list
	OnlineMap map[string]*User
	maplock   sync.RWMutex

	// channel for msg broadcast
	Message chan string
}

// create a server api
func NewServer(ip string, port int) *Server {
	new_server := &Server{
		Ip:        ip,
		Port:      port,
		OnlineMap: make(map[string]*User),
		Message:   make(chan string),
	}
	return new_server
}

// broadcast
func (this *Server) BroadCast(user *User, msg string) {
	sendMsg := "[" + user.addr + "]" + user.Name + ":" + msg
	this.Message <- sendMsg
	fmt.Println(sendMsg)
}

// listen broadcast
func (this *Server) ListenMessage() {
	for {
		msg := <-this.Message

		//send msg to all user
		this.maplock.Lock()
		for _, cli := range this.OnlineMap {
			cli.C <- msg
		}
		this.maplock.Unlock()
	}
}

// handler
func (this *Server) Handler(conn net.Conn) {

	user := NewUser(conn)

	// add user to onlinelist
	this.maplock.Lock()
	this.OnlineMap[user.Name] = user
	this.maplock.Unlock()

	//broadcast msg
	this.BroadCast(user, "is online now")

	// block
	select {}
}

// start server
func (this *Server) Start() {
	// socket listen
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", this.Ip, this.Port))
	if err != nil {
		fmt.Println("net.listen error:", err)
		return
	}
	//close listen socket
	defer listener.Close()

	// start listen message
	go this.ListenMessage()

	for {
		//accept
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("accept error:", err)
			continue
		}

		//handler
		go this.Handler(conn)
	}
}
