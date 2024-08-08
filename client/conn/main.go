package conn

import (
	"bufio"
	"fmt"
	"net"

	"github.com/ctheil/conc-chat-app-client/tui"
)

func NewConnection() (net.Conn, error) {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		panic(err)
	}

	return conn, err
}

func Listen(conn net.Conn, msgs chan<- string, tui *tui.TUI) {
	fmt.Println("[Listen]: go routine started: listening")
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		message := scanner.Text()
		// msgs <- message
		fmt.Fprintf(&tui.ChatView, "[yellow::b]%s\n", message)
		//
	}

	if err := scanner.Err(); err != nil {
		tui.Error(err.Error())
		panic(err)
	}
}
