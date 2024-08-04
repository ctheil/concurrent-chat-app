package service

import "fmt"

type Job struct {
	RoomID uint32
	Msg    string
}

var (
	jobQueue   = make(chan Job, 100)
	numWorkers = 5
)

func worker() {
	for job := range jobQueue {
		// room := RIDs[job.RoomID]
		room := getRoomById(job.RoomID)
		if room == nil {
			fmt.Println("room not found")
			return
		}
		room.Mutex.Lock()
		for _, client := range room.Clients {
			client.Send <- job.Msg
		}
		room.Mutex.Unlock()
	}
}

func StartWorkerPool() {
	for i := 0; i < numWorkers; i++ {
		go worker()
	}
}
