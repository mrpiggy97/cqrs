package main

import (
	"fmt"

	"github.com/gorilla/websocket"
	"github.com/mrpiggy97/cqrs/repository"
)

type Client struct {
	Id       string
	Socket   *websocket.Conn
	Outbound chan []byte
}

func (client *Client) Write() {
	for {
		select {
		case message, ok := <-client.Outbound:
			if !ok {
				client.Socket.WriteMessage(websocket.CloseMessage, []byte{})
				repository.DeRegisterClient(client)
				return
			}
			fmt.Println("sending message from ", client.Id)
			var activeClients int = repository.GetNumberOfActiveClients()
			var log string = fmt.Sprintf("number of active clients %d", activeClients)
			fmt.Println(log)
			var sendingErr error = client.Socket.WriteMessage(websocket.TextMessage, message)
			if sendingErr != nil {
				fmt.Println(sendingErr.Error())
				close(client.Outbound)
			}
		}
	}
}

func (client *Client) GetId() string {
	return client.Id
}

func (client *Client) CloseConnection() {
	client.Socket.Close()
}

func (client *Client) Send(data []byte) {
	client.Outbound <- data
}

func NewClient(socket *websocket.Conn) *Client {
	return &Client{
		Socket:   socket,
		Outbound: make(chan []byte, 1),
		Id:       socket.RemoteAddr().String(),
	}
}
