package main

import (
	"bufio"
	"fmt"
	"net"
	"sync"
)

type Client struct {
	conn  net.Conn
	send  chan string
	uname string
}

var (
	clients    = make(map[net.Conn]Client)
	unregister = make(chan Client)
	register   = make(chan Client)
	broadcast  = make(chan string)
	mutex      sync.Mutex
)

func main() {
	lr, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer lr.Close()

	go handleMessages()

	for {
		conn, err := lr.Accept()
		if err != nil {
			fmt.Println("error connection to client: ", err)
			continue
		}
		client := Client{conn: conn, send: make(chan string)}
		mutex.Lock()
		clients[conn] = client
		mutex.Unlock()
		register <- client
		go handleClient(client)
	}
}

func handleMessages() {
	for {
		select {
		case msg := <-broadcast:
			mutex.Lock()
			for _, c := range clients {
				c.send <- msg
			}
			mutex.Unlock()
		case c := <-register:
			fmt.Println("Hello ", c.conn.RemoteAddr())
			c.send <- "Hello! Please provide a screen name: "
		case c := <-unregister:
			mutex.Lock()
			delete(clients, c.conn)
			mutex.Unlock()
			fmt.Println("Goodbye ", c.uname)
		}
	}
}

func handleClient(client Client) {
	defer func() {
		broadcast <- fmt.Sprintf("\n%s has left the chat\n", client.uname)
		unregister <- client
		client.conn.Close()
	}()

	registered := false

	go sendMesagesToClient(client)

	scanner := bufio.NewScanner(client.conn)
	for scanner.Scan() {
		message := scanner.Text()
		if !registered {
			// first message
			client.uname = message
			registered = true
		}

		if message == client.uname {
			broadcast <- fmt.Sprintf("\n%s has joined the chat!\n", client.uname)
		} else {
			broadcast <- fmt.Sprintf("%s: %s", client.uname, message)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("error reading from client", err)
	}
}

func sendMesagesToClient(client Client) {
	for msg := range client.send {
		fmt.Fprintln(client.conn, msg)
	}
}
