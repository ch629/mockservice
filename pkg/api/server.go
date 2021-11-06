package api

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/ch629/mockservice/pkg/config"
	"github.com/ch629/mockservice/pkg/stub"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type Server interface {
	Start(ctx context.Context, cfg config.API)
}

type server struct {
	router *mux.Router
	logger *zap.Logger

	stubService stub.Service
}

func New(logger *zap.Logger, stubService stub.Service) Server {
	s := &server{
		router:      mux.NewRouter(),
		logger:      logger,
		stubService: stubService,
	}
	s.registerRoutes()
	return s
}

func (s *server) Start(ctx context.Context, cfg config.API) {
	server := http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Port),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  10 * time.Second,
		Handler:      s.router,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil {
			s.logger.Warn("server shutdown", zap.Error(err))
		}
	}()
	<-ctx.Done()
	// TODO: Deadline on shutdown ctx?
	if err := server.Shutdown(context.Background()); err != nil {
		s.logger.Error("failed to shutdown", zap.Error(err))
	}
}
