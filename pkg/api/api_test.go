package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"web/pkg/storage"
)

const (
	DBHost     = "localhost"
	DBPort     = "5433"
	DBName     = "postgres"
	DBUser     = "postgres"
	DBPassword = "ZAQzaqzaq97"
)

func TestAPI_posts(t *testing.T) {
	dbase, err := storage.New(fmt.Sprintf("postgres://%s:%s@%s:%s/%s", DBUser, DBPassword, DBHost, DBPort, DBName))
	if err != nil {
		t.Fatal("Не удалось подключиться к базе")
	}
	api := New(dbase)
	req := httptest.NewRequest(http.MethodGet, "/news/5", nil)
	rr := httptest.NewRecorder()
	api.r.ServeHTTP(rr, req)
	if !(rr.Code == http.StatusOK) {
		t.Errorf("Получили %d, а хотели %d", rr.Code, http.StatusOK)
	}
	b, err := io.ReadAll(rr.Body)
	if err != nil {
		t.Fatalf("Не удалось прочитать ответ сервера: %v", err)
	}
	var data []storage.Post
	err = json.Unmarshal(b, &data)
	if err != nil {
		t.Fatalf("Не удалось раскодировать сообщение сервера: %v", err)
	}
	const wantLen = 5
	if len(data) != wantLen {
		t.Fatalf("Получено %d записей, ожидалось %d", len(data), wantLen)
	}
	t.Log(string(b))
}
