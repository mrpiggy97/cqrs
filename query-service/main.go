package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/kelseyhightower/envconfig"
	"github.com/mrpiggy97/cqrs/database"
	"github.com/mrpiggy97/cqrs/events"
	"github.com/mrpiggy97/cqrs/repository"
	"github.com/mrpiggy97/cqrs/search"
)

type Config struct {
	PostgresDb          string `envconfig:"POSTGRES_DB"`
	PostgresUser        string `envconfig:"POSTGRES_USER"`
	PostgresPassword    string `envconfig:"POSTGRES_PASSWORD"`
	NatsAddress         string `envconfig:"NATS_ADDRESS"`
	ElasticSearchAdress string `envconfig:"ELASTICSEARCH_ADDRESS"`
}

func newRouter() *mux.Router {
	var newRouter *mux.Router = mux.NewRouter()
	newRouter.HandleFunc("/feeds", ListFeedsHandler).Methods(http.MethodGet)
	newRouter.HandleFunc("/search", SearchHandler).Methods(http.MethodGet)
	return newRouter
}

func main() {
	var newConfig *Config = new(Config)
	err := envconfig.Process("", newConfig)
	if err != nil {
		log.Fatalf(err.Error())
	}
	var address string = fmt.Sprintf(
		"postgres://%s:%s@postgres/%s?sslmode=disable",
		newConfig.PostgresUser, newConfig.PostgresPassword, newConfig.PostgresDb,
	)

	postgresDb, err := database.NewPostgresRepository(address)
	if err != nil {
		log.Fatalf(err.Error())
	}
	repository.SetRepository(postgresDb)

	var elasticSearchUrl string = fmt.Sprintf("http://%s", newConfig.ElasticSearchAdress)
	elasticSearch, err := search.NewElastic(elasticSearchUrl)
	if err != nil {
		log.Fatalf(err.Error())
	}

	search.SetSearchRepository(elasticSearch)
	defer search.Close()

	nats, err := events.NewNats(fmt.Sprintf("nats://%s", newConfig.NatsAddress))
	if err != nil {
		log.Fatalf(err.Error())
	}

	events.SetEventStore(nats)
	defer events.Close()

	err = nats.StartSubscribing(indexIncomingFeed)
	if err != nil {
		log.Fatalf(err.Error())
	}

	router := newRouter()
	err = http.ListenAndServe(":8080", router)
	if err != nil {
		log.Fatalf(err.Error())
	}
}
