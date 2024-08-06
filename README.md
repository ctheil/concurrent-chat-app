# Conc Chat App

*Small project to apply go concurrency patterns via real world applications*
![preview](https://github.com/user-attachments/assets/29c782a7-643c-4440-ad9e-f4cf7e84ce89)

## Client API Interaction

Clients mainly utilize the application to connect to the TCP port, enter a room and begin chatting. However, once in a room, they need a way to escape the chat (a standalone message prefixes by a colon, the command, and additional args) and interact with the api. Some of my favorites:

### `:erm`

#### Extra-Room Messaging

Allows you to quickly send messages
to another room without leaving your current
room.
e.g. <br/>`:era <room name> send this message to anohter room!`

```go
func (c *Client) ExtraRoomMessage(room, msg string) error {
  r, ok := GetRoom(room, 0, false) // GetRoom(rName string, rID int, createIfNone bool) (r *Room, ok bool)
  if !ok {
    return fmt.Errorf("%s room not found", room)
  }
  r.Broadcast <- fmt.Sprintf("[%s] from room %s says: %s", c.Uname, c.Room.Name, msg)
  return nil
}
```

### `:cr` && `:n`

#### Change Room and New Room
>
>`:cr` or `:change-room` or `:changeroom` or `:change_room`
>allows you to switch out of your current room
>and into a different room.
><br/> e.g. `:cr <room name>`

>`:new` or `:n` of `:nr`
>allows you to switch out of your current room
>and into a brand new room
><br/>e.g. `:n <room name>`

```go
/* The only real different between cr and n is the behavior of GetRoom
* GetRoom should be provided true for `createIfNone` 
* and return false since no room should exist
* 
* This is contrasted by cr 
* which should be provided false for `createIfNone` 
* and return an error if no room is found. 
*/

func (c *Client) ChangeRoom(rname string) error {
  r, ok := GetRoom(rname, 0, false)
  if !ok {
    return fmt.Errorf("no room found with name %s", rname)
  }
 c.ExitRoom()
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
```

### Additional API Methods

- `:h` or `:help` returns the help menu
- `:e` or `:exit` unregisters the client from the room and the TCP connection. When using with `nc` it will end the process.

#### Future methods

I want to add simple auth to rooms where the first person in a room is the host which would expose additional methods for the host relative to the room, maybe `rcmds`, like `hostadd` to add new hosts to the room, and `:passwd` to set a password for the room required to users who want to enter the room or access the ERM method for that room. This introduces interesting complexity to the way clients interact with rooms. This is intriguing because it heavily relies on keeping some values and methods private to the room api.

> [!error] Blocking!
> Because I am using nc as the clinet, I do not have methods to, for example, obfuscate the password input when setting it.

## For-Select Loop && Job Queue Concurrency

### Job Queue

This project utilizes a concurrent job queue to listen to broadcasted messages across room, each with a specific number of clients.
NOTE: I ran into a few issues reading from the map concurrently, so I abstracted reads to the global room map into a function which controls the room mutex (RMastMutex), which differs from each rooms scoped mutex controlling its Client map for broadcasting.

```go
// ...
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
// ...
for job := range jobQueue {
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
// ...
```

### For-Select Loop

This example leverages the `for-select loop` to manage client registration, unregistration, and message broadcasting.

```go
// ...
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
// ...
```

## Running it locally

First, build the server, or execute `go run` to spin up the tcp server:

```bash
cd server
go build -o tmp
./tmp/conc-chat-app-server
```

The `./client` directory is dead as I've simply utilized `$ nc localhost 8080` across different terminals like so:

1. Leveraging a terminal-multiplexer with 3 windows (or by just opening three different terminal windows), execute the following in each:

```sh
nc localhost 8080
```

2. Upon client registration, the server will respond with:

```
❯ nc localhost 8080
Hello! Please provide a screen name:
```

3. For the first type `room`, then for the remaining two windows, provide different screen names: `foo` and `bar`

```
room has joined the chat!
foo has joined the chat!
bar has joined the chat!
```

And that's it! Now the `room` user will show the complese chat history between `foo` and `bar`, and will broadcast each client's registration process as well.
