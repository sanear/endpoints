// Package endpoints provides wrapper structs for standard net/http hooks
// Goal is to allow simple creation of CRUD webservices against SQL dbs

// TODO: Restructuring, probably 'handlers.go' 'service.go' etc.

// Huh, I have less of a clear idea of what to do here than I wanted.
// As I'm remembering what I did at Helix, it's basically all struct
// wrappers for existing stuff - ServeMux would have done the trick for
// that Webservice type that had "AddEndpoint" and all that. Reinvented
// the wheel somewhat.

// Maybe all I need here is a set of standard Handlers that hook into
// the db interfaces, and a struct to hang them all off of...

package endpoints

import (
	"net/http"
	"database/sql"
	"log"
)

// CrudService exposes a db connection as a REST-ish service,
// backed by a standard server mutex
type CrudService struct {
	// Configurables
	RootPath string
	Port string
	// TODO: Auth Strategy
	// TODO: Pagination config

	// Internals
	paths map[string]http.Handler
	db *sql.DB
	mux *http.ServeMux
}

func NewCrudService(path, port string, db *sql.DB) CrudService {
	return CrudService {
		path,
		port,
		map[string]http.Handler{},
		db,
		http.NewServeMux(),
	}
}

func (s *CrudService) ListenAndServe() {
	log.Fatal(http.ListenAndServe(s.Port, s.mux))
}

func (s *CrudService) AddEndpoint(path string, f http.HandlerFunc) {
	s.mux.Handle(path, f)
}
