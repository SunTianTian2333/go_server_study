package main

import (
	"net"
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

	//broadcast msg
	this.server.BroadCast(this, "is offline now")
}

// deal msg
func (this *User) DoMessage(msg string) {
	//receive msg from client
	this.server.BroadCast(this, msg)
}
