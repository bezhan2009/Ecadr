package server

import (
	"context"
	"net/http"
	"time"
)

type Server struct {
	httpServer *http.Server
}

func (s *Server) Run(port string, handler http.Handler, certFile, keyFile string) error {
	s.httpServer = &http.Server{
		Addr:           ":" + port,
		Handler:        handler,
		MaxHeaderBytes: 1 << 20,
		ReadTimeout:    30 * time.Second,
		WriteTimeout:   30 * time.Second,
	}

	if certFile != "" && keyFile != "" {
		// HTTPS
		return s.httpServer.ListenAndServeTLS(certFile, keyFile)
	}

	// HTTP fallback
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
