package server

import (
	"context"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
	"github.com/nagdahimanshu/ethereum-block-scanner/internal/logger"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Server struct {
	ctx    context.Context
	server *http.Server
	logger logger.Logger
	health bool
	mu     sync.RWMutex
	Port   string
}

func NewServer(ctx context.Context, logger logger.Logger, port string) *Server {
	return &Server{
		ctx:    ctx,
		logger: logger,
		health: true,
		Port:   port,
	}
}

func (s *Server) Start(ctx context.Context) {
	s.logger.Infof("Starting server at port %s", s.Port)
	go s.startServer()

	s.logger.Infof("HTTP server started at port %s", s.Port)

	<-ctx.Done()
	s.Stop()
}

func (s *Server) startServer() {
	router := mux.NewRouter()

	router.HandleFunc("/health", s.healthHandler)
	router.Handle("/metrics", promhttp.Handler())

	s.server = &http.Server{
		Addr:    "0.0.0.0:" + s.Port,
		Handler: router,
	}

	if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		s.logger.Error("server failed", err)
	}
}

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var err error
	if s.health {
		w.WriteHeader(http.StatusOK)
		_, err = w.Write([]byte(`{"status": "healthy"}`))
	} else {
		w.WriteHeader(http.StatusServiceUnavailable)
		_, err = w.Write([]byte(`{"status": "unhealthy"}`))
	}

	if err != nil {
		s.logger.Errorf("failed to write health response: %v", err)
	}
}

func (s *Server) SetHealth(health bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.health = health
}

func (s *Server) Stop() {
	if s.server != nil {
		if err := s.server.Shutdown(s.ctx); err != nil {
			s.logger.Errorf("server shutdown error: %v", err)
		}
	}
}
