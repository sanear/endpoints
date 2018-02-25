package endpoints

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type CrudHandler struct {
	table, schema string
	fieldToColumn map[string]string
	db            *sql.DB
}

func (h *CrudHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		result, err := h.get(r.URL.Query(), h.db)
		if err != nil {
			log.Printf("Failed to perform SELECT! %s", err)
			failResponse(500, "Server error", &w)
		} else {
			b, _ := json.Marshal(result) // We know the return from get()
			w.WriteHeader(200)
			w.Write(b)
		}
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

// For a SELECT, all given values are for the WHERE clause
func (h *CrudHandler) get(queryStr url.Values, db *sql.DB) ([]map[string]interface{}, error) {
	log.Printf("Handling GET with values: %v", queryStr)
	vals := make([]interface{}, 0, len(queryStr))
	where := make([]string, 0, len(queryStr))
	i := 1
	for k, v := range queryStr {
		if col, ok := h.fieldToColumn[k]; ok {
			switch len(v) {
			case 0:
				log.Printf("Got empty value in query string for key %s. Ignoring", k)
			case 1:
				vals = append(vals, v[0])
				where = append(where, fmt.Sprintf("%s = $%d", col, i))
				i++
			default:
				log.Printf("Got multivalued arg in queryString, %s:%v. Using IN list.", k, v)
				// Unpack into values array
				inList := make([]string, 0, len(v))
				for _, el := range v {
					vals = append(vals, el)
					inList = append(inList, "$"+el)
					i++
				}
				where = append(where, fmt.Sprintf("%s IN (%s)", col, strings.Join(inList, ",")))
			}
		} else {
			log.Printf("Given field %s not recognized. Omitting", k)
		}
	}

	var stmt string
	if len(where) > 0 {
		stmt = fmt.Sprintf("SELECT %s FROM %s.%s WHERE %s",
			strings.Join(h.columnList(), ","), h.schema, h.table, strings.Join(where, " AND "))
	} else {
		stmt = fmt.Sprintf("SELECT %s FROM %s.%s",
			strings.Join(h.columnList(), ","), h.schema, h.table)
	}
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}
	s, err := tx.Prepare(stmt)
	log.Printf("Statement: %s", stmt)
	if err != nil {
		return nil, err
	}
	rows, err := s.Query(vals...)
	defer rows.Close()
	if err != nil {
		return nil, err
	}

	cols, _ := rows.Columns()
	result := make([]map[string]interface{}, 0)
	for i := 0; rows.Next(); i++ {
		columns := make([]interface{}, len(cols))
		colPointers := make([]interface{}, len(cols))
		for i, _ := range columns {
			colPointers[i] = &columns[i]
		}
		err = rows.Scan(colPointers...)
		if err != nil {
			return nil, err
		}
		m := make(map[string]interface{})
		for i, c := range columns {
			m[cols[i]] = c
		}
		result = append(result, m)
	}
	return result, nil
}

func failResponse(code int, reason string, w *http.ResponseWriter) {
	(*w).WriteHeader(code)
	body := map[string]string{
		"result": "failure",
		"reason": reason,
	}
	b, _ := json.Marshal(body) // We know this works, we just made the map
	(*w).Write(b)
	return
}

func (h *CrudHandler) columnList() []string {
	result := make([]string, 0, len(h.fieldToColumn))
	for _, col := range h.fieldToColumn {
		result = append(result, col)
	}
	return result
}
