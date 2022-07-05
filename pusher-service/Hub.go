package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/mrpiggy97/cqrs/repository"
)

type Hub struct {
	Clients    []repository.IClient
	Register   chan repository.IClient
	Unregister chan repository.IClient
	Number     chan int
	Locker     *sync.Mutex
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// flow of code will be the following hub.run will be waiting for
// register channel to send any new *Client type, when a request is sent
// to handlers.WebSocketHandler register channel will send new *Client created
// and will run client.Write which will have a channel waiting for a byte message
// once we get another request to an eligible endpoint, that same endpoint will
// send a message through hub.BroadCast, which will send a byte message
// that will be recieved by every client in its array of clients
// client.Write will then send that same byte message
// through websocket connection

func (hub *Hub) OnConnect(client repository.IClient) {
	fmt.Println("client has connected ", client.GetId())
	hub.Locker.Lock()
	defer hub.Locker.Unlock()
	hub.Clients = append(hub.Clients, client)
}

func (hub *Hub) OnDisconnect(client repository.IClient) {
	fmt.Println("client has disconnected ", client.GetId())
	client.CloseConnection()
	hub.Locker.Lock()
	var newSlice []repository.IClient = make([]repository.IClient, 0)
	defer hub.Locker.Unlock()
	for _, currentClient := range hub.Clients {
		if currentClient.GetId() != client.GetId() {
			newSlice = append(newSlice, currentClient)
		}
	}
	hub.Clients = newSlice
}

func (hub *Hub) GetNumberOfActiveClients() int {
	return len(hub.Clients)
}

func (hub *Hub) RegisterClient(client repository.IClient) {
	hub.Register <- client
}

func (hub *Hub) DeRegisterClient(client repository.IClient) {
	hub.Unregister <- client
}

func (hub *Hub) Run() {
	for {
		fmt.Println("waiting")
		select {
		case client := <-hub.Register:
			fmt.Println("registering client")
			hub.OnConnect(client)
		case client := <-hub.Unregister:
			fmt.Println("disconnecting client")
			hub.OnDisconnect(client)
		}
	}
}

func (hub *Hub) BroadCast(message interface{}, ignoreClient string) {
	data, _ := json.Marshal(message)
	for _, client := range hub.Clients {
		if client.GetId() != ignoreClient {
			client.Send(data)
		}
	}
}

func NewHub() *Hub {
	return &Hub{
		Clients:    make([]repository.IClient, 0),
		Register:   make(chan repository.IClient, 100),
		Unregister: make(chan repository.IClient, 100),
		Locker:     new(sync.Mutex),
		Number:     make(chan int, 1),
	}
}
