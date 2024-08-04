package main

import (
	"fmt"
	"net"

	"github.com/ctheil/conc-chat-app-server/service"
)

func main() {
	lr, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer lr.Close()

	go service.MonitorRooms()
	service.StartWorkerPool()

	cs := service.Clients{
		Clients:    make(map[net.Conn]service.Client),
		Unregister: make(chan service.Client),
		Register:   make(chan service.Client),
		// Broadcast:  make(chan string),
	}
	go cs.HandleRegistration()

	// go service.HandleRooms()

	for {
		conn, err := lr.Accept()
		if err != nil {
			fmt.Println("error connection to client: ", err)
			continue
		}
		client := service.Client{Conn: conn, Send: make(chan string)}
		cs.Mutex.Lock()
		cs.Clients[conn] = client
		cs.Mutex.Unlock()
		cs.Register <- client
		go cs.HandleClient(&client)
	}
}
