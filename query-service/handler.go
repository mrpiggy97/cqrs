package main

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/mrpiggy97/cqrs/models"
	"github.com/mrpiggy97/cqrs/repository"
	"github.com/mrpiggy97/cqrs/search"
	"github.com/nats-io/nats.go"
)

// creates a new models.Feed object from message.Data
// and then indexes that model
func indexIncomingFeed(message *nats.Msg) {
	var newFeed *models.Feed = new(models.Feed)
	var bufferer *bytes.Buffer = new(bytes.Buffer)
	bufferer.Write(message.Data)
	json.NewDecoder(bufferer).Decode(newFeed)

	var err error = search.IndexFeed(context.Background(), newFeed)
	if err != nil {
		log.Printf("failed to index feed: %s", err.Error())
	}
}

func ListFeedsHandler(writer http.ResponseWriter, req *http.Request) {
	var cxt context.Context = req.Context()
	var err error
	feeds, err := repository.ListFeeds(cxt)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	writer.Header().Add("Content-type", "application/json")
	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(feeds)
}

func SearchHandler(writer http.ResponseWriter, req *http.Request) {
	var cxt context.Context = req.Context()
	var err error
	query := req.URL.Query().Get("q")
	if len(query) == 0 {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	feeds, err := search.SearchFeed(cxt, query)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	writer.Header().Add("Content-type", "application/json")
	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(feeds)
}
