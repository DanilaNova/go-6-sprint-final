package server

import (
	"log"
	"net/http"

	"github.com/DanilaNova/go-6-sprint-final/internal/handlers"
)

type Server struct {
	Logger *log.Logger
	Server http.Server
}

// Creates new http server with this logger and port 8080
func Create(logger *log.Logger) *Server {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /", handlers.Root(logger))
	mux.HandleFunc("POST /upload", handlers.Upload(logger))

	return &Server{
		Logger: logger,
		Server: http.Server{
			Addr:    ":8080",
			Handler: mux,
		},
	}
}
