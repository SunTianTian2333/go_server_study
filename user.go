package main

import (
	"net"
	"strings"
)

type User struct {
	Name   string
	addr   string
	C      chan string
	conn   net.Conn
	server *Server
}

// create a user api
func NewUser(conn net.Conn, server *Server) *User {
	userAddr := conn.RemoteAddr().String()
	user := &User{
		Name:   userAddr,
		addr:   userAddr,
		C:      make(chan string),
		conn:   conn,
		server: server,
	}
	go user.ListenMessage()
	return user
}

// listen message
func (this *User) ListenMessage() {
	for {
		msg := <-this.C
		// send to client
		this.conn.Write([]byte(msg + "\n"))
	}

}

// user online
func (this *User) Online() {

	// add user to onlinelist
	this.server.maplock.Lock()
	this.server.OnlineMap[this.Name] = this
	this.server.maplock.Unlock()

	//broadcast msg
	this.server.BroadCast(this, "is online now")

}

// user offline
func (this *User) Offline() {
	// add user to onlinelist
	this.server.maplock.Lock()
	delete(this.server.OnlineMap, this.Name)
	this.server.maplock.Unlock()

	//clear
	close(this.C)
	this.conn.Close()
	//broadcast msg
	this.server.BroadCast(this, "is offline now")
}

// send msg
func (this *User) SendMessage(msg string) {
	this.conn.Write([]byte(msg + "\n"))
}

// deal mode
func (this *User) deal_mode() {

}

// deal msg
func (this *User) DoMessage(msg string) {

	if msg == "who" {
		this.server.maplock.Lock()
		for _, user := range this.server.OnlineMap {
			onlineMsg := "[" + user.addr + "]" + user.Name + ":" + "online \n"
			this.SendMessage(onlineMsg)
		}
		this.server.maplock.Unlock()
	} else if len(msg) > 7 && msg[:7] == "rename:" {
		newName := strings.Split(msg, ":")[1]
		_, ok := this.server.OnlineMap[newName]
		if ok {
			this.SendMessage("this name already exist")
		} else {
			this.server.maplock.Lock()
			delete(this.server.OnlineMap, this.Name)
			this.Name = newName
			this.server.OnlineMap[newName] = this
			this.SendMessage("name changed! your new name is : " + this.Name)
			this.server.maplock.Unlock()
		}
	} else if len(msg) > 4 && msg[:3] == "to|" {
		rName := strings.Split(msg, "|")[1]
		if rName == "" {
			this.SendMessage("wrong format! use \"to|someone|message\"")
			return
		}
		rUser, ok := this.server.OnlineMap[rName]
		if !ok {
			this.SendMessage("user name doesnt exist!")
			return
		}
		content := strings.Split(msg, "|")[2]
		if content == "" {
			this.SendMessage("empty message!")
			return
		}
		rUser.SendMessage(this.Name + " said to you: " + content)
		this.SendMessage("message sent")

	}
}
