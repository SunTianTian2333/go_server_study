package main

import "net"

type User struct {
	Name string
	addr string
	C    chan string
	conn net.Conn
}

// create a user api
func NewUser(conn net.Conn) *User {
	userAddr := conn.RemoteAddr().String()
	user := &User{
		Name: userAddr,
		addr: userAddr,
		C:    make(chan string),
		conn: conn,
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
