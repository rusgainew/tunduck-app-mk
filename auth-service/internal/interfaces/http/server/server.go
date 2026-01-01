package server

import (
	"context"
	"log"
	"net/http"
	"time"
)

// HTTPServer wraps the net/http server with sensible timeouts.
type HTTPServer struct {
	server *http.Server
}

// NewHTTPServer creates a configured HTTP server.
func NewHTTPServer(port string, handler http.Handler) *HTTPServer {
	return &HTTPServer{
		server: &http.Server{
			Addr:         ":" + port,
			Handler:      handler,
			ReadTimeout:  15 * time.Second,
			WriteTimeout: 15 * time.Second,
			IdleTimeout:  60 * time.Second,
		},
	}
}

// Start - запустить HTTP сервер
func (s *HTTPServer) Start() error {
	log.Printf("Starting HTTP server on %s\n", s.server.Addr)
	return s.server.ListenAndServe()
}

// Shutdown - корректно остановить сервер
func (s *HTTPServer) Shutdown(ctx context.Context) error {
	log.Println("Shutting down HTTP server...")
	return s.server.Shutdown(ctx)
}
