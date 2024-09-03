package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type ApiServer struct {
	addr string
}
type ApiErrType struct {
	Error string `json:"error"`
}

type ApiHealthz struct {
	Status string `json:"status"`
}

func NewApiServer(addr string) *ApiServer {
	return &ApiServer{
		addr: addr,
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
