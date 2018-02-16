package endpoints

import (
	"log"
	"net/http"
	"net/url"
	"database/sql"
)

type CrudHandler struct {
	table, schema string
	fieldToColumn map[string]string
	db *sql.DB
}

func (h *CrudHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
	case "PUT":
	case "POST":
	case "DELETE":
	case "OPTIONS":
		w.WriteHeader(200)
	default:
		w.WriteHeader(405)
	}
}

func NewCrudHandler(table, schema string, fieldToColumn map[string]string, db *sql.DB) *CrudHandler {
	return &CrudHandler{
		table,
		schema,
		fieldToColumn,
		db,
	}
}

func get(v url.Values, db *sql.DB) ([]map[string]interface{}, int, error) {
	log.Printf("Handling GET with values: %v", v)
	// TODO
	return make([]map[string]interface{}, 0), 200, nil
}
