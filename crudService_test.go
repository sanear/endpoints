package endpoints

import (
	"database/sql"
	"encoding/json"
	_ "github.com/lib/pq"
	"io/ioutil"
	"log"
	"net/http"
	"testing"
	"time"
)

func init() {
	connStr := "postgresql://andrews:test123@localhost:5432/andrews"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Error opening DB connection: %s", err)
	}

	service := NewCrudService("/", ":8080", db)
	handler := NewCrudHandler("test", "andrews",
		map[string]string{
			"ID":   "id",
			"NAME": "name",
		}, db)
	service.AddEndpoint("/test", handler)
	ch := make(chan bool)
	go func(ch chan bool) {
		service.ListenAndServe()
		ch <- true
	}(ch)
	time.Sleep(1)
}

func TestCrudService(t *testing.T) {
	resp, err := http.Get("http://localhost:8080/test?ID=1")
	if err != nil {
		t.Fatalf("Failed to perform http GET: %s", err)
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %s", err)
	}
	if resp.StatusCode != 200 {
		t.Errorf("Got response code %d", resp.StatusCode)
	}

	var result interface{}
	err = json.Unmarshal(b, &result)
	if err != nil {
		t.Fatalf("Failed to parse response body as json: %s", err)
	}
	if r, ok := result.([]interface{}); !ok || len(r) == 0 {
		t.Error("Got back an empty result set. Is that what we expected?")
	}

	t.Logf("Got response code %d and body %s", resp.StatusCode, b)
}

func TestCrudServiceAgain(t *testing.T) {
	resp, err := http.Get("http://localhost:8080/test?ID=1&ID=2")
	if err != nil {
		t.Fatalf("Failed to perform http GET: %s", err)
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %s", err)
	}
	if resp.StatusCode != 200 {
		t.Errorf("Got response code %d", resp.StatusCode)
	}

	var result interface{}
	err = json.Unmarshal(b, &result)
	if err != nil {
		t.Fatalf("Failed to parse response body as json: %s", err)
	}
	if r, ok := result.([]interface{}); !ok || len(r) == 0 {
		t.Error("Got back an empty result set. Is that what we expected?")
	}

	t.Logf("Got response code %d and body %s", resp.StatusCode, b)
}
