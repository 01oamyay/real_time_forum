package server

import (
	"log"
	"net/http"

	"rlf/pkg/config"
)

type Server struct {
	httpServer *http.Server
}

func (s *Server) Run(c *config.API, handler http.Handler) error {
	s.httpServer = &http.Server{
		Addr:    ":" + c.Port,
		Handler: handler,
	}
	log.Printf("\033[32mServer is running...🚀\nLink: 🌐 http://%s%s", c.Host, s.httpServer.Addr)
	return s.httpServer.ListenAndServe()
}
