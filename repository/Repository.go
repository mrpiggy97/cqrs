package repository

import (
	"context"

	"github.com/mrpiggy97/cqrs/models"
)

type Repository interface {
	Close()
	InsertFeed(cxt context.Context, feed *models.Feed) error
	ListFeeds(cxt context.Context) ([]*models.Feed, error)
}

var repository Repository

func SetRepository(repo Repository) {
	repository = repo
}

func Close() {
	repository.Close()
}

func InsertFeed(cxt context.Context, feed *models.Feed) error {
	return repository.InsertFeed(cxt, feed)
}

func ListFeeds(cxt context.Context) ([]*models.Feed, error) {
	return repository.ListFeeds(cxt)
}
