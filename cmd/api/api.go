package api

import (
	"github.com/gorilla/mux"
)

type APIServer struct {
	addr string
}

func NewAPIServer(addr string) *APIServer {
	return &APIServer{
		addr: addr,
	}
}

func (s *APIServer) Run() {
	router := mux.NewRouter()

	getRouter := router.Methods("GET").Subrouter()
	getRouter.HandleFunc("/login", nil)
}
