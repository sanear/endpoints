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
	"database/sql"
	"log"
	"net/http"
)

// CrudService exposes a db connection as a REST-ish service,
// backed by a standard server mutex
type CrudService struct {
	// Configurables
	RootPath string
	// TODO: Auth Strategy
	// TODO: Pagination config

	// Internals
	db  *sql.DB
	mux *http.ServeMux // Exposed for convenience
	srv *http.Server
}

func NewCrudService(path, port string, db *sql.DB) CrudService {
	mux := http.NewServeMux()
	return CrudService{
		path,
		db,
		mux,
		&http.Server{Addr: port, Handler: mux},
	}
}

func (s *CrudService) ListenAndServe() {
	defer func() {
		err := s.db.Close()
		if err != nil {
			log.Printf("Unable to close DB connection! %s", err)
		}
	}()
	log.Fatal(s.srv.ListenAndServe())
}

// func (s *CrudService) Close() {
// 	log.Fatal(s.srv.Close())
// }

func (s *CrudService) AddEndpoint(path string, h http.Handler) {
	s.mux.Handle(path, h)
}
