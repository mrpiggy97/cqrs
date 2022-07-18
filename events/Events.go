package events

import (
	"context"

	"github.com/mrpiggy97/cqrs/models"
	"github.com/nats-io/nats.go"
)

type EventStore interface {
	Close()
	PublishCreatedFeed(cxt context.Context, feed *models.Feed) error
	StartSubscribing(natsMessageHandler nats.MsgHandler) error
}

var eventStore EventStore

func SetEventStore(store EventStore) {
	eventStore = store
}

func Close() {
	eventStore.Close()
}

func PublishCreatedFeed(cxt context.Context, feed *models.Feed) error {
	return eventStore.PublishCreatedFeed(cxt, feed)
}

func StartSubscribing(natsMessageHandler nats.MsgHandler) error {
	return eventStore.StartSubscribing(natsMessageHandler)
}
