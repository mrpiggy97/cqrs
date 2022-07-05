package main

import (
	"net/http"

	"github.com/mrpiggy97/cqrs/repository"
)

func WebSocketHandler(writer http.ResponseWriter, req *http.Request) {
	socket, err := upgrader.Upgrade(writer, req, nil)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	var newClient *Client = NewClient(socket)
	repository.RegisterClient(newClient)
	go newClient.Write()
}
