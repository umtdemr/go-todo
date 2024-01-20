package server

import (
	"encoding/json"
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
	http.ListenAndServe(s.ListenAddr, RequestLoggerMiddleware(s.Router))
}

func DecodeBody(r *http.Request, v interface{}) error {
	decoder := json.NewDecoder(r.Body)

	return decoder.Decode(v)
}
