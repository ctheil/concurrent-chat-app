package service

import (
	"fmt"
	"net"
	"sync"
)

type (
	ClientsMap map[net.Conn]Client
	Clients    struct {
		Clients ClientsMap
		// Broadcast  chan string
		Register   chan Client
		Unregister chan Client
		Mutex      sync.Mutex
	}
)

func (c *Clients) HandleClient(client *Client) {
	defer func() {
		if client.Room != nil {
			client.Room.Broadcast <- fmt.Sprintf("\n%s has left the chat\n", client.Uname)
		}
		c.Unregister <- *client
	}()

	go client.SendMessage()
	client.ListenForMessages()
}

func (c *Clients) HandleRegistration() {
	for {
		select {
		case client := <-c.Register:
			// fmt.Printf("new client connected\n", client.Conn.RemoteAddr())
			client.Send <- "Hello! Please provide a screen name & room like so: 'my_username:my_room'"
		case client := <-c.Unregister:
			c.Mutex.Lock()
			delete(c.Clients, client.Conn)
			// c.HandleUnregister()
			c.Mutex.Unlock()
			fmt.Printf("\n%s:%s unregistered\n", client.Conn.RemoteAddr(), client.Uname)
		}
	}
}
