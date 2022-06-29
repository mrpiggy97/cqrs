package search

import (
	"context"

	"github.com/mrpiggy97/cqrs/models"
)

type SearchRepository interface {
	Close()
	IndexFeed(cxt context.Context, feed *models.Feed) error
	SearchFeed(cxt context.Context, query string) ([]models.Feed, error)
}

var repo SearchRepository

func SetSearchRepository(searchRepo SearchRepository) {
	repo = searchRepo
}

func Close() {
	repo.Close()
}

func IndexFeed(cxt context.Context, feed *models.Feed) error {
	return repo.IndexFeed(cxt, feed)
}

func SearchFeed(cxt context.Context, query string) ([]models.Feed, error) {
	return repo.SearchFeed(cxt, query)
}
