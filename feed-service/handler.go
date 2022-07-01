package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/mrpiggy97/cqrs/events"
	"github.com/mrpiggy97/cqrs/models"
	"github.com/mrpiggy97/cqrs/repository"
	"github.com/segmentio/ksuid"
)

type createdFeedRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

func CreatedFeedHandler(writer http.ResponseWriter, req *http.Request) {
	var decodedRequest *createdFeedRequest = new(createdFeedRequest)
	var decodingErr error = json.NewDecoder(req.Body).Decode(decodedRequest)
	if decodingErr != nil {
		http.Error(writer, decodingErr.Error(), http.StatusBadRequest)
		return
	}
	var newFeed *models.Feed = &models.Feed{
		Id:          ksuid.New().String(),
		CreatedAt:   time.Now().UTC(),
		Title:       decodedRequest.Title,
		Description: decodedRequest.Description,
	}

	var databaseError error = repository.InsertFeed(
		context.Background(),
		newFeed,
	)
	if databaseError != nil {
		http.Error(writer, databaseError.Error(), http.StatusInternalServerError)
		return
	}

	var eventError error = events.PublishCreatedFeed(context.Background(), newFeed)
	if eventError != nil {
		fmt.Println("failed ot publish created feed")
	}

	writer.Header().Add("content-type", "application/json")
	writer.WriteHeader(http.StatusCreated)
	json.NewEncoder(writer).Encode(newFeed)
}
