package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func send_message(msg string, conn net.Conn) {
	_, err := conn.Write([]byte(msg))
	if err != nil {
		fmt.Println("error writing to server: ", err)
		return
	} else {
		fmt.Println("sent!")
	}

	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("error reading from the server", err)
		return
	}

	fmt.Println("recieved from the server:", string(buffer[:n]))
}

func main() {
	conn, err := net.Dial("tcp", "10.0.0.218:8080")
	if err != nil {
		fmt.Println("error connection: ", err)
		os.Exit(1)
	}
	defer conn.Close()

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Println("Compose and Send message: ")
		if scanner.Scan() {
			text := scanner.Text()
			if len(text) != 0 {
				if text == "DISCONN" {
					conn.Close()
					fmt.Println("Disconnecting...")
					os.Exit(0)
				}
				fmt.Println("sending...")
				send_message(text, conn)
			}
		} else {
			fmt.Println("Error reading input:", scanner.Err())
			break
		}
	}
}
