package api

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/ch629/mockservice/pkg/config"
	"github.com/ch629/mockservice/pkg/recorder"
	"github.com/ch629/mockservice/pkg/stub"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type Server interface {
	Start(ctx context.Context, cfg config.API)
}

type server struct {
	router *mux.Router
	log    *zap.Logger

	stubService     stub.Service
	recorderService recorder.Service
}

func New(log *zap.Logger, stubService stub.Service, recorderService recorder.Service) Server {
	s := &server{
		router:          mux.NewRouter(),
		log:             log,
		stubService:     stubService,
		recorderService: recorderService,
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
			s.log.Warn("server shutdown", zap.Error(err))
		}
	}()
	<-ctx.Done()
	// TODO: Deadline on shutdown ctx?
	if err := server.Shutdown(context.Background()); err != nil {
		s.log.Error("failed to shutdown", zap.Error(err))
	}
}
