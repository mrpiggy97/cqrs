package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/mrpiggy97/cqrs/events"
	"github.com/mrpiggy97/cqrs/models"
	"github.com/mrpiggy97/cqrs/repository"
	"github.com/mrpiggy97/cqrs/search"
)

func onCreatedFeed(message events.CreatedFeedMessage) {
	var newFeed models.Feed = models.Feed{
		Id:          message.Id,
		Title:       message.Title,
		Description: message.Description,
		CreatedAt:   message.CreatedAt,
	}

	var err error = search.IndexFeed(context.Background(), &newFeed)
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
