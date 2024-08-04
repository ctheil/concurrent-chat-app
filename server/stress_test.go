package main

import (
	"fmt"
	"net"
	"sync"
	"testing"
	"time"
)

// func startServer() {
// 	go service.MonitorRooms()
// 	service.StartWorkerPool()
//
// 	for i := 1; i <= 5; i++ {
// 		service.NewRoom(fmt.Sprintf("room%d", i))
// 	}
// }

func createClient(roomName, clientName string, wg *sync.WaitGroup, msgCount int, delay time.Duration) bool {
	defer wg.Done()

	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Error connecting:", err)
		return false
	}

	defer conn.Close()

	// SET CLIENT NAME && ROOM
	conn.Write([]byte(fmt.Sprintf("%s:%s", clientName, roomName)))

	for i := 0; i < msgCount; i++ { // adjust the number of messages per client as needed
		msg := fmt.Sprintf("message %d from client in %s", i, roomName)
		conn.Write([]byte(msg))
		time.Sleep(delay)
	}
	return true
}

func TestStress(t *testing.T) {
	var wg sync.WaitGroup

	numClients := 2000
	numRooms := 1000
	messageCount := 1000
	delay := 10 * time.Millisecond
	fmt.Printf("Stress Testing Server with:\n\n    %d Clients\n    Across %d Rooms\n    Each sending %d Messages\n    With a delay of %d MS", numClients, numRooms, messageCount, delay/100000)

	for i := 0; i < numClients; i++ {
		wg.Add(1)
		rName := fmt.Sprintf("r%d", (i%numRooms)+1)
		cName := fmt.Sprintf("c%d", i)
		go createClient(rName, cName, &wg, messageCount, delay)
	}

	wg.Wait()
}
