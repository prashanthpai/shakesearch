package server

import (
	"context"
	"log"
	"net"
	"net/http"
	"time"
)

type Config struct {
	Addr string
}

type Server struct {
	srv *http.Server
}

func New(cfg *Config, handler http.Handler) *Server {
	return &Server{
		srv: &http.Server{
			Addr:         cfg.Addr,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  120 * time.Second,
			Handler:      handler,
		},
	}
}

func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.srv.Addr)
	if err != nil {
		return err
	}

	go func() {
		if err := s.srv.Serve(ln); err != http.ErrServerClosed {
			log.Fatalf("srv.Serve() failed: %s", err.Error())
		}
	}()
	log.Printf("server started; listening at: %s", s.srv.Addr)

	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}
