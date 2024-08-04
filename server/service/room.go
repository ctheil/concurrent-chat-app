package service

import (
	"fmt"
	"sync"
)

var (
	RIDs                  = make(RoomIDS)
	RNames                = make(RoomNames)
	IDCounter  uint32     = 0
	Rooms                 = make(chan Room)
	RMastMutex sync.Mutex // Master mutex for RIDs and RNames
)

type (
	RoomIDS   map[uint32]*Room
	RoomNames map[string]*Room
	Room      struct {
		Broadcast chan string
		ID        uint32
		Name      string
		Clients   ClientsMap
		Mutex     sync.Mutex
	}
)

func (r *Room) AddClient(client *Client) {
	r.Mutex.Lock()
	r.Clients[client.Conn] = *client
	r.Mutex.Unlock()
	r.Broadcast <- fmt.Sprintf("%s has joined the room\n", client.Uname)
}

func (r *Room) RemoveClient(client *Client) {
	r.Mutex.Lock()
	delete(r.Clients, client.Conn)
	r.Mutex.Unlock()
	r.Broadcast <- fmt.Sprintf("%s has left the room\n", client.Uname)
}

func NewRoom(name string) *Room {
	IDCounter++ // WHY GO, WHY???
	r := Room{
		Broadcast: make(chan string),
		ID:        IDCounter, // WHY NOT IDCounter++???
		Name:      name,
		Clients:   make(ClientsMap),
	}
	RMastMutex.Lock()
	RIDs[r.ID] = &r
	RNames[r.Name] = &r
	RMastMutex.Unlock()

	return &r
}

func GetRoom(name string, id uint32, createNewIfNone bool) (room *Room, ok bool) {
	var r *Room

	RMastMutex.Lock()

	if name != "" {
		r = RNames[name]
	}
	if id != 0 {
		r = RIDs[id]
	}

	if r == nil && name != "" && createNewIfNone {
		RMastMutex.Unlock()
		return NewRoom(name), false
	}

	RMastMutex.Unlock()
	return r, r != nil
}

func getRoomById(id uint32) *Room {
	RMastMutex.Lock()
	defer RMastMutex.Unlock()
	return RIDs[id]
}

func MonitorRooms() {
	for {
		RMastMutex.Lock()
		for _, room := range RIDs {
			select {
			case msg := <-room.Broadcast:
				jobQueue <- Job{RoomID: room.ID, Msg: fmt.Sprintf("[%s] %s", room.Name, msg)}
			default:
				// no message, continue to next room
			}
		}
		RMastMutex.Unlock()
	}
}
