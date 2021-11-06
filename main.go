package main

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/ch629/mockservice/pkg/api"
	"github.com/ch629/mockservice/pkg/config"
	"github.com/ch629/mockservice/pkg/stub"
	"go.uber.org/zap"
)

func main() {
	devCfg := zap.NewDevelopmentConfig()
	devCfg.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	log, _ := devCfg.Build()
	zap.ReplaceGlobals(log)
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()
	api := api.New(log, stub.NewService(log))
	log.Info("starting HTTP server")
	api.Start(ctx, config.API{
		Port: 8080,
	})
}
