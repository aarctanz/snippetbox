package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type ApiServer struct {
	addr string
	db   *sql.DB
}
type ApiErrType struct {
	Error string `json:"error"`
}

type ApiHealthz struct {
	Status string `json:"status"`
}

func NewApiServer(addr string, db *sql.DB) *ApiServer {

	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			name VARCHAR(128) NOT NULL,
			email VARCHAR(128) UNIQUE NOT NULL,
			password VARCHAR(256) NOT NULL,
			created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS snippets(
			id SERIAL PRIMARY KEY,
			title VARCHAR(256) NOT NULL,
			is_private BOOLEAN NOT NULL,
			created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
			
			user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE
		)
	`)

	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS files(
			id SERIAL PRIMARY KEY,
			name VARCHAR(128),
			content VARCHAR(1024),
			created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,

			snippet_id INTEGER NOT NULL REFERENCES snippets(id) ON DELETE CASCADE
		)
	`)
	if err != nil {
		log.Fatal(err)
	}

	return &ApiServer{
		addr: addr,
		db:   db,
	}
}

func (s *ApiServer) Run() error {
	router := mux.NewRouter()

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, World\n")
	})

	router.Handle("/healthz", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ResponseWithhJSON(w, 200, ApiHealthz{"ok"})
	})).Methods("GET")

	log.Printf("Server started on port %s\n", s.addr)
	return http.ListenAndServe(s.addr, router)

}

func ResponseWithhJSON(w http.ResponseWriter, code int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	return json.NewEncoder(w).Encode(v)
}

func ResponseWithError(w http.ResponseWriter, code int, err string) error {
	return ResponseWithhJSON(w, code, ApiErrType{Error: err})
}
