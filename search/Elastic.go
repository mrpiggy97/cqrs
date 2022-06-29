package search

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/mrpiggy97/cqrs/models"
)

type ElasticSearchRepository struct {
	client *elasticsearch.Client
}

func NewElastic(url string) (*ElasticSearchRepository, error) {
	client, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{url},
	})

	if err != nil {
		return nil, err
	}

	var newRepo *ElasticSearchRepository = &ElasticSearchRepository{
		client: client,
	}
	return newRepo, nil
}

func (elasticRepo *ElasticSearchRepository) Close() {
	//
}

func (elasticRepo *ElasticSearchRepository) IndexFeed(cxt context.Context, feed *models.Feed) error {
	body, _ := json.Marshal(feed)
	var idAsString string = fmt.Sprintf("%d", feed.Id)
	_, err := elasticRepo.client.Index(
		"feeds",
		bytes.NewReader(body),
		elasticRepo.client.Index.WithDocumentID(idAsString),
		elasticRepo.client.Index.WithContext(cxt),
		elasticRepo.client.Index.WithRefresh("wait_for"),
	)
	return err
}

func (elasticRepo *ElasticSearchRepository) SearchFeed(cxt context.Context, query string) ([]models.Feed, error) {
	var bufferer *bytes.Buffer = &bytes.Buffer{}
	var searchQuery map[string]interface{} = make(map[string]interface{})
	searchQuery[query] = map[string]interface{}{
		"multi_match": map[string]interface{}{
			"query":            query,
			"fields":           []string{"title", "description"},
			"fuzziness":        3,
			"cutoff_frequency": 0.0001,
		},
	}

	var err error = json.NewEncoder(bufferer).Encode(searchQuery)
	if err != nil {
		return nil, err
	}
	res, err := elasticRepo.client.Search(
		elasticRepo.client.Search.WithContext(cxt),
		elasticRepo.client.Search.WithIndex("feeds"),
		elasticRepo.client.Search.WithBody(bufferer),
		elasticRepo.client.Search.WithTrackTotalHits(true),
	)

	if err != nil {
		return nil, err
	}
	defer func() {
		var closingErr error = res.Body.Close()
		panic(closingErr)
	}()

	if res.IsError() {
		return nil, errors.New(res.String())
	}

	var eRes map[string]interface{}
	err = json.NewDecoder(res.Body).Decode(eRes)
	if err != nil {
		return nil, err
	}

	var feeds []models.Feed
	for _, hit := range eRes["hits"].(map[string]interface{})["hits"].([]interface{}) {
		source := hit.(map[string]interface{})["_source"]
		jsonSource, marshalingErr := json.Marshal(source)
		if marshalingErr != nil {
			return nil, marshalingErr
		}
		var feed *models.Feed = new(models.Feed)
		var unmarshalingErr error = json.Unmarshal(jsonSource, feed)
		if unmarshalingErr == nil {
			feeds = append(feeds, *feed)
		}
	}

	return feeds, nil
}
