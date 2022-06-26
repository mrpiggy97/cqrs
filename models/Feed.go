package models

import "time"

type Feed struct {
	Id          int64  `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	CreatedAt   time.Time
}
