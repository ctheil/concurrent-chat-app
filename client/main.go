package main

import (
	"github.com/ctheil/conc-chat-app-client/conn"
	"github.com/ctheil/conc-chat-app-client/tui"
)

func main() {
	tui := tui.NewApp()
	// go tui.Listen()
	//
	//
	//
	c, err := conn.NewConnection()
	if err != nil {
		tui.Error("could not read from server.")
		tui.Run()
		return
	}
	defer c.Close()
	go conn.Listen(c, tui.Incoming, tui)
	if err := tui.Run(); err != nil {
		panic(err)
	}
}
