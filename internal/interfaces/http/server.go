package httpapi

import (
	"log"
	"net/http"
)

type Server struct {
	mux *http.ServeMux
}

func NewServer() *Server {
	mux := http.NewServeMux()

	return &Server{
		mux: mux,
	}
}

func (s *Server) RegisterRoutes(register func(mux *http.ServeMux)) {
	register(s.mux)
}

func (s *Server) Start(addr string) error {
	log.Printf("HTTP server running on %s\n", addr)
	return http.ListenAndServe(addr, s.mux)
}
