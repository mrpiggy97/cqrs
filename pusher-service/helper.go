package main

import "time"

type CreatedFeedMessage struct {
	Type        string    `json:"type"`
	Id          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

func newCreatedFeedMessage(id string, title string, description string, created time.Time) *CreatedFeedMessage {
	return &CreatedFeedMessage{
		Type:        "created_feed",
		Id:          id,
		Title:       title,
		Description: description,
		CreatedAt:   created,
	}
}
