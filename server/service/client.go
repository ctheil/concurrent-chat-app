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
		isCmd, cmd := messageIsCommand(msg)
		if isCmd {
			c.ExecuteCommand(cmd)
			continue
		}

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

	r, _ := GetRoom(rname, 0, true)
	// if !ok {
	// 	return fmt.Errorf("error searching for or creating new room")
	// }
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

func (c *Client) JoinRoomByName(name string, rooms RoomNames) {
	r := rooms[name]
	c.Room = r
	r.AddClient(c)
}

func (c *Client) ExitRoom() {
	c.Room.RemoveClient(c)
	c.Room = nil
}

func messageIsCommand(s string) (bool, []string) {
	fmt.Println("checking is cmd:", s)
	return strings.HasPrefix(s, ":"), strings.Split(strings.TrimPrefix(s, ":"), " ")
}

func (c *Client) ChangeRoom(rname string) error {
	c.ExitRoom()
	r, ok := GetRoom(rname, 0, false)
	if !ok {
		return fmt.Errorf("no room found with name %s", rname)
	}
	c.Room = r
	r.AddClient(c)
	return nil
}

func (c *Client) NewRoom(rname string) error {
	r, ok := GetRoom(rname, 0, true)
	if ok {
		return fmt.Errorf("%s already exists", rname)
	}
	c.ExitRoom()
	r.AddClient(c)
	c.Room = r
	return nil
}

func (c *Client) ExtraRoomMessage(room, msg string) error {
	r, ok := GetRoom(room, 0, false)
	if !ok {
		return fmt.Errorf("%s room not found", room)
	}
	r.Broadcast <- fmt.Sprintf("[%s] from room %s says: %s", c.Uname, c.Room.Name, msg)
	return nil
}

func (c *Client) ExecuteCommand(cmd []string) {
	switch cmd[0] {
	case "help", "h":
		fmt.Println("Client issued help")
		c.Send <- help
	// c.Send <- //help message
	case "exit", "e":
		fmt.Println("Client issued exit")
		c.Conn.Close()
	case "cr", "change-room", "changeroom", "change_room":
		if cmd[1] == "" {
			c.Send <- "change room requires the name of the room you wish to switch into, like so ':cr room_2'"
			return
		}
		if err := c.ChangeRoom(cmd[1]); err != nil {
			c.Send <- fmt.Sprintf("error: %s", err.Error())
		}
		c.Send <- fmt.Sprintf("changed room to %s", cmd[1])
	case "new", "n", "nr":
		if cmd[1] == "" {
			c.Send <- "new room requires the name of the room you wish to create, like so ':n room_2'"
			return
		}
		if err := c.NewRoom(cmd[1]); err != nil {
			c.Send <- fmt.Sprintf("error: %s. If the room already exists, switch into it using ':cr %s", err.Error(), cmd[1])
		}
		c.Send <- fmt.Sprintf("%s created new room %s", c.Uname, cmd[1])
	case "erm":
		fmt.Println("client issued erm cmd")
		if cmd[1] == "" || cmd[2] == "" {
			c.Send <- "erm requires the name of the room you with to message, followed by the message, like so ':erm room_2 some message'"
			return
		}
		if err := c.ExtraRoomMessage(cmd[1], cmd[2]); err != nil {
			c.Send <- err.Error()
		}
		c.Send <- fmt.Sprintf("message sent to %s", cmd[1])
		return
	default:
		c.Send <- fmt.Sprintf("%s is an unknown command. type :h for help.", cmd[0])

	}
}

const help = `CONC Chat App: 

  :h or :help will always issue this command.

  :e or :exit will end your current session.

  :cr or :change-room or :changeroom or :change_room 
  allows you to switch out of your current room
  and into a different room. 
  e.g. ':cr <room name>'

  :new or :n of :nr 
  allows you to switch out of your current room 
  and into a brand new room
  e.g. ':n <room name>
  :erm
  Stands for extra-room message, 
  allowing you to quickly send messages 
  to another room without leaving your current 
  room. 
  e.g. ':era <room name> send this message to anohter room!'

  Enjoy!
  `
