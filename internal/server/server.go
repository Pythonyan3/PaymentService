package server

import (
	"net/http"
)

type Server struct {
	httpServer *http.Server
}

func NewServer(port string, handler http.Handler) *Server {
	return &Server{httpServer: &http.Server{
		Addr:    ":" + port,
		Handler: handler,
	}}
}

func (server *Server) Run() error {
	return server.httpServer.ListenAndServe()
}
