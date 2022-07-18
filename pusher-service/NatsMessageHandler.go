package main

import (
	"bytes"
	"encoding/json"

	"github.com/mrpiggy97/cqrs/repository"
	"github.com/nats-io/nats.go"
)

func NatsMessageHandler(message *nats.Msg) {
	var feedMessage *CreatedFeedMessage = new(CreatedFeedMessage)
	var bufferrer *bytes.Buffer = new(bytes.Buffer)
	bufferrer.Write(message.Data)
	json.NewDecoder(bufferrer).Decode(feedMessage)
	bufferrer = nil
	repository.BroadCast(*feedMessage, "")
}
