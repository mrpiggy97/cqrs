package database

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/lib/pq"
	"github.com/mrpiggy97/cqrs/models"
)

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(url string) (*PostgresRepository, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}
	var repo *PostgresRepository = &PostgresRepository{
		db: db,
	}

	return repo, nil
}

func (repo *PostgresRepository) Close() {
	repo.db.Close()
}

func (repo *PostgresRepository) InsertFeed(cxt context.Context, feed *models.Feed) error {
	_, err := repo.db.ExecContext(
		cxt,
		"INSERT INTO feeds(id,title,description)values($1,$2,$3);", feed.Id, feed.Title, feed.Description,
	)
	return err
}

func (repo *PostgresRepository) ListFeeds(cxt context.Context) ([]*models.Feed, error) {
	var query string = "SELECT * FROM feeds;"
	var feeds []*models.Feed = []*models.Feed{}
	rows, err := repo.db.QueryContext(cxt, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var id *string = new(string)
		var title *string = new(string)
		var description *string = new(string)
		var createdAt *time.Time = new(time.Time)
		var scanningErr error = rows.Scan(id, title, description, createdAt)
		if scanningErr != nil {
			return nil, scanningErr
		}
		var newFeed *models.Feed = &models.Feed{
			Id:          *id,
			Title:       *title,
			Description: *description,
			CreatedAt:   *createdAt,
		}
		feeds = append(feeds, newFeed)
	}
	return feeds, nil
}
