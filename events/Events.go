package events

import (
	"context"

	"github.com/mrpiggy97/cqrs/models"
)

type EventStore interface {
	Close()
	PublishCreatedFeed(cxt context.Context, feed *models.Feed) error
	SubscribeCreatedFeed(cxt context.Context) (<-chan CreatedFeedMessage, error)
	OnCreateFeed(callback func(CreatedFeedMessage)) error
}

var eventStore EventStore

func Close() {
	eventStore.Close()
}

func PublishCreatedFeed(cxt context.Context, feed *models.Feed) error {
	return eventStore.PublishCreatedFeed(cxt, feed)
}

func SubscribeCreatedFeed(cxt context.Context) (<-chan CreatedFeedMessage, error) {
	return eventStore.SubscribeCreatedFeed(cxt)
}

func OnCreateFeed(callback func(CreatedFeedMessage)) error {
	return eventStore.OnCreateFeed(callback)
}
