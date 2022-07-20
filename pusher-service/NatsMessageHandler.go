package main

import (
	"bytes"
	"encoding/gob"

	"github.com/mrpiggy97/cqrs/events"
	"github.com/mrpiggy97/cqrs/repository"
	"github.com/nats-io/nats.go"
)

func NatsMessageHandler(message *nats.Msg) {
	var feedMessage *events.CreatedFeedMessage = new(events.CreatedFeedMessage)
	var bufferrer *bytes.Buffer = new(bytes.Buffer)
	bufferrer.Write(message.Data)
	gob.NewDecoder(bufferrer).Decode(feedMessage)
	bufferrer = nil
	var createdMessage *CreatedFeedMessage = &CreatedFeedMessage{
		Id:          feedMessage.Id,
		Description: feedMessage.Description,
		CreatedAt:   feedMessage.CreatedAt,
		Type:        "created_feed",
		Title:       feedMessage.Title,
	}
	repository.BroadCast(*createdMessage, "")
}
