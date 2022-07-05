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
)

func newRouter() *mux.Router {
	var newRouter *mux.Router = mux.NewRouter()
	newRouter.HandleFunc("/feeds", CreatedFeedHandler).Methods(http.MethodPost)
	return newRouter
}

type Config struct {
	PostgresDb       string `envconfig:"POSTGRES_DB"`
	PostgresUser     string `envconfig:"POSTGRES_USER"`
	PostgresPassword string `envconfig:"POSTGRES_PASSWORD"`
	NatsAddress      string `envconfig:"NATS_ADDRESS"`
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

	nats, err := events.NewNats(fmt.Sprintf("nats://%s", newConfig.NatsAddress))
	if err != nil {
		log.Fatalf(err.Error())
	}

	events.SetEventStore(nats)
	defer events.Close()
	router := newRouter()
	err = http.ListenAndServe(":8080", router)
	if err != nil {
		log.Fatalf(err.Error())
	}
}
