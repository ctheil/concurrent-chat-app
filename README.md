# Conc Chat App

### Small project to apply go concurrency patterns via real world applications.

![preview](https://github.com/user-attachments/assets/29c782a7-643c-4440-ad9e-f4cf7e84ce89)


## For-Select Loop

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
