package events

import (
	"bytes"
	"context"
	"encoding/gob"

	"github.com/mrpiggy97/cqrs/models"
	natsio "github.com/nats-io/nats.go"
)

type NatsEventStore struct {
	Connection      *natsio.Conn
	FeedCreatedSub  *natsio.Subscription
	FeedCreatedChan chan CreatedFeedMessage
}

func NewNats(url string) (*NatsEventStore, error) {
	conn, err := natsio.Connect(url)
	if err != nil {
		return nil, err
	}

	var newNats *NatsEventStore = &NatsEventStore{
		Connection: conn,
	}

	return newNats, nil
}

func (nats *NatsEventStore) Close() {
	if nats.Connection != nil {
		nats.Connection.Close()
	}
	if nats.FeedCreatedSub != nil {
		nats.FeedCreatedSub.Unsubscribe()
	}
	close(nats.FeedCreatedChan)
}

func (nats *NatsEventStore) encodeMessage(message Message) ([]byte, error) {
	var bufferer *bytes.Buffer = &bytes.Buffer{}
	var err error = gob.NewEncoder(bufferer).Encode(message)
	if err != nil {
		return nil, err
	}
	return bufferer.Bytes(), nil
}

func (nats *NatsEventStore) decodeMessage(data []byte, message interface{}) error {
	var bufferer *bytes.Buffer = &bytes.Buffer{}
	bufferer.Write(data)
	return gob.NewDecoder(bufferer).Decode(message)
}

func (nats *NatsEventStore) PublishCreatedFeed(cxt context.Context, feed *models.Feed) error {
	var message *CreatedFeedMessage = &CreatedFeedMessage{
		Id:          feed.Id,
		Title:       feed.Title,
		Description: feed.Description,
		CreatedAt:   feed.CreatedAt,
	}
	data, err := nats.encodeMessage(message)
	if err != nil {
		return err
	}
	return nats.Connection.Publish(message.Type(), data)
}

func (nats *NatsEventStore) OnCreateFeed(callback func(CreatedFeedMessage)) (err error) {
	var message *CreatedFeedMessage = new(CreatedFeedMessage)
	nats.FeedCreatedSub, err = nats.Connection.Subscribe(message.Type(), func(msg *natsio.Msg) {
		nats.decodeMessage(msg.Data, message)
		callback(*message)
	})
	return
}

func (nats *NatsEventStore) SubscribeCreatedFeed(cxt context.Context) (<-chan CreatedFeedMessage, error) {
	var message *CreatedFeedMessage = new(CreatedFeedMessage)
	nats.FeedCreatedChan = make(chan CreatedFeedMessage, 64)
	var natsioChannel chan *natsio.Msg = make(chan *natsio.Msg, 64)
	var err error
	nats.FeedCreatedSub, err = nats.Connection.ChanSubscribe(message.Type(), natsioChannel)
	if err != nil {
		return nil, err
	}
	go func() {
		for {
			select {
			case m := <-natsioChannel:
				nats.decodeMessage(m.Data, message)
				nats.FeedCreatedChan <- *message
			}
		}
	}()
	return (<-chan CreatedFeedMessage)(nats.FeedCreatedChan), nil
}
