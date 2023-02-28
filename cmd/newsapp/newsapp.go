package main

import (
	"encoding/json"
	"fmt"
	"os"

	"log"
	"net/http"
	"time"

	"web/pkg/api"
	"web/pkg/rss"
	"web/pkg/storage"
)

const (
	DBHost     = "localhost"
	DBPort     = "5433"
	DBName     = "postgres"
	DBUser     = "postgres"
	DBPassword = "ZAQzaqzaq97"
)

type config struct {
	Sources []string `json:"rss"`
	Period  int      `json:"period"`
}

func main() {

	db, err := storage.New(fmt.Sprintf("postgres://%s:%s@%s:%s/%s", DBUser, DBPassword, DBHost, DBPort, DBName))
	if err != nil {
		log.Fatal(err)
	}
	api := api.New(db)

	b, err := os.ReadFile("./cmd/newsapp/config.json")
	if err != nil {
		log.Fatal(err)
	}
	var config config
	err = json.Unmarshal(b, &config)
	if err != nil {
		log.Fatal(err)
	}

	chPosts := make(chan []storage.Post)
	chErrs := make(chan error)
	for _, url := range config.Sources {
		go parseURL(url, db, chPosts, chErrs, config.Period)
	}

	go func() {
		for posts := range chPosts {
			db.StorePosts(posts)
		}
	}()

	go func() {
		for err := range chErrs {
			log.Println(err)
		}
	}()

	err = http.ListenAndServe(":80", api.Router())
	if err != nil {
		log.Fatal(err)
	}
}

func parseURL(url string, db *storage.DB, posts chan<- []storage.Post, errs chan<- error, period int) {
	for {
		news, err := rss.Parse(url)
		if err != nil {
			errs <- err
			continue
		}
		posts <- news
		time.Sleep(time.Minute * time.Duration(period))
	}
}
