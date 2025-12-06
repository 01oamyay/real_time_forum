package server

import (
	"log"
	"net/http"

	"rlf/pkg/config"
)

// Server wraps the standard http.Server to keep configuration centralized.
type Server struct {
	httpServer *http.Server
}

// Run starts listening on the configured port using the provided handler.
func (s *Server) Run(c *config.API, handler http.Handler) error {
	s.httpServer = &http.Server{
		Addr:    ":" + c.Port,
		Handler: handler,
	}
	log.Printf("\033[32mServer is running...🚀\nLink: 🌐 http://%s%s", c.Host, s.httpServer.Addr)
	return s.httpServer.ListenAndServe()
}
