package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/kelseyhightower/envconfig"
	"github.com/mrpiggy97/cqrs/events"
	"github.com/mrpiggy97/cqrs/repository"
)

type Config struct {
	NatsAddress string `envconfig:"NATS_ADDRESS"`
}

func main() {
	var newConfig *Config = new(Config)
	err := envconfig.Process("", newConfig)
	if err != nil {
		log.Fatalf(err.Error())
	}

	var hub *Hub = NewHub()
	repository.SetAppHub(hub)
	nats, err := events.NewNats(fmt.Sprintf("nats://%s", newConfig.NatsAddress))
	if err != nil {
		log.Fatalf(err.Error())
	}
	events.SetEventStore(nats)
	defer events.Close()
	err = nats.OnCreateFeed(func(m events.CreatedFeedMessage) {
		hub.BroadCast(newCreatedFeedMessage(m.Id, m.Title, m.Description, m.CreatedAt), "")
	})

	if err != nil {
		log.Fatalf(err.Error())
	}

	go repository.Run()

	http.HandleFunc("/ws", WebSocketHandler)
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf(err.Error())
	}
}
