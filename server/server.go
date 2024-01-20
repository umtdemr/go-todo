package server

import (
	"github.com/gorilla/mux"
	"net/http"
)

type APIServer struct {
	ListenAddr string
	Router     *mux.Router
}

func NewAPIServer(listenAddr string) *APIServer {
	return &APIServer{ListenAddr: listenAddr, Router: mux.NewRouter()}
}

func (s *APIServer) Run() {
	http.ListenAndServe(s.ListenAddr, requestLogger(s.Router))
}
