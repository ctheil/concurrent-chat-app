package service

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

type Client struct {
	Conn  net.Conn
	Send  chan string
	Uname string

	Room *Room
}

func (c *Client) ListenForMessages() {
	scanner := bufio.NewScanner(c.Conn)
	for scanner.Scan() {
		msg := scanner.Text()
		// handle registration
		if !c.IsRegistered() {
			if err := c.HandleRegistration(msg); err != nil {
				fmt.Println("Error registering client!", err)
				c.Send <- "Invalid format. Please enter your username, a colon, followed by the room you would like to join"
				continue
			}
			c.Room.Broadcast <- fmt.Sprintf("%s has joined %s", c.Uname, c.Room.Name)
			continue
		}

		c.Room.Broadcast <- fmt.Sprintf("%s: %s", c.Uname, msg)
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("error reading from client", err)
	}
}

func (c *Client) IsRegistered() bool {
	return c.Uname != "" && c.Room != nil
}

func (c *Client) HandleRegistration(data string) error {
	strs := strings.Split(data, ":")
	if len(strs) != 2 {
		return fmt.Errorf("missing username or roomname")
	}

	c.Uname = strs[0]
	// fmt.Println("client registered:", c.Uname)
	rname := strs[1]

	r, ok := GetRoom(rname, 0, true)
	if !ok {
		return fmt.Errorf("error searching for or creating new room")
	}
	r.AddClient(c)
	c.Room = r
	fmt.Printf("\n %s registered to room %s \n", c.Uname, c.Room.Name)

	return nil
}

func (c *Client) GetMsgPrefix() string {
	if !c.IsRegistered() {
		return ""
	}
	return fmt.Sprintf("%s : %s    ", c.Room.Name, c.Uname)
}

func (c *Client) SendMessage() {
	for msg := range c.Send {
		fmt.Fprintln(c.Conn, msg)
	}
}

func (c *Client) JoinRoomById(id uint32, rooms RoomIDS) {
	r := rooms[id]
	c.Room = r
	r.AddClient(c)
}

func (c *Client) JoinRoomByName(name string, rooms RoomNames) {
	r := rooms[name]
	c.Room = r
	r.AddClient(c)
}
