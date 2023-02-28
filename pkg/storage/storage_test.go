package storage

import (
	"fmt"
	"math/rand"
	"strconv"
	"testing"
	"time"
)

const (
	DBHost     = "localhost"
	DBPort     = "5433"
	DBName     = "postgres"
	DBUser     = "postgres"
	DBPassword = "ZAQzaqzaq97"
)

func TestDB_News(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	posts := []Post{
		{
			Title: "Test Post",
			URL:   strconv.Itoa(rand.Intn(1_000_000_000)),
		},
	}
	db, err := New(fmt.Sprintf("postgres://%s:%s@%s:%s/%s", DBUser, DBPassword, DBHost, DBPort, DBName))
	if err != nil {
		t.Fatal(err)
	}
	err = db.StorePosts(posts)
	if err != nil {
		t.Fatal(err)
	}
	news, err := db.News(2)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", news)
}

func TestNew(t *testing.T) {
	_, err := New(fmt.Sprintf("postgres://%s:%s@%s:%s/%s", DBUser, DBPassword, DBHost, DBPort, DBName))
	if err != nil {
		t.Fatal(err)
	}
}
